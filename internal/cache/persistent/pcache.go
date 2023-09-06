package persistent

import (
	"encoding/gob"
	"io"
	"io/fs"
	"os"
	"sync"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/zeebo/xxh3"
)

type PCache interface {
	Get(string) (struct{}, bool, error)
	Set(string, struct{}) error
	Delete(string) error
	Close() error
}

type Shard interface {
	Get(string) (struct{}, bool, error)
	Set(string, struct{}) error
	Delete(string) error
	Close() error
}

var _ PCache = (*pcache)(nil)
var _ Shard = (*shard)(nil)

type pcache struct {
	shards [slen]Shard
}

type shard struct {
	path string
	dl   int
	l    int
	mu   sync.Mutex
	perm fs.FileMode
}

const (
	// slen is shards length.
	slen = 512
	// slen = 4096
	// mask is slen-1 Hex value.
	mask = 0x1FF
	// mask = 0xFFF.
)

func NewPCache(basePath string) (PCache, error) {
	var shards [slen]Shard
	for i := range shards {
		s, err := newShard(basePath)
		if err != nil {
			return nil, err
		}
		shards[i] = s
	}
	return &pcache{
		shards: shards,
	}, nil
}

// New returns the pcache that satisfies the PCache interface.
func (p *pcache) Get(key string) (struct{}, bool, error) {
	data, ok, err := p.shards[getShardID(key)].Get(key)
	if err != nil {
		return data, false, err
	}
	if !ok {
		return data, false, nil
	}

	return data, true, nil
}

func (p *pcache) Set(key string, data struct{}) error {
	return p.shards[getShardID(key)].Set(key, data)
}

func (p *pcache) Delete(key string) error {
	return p.shards[getShardID(key)].Delete(key)
}

func (p *pcache) Close() error {
	for _, s := range p.shards {
		err := s.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func newShard(basePath string) (*shard, error) {
	f, err := os.CreateTemp(basePath, "pcache-*")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return &shard{
		perm: 0600,
		path: f.Name(),
	}, nil
}

func (s *shard) Get(key string) (data struct{}, ok bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := file.Open(s.path, os.O_RDWR, s.perm)
	if err != nil {
		return
	}
	defer f.Close()

	m := make(map[string]struct{}, s.l)
	err = gob.NewDecoder(f).Decode(&m)
	if err != nil {
		// empty shard file returns EOF
		if errors.Is(err, io.EOF) {
			return data, false, nil
		}
		return data, false, err
	}

	data, ok = m[key]

	return data, ok, nil
}

func (s *shard) Set(key string, data struct{}) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := file.Open(s.path, os.O_RDWR, s.perm)
	if err != nil {
		return err
	}
	defer f.Close()

	m := make(map[string]struct{}, s.l)
	if s.dl != 0 {
		err = gob.NewDecoder(f).Decode(&m)
		if err != nil {
			return err
		}
	}

	m[key] = data

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	err = gob.NewEncoder(f).Encode(m)
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	s.dl = int(fi.Size())
	s.l++

	return f.Sync()
}

func (s *shard) Delete(key string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := file.Open(s.path, os.O_RDWR, s.perm)
	if err != nil {
		return err
	}
	defer f.Close()

	m := make(map[string]struct{}, s.l)
	err = gob.NewDecoder(f).Decode(&m)
	if err != nil {
		return
	}

	delete(m, key)

	// Write the updated data to the file
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	err = gob.NewEncoder(f).Encode(m)
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	s.dl = int(fi.Size())
	s.l--

	return f.Sync()
}

func (s *shard) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.Remove(s.path); err != nil {
		return err
	}
	return nil
}

func getShardID(key string) (id uint64) {
	if len(key) > 128 {
		return xxh3.HashString(key[:128]) & mask
	}
	return xxh3.HashString(key) & mask
}

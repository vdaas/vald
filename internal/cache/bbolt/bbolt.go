package bbolt

import (
	"fmt"
	"os"

	"github.com/vdaas/vald/internal/errors"
	bolt "go.etcd.io/bbolt"
)

type Bbolt struct {
	db   *bolt.DB
	file string
}

const bucket = "vald-bbolt-bucket"

func New(filepath string) (*Bbolt, error) {
	// TODO: 初期化をここでするか、DIするか。ライフタイムを管理するのだるいのでDIの方がいいかも
	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		return nil
	})
	return &Bbolt{
		db:   db,
		file: filepath,
	}, nil
}

func (b *Bbolt) Set(key string, val []byte) error {
	if err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(val)
		err := b.Put([]byte(key), nil)
		return err
	}); err != nil {
		return err
	}

	return nil
}

func (b *Bbolt) Get(key string) ([]byte, bool, error) {
	var val []byte
	if err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		copy(val, b.Get([]byte(key)))
		return nil
	}); err != nil {
		return nil, false, err
	}

	if val == nil {
		return nil, false, nil
	}

	return val, true, nil
}

func (b *Bbolt) Close() (err error) {
	if cerr := b.db.Close(); cerr != nil {
		err = cerr
	}

	if rerr := os.RemoveAll(b.file); rerr != nil {
		err = errors.Wrap(rerr, err.Error())
	}

	return err
}


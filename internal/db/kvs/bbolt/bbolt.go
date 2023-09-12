package bbolt

import (
	"fmt"
	"os"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
	bolt "go.etcd.io/bbolt"
)

type Bbolt struct {
	db     *bolt.DB
	file   string
	bucket string
}

const defaultBucket = "vald-bbolt-bucket"

// New returns a new Bbolt instance.
// If file does not exist, it creates a new file. If bucket is empty, it uses default_bucket.
// If opts is nil, it uses default options.
func New(file, bucket string, opts *bolt.Options) (*Bbolt, error) {
	db, err := bolt.Open(file, 0600, opts)
	if err != nil {
		return nil, err
	}

	if bucket == "" {
		bucket = defaultBucket
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		return nil
	})
	return &Bbolt{
		db:     db,
		file:   file,
		bucket: bucket,
	}, nil
}

func (b *Bbolt) Set(key, val []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		err := b.Put(key, val)
		return err
	})
}

func (b *Bbolt) Get(key []byte) ([]byte, bool, error) {
	var val []byte
	if err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		ret := b.Get(key)
		if ret == nil {
			// key not found. just return without copying anything to val
			return nil
		}

		// key found. copy the value to val because ret is only valid in this scope
		val = make([]byte, len(ret))
		copy(val, ret)
		return nil
	}); err != nil {
		return nil, false, err
	}

	if val == nil {
		return nil, false, nil
	}

	return val, true, nil
}

// AsyncSet sets the key and value asynchronously for better write performance.
// It accumulates the keys and values until the batch size is reached or the timeout comes, then
// writes them all at once. Wait for the errgroup to make sure all the batches finished if required.
func (b *Bbolt) AsyncSet(eg *errgroup.Group, key, val []byte) error {
	if eg == nil {
		return errors.ErrNilErrGroup
	}
	(*eg).Go(func() error {
		b.db.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(b.bucket))
			err := b.Put(key, val)
			return err
		})
		return nil
	})

	return nil
}

// Close closes the database and removes the file if remove is true.
func (b *Bbolt) Close(remove bool) (err error) {
	if cerr := b.db.Close(); cerr != nil {
		err = cerr
	}

	if remove {
		if rerr := os.RemoveAll(b.file); rerr != nil {
			err = errors.Wrap(rerr, err.Error())
		}
	}

	return err
}

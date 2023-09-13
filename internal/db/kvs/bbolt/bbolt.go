// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package bbolt

import (
	"fmt"
	"os"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
	bolt "go.etcd.io/bbolt"
)

type Bbolt interface {
	Set(key, val []byte) error
	Get(key []byte) ([]byte, bool, error)
	AsyncSet(eg errgroup.Group, key, val []byte) error
	Close(remove bool) error
}

type bbolt struct {
	db     *bolt.DB
	file   string
	bucket []byte
}

const defaultBucket = "vald-bbolt-bucket"

// New returns a new Bbolt instance.
// If file does not exist, it creates a new file. If bucket is empty, it uses default_bucket.
// If opts is nil, it uses default options.
func New(file, bucket string, opts *bolt.Options) (Bbolt, error) {
	db, err := bolt.Open(file, 0o600, opts)
	if err != nil {
		return nil, err
	}

	bk := []byte(defaultBucket)
	if bucket != "" {
		bk = []byte(bucket)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bk)
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		return nil
	})
	return &bbolt{
		db:     db,
		file:   file,
		bucket: bk,
	}, nil
}

func (b *bbolt) Set(key, val []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(b.bucket)
		err := b.Put(key, val)
		return err
	})
}

func (b *bbolt) Get(key []byte) (val []byte, ok bool, err error) {
	if err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(b.bucket)
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
func (b *bbolt) AsyncSet(eg errgroup.Group, key, val []byte) error {
	eg.Go(func() error {
		b.db.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket(b.bucket)
			err := b.Put(key, val)
			return err
		})
		return nil
	})

	return nil
}

// Close closes the database and removes the file if remove is true.
func (b *bbolt) Close(remove bool) (err error) {
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

// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package pogreb

import (
	"context"
	"os"
	"reflect"

	"github.com/akrylysov/pogreb"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
)

// Pogreb represents an interface for operating the pogreb database.
type Pogreb interface {
	Set(key string, val []byte) error
	Get(key string) ([]byte, bool, error)
	Delete(key string) error
	Range(ctx context.Context, f func(key string, val []byte) bool) error
	Len() uint32
	Close(remove bool) error
}

type db struct {
	db   *pogreb.DB
	opts *pogreb.Options
	path string
}

// New returns a new pogreb instance.
// If the directory path does not exist, it creates a directory for database.
// If opts is nil, it uses default options.
func New(opts ...Option) (_ Pogreb, err error) {
	db := new(db)
	for _, opt := range append(deafultOpts, opts...) {
		if err := opt(db); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	// If db.opts is nil, an default value is used.
	db.db, err = pogreb.Open(db.path, db.opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Set sets the value for a key.
func (d *db) Set(key string, val []byte) error {
	return d.db.Put(conv.Atob(key), val)
}

// Get returns the value stored in the database for a key.
// The ok result indicates whether value was found in the database.
func (d *db) Get(key string) ([]byte, bool, error) {
	val, err := d.db.Get(conv.Atob(key))
	if err != nil {
		return nil, false, err
	}
	// If val is nil, it means that there is no value associated with key, so false is returned.
	if val == nil {
		return nil, false, nil
	}
	return val, true, nil
}

// Delete deletes the given key from the database.
func (d *db) Delete(key string) error {
	// NOTE: Even if the key does not exist in the database, no error occurs.
	// Depending on the future use case, it may be necessary to check for the existence of the key, in which case the `Has` method can be used.
	return d.db.Delete(conv.Atob(key))
}

// Range calls f sequentially for each key and value present in the database.
// If f returns false, range stops the iteration.
func (d *db) Range(ctx context.Context, f func(key string, val []byte) bool) error {
	it := d.db.Items()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			key, val, err := it.Next()
			if err != nil {
				if errors.Is(err, pogreb.ErrIterationDone) {
					return nil
				}
				return err
			}
			if !f(conv.Btoa(key), val) {
				return nil
			}
		}
	}
}

// Len returns the number of keys in the DB.
func (d *db) Len() uint32 {
	return d.db.Count()
}

// Close closes the database and removes the file if remove is true.
func (d *db) Close(remove bool) (err error) {
	if serr := d.db.Sync(); serr != nil {
		err = serr
	}
	if cerr := d.db.Close(); cerr != nil {
		err = errors.Join(err, cerr)
	}
	if remove {
		if rerr := os.RemoveAll(d.path); rerr != nil {
			err = errors.Join(err, rerr)
		}
	}
	return err
}

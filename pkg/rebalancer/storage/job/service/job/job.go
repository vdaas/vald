//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package service manages the main logic of server.
package job

import (
	"archive/tar"
	"context"
	"encoding/gob"
	"io"
	"reflect"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	ctxio "github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/rebalancer/storage/job/service/storage"
)

type Rebalancer interface {
	Start(context.Context) (<-chan error, error)
}

type rebalancer struct {
	eg              errgroup.Group
	targetAgentName string
	rate            float64
	gatewayHost     string
	gatewayPort     int
	storage         storage.Storage
}

const (
	kvsFileName = "ngt-meta.kvsdb"
)

func New(opts ...Option) (dsc Rebalancer, err error) {
	r := new(rebalancer)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return r, nil
}

func (r *rebalancer) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 2)

	pr, pw := io.Pipe()
	defer pr.Close()

	r.eg.Go(func() error {
		// Download tar gz file
		r.eg.Go(safety.RecoverFunc(func() (err error) {
			defer pw.Close()

			// TODO consider the error handling (why we need channel?)
			defer func() {
				if err != nil {
					select {
					case <-ctx.Done():
						ech <- errors.Wrap(err, ctx.Err().Error())
					case ech <- err:
					}
				}
			}()

			sr, err := r.storage.Reader(ctx)
			if err != nil {
				return err
			}

			sr, err = ctxio.NewReadCloserWithContext(ctx, sr)
			if err != nil {
				return err
			}
			defer func() {
				err = sr.Close()
				if err != nil {
					log.Errorf("error on closing blob-storage reader: %s", err)
				}
			}()

			_, err = io.Copy(pw, sr)
			if err != nil {
				return err
			}
			return nil
		}))

		// Unpack tar.gz file and Decode kvsdb file to get the vector ids
		idm, err := r.loadKVS(ctx, pr)
		if err != nil {
			// TODO: should we return here?
			select {
			case <-ctx.Done():
				// loadKVSでcontext.Errが返ってきたら重複してwrapされるので別途考えた方がいいかもしれない
				ech <- errors.Wrap(err, ctx.Err().Error())
			case ech <- err:
			}
		}

		// Calculate to process data from the above data
		amntData := r.rate * float64(len(idm))

		// Rebalance
		data := map[string][]float64
		// TODO: get gateway address and set option
		// e.g. https://github.com/vdaas/vald/blob/master/charts/vald/templates/gateway/backup/configmap.yaml#L54
		client, err := vald.New()
		if err != nil {
			return err
		}
		for i, id := range idm {
			// get vecotr by id
			// apppned data
			if amntData - 1 == i {
				break
			}
		}
		// request multi update using v1 client

		return nil
	})

	return ech, nil
}

func (r *rebalancer) loadKVS(ctx context.Context, reader io.Reader) (map[string]uint32, error) {
	tr := tar.NewReader(reader)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var header *tar.Header
		header, err := tr.Next()
		// TODO: error handling
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}

			return nil, err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			if header.Name != kvsFileName {
				continue
			}

			gob.Register(map[string]uint32{})
			idm := make(map[string]uint32)

			err = gob.NewDecoder(tr).Decode(&idm)
			if err != nil {
				return nil, err
			}

			return idm, nil
		default:
			// TODO: define in errors package
			return nil, errors.New("invalid file type")
		}

	}

	// TODO; define in errors package
	return nil, errors.New("kvsdb file not found")
}

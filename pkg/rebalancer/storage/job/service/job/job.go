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
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
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
	storage         storage.Storage
	kvsdbStorage    storage.Storage
	client          vald.Client
	parallelism     int
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
	ech := make(chan error, 3)
	cech, err := r.client.Start(ctx)
	if err != nil {
		log.Errorf("[job debug] failed start connection monitor")
		return nil, err
	}
	r.eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case ech <- <-cech:
			}
		}
	})
	r.eg.Go(func() (err error) {
		log.Infof("[job debug] Start rebalance job service: %s", r.targetAgentName)

		pr, pw := io.Pipe()
		defer pr.Close()
		defer func() {
			p, perr := os.FindProcess(os.Getpid())
			if perr != nil {
				log.Error(perr)
				return
			}
			if err != nil {
				p.Signal(syscall.SIGKILL) // TODO: #403
			} else {
				p.Signal(syscall.SIGTERM)
			}
		}()

		// Download tar gz file
		log.Info("[job debug] download s3 backup file")
		r.eg.Go(safety.RecoverFunc(func() (err error) {
			defer pw.Close()

			defer func() {
				if err != nil {
					select {
					case <-ctx.Done():
						ech <- errors.Wrap(err, ctx.Err().Error())
					case ech <- err:
					}
				}
			}()

			log.Info("[job debug] read buffer from download s3 backup file")
			sr, err := r.kvsdbStorage.Reader(ctx)
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
		log.Info("[job debug] unpack tar.gz file and decode kvsdb file")
		idm, err := r.loadKVS(ctx, pr)
		log.Infof("[job debug] finish unpack: len kvs: %d, err: %#v ", len(idm), err)
		if err != nil {
			select {
			case <-ctx.Done():
				if err != context.Canceled {
					ech <- errors.Wrap(err, ctx.Err().Error())
				} else {
					ech <- err
				}
			case ech <- err:
			}
			return err
		}

		// Calculate to process data from the above data
		amntData := int(r.rate * float64(len(idm)))
		if amntData == 0 {
			log.Info("[job] no data to rebalance")
			return nil
		}
		var cnt int64 = 0

		log.Infof("Start rebalance data: %d", amntData)
		var mu sync.Mutex

		eg, egctx := errgroup.New(ctx)
		eg.Limitation(r.parallelism)
		for id, _ := range idm {
			id := id
			select {
			case <-egctx.Done():
				break
			default:
			}
			eg.Go(func() error {
				// get vecotr by id
				log.Infof("[job debug] Get object data: %s", id)
				vec, gerr := r.client.GetObject(egctx, &payload.Object_VectorRequest{
					Id: &payload.Object_ID{
						Id: id,
					},
				})
				if gerr != nil {
					log.Errorf("failed to send GetObject request, uuid: %s, err: %s", id, gerr.Error())
					mu.Lock()
					err = errors.Wrap(err, gerr.Error())
					mu.Unlock()
					return nil
				}

				// update data
				// TODO: use stream or upsert?
				log.Infof("[job debug] Update object data: %s", id)
				_, uerr := r.client.Update(egctx, &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     vec.GetId(),
						Vector: vec.GetVector(),
					},
				})
				if uerr != nil {
					log.Errorf("failed to send Update request, uuid: %s, err: %s", id, uerr.Error())
					mu.Lock()
					err = errors.Wrap(err, uerr.Error())
					mu.Unlock()
					return nil
				}

				n := atomic.AddInt64(&cnt, 1)
				log.Infof("[job debug] Success Rebalance data: success amount data = %d", n)
				return nil
			})

			if amntData--; amntData == 0 {
				log.Infof("[job debug] Finish Rebalance data: success amount data = %d", atomic.LoadInt64(&cnt))
				break
			}
		}
		eg.Wait()

		if err != nil {
			log.Errorf("failed to rebalance data: %s", err.Error())
			return err
		}

		// rename backup file
		if r.rate == float64(1) {
			log.Info("Start create backup file")
			if err := renameBackupFile(ctx, r.storage); err != nil {
				log.Error("Failed to rename full backup file: %s", err.Error())
				return err
			}
			if err := renameBackupFile(ctx, r.kvsdbStorage); err != nil {
				log.Error("Failed to rename kvsdb backup file: %s", err.Error())
				return err
			}
			log.Info("Finish delete original backup file")
		}

		// request multi update using v1 client
		log.Infof("Finish rebalance data: %d, remaining data: %d", cnt, amntData)

		return nil
	})

	return ech, nil
}

func renameBackupFile(ctx context.Context, storage storage.Storage) (err error) {
	err = storage.Backup(ctx)
	if err != nil {
		log.Errorf("failed to create backup file: %s", err.Error())
		return err
	}
	log.Info("Finish create backup file")

	log.Info("Start delete original backup file")
	err = storage.Delete(ctx)
	if err != nil {
		log.Errorf("failed to delete original backup file: %s", err.Error())
		return err
	}
	return nil
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
		if err != nil {
			log.Debugf("[job debug] err: %s", err)
			if err == io.EOF {
				// TODO; define in errors package (after controller PR merged)
				return nil, errors.New("kvsdb file not found")
			}
			return nil, err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			log.Infof("[job debug] header.Name: %s, kvsFileName: %s", header.Name, kvsFileName)
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
		}
	}
}

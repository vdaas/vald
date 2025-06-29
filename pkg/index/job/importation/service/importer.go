// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package service

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/db/kvs/pogreb"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"google.golang.org/grpc"
)

const (
	apiName        = "vald/index/job/import"
	grpcMethodName = "vald.v1.StreamUpsert/" + vald.StreamUpsertRPCName
)

// Importer represents an interface for importing.
type Importer interface {
	StartClient(ctx context.Context) (<-chan error, error)
	Start(ctx context.Context) error
	PreStop(ctx context.Context) error
}

type importer struct {
	eg           errgroup.Group
	gateway      vc.Client
	storedVector pogreb.DB

	streamListConcurrency        int
	backgroundSyncInterval       time.Duration
	backgroundCompactionInterval time.Duration
	indexPath                    string
	forceUpdate                  bool
}

// New returns Importer object if no error occurs.
func New(opts ...Option) (Importer, error) {
	i := new(importer)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(i); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(err)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}

	db, err := pogreb.New(pogreb.WithPath(i.indexPath))
	if err != nil {
		log.Errorf("failed to open checked List kvs DB %s", i.indexPath)
		return nil, err
	}
	i.storedVector = db
	return i, nil
}

// StartClient starts the gRPC client.
func (i *importer) StartClient(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 1)
	gch, err := i.gateway.Start(ctx)
	if err != nil {
		return nil, err
	}
	i.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.Canceled) {
					log.Warn("context canceled when starting client")
					return ctx.Err()
				}
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					log.Warn("context deadline exceeded when starting client")
					return ctx.Err()
				}
				return ctx.Err()
			case err = <-gch:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ech <- err:
				}
			}
		}
	}))
	return ech, nil
}

func (i *importer) Start(ctx context.Context) error {
	err := i.doImportIndex(ctx)
	return err
}

func (i *importer) doImportIndex(
	ctx context.Context,
) (errs error) {
	log.Info("starting doImportIndex")
	ctx, span := trace.StartSpan(igrpc.WrapGRPCMethod(ctx, grpcMethodName), apiName+"/service/index.doImportIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	eg, egctx := errgroup.WithContext(ctx)
	eg.SetLimit(i.streamListConcurrency)

	vch := make(chan *payload.Object_Vector, 5*5)
	ctx, cancel := context.WithCancel(egctx)

	eg.Go(safety.RecoverFunc(func() (err error) {
		defer cancel()
		defer close(vch)
		i.storedVector.Range(egctx, func(key string, value []byte) bool {
			select {
			case <-egctx.Done():
				return false
			case vch <- func() *payload.Object_Vector {
				objVec := new(payload.Object_Vector)
				if err := objVec.UnmarshalVT(value); err != nil {
					log.Error("failed to Unmarshal proto to payload.Object_Vector: %v", err)
					return nil
				}
				return objVec
			}():
			}
			return true
		})
		return nil
	}))

	i.gateway.GRPCClient().RangeConcurrent(egctx, 1, func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		cs, err := vc.NewValdClient(conn).StreamUpsert(egctx, copts...)
		if err != nil {
			log.Errorf("failed to create stream upsert client: %v", err)
			return nil
		}

		err = igrpc.BidirectionalStreamClient(cs, 1, func() (res *payload.Upsert_Request, ok bool) {
			select {
			case <-ctx.Done():
				return nil, false
			case vec, ok := <-vch:
				fmt.Println("get vector", vec, ok)
				if !ok {
					return nil, false
				}
				return &payload.Upsert_Request{
					Vector: vec,
					Config: &payload.Upsert_Config{
						SkipStrictExistCheck: i.forceUpdate,
						Timestamp:            vec.GetTimestamp(),
					},
				}, true
			}
		}, func(loc *payload.Object_StreamLocation, err error) bool {
			if err != nil {
				log.Errorf("stream location error: %v", err)
				return false
			}
			return true
		})
		if err != nil {
			log.Errorf("failed to range gateway: %v", err)
			return nil
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		log.Errorf("importer returned error status errgroup returned error: %v", ctx.Err())
		return err
	}

	log.Infof("importer finished")
	return nil
}

func (i *importer) PreStop(ctx context.Context) error {
	log.Info("removing lock.")
	return i.storedVector.Close(false)
}

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
package service

import (
	"context"
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
	apiName        = "vald/index/job/export"
	grpcMethodName = "vald.v1.StreamListObject/" + vald.StreamListObjectRPCName
)

// Exporter represents an interface for exporting.
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
	log.Info("starting doImporterIndex")
	ctx, span := trace.StartSpan(igrpc.WrapGRPCMethod(ctx, grpcMethodName), apiName+"/service/index.doExportIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	eg, egctx := errgroup.WithContext(ctx)
	eg.SetLimit(i.streamListConcurrency)
	ctx = context.WithoutCancel(egctx)
	gatewayAddrs := i.gateway.GRPCClient().ConnectedAddrs()
	if len(gatewayAddrs) == 0 {
		log.Errorf("Active gateway is not found.: %v ", ctx.Err())
	}

	vch := make(chan *payload.Object_Vector, 5*len(gatewayAddrs))
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

	i.gateway.GRPCClient().RangeConcurrent(egctx, len(gatewayAddrs), func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		cs, err := vc.NewValdClient(conn).StreamUpsert(egctx, copts...)
		if err != nil {
			log.Errorf("failed to create stream upsert client: %v", err)
			return nil
		}

		err = igrpc.BidirectionalStreamClient(cs, func() (res *payload.Upsert_Request) {
			select {
			case <-ctx.Done():
				select {
				case vec, ok := <-vch:
					if !ok {
						return nil
					}
					return &payload.Upsert_Request{
						Vector: vec,
						Config: &payload.Upsert_Config{
							SkipStrictExistCheck: i.forceUpdate,
							Timestamp:            vec.GetTimestamp(),
						},
					}
				case <-time.After(time.Second):
					return nil
				}
			case vec, ok := <-vch:
				if !ok {
					return nil
				}
				return &payload.Upsert_Request{
					Vector: vec,
					Config: &payload.Upsert_Config{
						SkipStrictExistCheck: i.forceUpdate,
						Timestamp:            vec.GetTimestamp(),
					},
				}
			}
		}, func(loc *payload.Object_StreamLocation, err error) {})
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

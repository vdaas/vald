// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"strconv"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/db/kvs/pogreb"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"google.golang.org/grpc"
)

const (
	apiName        = "vald/index/job/export"
	grpcMethodName = "vald.v1.StreamListObject/" + vald.StreamListObjectRPCName
)

// Exporter represents an interface for exporting.
type Exporter interface {
	StartClient(ctx context.Context) (<-chan error, error)
	Start(ctx context.Context) error
	PreStop(ctx context.Context) error
}

type export struct {
	eg                           errgroup.Group
	gateway                      vc.Client
	storedVector                 pogreb.DB
	indexPath                    string
	streamListConcurrency        int
	backgroundSyncInterval       time.Duration
	backgroundCompactionInterval time.Duration
}

// New returns Exporter object if no error occurs.
func New(opts ...Option) (Exporter, error) {
	e := new(export)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(e); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(err)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}

	if err := file.MkdirAll(e.indexPath, os.ModePerm); err != nil {
		log.Errorf("failed to create dir %s", e.indexPath)
		return nil, errors.Wrap(err, "failed to create index path directory")
	}

	path := file.Join(e.indexPath, fmt.Sprintf("%s.db", strconv.FormatInt(time.Now().Unix(), 10)))
	db, err := pogreb.New(pogreb.WithPath(path),
		pogreb.WithBackgroundCompactionInterval(e.backgroundCompactionInterval),
		pogreb.WithBackgroundSyncInterval(e.backgroundSyncInterval))
	if err != nil {
		log.Errorf("failed to open checked List kvs DB %s", path)
		return nil, err
	}
	e.storedVector = db
	return e, nil
}

// StartClient starts the gRPC client.
func (e *export) StartClient(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 1)
	gch, err := e.gateway.Start(ctx)
	if err != nil {
		return nil, err
	}
	e.eg.Go(safety.RecoverFunc(func() (err error) {
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

func (e *export) Start(ctx context.Context) error {
	err := e.doExportIndex(ctx)
	return err
}

func (e *export) doExportIndex(ctx context.Context) (err error) {
	ctx, span := trace.StartSpan(igrpc.WrapGRPCMethod(ctx, grpcMethodName), apiName+"/service/index.doExportIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	emptyReq := new(payload.Object_List_Request)

	grpcCallOpts := []grpc.CallOption{
		grpc.WaitForReady(true),
	}

	stream, err := e.gateway.StreamListObject(ctx, emptyReq, grpcCallOpts...)
	if err != nil || stream == nil {
		return err
	}

	eg, egctx := errgroup.WithContext(ctx)
	eg.SetLimit(e.streamListConcurrency)
	ctx, cancel := context.WithCancelCause(egctx)
	defer cancel(nil)

	var (
		emu  sync.Mutex
		errs = make([]error, 0, e.streamListConcurrency*2) // concurrency * recv+send
	)

	finalize := func() (err error) {
		err = eg.Wait()
		if err != nil {
			emu.Lock()
			errs = append(errs, err)
			emu.Unlock()
		}
		errs := errors.RemoveDuplicates(errs)
		emu.Lock()
		err = errors.Join(errs...)
		emu.Unlock()
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse BidirectionalStream final gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, msg)
		}

		if err != nil {
			log.Errorf("exporter returned error status errgroup returned error: %v", ctx.Err())
			return err
		}

		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return finalize()
		default:
			res, err := stream.Recv()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					err = errors.Wrap(err, "BidirectionalStream Recv returned error")
					emu.Lock()
					errs = append(errs, err)
					emu.Unlock()
					log.Errorf("failed to receive stream message: %v", err)
				}
				return finalize()
			}

			if res != nil && res.GetVector() != nil && res.GetVector().GetId() != "" {
				eg.Go(safety.RecoverFunc(func() (err error) {
					objVec := res.GetVector()
					log.Infof("received object vector id: %s, timestamp: %d", objVec.GetId(), objVec.GetTimestamp())

					storedBinary, ok, err := e.storedVector.Get(objVec.GetId())
					if err != nil {
						log.Errorf("failed to perform Get from check list but still try to finish processing without cache: %v", err)
						return err
					}

					var storedObjVec payload.Object_Vector
					if ok {
						if err := storedObjVec.UnmarshalVT(storedBinary); err != nil {
							log.Errorf("failed to Unmarshal proto to payload.Object_Vector: %v", err)
							return err
						}
					}

					isUpsertVector := !ok || storedObjVec.GetTimestamp() < objVec.GetTimestamp()
					if isUpsertVector {
						dAtA, err := objVec.MarshalVT()
						if err != nil {
							return err
						}
						if err := e.storedVector.Set(objVec.GetId(), dAtA); err != nil {
							log.Errorf("failed to perform Set: %v", err)
							return err
						}
					}
					return nil
				}))
			}
		}
	}
}

func (e *export) PreStop(ctx context.Context) error {
	return e.storedVector.Close(false)
}

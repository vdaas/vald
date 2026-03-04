//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

func (s *server) GetObjectWithMetadata(
	ctx context.Context, req *payload.Object_VectorRequest,
) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, vald.MetadataSpanName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec, err = s.GetObject(ctx, req)
	if err != nil {
		return vec, err
	}
	meta, err := s.metadataClient.Get(ctx, []byte(req.Id.Id))
	if err != nil {
		st, _ := status.FromError(err)
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec.Metadata = meta
	return vec, nil
}

func (s *server) StreamGetObjectWithMetadata(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamGetObjectRPCName+vald.MetadataSpanName),
		apiName+"/"+vald.StreamGetObjectRPCName+vald.MetadataSpanName,
	)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Object_VectorRequest) (*payload.Object_StreamVector, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamGetObjectRPCName+"/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.GetObjectWithMetadata(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				if st != nil && sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamVector{
					Payload: &payload.Object_StreamVector_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamVector{
				Payload: &payload.Object_StreamVector_Vector{
					Vector: res,
				},
			}, nil
		})
	if err != nil {
		st, _ := status.FromError(err)
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())

		}
		return err
	}
	return nil
}

func (s *server) StreamListObjectWithMetadata(
	req *payload.Object_List_Request, stream vald.Object_StreamListObjectServer,
) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamListObjectRPCName), apiName+"/"+vald.StreamListObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var rmu, smu sync.Mutex
	err := s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.StreamListObjectRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()

		client, err := vc.StreamListObject(ctx, req, copts...)
		if err != nil {
			log.Errorf("failed to get StreamListObject client for agent(%s): %v", target, err)
			return err
		}

		eg, ctx := errgroup.WithContext(ctx)
		ectx, ecancel := context.WithCancel(ctx)
		defer ecancel()
		eg.SetLimit(s.streamConcurrency)

		for {
			select {
			case <-ectx.Done():
				var err error
				if !errors.Is(ctx.Err(), context.Canceled) {
					err = errors.Join(err, ctx.Err())
				}
				if egerr := eg.Wait(); err != nil {
					err = errors.Join(err, egerr)
				}
				return err
			default:
				eg.Go(safety.RecoverFunc(func() error {
					rmu.Lock()
					res, err := client.Recv()
					rmu.Unlock()
					if err != nil {
						if errors.Is(err, io.EOF) {
							ecancel()
							return nil
						}
						return errors.ErrServerStreamClientRecv(err)
					}

					vec := res.GetVector()
					if vec == nil {
						st := res.GetStatus()
						log.Warnf("received empty vector: code %v: details %v: message %v",
							st.GetCode(),
							st.GetDetails(),
							st.GetMessage(),
						)
						return nil
					}

					vec.Metadata, err = s.metadataClient.Get(ctx, []byte(vec.Id))
					if err != nil {
						return errors.ErrServerStreamClientRecv(err)
					}

					smu.Lock()
					err = stream.Send(res)
					smu.Unlock()
					if err != nil {
						if sspan != nil {
							st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamListObject send gRPC error response")
							sspan.RecordError(err)
							sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
							sspan.SetStatus(trace.StatusError, err.Error())
						}
						return errors.ErrServerStreamServerSend(err)
					}

					return nil
				}))
			}
		}
	})
	return err
}

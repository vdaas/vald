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
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
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
	vec, err = s.gateway.GetObject(ctx, req, s.copts...)
	if err != nil {
		return vec, err
	}
	meta, err := s.metadataClient.Get(ctx, []byte(req.GetId().GetId()))
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

func (s *server) StreamGetObjectWithMetadata(
	stream vald.ObjectWithMetadata_StreamGetObjectWithMetadataServer,
) (err error) {
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
	req *payload.Object_List_Request,
	stream vald.ObjectWithMetadata_StreamListObjectWithMetadataServer,
) error {
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamListObjectRPCName+vald.MetadataSpanName),
		apiName+"/"+vald.StreamListObjectRPCName+vald.MetadataSpanName,
	)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	client, err := s.gateway.StreamListObject(ctx, req, s.copts...)
	if err != nil {
		st, _ := status.FromError(err)
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	for {
		res, err := client.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			st, _ := status.FromError(err)
			if st != nil && span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return err
		}
		vec := res.GetVector()
		if vec != nil && vec.GetId() != "" {
			meta, merr := s.metadataClient.Get(ctx, []byte(vec.GetId()))
			if merr != nil {
				st, _ := status.FromError(merr)
				if st != nil && span != nil {
					span.RecordError(merr)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
					span.SetStatus(trace.StatusError, merr.Error())
				}
				return merr
			}
			vec.Metadata = meta
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
}

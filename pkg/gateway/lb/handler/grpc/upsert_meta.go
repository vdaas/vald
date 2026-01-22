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
	"fmt"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func (s *server) UpsertWithMetadata(
	ctx context.Context, req *payload.Upsert_Request,
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, vald.MetadataSpanName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	loc, err = s.Upsert(ctx, req)
	if err != nil {
		return loc, err
	}
	err = s.metadataClient.Put(ctx, []byte(req.Vector.Id), req.Vector.Metadata)
	if err != nil {
		st, _ := status.FromError(err)
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpsertWithMetadata(stream vald.Upsert_StreamUpsertServer) (err error) {
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.StreamUpsertRPCName+vald.MetadataSpanName),
		apiName+"/"+vald.StreamUpsertRPCName+vald.MetadataSpanName,
	)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Upsert_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamUpsertRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.UpsertWithMetadata(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				if st != nil && sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
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

func (s *server) MultiUpsertWithMetadata(
	ctx context.Context, reqs *payload.Upsert_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.MultiUpsertRPCName+vald.MetadataSpanName),
		apiName+"/"+vald.MultiUpsertRPCName+vald.MetadataSpanName,
	)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var (
		emu sync.Mutex
		lmu sync.Mutex
	)
	eg, ectx := errgroup.New(ctx)
	eg.SetLimit(s.multiConcurrency)
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	for i, r := range reqs.GetRequests() {
		if r != nil && r.GetVector() != nil && len(r.GetVector().GetVector()) >= algorithm.MinimumVectorDimensionSize && r.GetVector().GetId() != "" {
			idx := i
			req := r
			eg.Go(safety.RecoverFunc(func() (err error) {
				ectx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ectx, "eg.Go"), apiName+"/"+vald.MultiUpsertRPCName+"/id-"+req.GetVector().GetId())
				defer func() {
					if sspan != nil {
						sspan.End()
					}
				}()
				res, err := s.UpsertWithMetadata(ectx, req)
				if err != nil {
					st, _ := status.FromError(err)
					if st != nil && sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					emu.Lock()
					if errs == nil {
						errs = err
					} else {
						errs = errors.Join(errs, err)
					}
					emu.Unlock()
				} else if res != nil && res.GetUuid() == req.GetVector().GetId() && res.GetIps() != nil {
					lmu.Lock()
					locs.Locations[idx] = res
					lmu.Unlock()
				}
				return nil
			}))
		} else {
			var (
				err   error
				field string
			)
			switch {
			case r.GetVector() == nil, len(r.GetVector().GetVector()) < algorithm.MinimumVectorDimensionSize:
				err = errors.ErrInvalidDimensionSize(len(r.GetVector().GetVector()), 0)
				field = "vector"
			case r.GetVector().GetId() == "":
				err = errors.ErrInvalidUUID(r.GetVector().GetId())
				field = "uuid"
			}
			err = status.WrapWithInvalidArgument(vald.MultiUpsertRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   r.GetVector().GetId(),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       field,
							Description: err.Error(),
						},
					},
				}, info.Get())
			emu.Lock()
			if errs == nil {
				errs = err
			} else {
				errs = errors.Join(errs, err)
			}
			emu.Unlock()

		}
	}
	err := eg.Wait()
	if err != nil {
		emu.Lock()
		if errs == nil {
			errs = err
		} else {
			errs = errors.Join(errs, err)
		}
		emu.Unlock()
	}

	if errs != nil {
		st, _ := status.FromError(err)
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		errs = err
	}

	return locs, errs
}

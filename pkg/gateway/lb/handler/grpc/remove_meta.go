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
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
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

func (s *server) RemoveWithMetadata(
	ctx context.Context, req *payload.Remove_Request,
) (locs *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, vald.MetadataSpanName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs, err = s.Remove(ctx, req)
	if err != nil {
		return locs, err
	}
	err = s.metadataClient.Delete(ctx, []byte(req.Id.Id))
	if err != nil {
		st, _ := status.FromError(err)
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return locs, nil
}

func (s *server) StreamRemoveWithMetadata(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.StreamRemoveRPCName+vald.MetadataSpanName),
		apiName+"/"+vald.StreamRemoveRPCName+vald.MetadataSpanName,
	)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Remove_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamRemoveRPCName+"/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.RemoveWithMetadata(ctx, req)
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

func (s *server) MultiRemoveWithMetadata(
	ctx context.Context, reqs *payload.Remove_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName+vald.MetadataSpanName),
		apiName+"/"+vald.MultiRemoveRPCName+vald.MetadataSpanName,
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
		if r != nil && r.GetId().GetId() != "" {
			idx := i
			req := r
			eg.Go(safety.RecoverFunc(func() (err error) {
				ectx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ectx, "eg.Go"), apiName+"/"+vald.MultiRemoveRPCName+"/id-"+req.GetId().GetId())
				defer func() {
					if sspan != nil {
						sspan.End()
					}
				}()
				res, err := s.RemoveWithMetadata(ectx, req)
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
				} else if res != nil && res.GetUuid() == req.GetId().GetId() && res.GetIps() != nil {
					lmu.Lock()
					locs.Locations[idx] = res
					lmu.Unlock()
				}
				return nil
			}))
		} else {
			err := errors.ErrInvalidUUID(r.GetId().GetId())
			err = status.WrapWithInvalidArgument(vald.MultiRemoveRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   r.GetId().GetId(),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid",
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
		st, _ := status.FromError(errs)
		if st != nil && span != nil {
			span.RecordError(errs)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, errs.Error())
		}
	}

	return locs, errs
}

func (s *server) RemoveByTimestampWithMetadata(
	ctx context.Context, req *payload.Remove_TimestampRequest,
) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(ctx, vald.MetadataSpanName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs, err := s.RemoveByTimestamp(ctx, req)
	if err != nil {
		return locs, err
	}

	var wg sync.WaitGroup
	var emu sync.Mutex
	for i, loc := range locs.Locations {
		idx, id := i, loc.Uuid
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/Metadata/"+ti)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			err := s.metadataClient.Delete(ctx, []byte(id))
			if err != nil {
				st, _ := status.FromError(err)
				if st != nil && sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				emu.Lock()
				if errs != nil {
					errs = errors.Join(errs, err)
				} else {
					errs = err
				}
				emu.Unlock()
				return nil
			}
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, _ := status.FromError(errs)
		if st != nil && span != nil {
			span.RecordError(errs)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, errs.Error())
		}
		return locs, errs
	}
	return locs, nil
}

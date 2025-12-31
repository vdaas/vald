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
	"slices"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

func (s *server) Update(
	ctx context.Context, req *payload.Update_Request,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetVector().GetId()
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid uuid", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec := req.GetVector().GetVector()
	vl := len(vec)
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid vector argument", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	if req.GetConfig().GetDisableBalancedUpdate() {
		var (
			mu      sync.RWMutex
			aeCount atomic.Uint64
			updated atomic.Uint64
			ls      = make([]string, 0, s.replica)
			visited = make(map[string]bool, s.replica)
			locs    = &payload.Object_Location{
				Uuid: uuid,
				Ips:  make([]string, 0, s.replica),
			}
		)
		err = s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpdateRPCName+"/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			loc, err := vc.Update(ctx, req, copts...)
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					if st.Code() != codes.AlreadyExists &&
						st.Code() != codes.Canceled &&
						st.Code() != codes.DeadlineExceeded &&
						st.Code() != codes.InvalidArgument &&
						st.Code() != codes.NotFound &&
						st.Code() != codes.OK &&
						st.Code() != codes.Unimplemented {
						if span != nil {
							span.RecordError(err)
							span.SetAttributes(trace.FromGRPCStatus(st.Code(), fmt.Sprintf("Update operation for Agent %s failed,\terror: %v", target, err))...)
							span.SetStatus(trace.StatusError, err.Error())
						}
						return err
					}
					if st.Code() == codes.AlreadyExists {
						host, _, err := net.SplitHostPort(target)
						if err != nil {
							host = target
						}
						aeCount.Add(1)
						mu.Lock()
						visited[target] = true
						locs.Ips = append(locs.GetIps(), host)
						ls = append(ls, host)
						mu.Unlock()

					}
				}
				return nil
			}
			if loc != nil {
				updated.Add(1)
				mu.Lock()
				visited[target] = true
				locs.Ips = append(locs.GetIps(), loc.GetIps()...)
				ls = append(ls, loc.GetName())
				mu.Unlock()
			}
			return nil
		})
		switch {
		case err != nil:
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		case len(locs.Ips) <= 0:
			err = errors.ErrIndexNotFound
			err = status.WrapWithNotFound(vald.UpdateRPCName+" API update target not found", err, reqInfo, resInfo)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		case updated.Load()+aeCount.Load() < uint64(s.replica):
			shortage := s.replica - int(updated.Load()+aeCount.Load())
			err = s.gateway.DoMulti(ctx, shortage, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
				mu.RLock()
				tf, ok := visited[target]
				mu.RUnlock()
				if tf && ok {
					return errors.Errorf("target: %s already inserted will skip", target)
				}
				ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "DoMulti/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
				defer func() {
					if span != nil {
						span.End()
					}
				}()
				loc, err := vc.Insert(ctx, &payload.Insert_Request{
					Vector: req.GetVector(),
					Config: &payload.Insert_Config{
						SkipStrictExistCheck: true,
						Filters:              req.GetConfig().GetFilters(),
						Timestamp:            req.GetConfig().GetTimestamp(),
					},
				}, copts...)
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil && span != nil {
						span.RecordError(err)
						span.SetAttributes(trace.FromGRPCStatus(st.Code(), fmt.Sprintf("Shortage index Insert for Update operation for Agent %s failed,\terror: %v", target, err))...)
						span.SetStatus(trace.StatusError, err.Error())
					}
					return err
				}
				if loc != nil {
					updated.Add(1)
					mu.Lock()
					locs.Ips = append(locs.GetIps(), loc.GetIps()...)
					ls = append(ls, loc.GetName())
					mu.Unlock()
				}
				return nil
			})
		case updated.Load() == 0 && aeCount.Load() > 0:
			err = errors.ErrSameVectorAlreadyExists(uuid, vec, vec)
			err = status.WrapWithAlreadyExists(vald.UpdateRPCName+" API update target same vector already exists", err, reqInfo, resInfo)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err

		}
		slices.Sort(ls)
		locs.Name = strings.Join(ls, ",")
		return locs, nil
	}

	if !req.GetConfig().GetSkipStrictExistCheck() {
		vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if err != nil {
			st, _ := status.FromError(err)
			if span != nil && st != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
			if vec.GetTimestamp() < req.GetVector().GetTimestamp() {
				return s.UpdateTimestamp(ctx, &payload.Update_TimestampRequest{
					Id:        uuid,
					Timestamp: req.GetVector().GetTimestamp(),
				})
			}
			err = errors.ErrSameVectorAlreadyExists(uuid, vec.GetVector(), req.GetVector().GetVector())
			st, msg, err := status.ParseError(err, codes.AlreadyExists,
				"error "+vald.UpdateRPCName+" API ID = "+uuid+"'s same vector data already exists",
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				}, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Update_Config{SkipStrictExistCheck: true}
		}
	}
	var now int64
	if req.GetConfig().GetTimestamp() > 0 {
		now = req.GetConfig().GetTimestamp()
	} else {
		now = time.Now().UnixNano()
	}

	rreq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: uuid,
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: true,
			Timestamp:            now,
		},
	}
	res, err = s.Remove(ctx, rreq)
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(vald.RemoveRPCName+" for "+vald.UpdateRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(rreq),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.RemoveRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		st, _ := status.FromError(err)
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	now++
	ireq := &payload.Insert_Request{
		Vector: req.GetVector(),
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
			Filters:              req.GetConfig().GetFilters(),
			Timestamp:            now,
		},
	}
	res, err = s.Insert(ctx, ireq)
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(ireq),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		st, _ := status.FromError(err)
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.StreamUpdateRPCName), apiName+"/"+vald.StreamUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Update_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamUpdateRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Update(ctx, req)
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

func (s *server) MultiUpdate(
	ctx context.Context, reqs *payload.Update_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
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
				ectx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ectx, "eg.Go"), apiName+"/"+vald.MultiUpdateRPCName+"/id-"+req.GetVector().GetId())
				defer func() {
					if sspan != nil {
						sspan.End()
					}
				}()
				res, err := s.Update(ectx, req)
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
			err = status.WrapWithInvalidArgument(vald.MultiUpdateRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   r.GetVector().GetId(),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName,
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
		st, _ := status.FromError(errs)
		if st != nil && span != nil {
			span.RecordError(errs)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, errs.Error())
		}
	}

	return locs, errs
}

func (s *server) UpdateTimestamp(
	ctx context.Context, req *payload.Update_TimestampRequest,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateTimestampRPCName), apiName+"/"+vald.UpdateTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetId()
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateTimestampRPCName + "." + vald.GetObjectRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.UpdateTimestampRPCName+" API invalid uuid", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	ts := req.GetTimestamp()
	if ts < 0 {
		err = errors.ErrInvalidTimestamp(ts)
		err = status.WrapWithInvalidArgument(vald.UpdateTimestampRPCName+" API invalid vector argument", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "timestamp",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	var (
		mu      sync.RWMutex
		aeCount atomic.Uint64
		updated atomic.Uint64
		ls      = make([]string, 0, s.replica)
		visited = make(map[string]bool, s.replica)
		locs    = &payload.Object_Location{
			Uuid: uuid,
			Ips:  make([]string, 0, s.replica),
		}
	)
	err = s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpdateTimestampRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.UpdateTimestamp(ctx, req, copts...)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				if st.Code() != codes.AlreadyExists &&
					st.Code() != codes.Canceled &&
					st.Code() != codes.DeadlineExceeded &&
					st.Code() != codes.InvalidArgument &&
					st.Code() != codes.NotFound &&
					st.Code() != codes.OK &&
					st.Code() != codes.Unimplemented {
					if span != nil {
						span.RecordError(err)
						span.SetAttributes(trace.FromGRPCStatus(st.Code(), fmt.Sprintf("UpdateTimestamp operation for Agent %s failed,\terror: %v", target, err))...)
						span.SetStatus(trace.StatusError, err.Error())
					}
					return err
				}
				if st.Code() == codes.AlreadyExists {
					host, _, err := net.SplitHostPort(target)
					if err != nil {
						host = target
					}
					aeCount.Add(1)
					mu.Lock()
					visited[target] = true
					locs.Ips = append(locs.GetIps(), host)
					ls = append(ls, host)
					mu.Unlock()

				}
			}
			return nil
		}
		if loc != nil {
			updated.Add(1)
			mu.Lock()
			visited[target] = true
			locs.Ips = append(locs.GetIps(), loc.GetIps()...)
			ls = append(ls, loc.GetName())
			mu.Unlock()
		}
		return nil
	})
	switch {
	case err != nil:
		st, _ := status.FromError(err)
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case len(locs.Ips) <= 0:
		err = errors.ErrIndexNotFound
		err = status.WrapWithNotFound(vald.UpdateTimestampRPCName+" API update target not found", err, reqInfo, resInfo)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case updated.Load()+aeCount.Load() < uint64(s.replica):
		shortage := s.replica - int(updated.Load()+aeCount.Load())
		vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if err != nil {
			st, _ := status.FromError(err)
			if st != nil && span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		err = s.gateway.DoMulti(ctx, shortage, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
			mu.RLock()
			tf, ok := visited[target]
			mu.RUnlock()
			if tf && ok {
				return errors.Errorf("target: %s already inserted will skip", target)
			}
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "DoMulti/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			loc, err := vc.Insert(ctx, &payload.Insert_Request{
				Vector: vec,
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
					Timestamp:            ts,
				},
			}, copts...)
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil && span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), fmt.Sprintf("Shortage index Insert for Update operation for Agent %s failed,\terror: %v", target, err))...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return err
			}
			if loc != nil {
				updated.Add(1)
				mu.Lock()
				locs.Ips = append(locs.GetIps(), loc.GetIps()...)
				ls = append(ls, loc.GetName())
				mu.Unlock()
			}
			return nil
		})
		if err != nil {
			st, _ := status.FromError(err)
			if st != nil && span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
	case updated.Load() == 0 && aeCount.Load() > 0:
		err = status.WrapWithAlreadyExists(vald.UpdateTimestampRPCName+" API update target same vector already exists", errors.ErrSameVectorAlreadyExists(uuid, nil, nil), reqInfo, resInfo)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err

	}
	slices.Sort(ls)
	locs.Name = strings.Join(ls, ",")
	return locs, nil
}

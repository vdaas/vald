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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

func (s *server) exists(ctx context.Context, uuid string) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "exists"), apiName+"/"+vald.ExistsRPCName+"/exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	ich := make(chan *payload.Object_ID, 1)
	ech := make(chan error, 1)
	doneErr := errors.New("done exists")
	ctx, cancel := context.WithCancelCause(ctx)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ich)
		defer close(ech)
		var once sync.Once
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/exists/BroadCast/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			meta := &payload.Object_ID{
				Id: uuid,
			}
			oid, err := vc.Exists(sctx, meta, copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.ExistsRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.ExistsRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.ExistsRPCName+" API meta "+uuid+"'s uuid not found",
						&errdetails.RequestInfo{
							RequestId:   uuid,
							ServingData: errdetails.Serialize(meta),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.ExistsRPCName,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				if err != nil && st != nil &&
					code != codes.Canceled &&
					code != codes.DeadlineExceeded &&
					code != codes.InvalidArgument &&
					code != codes.NotFound &&
					code != codes.OK &&
					code != codes.Unimplemented {
					return err
				}
				return nil
			}
			if oid != nil && oid.GetId() != "" {
				once.Do(func() {
					ich <- oid
					cancel(doneErr)
				})
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
		if errors.Is(err, context.Canceled) && errors.Is(context.Cause(ctx), doneErr) {
			select {
			case id = <-ich:
				if id == nil || id.GetId() == "" {
					err = errors.ErrObjectIDNotFound(uuid)
				} else {
					err = nil
				}
			default:
			}
		}
	case id = <-ich:
		if id == nil || id.GetId() == "" {
			err = errors.ErrObjectIDNotFound(uuid)
		}
	case err = <-ech:
	}
	if err == nil && (id == nil || id.GetId() == "") {
		err = errors.ErrObjectIDNotFound(uuid)
	}
	if err != nil {
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return id, nil
}

func (s *server) Exists(
	ctx context.Context, meta *payload.Object_ID,
) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := meta.GetId()
	id, err = s.exists(ctx, uuid)
	if err == nil {
		return id, nil
	}
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(meta),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.ExistsRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	var attrs trace.Attributes
	switch {
	case errors.Is(err, errors.ErrInvalidUUID(uuid)):
		err = status.WrapWithInvalidArgument(vald.ExistsRPCName+" API invalid argument for uuid \""+uuid+"\" detected", err, reqInfo, resInfo, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       "uuid",
					Description: err.Error(),
				},
			},
		})
		attrs = trace.StatusCodeInvalidArgument(err.Error())
	case errors.Is(err, errors.ErrObjectIDNotFound(uuid)):
		err = status.WrapWithNotFound(vald.ExistsRPCName+" API id "+uuid+"'s data not found", err, reqInfo, resInfo)
		attrs = trace.StatusCodeNotFound(err.Error())
	case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
		err = status.WrapWithInternal(vald.ExistsRPCName+" API connection not found", err, reqInfo, resInfo)
		attrs = trace.StatusCodeInternal(err.Error())
	case errors.Is(err, context.Canceled):
		err = status.WrapWithCanceled(vald.ExistsRPCName+" API canceled", err, reqInfo, resInfo)
		attrs = trace.StatusCodeCancelled(err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API deadline exceeded", err, reqInfo, resInfo)
		attrs = trace.StatusCodeDeadlineExceeded(err.Error())
	default:
		var (
			st   *status.Status
			code codes.Code
			msg  string
		)
		st, _ = status.FromError(err)
		if st != nil {
			code = st.Code()
			msg = uuid + "'s object id:" + vald.ExistsRPCName + " API uuid " + uuid + "'s request returned error\t" + st.String()
		} else {
			code = codes.Unknown
			msg = uuid + "'s object id:" + vald.ExistsRPCName + " API uuid " + uuid + "'s request returned error"
		}
		attrs = trace.FromGRPCStatus(code, msg)
	}
	log.Debug(err)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(attrs...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) getObject(
	ctx context.Context, uuid string,
) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "getObject"), apiName+"/"+vald.GetObjectRPCName+"/getObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vch := make(chan *payload.Object_Vector, 1)
	ech := make(chan error, 1)
	doneErr := errors.New("done getObject")
	ctx, cancel := context.WithCancelCause(ctx)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(vch)
		defer close(ech)
		var once sync.Once
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/getObject/BroadCast/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			req := &payload.Object_VectorRequest{
				Id: &payload.Object_ID{
					Id: uuid,
				},
			}
			ovec, err := vc.GetObject(sctx, req, copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.GetObjectRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.GetObjectRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.GetObjectRPCName+" API meta "+uuid+"'s uuid not found",
						&errdetails.RequestInfo{
							RequestId:   uuid,
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				if err != nil && st != nil &&
					code != codes.Canceled &&
					code != codes.DeadlineExceeded &&
					code != codes.InvalidArgument &&
					code != codes.NotFound &&
					code != codes.OK &&
					code != codes.Unimplemented {
					return err
				}
				return nil
			}
			if ovec != nil && ovec.GetId() != "" && ovec.GetVector() != nil {
				once.Do(func() {
					vch <- ovec
					cancel(doneErr)
				})
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
		if errors.Is(err, context.Canceled) && errors.Is(context.Cause(ctx), doneErr) {
			select {
			case vec = <-vch:
				if vec == nil || vec.GetId() == "" || vec.GetVector() == nil {
					err = errors.ErrObjectNotFound(nil, uuid)
				} else {
					err = nil
				}
			default:
			}
		}
	case vec = <-vch:
		if vec == nil || vec.GetId() == "" || vec.GetVector() == nil {
			err = errors.ErrObjectNotFound(nil, uuid)
		}
	case err = <-ech:
	}
	if err == nil && (vec == nil || vec.GetId() == "" || vec.GetVector() == nil) {
		err = errors.ErrObjectNotFound(nil, uuid)
	}
	if err != nil {
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	return vec, nil
}

func (s *server) GetObject(
	ctx context.Context, req *payload.Object_VectorRequest,
) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetId().GetId()
	vec, err = s.getObject(ctx, uuid)
	if err == nil {
		if vec != nil && vec.GetId() != "" && vec.GetVector() != nil {
			return vec, nil
		}
		err = errors.ErrObjectNotFound(nil, uuid)
	}
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	var attrs trace.Attributes
	switch {
	case errors.Is(err, errors.ErrInvalidUUID(uuid)):
		err = status.WrapWithInvalidArgument(vald.GetObjectRPCName+" API invalid argument for uuid \""+uuid+"\" detected", err, reqInfo, resInfo, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       "uuid",
					Description: err.Error(),
				},
			},
		})
		attrs = trace.StatusCodeInvalidArgument(err.Error())
	case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
		err = status.WrapWithNotFound(vald.GetObjectRPCName+" API id "+uuid+"'s object not found", err, reqInfo, resInfo)
		attrs = trace.StatusCodeNotFound(err.Error())
	case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
		err = status.WrapWithInternal(vald.GetObjectRPCName+" API connection not found", err, reqInfo, resInfo)
		attrs = trace.StatusCodeInternal(err.Error())
	case errors.Is(err, context.Canceled):
		err = status.WrapWithCanceled(vald.GetObjectRPCName+" API canceled", err, reqInfo, resInfo)
		attrs = trace.StatusCodeCancelled(err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.WrapWithDeadlineExceeded(vald.GetObjectRPCName+" API deadline exceeded", err, reqInfo, resInfo)
		attrs = trace.StatusCodeDeadlineExceeded(err.Error())
	default:
		st, _ := status.FromError(err)
		if st != nil {
			attrs = trace.FromGRPCStatus(st.Code(), st.Message())
		}
	}
	log.Debug(err)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(attrs...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamGetObjectRPCName), apiName+"/"+vald.StreamGetObjectRPCName)
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
			res, err := s.GetObject(ctx, req)
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

func (s *server) StreamListObject(
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

func (s *server) GetTimestamp(
	ctx context.Context, req *payload.Object_TimestampRequest,
) (ts *payload.Object_Timestamp, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetTimestampRPCName), apiName+"/"+vald.GetTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetId().GetId()
	tch := make(chan *payload.Object_Timestamp, 1)
	ech := make(chan error, 1)
	doneErr := errors.New("done getTimestamp")
	ctx, cancel := context.WithCancelCause(ctx)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(tch)
		defer close(ech)
		var once sync.Once
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/getTimestamp/BroadCast/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			req := &payload.Object_TimestampRequest{
				Id: &payload.Object_ID{
					Id: uuid,
				},
			}
			ots, err := vc.GetTimestamp(sctx, req, copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.GetTimestampRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.GetTimestampRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.GetTimestampRPCName+" API meta "+uuid+"'s uuid not found",
						&errdetails.RequestInfo{
							RequestId:   uuid,
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetTimestampRPCName,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				if err != nil && st != nil &&
					code != codes.Canceled &&
					code != codes.DeadlineExceeded &&
					code != codes.InvalidArgument &&
					code != codes.NotFound &&
					code != codes.OK &&
					code != codes.Unimplemented {
					return err
				}
				return nil
			}
			if ots != nil && ots.GetId() != "" {
				once.Do(func() {
					tch <- ots
					cancel(doneErr)
				})
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
		if errors.Is(err, context.Canceled) && errors.Is(context.Cause(ctx), doneErr) {
			select {
			case ts = <-tch:
				if ts == nil || ts.GetId() == "" {
					err = errors.ErrObjectNotFound(nil, uuid)
				} else {
					err = nil
				}
			default:
			}
		}
	case ts = <-tch:
		if ts == nil || ts.GetId() == "" {
			err = errors.ErrObjectNotFound(nil, uuid)
		}
	case err = <-ech:
	}
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   uuid,
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetTimestampRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrInvalidUUID(uuid)):
			err = status.WrapWithInvalidArgument(vald.GetTimestampRPCName+" API invalid argument for uuid \""+uuid+"\" detected", err, reqInfo, resInfo, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "uuid",
						Description: err.Error(),
					},
				},
			})
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
			err = status.WrapWithNotFound(vald.GetTimestampRPCName+" API id "+uuid+"'s object not found", err, reqInfo, resInfo)
			attrs = trace.StatusCodeNotFound(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.GetTimestampRPCName+" API connection not found", err, reqInfo, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.GetTimestampRPCName+" API canceled", err, reqInfo, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.GetTimestampRPCName+" API deadline exceeded", err, reqInfo, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			st, _ := status.FromError(err)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return ts, nil
}

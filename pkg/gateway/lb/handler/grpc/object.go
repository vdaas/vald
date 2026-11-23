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
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

// handleObjectBroadCastError classifies a per-target BroadCast error for the
// object read RPCs (Exists/GetObject/GetTimestamp): it records the error on the
// per-target span and returns a non-nil error only for codes that must fail the
// whole BroadCast, nil for the tolerable per-target outcomes (a missing uuid on
// one agent is expected). Extracted from the three verbatim-identical per-target
// blocks (they differed only in rpcName and the request payload serialized into
// RequestInfo).
func (s *server) handleObjectBroadCastError(
	sspan trace.Span, target, rpcName, uuid string, servingData any, err error,
) error {
	var (
		attrs trace.Attributes
		st    *status.Status
		msg   string
		code  codes.Code
	)
	switch {
	case errors.Is(err, context.Canceled):
		if sspan != nil {
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix +
					"/vald.v1." + rpcName + ".BroadCast/" +
					target + " canceled: " + err.Error())
		}
		code = codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		if sspan != nil {
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix +
					"/vald.v1." + rpcName + ".BroadCast/" +
					target + " deadline_exceeded: " + err.Error())
		}
		code = codes.DeadlineExceeded
	default:
		var ok bool
		st, ok = status.FromError(err)
		if ok && st != nil {
			code = st.Code()
			msg = st.Message()
		}
		if ok && st != nil && (code == codes.NotFound ||
			code == codes.Canceled ||
			code == codes.DeadlineExceeded ||
			code == codes.InvalidArgument ||
			code == codes.OK ||
			code == codes.Unimplemented) {
			if sspan != nil {
				attrs = trace.FromGRPCStatus(code, msg)
			}
		} else {
			st, msg, err = status.ParseError(err, codes.NotFound, "error "+rpcName+" API meta "+uuid+"'s uuid not found",
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(servingData),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + rpcName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				})
			if st != nil {
				code = st.Code()
			} else {
				code = codes.NotFound
			}
			if sspan != nil {
				attrs = trace.FromGRPCStatus(code, msg)
			}
		}
	}
	if sspan != nil {
		errhandler.RecordSpanAttrs(sspan, attrs, err)
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

var errBroadcastFirstDone = errors.New("broadcastFirst done")

// broadcastFirst broadcasts call to every READ target and returns the first
// result accepted by hit, canceling the remaining targets once it arrives. name
// labels the per-target spans (apiName/<name>/BroadCast/<target>). When
// probeCover is true and replication exceeds 1, it first probes a minimal
// cyclic subset of size len(addrs)-replica+1 and falls back to the remaining
// targets only if needed. call must return the zero T together with any error;
// a nil error with a missing object is reported by returning a result rejected
// by hit. ok is false (with a nil err) when no target had the object.
func broadcastFirst[T any](
	ctx context.Context,
	s *server,
	name string,
	probeCover bool,
	hit func(T) bool,
	call func(ctx context.Context, sspan trace.Span, target string, vc vald.Client, copts ...grpc.CallOption) (T, error),
) (res T, ok bool, err error) {
	bctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	var (
		hitVal T
		hasHit atomic.Bool
	)
	cb := func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		sctx := ctx
		var sspan trace.Span
		if trace.Enabled() {
			sctx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+name+"/BroadCast/"+target)
			defer trace.End(sspan)
		}
		r, cerr := call(sctx, sspan, target, vc, copts...)
		if cerr != nil {
			return cerr
		}
		if hit(r) && !hasHit.Swap(true) {
			hitVal = r
			cancel(errBroadcastFirstDone)
		}
		return nil
	}
	if probeCover && s.replica > 1 {
		if addrs := s.gateway.Addrs(ctx); len(addrs) > s.replica {
			cover := len(addrs) - s.replica + 1
			offset := s.gateway.NextRotateOffset(len(addrs))
			err = s.gateway.BroadCastSlice(bctx, service.READ, addrs, offset, cover, cb)
			if hasHit.Load() {
				return hitVal, true, nil
			}
			if ctx.Err() != nil {
				if err == nil {
					err = ctx.Err()
				}
				var zero T
				return zero, false, err
			}
			if err != nil && errors.Is(err, errGRPCClientConnNotFoundStar) {
				var zero T
				return zero, false, err
			}
			remOffset := (offset + cover) % len(addrs)
			err2 := s.gateway.BroadCastSlice(bctx, service.READ, addrs, remOffset, len(addrs)-cover, cb)
			if hasHit.Load() {
				return hitVal, true, nil
			}
			if err == nil {
				err = err2
			} else if err2 != nil {
				err = errors.Join(err, err2)
			}
			if err == nil {
				err = ctx.Err()
			}
			var zero T
			return zero, false, err
		}
	}
	err = s.gateway.BroadCast(bctx, service.READ, cb)
	if hasHit.Load() {
		return hitVal, true, nil
	}
	if err == nil {
		err = ctx.Err()
	}
	var zero T
	return zero, false, err
}

func (s *server) existsWithID(
	ctx context.Context, uuid string, probeCover bool, meta *payload.Object_ID,
) (id *payload.Object_ID, ok bool, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "exists"), apiName+"/"+vald.ExistsRPCName+"/exists")
		defer trace.End(span)
	}

	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		log.Warn(err)
		errhandler.RecordSpanError(span, codes.InvalidArgument, err)
		return nil, false, err
	}
	if meta == nil {
		meta = &payload.Object_ID{
			Id: uuid,
		}
	}
	id, ok, err = broadcastFirst(ctx, s, "exists", probeCover,
		func(oid *payload.Object_ID) bool { return oid.GetId() != "" },
		func(ctx context.Context, sspan trace.Span, target string, vc vald.Client, copts ...grpc.CallOption) (*payload.Object_ID, error) {
			oid, rerr := vc.Exists(ctx, meta, copts...)
			if rerr != nil {
				return nil, s.handleObjectBroadCastError(sspan, target, vald.ExistsRPCName, uuid, meta, rerr)
			}
			return oid, nil
		})
	if err != nil {
		errhandler.RecordSpanError(span, spanErrorCode(err, codes.NotFound), err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	return id, true, nil
}

func spanErrorCode(err error, fallback codes.Code) codes.Code {
	switch {
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, errGRPCClientConnNotFoundStar):
		return codes.Internal
	default:
		if st, ok := status.FromError(err); ok && st != nil && st.Code() != codes.OK {
			return st.Code()
		}
		return fallback
	}
}

func (s *server) wrapObjectReadError(
	ctx context.Context,
	span trace.Span,
	rpcName, uuid string,
	req any,
	notFoundSuffix string,
	err error,
) error {
	details := s.errDetails(ctx, rpcName, uuid, req)
	var attrs trace.Attributes
	switch {
	case errors.Is(err, errors.ErrInvalidUUID(uuid)):
		err = status.WrapWithInvalidArgument(rpcName+" API invalid argument for uuid \""+uuid+"\" detected", err, append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       uuidField,
					Description: err.Error(),
				},
			},
		})...)
		if span != nil {
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		}
	case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
		err = status.WrapWithNotFound(rpcName+" API id "+uuid+"'s "+notFoundSuffix, err, details...)
		if span != nil {
			attrs = trace.StatusCodeNotFound(err.Error())
		}
	case errors.Is(err, errGRPCClientConnNotFoundStar):
		err = status.WrapWithInternal(rpcName+" API connection not found", err, details...)
		if span != nil {
			attrs = trace.StatusCodeInternal(err.Error())
		}
	case errors.Is(err, context.Canceled):
		err = status.WrapWithCanceled(rpcName+" API canceled", err, details...)
		if span != nil {
			attrs = trace.StatusCodeCancelled(err.Error())
		}
	case errors.Is(err, context.DeadlineExceeded):
		err = status.WrapWithDeadlineExceeded(rpcName+" API deadline exceeded", err, details...)
		if span != nil {
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		}
	default:
		if span != nil {
			st, _ := status.FromError(err)
			if rpcName == vald.ExistsRPCName {
				var (
					code codes.Code
					msg  string
				)
				if st != nil {
					code = st.Code()
					msg = uuid + "'s object id:" + vald.ExistsRPCName + " API uuid " + uuid + "'s request returned error\t" + st.String()
				} else {
					code = codes.Unknown
					msg = uuid + "'s object id:" + vald.ExistsRPCName + " API uuid " + uuid + "'s request returned error"
				}
				attrs = trace.FromGRPCStatus(code, msg)
			} else if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
	}
	log.Debug(err)
	if span != nil {
		errhandler.RecordSpanAttrs(span, attrs, err)
	}
	return err
}

func (s *server) existsCheck(ctx context.Context, uuid string) (bool, error) {
	_, ok, err := s.existsWithID(ctx, uuid, false, nil)
	return ok, err
}

func (s *server) exists(ctx context.Context, uuid string) (id *payload.Object_ID, err error) {
	id, ok, err := s.existsWithID(ctx, uuid, true, nil)
	if err == nil && !ok {
		err = errors.ErrObjectIDNotFound(uuid)
	}
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *server) Exists(
	ctx context.Context, meta *payload.Object_ID,
) (id *payload.Object_ID, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
		defer trace.End(span)
	}
	uuid := meta.GetId()
	id, ok, err := s.existsWithID(ctx, uuid, true, meta)
	if err == nil && !ok {
		err = errors.ErrObjectIDNotFound(uuid)
	}
	if err == nil {
		return id, nil
	}
	return nil, s.wrapObjectReadError(ctx, span, vald.ExistsRPCName, uuid, meta, "data not found", err)
}

func (s *server) getObjectWithReq(
	ctx context.Context, uuid string, req *payload.Object_VectorRequest,
) (vec *payload.Object_Vector, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "getObject"), apiName+"/"+vald.GetObjectRPCName+"/getObject")
		defer trace.End(span)
	}
	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		log.Warn(err)
		errhandler.RecordSpanError(span, codes.InvalidArgument, err)
		return nil, err
	}
	if req == nil || req.GetFilters() != nil {
		req = &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		}
	}
	vec, ok, err := broadcastFirst(ctx, s, "getObject", true,
		func(ovec *payload.Object_Vector) bool { return ovec.GetId() != "" && ovec.GetVector() != nil },
		func(ctx context.Context, sspan trace.Span, target string, vc vald.Client, copts ...grpc.CallOption) (*payload.Object_Vector, error) {
			ovec, rerr := vc.GetObject(ctx, req, copts...)
			if rerr != nil {
				return nil, s.handleObjectBroadCastError(sspan, target, vald.GetObjectRPCName, uuid, req, rerr)
			}
			return ovec, nil
		})
	if err == nil && !ok {
		err = errors.ErrObjectNotFound(nil, uuid)
	}
	if err != nil {
		return errhandler.HandleError[payload.Object_Vector](span, spanErrorCode(err, codes.NotFound), err)
	}

	return vec, nil
}

func (s *server) getObject(ctx context.Context, uuid string) (*payload.Object_Vector, error) {
	return s.getObjectWithReq(ctx, uuid, nil)
}

func (s *server) GetObject(
	ctx context.Context, req *payload.Object_VectorRequest,
) (vec *payload.Object_Vector, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
		defer trace.End(span)
	}
	uuid := req.GetId().GetId()
	vec, err = s.getObjectWithReq(ctx, uuid, req)
	if err == nil {
		if vec != nil && vec.GetId() != "" && vec.GetVector() != nil {
			return vec, nil
		}
		err = errors.ErrObjectNotFound(nil, uuid)
	}
	return nil, s.wrapObjectReadError(ctx, span, vald.GetObjectRPCName, uuid, req, "object not found", err)
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	return grpc.HandleBidirectionalStream(
		stream,
		s.streamConcurrency,
		vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamGetObjectRPCName,
		apiName+"/"+vald.StreamGetObjectRPCName,
		func(req *payload.Object_VectorRequest) string { return "id-" + req.GetId().GetId() },
		s.GetObject,
		func(res *payload.Object_Vector) *payload.Object_StreamVector {
			return &payload.Object_StreamVector{
				Payload: &payload.Object_StreamVector_Vector{Vector: res},
			}
		},
		func(st *status.Status) *payload.Object_StreamVector {
			return &payload.Object_StreamVector{
				Payload: &payload.Object_StreamVector_Status{Status: st.Proto()},
			}
		},
	)
}

func (s *server) StreamListObject(
	req *payload.Object_List_Request, stream vald.Object_StreamListObjectServer,
) error {
	ctx := stream.Context()
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamListObjectRPCName), apiName+"/"+vald.StreamListObjectRPCName)
		defer trace.End(span)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// smu serializes stream.Send: BroadCast runs the per-target callbacks
	// concurrently and they all forward to the same caller stream. Each agent
	// client stream is only read by its own callback goroutine, so Recv needs no
	// lock and no per-message goroutine.
	var smu sync.Mutex
	err := s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		var sspan trace.Span
		if trace.Enabled() {
			ctx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.StreamListObjectRPCName+"/"+target)
			defer trace.End(sspan)
		}

		// Bind the client stream to a ctx canceled when this callback returns so
		// an early return (send failure) also tears down the agent stream.
		ctx, ccancel := context.WithCancel(ctx)
		defer ccancel()
		client, err := vc.StreamListObject(ctx, req, copts...)
		if err != nil {
			log.Errorf("failed to get StreamListObject client for agent(%s): %v", target, err)
			return err
		}

		// forwarded records whether any message of this agent stream has already
		// reached the caller. The gRPC client backoff re-invokes this callback on
		// retryable codes and the new agent stream restarts from the beginning, so a
		// failure after partial delivery is returned as Aborted (not retried) to
		// avoid sending duplicates to the caller.
		var forwarded bool
		fail := func(err error) error {
			if cerr := ctx.Err(); cerr != nil && !errors.Is(cerr, context.Canceled) {
				err = errors.Join(cerr, err)
			}
			if forwarded {
				return status.WrapWithAborted(vald.StreamListObjectRPCName+" API stream from "+target+" failed after partial delivery", err)
			}
			return err
		}
		for {
			res, err := client.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return fail(errors.ErrServerStreamClientRecv(err))
			}

			vec := res.GetVector()
			if vec == nil {
				st := res.GetStatus()
				log.Warnf("received empty vector: code %v: details %v: message %v",
					st.GetCode(),
					st.GetDetails(),
					st.GetMessage(),
				)
				continue
			}

			smu.Lock()
			err = stream.Send(res)
			smu.Unlock()
			if err != nil {
				cancel()
				if sspan != nil {
					st, msg, perr := status.ParseError(err, codes.Internal, "failed to parse StreamListObject send gRPC error response")
					errhandler.RecordSpanAttrs(sspan, trace.FromGRPCStatus(st.Code(), msg), perr)
				}
				return fail(errors.ErrServerStreamServerSend(err))
			}
			forwarded = true
		}
	})
	return err
}

func (s *server) GetTimestamp(
	ctx context.Context, req *payload.Object_TimestampRequest,
) (ts *payload.Object_Timestamp, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetTimestampRPCName), apiName+"/"+vald.GetTimestampRPCName)
		defer trace.End(span)
	}
	uuid := req.GetId().GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		log.Warn(err)
		return nil, s.wrapObjectReadError(ctx, span, vald.GetTimestampRPCName, uuid, req, "object not found", err)
	}
	treq := req
	ts, ok, err := broadcastFirst(ctx, s, "getTimestamp", true,
		func(ots *payload.Object_Timestamp) bool { return ots.GetId() != "" },
		func(ctx context.Context, sspan trace.Span, target string, vc vald.Client, copts ...grpc.CallOption) (*payload.Object_Timestamp, error) {
			ots, rerr := vc.GetTimestamp(ctx, treq, copts...)
			if rerr != nil {
				return nil, s.handleObjectBroadCastError(sspan, target, vald.GetTimestampRPCName, uuid, treq, rerr)
			}
			return ots, nil
		})
	if err == nil && !ok {
		err = errors.ErrObjectNotFound(nil, uuid)
	}
	if err != nil {
		return nil, s.wrapObjectReadError(ctx, span, vald.GetTimestampRPCName, uuid, req, "object not found", err)
	}
	return ts, nil
}

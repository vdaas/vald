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
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
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
	case errors.Is(err, context.Canceled),
		errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
		attrs = trace.StatusCodeCancelled(
			errdetails.ValdGRPCResourceTypePrefix +
				"/vald.v1." + rpcName + ".BroadCast/" +
				target + " canceled: " + err.Error())
		code = codes.Canceled
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
		attrs = trace.StatusCodeDeadlineExceeded(
			errdetails.ValdGRPCResourceTypePrefix +
				"/vald.v1." + rpcName + ".BroadCast/" +
				target + " deadline_exceeded: " + err.Error())
		code = codes.DeadlineExceeded
	default:
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
		attrs = trace.FromGRPCStatus(code, msg)
	}
	errhandler.RecordSpanAttrs(sspan, attrs, err)
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

// broadcastFirst broadcasts call to every READ target and returns the first
// result accepted by hit, canceling the remaining targets once it arrives. name
// labels the per-target spans (apiName/<name>/BroadCast/<target>). call must
// return the zero T together with any error; a nil error with a missing object
// is reported by returning a result rejected by hit. ok is false (with a nil
// err) when no target had the object.
func broadcastFirst[T any](
	ctx context.Context,
	s *server,
	name string,
	hit func(T) bool,
	call func(ctx context.Context, sspan trace.Span, target string, vc vald.Client, copts ...grpc.CallOption) (T, error),
) (res T, ok bool, err error) {
	rch := make(chan T, 1)
	ech := make(chan error, 1)
	doneErr := errors.New("done " + name)
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(rch)
		defer close(ech)
		var once sync.Once
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+name+"/BroadCast/"+target)
			defer trace.End(sspan)
			r, cerr := call(sctx, sspan, target, vc, copts...)
			if cerr != nil {
				return cerr
			}
			if hit(r) {
				once.Do(func() {
					rch <- r
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
			case res = <-rch:
				err = nil
			default:
			}
		}
	case r, received := <-rch:
		if received {
			res = r
		} else {
			// rch was closed without a hit: the broadcast result is already
			// buffered in ech, so report it instead of a blanket NotFound.
			err = <-ech
		}
	case err = <-ech:
		// The hit is buffered in rch before BroadCast returns, so prefer it over
		// the (nil or canceled) broadcast result.
		select {
		case r := <-rch:
			if hit(r) {
				res, err = r, nil
			}
		default:
		}
	}
	if err != nil || !hit(res) {
		var zero T
		return zero, false, err
	}
	return res, true, nil
}

func (s *server) exists(ctx context.Context, uuid string) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "exists"), apiName+"/"+vald.ExistsRPCName+"/exists")
	defer trace.End(span)

	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		log.Warn(err)
		return errhandler.HandleError[payload.Object_ID](span, codes.InvalidArgument, err)
	}
	meta := &payload.Object_ID{
		Id: uuid,
	}
	id, ok, err := broadcastFirst(ctx, s, "exists",
		func(oid *payload.Object_ID) bool { return oid.GetId() != "" },
		func(ctx context.Context, sspan trace.Span, target string, vc vald.Client, copts ...grpc.CallOption) (*payload.Object_ID, error) {
			oid, rerr := vc.Exists(ctx, meta, copts...)
			if rerr != nil {
				return nil, s.handleObjectBroadCastError(sspan, target, vald.ExistsRPCName, uuid, meta, rerr)
			}
			return oid, nil
		})
	if err == nil && !ok {
		err = errors.ErrObjectIDNotFound(uuid)
	}
	if err != nil {
		return errhandler.HandleError[payload.Object_ID](span, codes.NotFound, err)
	}
	return id, nil
}

func (s *server) Exists(
	ctx context.Context, meta *payload.Object_ID,
) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
	defer trace.End(span)
	uuid := meta.GetId()
	id, err = s.exists(ctx, uuid)
	if err == nil {
		return id, nil
	}
	details := s.errDetails(ctx, vald.ExistsRPCName, uuid, meta)
	var attrs trace.Attributes
	switch {
	case errors.Is(err, errors.ErrInvalidUUID(uuid)):
		err = status.WrapWithInvalidArgument(vald.ExistsRPCName+" API invalid argument for uuid \""+uuid+"\" detected", err, append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       uuidField,
					Description: err.Error(),
				},
			},
		})...)
		attrs = trace.StatusCodeInvalidArgument(err.Error())
	case errors.Is(err, errors.ErrObjectIDNotFound(uuid)):
		err = status.WrapWithNotFound(vald.ExistsRPCName+" API id "+uuid+"'s data not found", err, details...)
		attrs = trace.StatusCodeNotFound(err.Error())
	case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
		err = status.WrapWithInternal(vald.ExistsRPCName+" API connection not found", err, details...)
		attrs = trace.StatusCodeInternal(err.Error())
	case errors.Is(err, context.Canceled):
		err = status.WrapWithCanceled(vald.ExistsRPCName+" API canceled", err, details...)
		attrs = trace.StatusCodeCancelled(err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API deadline exceeded", err, details...)
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
	errhandler.RecordSpanAttrs(span, attrs, err)
	return nil, err
}

func (s *server) getObject(
	ctx context.Context, uuid string,
) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "getObject"), apiName+"/"+vald.GetObjectRPCName+"/getObject")
	defer trace.End(span)
	req := &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: uuid,
		},
	}
	vec, ok, err := broadcastFirst(ctx, s, "getObject",
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
		return errhandler.HandleError[payload.Object_Vector](span, codes.NotFound, err)
	}

	return vec, nil
}

func (s *server) GetObject(
	ctx context.Context, req *payload.Object_VectorRequest,
) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer trace.End(span)
	uuid := req.GetId().GetId()
	vec, err = s.getObject(ctx, uuid)
	if err == nil {
		if vec != nil && vec.GetId() != "" && vec.GetVector() != nil {
			return vec, nil
		}
		err = errors.ErrObjectNotFound(nil, uuid)
	}
	details := s.errDetails(ctx, vald.GetObjectRPCName, uuid, req)
	var attrs trace.Attributes
	switch {
	case errors.Is(err, errors.ErrInvalidUUID(uuid)):
		err = status.WrapWithInvalidArgument(vald.GetObjectRPCName+" API invalid argument for uuid \""+uuid+"\" detected", err, append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       uuidField,
					Description: err.Error(),
				},
			},
		})...)
		attrs = trace.StatusCodeInvalidArgument(err.Error())
	case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
		err = status.WrapWithNotFound(vald.GetObjectRPCName+" API id "+uuid+"'s object not found", err, details...)
		attrs = trace.StatusCodeNotFound(err.Error())
	case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
		err = status.WrapWithInternal(vald.GetObjectRPCName+" API connection not found", err, details...)
		attrs = trace.StatusCodeInternal(err.Error())
	case errors.Is(err, context.Canceled):
		err = status.WrapWithCanceled(vald.GetObjectRPCName+" API canceled", err, details...)
		attrs = trace.StatusCodeCancelled(err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.WrapWithDeadlineExceeded(vald.GetObjectRPCName+" API deadline exceeded", err, details...)
		attrs = trace.StatusCodeDeadlineExceeded(err.Error())
	default:
		st, _ := status.FromError(err)
		if st != nil {
			attrs = trace.FromGRPCStatus(st.Code(), st.Message())
		}
	}
	log.Debug(err)
	errhandler.RecordSpanAttrs(span, attrs, err)
	return nil, err
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamGetObjectRPCName), apiName+"/"+vald.StreamGetObjectRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Object_VectorRequest) (*payload.Object_StreamVector, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamGetObjectRPCName+"/id-"+req.GetId().GetId())
			defer trace.End(sspan)
			res, err := s.GetObject(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				errhandler.RecordSpanStatus(sspan, st, err)
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
	defer trace.End(span)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// smu serializes stream.Send: BroadCast runs the per-target callbacks
	// concurrently and they all forward to the same caller stream. Each agent
	// client stream is only read by its own callback goroutine, so Recv needs no
	// lock and no per-message goroutine.
	var smu sync.Mutex
	err := s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.StreamListObjectRPCName+"/"+target)
		defer trace.End(sspan)

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
				if sspan != nil {
					st, msg, perr := status.ParseError(err, codes.Internal, "failed to parse StreamListObject send gRPC error response")
					sspan.RecordError(perr)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, perr.Error())
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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetTimestampRPCName), apiName+"/"+vald.GetTimestampRPCName)
	defer trace.End(span)
	uuid := req.GetId().GetId()
	treq := &payload.Object_TimestampRequest{
		Id: &payload.Object_ID{
			Id: uuid,
		},
	}
	ts, ok, err := broadcastFirst(ctx, s, "getTimestamp",
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
		details := s.errDetails(ctx, vald.GetTimestampRPCName, uuid, req)
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrInvalidUUID(uuid)):
			err = status.WrapWithInvalidArgument(vald.GetTimestampRPCName+" API invalid argument for uuid \""+uuid+"\" detected", err, append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       uuidField,
						Description: err.Error(),
					},
				},
			})...)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
			err = status.WrapWithNotFound(vald.GetTimestampRPCName+" API id "+uuid+"'s object not found", err, details...)
			attrs = trace.StatusCodeNotFound(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.GetTimestampRPCName+" API connection not found", err, details...)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.GetTimestampRPCName+" API canceled", err, details...)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.GetTimestampRPCName+" API deadline exceeded", err, details...)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			st, _ := status.FromError(err)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	return ts, nil
}

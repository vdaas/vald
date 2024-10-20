package service

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	apiName        = "vald/index/job/delete"
	grpcMethodName = "vald.v1.Remove/" + vald.RemoveRPCName
)

// Deleter represents an interface for deleting.
type Deleter interface {
	StartClient(ctx context.Context) (<-chan error, error)
	Start(ctx context.Context) error
}

var defaultOpts = []Option{
	WithIndexingConcurrency(1),
}

type index struct {
	client        discoverer.Client
	targetAddrs   []string
	targetIndexID string

	concurrency int
}

// New returns Deleter object if no error occurs.
func New(opts ...Option) (Deleter, error) {
	idx := new(index)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(idx); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(err)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	idx.targetAddrs = delDuplicateAddrs(idx.targetAddrs)
	return idx, nil
}

func delDuplicateAddrs(targetAddrs []string) []string {
	addrs := make([]string, 0, len(targetAddrs))
	exist := make(map[string]bool)

	for _, addr := range targetAddrs {
		if !exist[addr] {
			addrs = append(addrs, addr)
			exist[addr] = true
		}
	}
	return addrs
}

// StartClient starts the gRPC client.
func (idx *index) StartClient(ctx context.Context) (<-chan error, error) {
	return idx.client.Start(ctx)
}

func (idx *index) Start(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, apiName+"/service/index.Delete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err := idx.doDeleteIndex(ctx,
		func(ctx context.Context, rc vald.RemoveClient, copts ...grpc.CallOption) (*payload.Object_Location, error) {
			return rc.Remove(ctx, &payload.Remove_Request{
				Id: &payload.Object_ID{
					Id: idx.targetIndexID,
				},
			}, copts...)
		},
	)
	if err != nil {
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.RemoveRPCName+" API connection not found", err,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCTargetAddrNotFound):
			err = status.WrapWithInternal(
				vald.RemoveRPCName+" API connection target address \""+strings.Join(idx.targetAddrs, ",")+"\" not found", err,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response",
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (idx *index) doDeleteIndex(
	ctx context.Context,
	fn func(_ context.Context, _ vald.RemoveClient, _ ...grpc.CallOption) (*payload.Object_Location, error),
) (errs error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, grpcMethodName), apiName+"/service/index.doDeleteIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	targetAddrs := idx.client.GetAddrs(ctx)
	if len(idx.targetAddrs) != 0 {
		// If target addresses is specified, that addresses are used in priority.
		for _, addr := range idx.targetAddrs {
			log.Infof("connect to target agent (%s)", addr)
			if _, err := idx.client.GetClient().Connect(ctx, addr); err != nil {
				return err
			}
		}
		targetAddrs = idx.targetAddrs
	}
	log.Infof("target agent addrs: %v", targetAddrs)

	var emu sync.Mutex
	err := idx.client.GetClient().OrderedRangeConcurrent(ctx, targetAddrs, idx.concurrency,
		func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "OrderedRangeConcurrent/"+target), vald.RemoveRPCName+"/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			_, err := fn(ctx, vald.NewRemoveClient(conn), copts...)
			if err != nil {
				var attrs trace.Attributes
				switch {
				case errors.Is(err, context.Canceled):
					err = status.WrapWithCanceled(
						vald.RemoveRPCName+" API canceld", err,
					)
					attrs = trace.StatusCodeCancelled(err.Error())
				case errors.Is(err, context.DeadlineExceeded):
					err = status.WrapWithCanceled(
						vald.RemoveRPCName+" API deadline exceeded", err,
					)
					attrs = trace.StatusCodeDeadlineExceeded(err.Error())
				case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
					err = status.WrapWithInternal(
						vald.RemoveRPCName+" API connection not found", err,
					)
					attrs = trace.StatusCodeInternal(err.Error())
				case errors.Is(err, errors.ErrTargetNotFound):
					err = status.WrapWithInvalidArgument(
						vald.RemoveRPCName+" API target not found", err,
					)
					attrs = trace.StatusCodeInternal(err.Error())
				default:
					var (
						st  *status.Status
						msg string
					)
					st, msg, err = status.ParseError(err, codes.Internal,
						"failed to parse "+vald.RemoveRPCName+" gRPC error response",
					)
					if st != nil && err != nil && st.Code() == codes.FailedPrecondition {
						log.Warnf("DeleteIndex of %s skipped, indexID:%s, message: %s, err: %v", target, idx.targetIndexID, st.Message(), errors.Join(st.Err(), err))
						return nil
					}
					if st != nil && err != nil && st.Code() == codes.NotFound {
						log.Warn("DeleteIndex of %s skipped, indexID: %s, message: %s, err: %v", target, idx.targetIndexID, st.Message(), errors.Join(st.Err(), err))
						return nil
					}
					attrs = trace.FromGRPCStatus(st.Code(), msg)
				}
				log.Warnf("an error occurred in (%s) deleting index(%s): %v", target, idx.targetIndexID, err)
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(attrs...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
			}
			return err
		},
	)
	return errors.Join(err, errs)
}

// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

const (
	apiName        = "vald/index/job/save"
	grpcMethodName = "core.v1.Agent/" + agent.SaveIndexRPCName
)

// Indexer represents an interface for indexing.
type Indexer interface {
	StartClient(ctx context.Context) (<-chan error, error)
	Start(ctx context.Context) error
}

type index struct {
	client         discoverer.Client
	targetAddrs    []string
	targetAddrList map[string]bool

	concurrency int
}

// New returns Indexer object if no error occurs.
func New(opts ...Option) (Indexer, error) {
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
	idx.targetAddrList = make(map[string]bool, len(idx.targetAddrs))
	for _, addr := range idx.targetAddrs {
		idx.targetAddrList[addr] = true
	}
	return idx, nil
}

// StartClient starts the gRPC client.
func (idx *index) StartClient(ctx context.Context) (<-chan error, error) {
	return idx.client.Start(ctx)
}

// Start starts indexing process.
func (idx *index) Start(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, apiName+"/service/index.Start")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err := idx.doSaveIndex(ctx,
		func(ctx context.Context, ac agent.AgentClient, copts ...grpc.CallOption) (*payload.Empty, error) {
			return ac.SaveIndex(ctx, &payload.Empty{}, copts...)
		},
	)
	if err != nil {
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				agent.SaveIndexRPCName+" API connection not found", err,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCTargetAddrNotFound):
			err = status.WrapWithInternal(
				agent.SaveIndexRPCName+" API connection target address \""+strings.Join(idx.targetAddrs, ",")+"\" not found", err,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+agent.SaveIndexRPCName+" gRPC error response",
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

func (idx *index) doSaveIndex(ctx context.Context, fn func(_ context.Context, _ agent.AgentClient, _ ...grpc.CallOption) (*payload.Empty, error)) (errs error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, grpcMethodName), apiName+"/service/index.doSaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	targetAddrs := idx.client.GetAddrs(ctx)
	if len(idx.targetAddrs) != 0 {
		targetAddrs = idx.extractTargetAddrs(targetAddrs)

		// If targetAddrs is empty, an invalid target addresses may be registered in targetAddrList.
		if len(targetAddrs) == 0 {
			return errors.ErrGRPCTargetAddrNotFound
		}
	}
	log.Infof("target agent addrs: %v", targetAddrs)

	var emu sync.Mutex
	err := idx.client.GetClient().OrderedRangeConcurrent(ctx, targetAddrs, idx.concurrency,
		func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "OrderedRangeConcurrent/"+target), agent.SaveIndexRPCName+"/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			_, err := fn(ctx, agent.NewAgentClient(conn), copts...)
			if err != nil {
				var attrs trace.Attributes
				switch {
				case errors.Is(err, context.Canceled):
					err = status.WrapWithCanceled(
						agent.SaveIndexRPCName+" API canceld", err,
					)
					attrs = trace.StatusCodeCancelled(err.Error())
				case errors.Is(err, context.DeadlineExceeded):
					err = status.WrapWithCanceled(
						agent.SaveIndexRPCName+" API deadline exceeded", err,
					)
					attrs = trace.StatusCodeDeadlineExceeded(err.Error())
				case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
					err = status.WrapWithInternal(
						agent.SaveIndexRPCName+" API connection not found", err,
					)
					attrs = trace.StatusCodeInternal(err.Error())
				case errors.Is(err, errors.ErrTargetNotFound):
					err = status.WrapWithInvalidArgument(
						agent.SaveIndexRPCName+" API target not found", err,
					)
					attrs = trace.StatusCodeInternal(err.Error())
				default:
					var (
						st  *status.Status
						msg string
					)
					st, msg, err = status.ParseError(err, codes.Internal,
						"failed to parse "+agent.SaveIndexRPCName+" gRPC error response",
					)
					attrs = trace.FromGRPCStatus(st.Code(), msg)
				}
				log.Warnf("an error occurred in (%s) during save indexing: %v", target, err)
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

// extractTargetAddresses filters and extracts target addresses registered in targetAddrList from the given address list.
func (idx *index) extractTargetAddrs(addrs []string) []string {
	res := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		if !idx.targetAddrList[addr] {
			log.Warnf("the gRPC target address not found: %s", addr)
		} else {
			res = append(res, addr)
		}
	}
	return res
}

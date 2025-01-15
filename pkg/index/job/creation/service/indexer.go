// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	apiName        = "vald/index/job/create"
	grpcMethodName = "core.v1.Agent/" + agent.CreateIndexRPCName
)

// Indexer represents an interface for indexing.
type Indexer interface {
	StartClient(ctx context.Context) (<-chan error, error)
	Start(ctx context.Context) error
}

type index struct {
	client      discoverer.Client
	targetAddrs []string

	creationPoolSize uint32
	concurrency      int
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

// Start starts indexing process.
func (idx *index) Start(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, apiName+"/service/index.Start")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err := idx.doCreateIndex(ctx,
		func(ctx context.Context, ac agent.AgentClient, copts ...grpc.CallOption) (*payload.Empty, error) {
			return ac.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
				PoolSize: idx.creationPoolSize,
			}, copts...)
		},
	)
	if err != nil {
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				agent.CreateIndexRPCName+" API connection not found", err,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCTargetAddrNotFound):
			err = status.WrapWithInternal(
				agent.CreateIndexRPCName+" API connection target address \""+strings.Join(idx.targetAddrs, ",")+"\" not found", err,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+agent.CreateIndexRPCName+" gRPC error response",
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

// skipcq: GO-R1005
func (idx *index) doCreateIndex(
	ctx context.Context,
	fn func(_ context.Context, _ agent.AgentClient, _ ...grpc.CallOption) (*payload.Empty, error),
) (errs error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, grpcMethodName), apiName+"/service/index.doCreateIndex")
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
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "OrderedRangeConcurrent/"+target), agent.CreateIndexRPCName+"/"+target)
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
						agent.CreateIndexRPCName+" API canceled", err,
					)
					attrs = trace.StatusCodeCancelled(err.Error())
				case errors.Is(err, context.DeadlineExceeded):
					err = status.WrapWithCanceled(
						agent.CreateIndexRPCName+" API deadline exceeded", err,
					)
					attrs = trace.StatusCodeDeadlineExceeded(err.Error())
				case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
					err = status.WrapWithInternal(
						agent.CreateIndexRPCName+" API connection not found", err,
					)
					attrs = trace.StatusCodeInternal(err.Error())
				case errors.Is(err, errors.ErrTargetNotFound):
					err = status.WrapWithInvalidArgument(
						agent.CreateIndexRPCName+" API target not found", err,
					)
					attrs = trace.StatusCodeInternal(err.Error())
				default:
					var (
						st  *status.Status
						msg string
					)
					st, msg, err = status.ParseError(err, codes.Internal,
						"failed to parse "+agent.CreateIndexRPCName+" gRPC error response",
					)
					if st != nil && err != nil && st.Code() == codes.FailedPrecondition {
						log.Warnf("CreateIndex of %s skipped, message: %s, err: %v", target, st.Message(), errors.Join(st.Err(), err))
						return nil
					}
					attrs = trace.FromGRPCStatus(st.Code(), msg)
				}
				log.Warnf("an error occurred in (%s) during indexing: %v", target, err)
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

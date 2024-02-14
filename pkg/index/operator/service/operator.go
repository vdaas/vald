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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	apiName        = "vald/index/operator"
	grpcMethodName = "core.v1.Agent/" + agent.CreateIndexRPCName
)

// Operator represents an interface for indexing.
type Operator interface {
	Start(ctx context.Context) error
}

type index struct {
	targetAddrs    []string
	targetAddrList map[string]bool

	creationPoolSize uint32
	concurrency      int
}

// New returns Indexer object if no error occurs.
func New(opts ...Option) (Operator, error) {
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

// Start starts indexing process.
func (idx *index) Start(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, apiName+"/service/index.Start")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	return nil
}

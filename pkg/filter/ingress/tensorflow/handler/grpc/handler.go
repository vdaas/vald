//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"fmt"

	"github.com/vdaas/vald/apis/grpc/v1/filter/ingress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/core/converter/tensorflow"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

type Server ingress.FilterServer

type server struct {
	*ingress.UnimplementedFilterServer
	tf tensorflow.TF
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) GenVector(ctx context.Context, req *payload.Object_Blob) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/.GetVector")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	f64vec, err := s.tf.GetVector(string(req.GetObject()))
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("GenVector API id %s's object could not vectorize", req.GetId()), err, info.Get())
	}

	vec = &payload.Object_Vector{
		Vector: make([]float32, 0, len(f64vec)),
	}
	for i, d := range f64vec {
		vec.Vector[i] = float32(d)
	}
	return nil, nil
}

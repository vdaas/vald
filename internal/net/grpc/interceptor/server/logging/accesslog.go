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

package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

const (
	grpcKindUnary  = "unary"
	grpcKindStream = "stream"

	rpcCompletedMessage = "rpc completed"
	rpcFailedMessage    = "rpc failed"
)

var grpcEntityCache sync.Map[string, *AccessLogGRPCEntity]

func getCachedGRPCEntity(kind, fullMethod string) *AccessLogGRPCEntity {
	key := kind + ":" + fullMethod
	if val, ok := grpcEntityCache.Load(key); ok {
		return val
	}
	service, method := parseMethod(fullMethod)
	entity := &AccessLogGRPCEntity{
		Kind:    kind,
		Service: service,
		Method:  method,
	}
	actual, _ := grpcEntityCache.LoadOrStore(key, entity)
	return actual
}

type AccessLogEntity struct {
	Error     error                `json:"error,omitempty"     yaml:"error"`
	GRPC      *AccessLogGRPCEntity `json:"grpc,omitempty"      yaml:"grpc"`
	TraceID   string               `json:"traceID,omitempty"   yaml:"traceID"`
	StartTime int64                `json:"startTime,omitempty" yaml:"startTime"`
	EndTime   int64                `json:"endTime,omitempty"   yaml:"endTime"`
	Latency   int64                `json:"latency,omitempty"   yaml:"latency"`
}

type AccessLogGRPCEntity struct {
	Kind    string `json:"kind,omitempty"    yaml:"kind"`
	Service string `json:"service,omitempty" yaml:"service"`
	Method  string `json:"method,omitempty"  yaml:"method"`
}

func (e AccessLogEntity) String() (str string) {
	var emsg string
	if e.Error != nil {
		st, ok := status.FromError(e.Error)
		if ok && st != nil {
			emsg = st.String()
		} else {
			emsg = e.Error.Error()
		}
	}
	eb, err := json.Marshal(e)
	if err != nil {
		str = fmt.Sprintf("%#v,\tfailed to json.Marshal(AccessLogEntity) error: %v", e, err)
	} else {
		str = string(eb)
	}

	if emsg != "" {
		return str + ",\terror message: " + emsg
	}
	return str
}

// lazyAccessLog defers AccessLogEntity.String (a JSON encoding) until the logger
// actually formats its arguments, so that it is not paid when the level is disabled.
// fmt-based backends render it through String, and the glg JSON format encodes it
// through MarshalText as the same JSON string the pre-rendered value produced.
// Being a defined type, it does not inherit AccessLogEntity's methods.
type lazyAccessLog AccessLogEntity

func (l *lazyAccessLog) String() string {
	return (*AccessLogEntity)(l).String()
}

func (l *lazyAccessLog) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func AccessLogInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		var traceID string

		span := trace.FromContext(ctx)
		if span != nil {
			traceID = span.SpanContext().TraceID().String()
		}

		start := time.Now()

		resp, err = handler(ctx, req)

		end := time.Now()

		grpcEntity := getCachedGRPCEntity(grpcKindUnary, info.FullMethod)

		entity := &AccessLogEntity{
			GRPC:      grpcEntity,
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Latency:   end.Sub(start).Nanoseconds(),
		}

		if traceID != "" {
			entity.TraceID = traceID
		}

		if err != nil {
			entity.Error = err
			log.Warn(rpcFailedMessage, (*lazyAccessLog)(entity))
		} else {
			log.Debug(rpcCompletedMessage, (*lazyAccessLog)(entity))
		}

		return resp, err
	}
}

func AccessLogStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		var traceID string

		span := trace.FromContext(ss.Context())
		if span != nil {
			traceID = span.SpanContext().TraceID().String()
		}

		start := time.Now()

		err := handler(srv, ss)

		end := time.Now()

		grpcEntity := getCachedGRPCEntity(grpcKindStream, info.FullMethod)

		entity := &AccessLogEntity{
			GRPC:      grpcEntity,
			StartTime: start.UnixNano(),
			EndTime:   end.UnixNano(),
			Latency:   end.Sub(start).Nanoseconds(),
		}

		if traceID != "" {
			entity.TraceID = traceID
		}

		if err != nil {
			entity.Error = err
			log.Warn(rpcFailedMessage, (*lazyAccessLog)(entity))
		} else {
			log.Debug(rpcCompletedMessage, (*lazyAccessLog)(entity))
		}

		return err
	}
}

func parseMethod(fullMethod string) (service, method string) {
	if len(fullMethod) == 0 {
		return "", ""
	}
	s := fullMethod
	if s[0] == '/' {
		s = s[1:]
	}
	idx := strings.LastIndexByte(s, '/')
	if idx < 0 {
		return "", s
	}
	return s[:idx], s[idx+1:]
}

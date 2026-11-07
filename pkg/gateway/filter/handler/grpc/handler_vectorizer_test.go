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
	"testing"
	"time"

	ingressapi "github.com/vdaas/vald/apis/grpc/v1/filter/ingress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/comparator"
)

// vectorizerIngressStub records the addresses passed to Target and vectorizes
// every object to genVec.
type vectorizerIngressStub struct {
	ingress.Client
	targets  [][]string
	released [][]string
	mu       sync.Mutex
}

type vectorizerFilterStub struct {
	ingressapi.FilterClient
}

func genVec() []float32 { return []float32{0.1, 0.2} }

const (
	confVectorizer = "conf:8080"
	testVectorizer = "vectorizer:8080"
)

func (v *vectorizerIngressStub) Target(
	_ context.Context, addrs ...string,
) (ingressapi.FilterClient, error) {
	v.mu.Lock()
	v.targets = append(v.targets, addrs)
	v.mu.Unlock()
	return &vectorizerFilterStub{}, nil
}

func (v *vectorizerIngressStub) Release(_ context.Context, addrs ...string) {
	v.mu.Lock()
	v.released = append(v.released, addrs)
	v.mu.Unlock()
}

func (*vectorizerFilterStub) GenVector(
	context.Context, *payload.Object_Blob, ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return &payload.Object_Vector{Vector: genVec()}, nil
}

func Test_server_vectorizerAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		vr         *payload.Filter_Target
		name       string
		configured string
		want       string
	}{
		{name: "request target wins over config", vr: &payload.Filter_Target{Host: "vec", Port: 9000}, configured: confVectorizer, want: "vec:9000"},
		{name: "request target without host defaults to localhost", vr: &payload.Filter_Target{Port: 9000}, want: "localhost:9000"},
		{name: "nil request target falls back to config", configured: confVectorizer, want: confVectorizer},
		{name: "request target without port falls back to config", vr: &payload.Filter_Target{Host: "vec"}, configured: confVectorizer, want: confVectorizer},
		{name: "neither set", vr: &payload.Filter_Target{Host: "portless"}, want: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			s := &server{Vectorizer: tc.configured}
			if got := s.vectorizerAddr(tc.vr); got != tc.want {
				tt.Errorf("vectorizerAddr() = %q, want %q", got, tc.want)
			}
		})
	}
}

// The *Object RPCs must use the configured ingress_filter.vectorizer when the
// request names no vectorizer, instead of rejecting the request.
func Test_server_searchObject_configuredVectorizer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		call       func(context.Context, *server, *payload.Search_ObjectRequest) error
		name       string
		configured string
		wantTarget [][]string
		wantCode   codes.Code
	}{
		{
			name: "SearchObject uses the configured vectorizer",
			call: func(ctx context.Context, s *server, req *payload.Search_ObjectRequest) error {
				_, err := s.SearchObject(ctx, req)
				return err
			},
			configured: testVectorizer,
			wantCode:   codes.OK,
			wantTarget: [][]string{{testVectorizer}},
		},
		{
			name: "LinearSearchObject uses the configured vectorizer",
			call: func(ctx context.Context, s *server, req *payload.Search_ObjectRequest) error {
				_, err := s.LinearSearchObject(ctx, req)
				return err
			},
			configured: testVectorizer,
			wantCode:   codes.OK,
			wantTarget: [][]string{{testVectorizer}},
		},
		{
			name: "SearchObject without any vectorizer is rejected",
			call: func(ctx context.Context, s *server, req *payload.Search_ObjectRequest) error {
				_, err := s.SearchObject(ctx, req)
				return err
			},
			wantCode: codes.InvalidArgument,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()
			ing := &vectorizerIngressStub{}
			gw := &recordSearchGatewayStub{}
			s := &server{
				gateway:    gw,
				ingress:    ing,
				Vectorizer: tc.configured,
			}
			done := make(chan error, 1)
			go func() {
				done <- tc.call(ctx, s, &payload.Search_ObjectRequest{Object: []byte("obj")})
			}()
			var err error
			select {
			case err = <-done:
			case <-ctx.Done():
				tt.Fatal("timed out waiting for the handler")
			}
			code := codes.OK
			if err != nil {
				st, _ := status.FromError(err)
				code = st.Code()
			}
			if code != tc.wantCode {
				tt.Fatalf("code = %v, want %v (err = %v)", code, tc.wantCode, err)
			}
			ing.mu.Lock()
			defer ing.mu.Unlock()
			if diff := comparator.Diff(ing.targets, tc.wantTarget); diff != "" {
				tt.Errorf("Target addrs mismatch (-got +want):\n%s", diff)
			}
			// Every target resolved for the request is released when it ends.
			if diff := comparator.Diff(ing.released, tc.wantTarget); diff != "" {
				tt.Errorf("Release addrs mismatch (-got +want):\n%s", diff)
			}
			if tc.wantCode != codes.OK {
				return
			}
			gw.mu.Lock()
			defer gw.mu.Unlock()
			if diff := comparator.Diff(gw.vecs, [][]float32{genVec()}); diff != "" {
				tt.Errorf("forwarded vectors mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

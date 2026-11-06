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
package service

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/db/kvs/pogreb"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/strings"
	"google.golang.org/grpc"
)

// fakeListStream returns one vector, then signals recvBlocked and blocks the
// next Recv until the stream ctx is canceled, returning the Canceled status
// grpc-go reports in that case.
type fakeListStream struct {
	grpc.ClientStream
	done        <-chan struct{}
	recvBlocked chan struct{}
	sent        bool
}

func (f *fakeListStream) Recv() (*payload.Object_List_Response, error) {
	if !f.sent {
		f.sent = true
		return &payload.Object_List_Response{
			Payload: &payload.Object_List_Response_Vector{
				Vector: &payload.Object_Vector{Id: "id-1", Vector: []float32{1}, Timestamp: 1},
			},
		}, nil
	}
	close(f.recvBlocked)
	<-f.done
	return nil, status.Error(codes.Canceled, context.Canceled.Error())
}

// fakeGateway implements only StreamListObject; any other call panics.
type fakeGateway struct {
	vc.Client
	recvBlocked chan struct{}
}

func (f *fakeGateway) StreamListObject(
	ctx context.Context, _ *payload.Object_List_Request, _ ...grpc.CallOption,
) (vald.Object_StreamListObjectClient, error) {
	return &fakeListStream{done: ctx.Done(), recvBlocked: f.recvBlocked}, nil
}

// failingSetDB fails Set, but only once the main loop is blocked in Recv, so
// the worker error always cancels a pending Recv.
type failingSetDB struct {
	pogreb.DB
	recvBlocked chan struct{}
	err         error
}

func (*failingSetDB) Get(string) ([]byte, bool, error) { return nil, false, nil }

func (f *failingSetDB) Set(string, []byte) error {
	select {
	case <-f.recvBlocked:
	case <-time.After(10 * time.Second):
	}
	return f.err
}

func TestDoExportIndex_WorkerErrorDecidesStatus(t *testing.T) {
	t.Parallel()
	recvBlocked := make(chan struct{})
	setErr := errors.New("pogreb set failed")
	e := &export{
		gateway:               &fakeGateway{recvBlocked: recvBlocked},
		storedVector:          &failingSetDB{recvBlocked: recvBlocked, err: setErr},
		streamListConcurrency: 2,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- e.doExportIndex(t.Context())
	}()
	var err error
	select {
	case err = <-errCh:
	case <-time.After(30 * time.Second):
		t.Fatal("doExportIndex did not return")
	}

	if err == nil {
		t.Fatal("doExportIndex must fail when a worker fails")
	}
	st, ok := status.FromError(err)
	if !ok || st == nil {
		t.Fatalf("doExportIndex must return a gRPC status error, got: %v", err)
	}
	if st.Code() != codes.Internal {
		t.Errorf("status code = %v, want %v (the Canceled Recv caused by our own cancellation must not decide it): %v", st.Code(), codes.Internal, err)
	}
	if strings.Contains(err.Error(), context.Canceled.Error()) {
		t.Errorf("error must not carry the self-inflicted cancellation, got: %v", err)
	}
}

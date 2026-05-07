//go:build e2e

//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package crud

import (
	"context"
	"testing"
	"time"

	embedderpb "github.com/vdaas/vald/apis/grpc/v1/embedder"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (r *runner) processEmbedding(t *testing.T, ctx context.Context, plan *config.Execution) error {
	t.Helper()
	if plan == nil {
		return errors.New("embedding plan is nil")
	}
	if plan.Text == "" {
		return errors.New("embedding text is empty")
	}
	res, err := grpc.RoundRobin(ctx, r.client.GRPCClient(), func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
		return embedderpb.NewEmbedderClient(conn).Embedding(ctx, &embedderpb.Text{Text: plan.Text}, copts...)
	})
	if err != nil {
		return err
	}
	vec, ok := res.(interface{ GetVector() []float32 })
	if !ok || vec == nil || len(vec.GetVector()) == 0 {
		return errors.New("embedding response vector is empty")
	}
	t.Logf("embedding vector length=%d", len(vec.GetVector()))
	return nil
}

func (r *runner) processEmbedderInsert(
	t *testing.T, ctx context.Context, plan *config.Execution,
) error {
	t.Helper()
	if plan == nil {
		return errors.New("embedder insert plan is nil")
	}
	if plan.Text == "" {
		return errors.New("embedder insert text is empty")
	}

	id := plan.DocumentID
	if id == "" {
		id = "e2e-embedder-" + time.Now().UTC().Format("20060102150405.000000000")
	}

	res, err := grpc.RoundRobin(ctx, r.client.GRPCClient(), func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
		return embedderpb.NewEmbedderClient(conn).Insert(ctx, &embedderpb.InsertRequest{
			Document: &embedderpb.Document{
				Id:        id,
				Text:      plan.Text,
				Timestamp: time.Now().UnixNano(),
			},
		}, copts...)
	})
	if err != nil {
		return err
	}
	loc, ok := res.(*payload.Object_Location)
	if !ok || loc == nil || (loc.GetUuid() == "" && loc.GetName() == "") {
		return errors.New("embedder insert response location is empty")
	}
	t.Logf("embedder insert location name=%s uuid=%s", loc.GetName(), loc.GetUuid())
	return nil
}

func (r *runner) processEmbedderMutation(
	t *testing.T, ctx context.Context, plan *config.Execution,
) error {
	t.Helper()
	if plan == nil {
		return errors.New("embedder mutation plan is nil")
	}
	id := plan.DocumentID
	if id == "" {
		id = "e2e-embedder-" + time.Now().UTC().Format("20060102150405.000000000")
	}
	metaAny, err := anypb.New(&wrapperspb.StringValue{Value: "e2e-meta"})
	if err != nil {
		return err
	}
	meta := &payload.Meta_Value{Value: metaAny}
	call := func(f func(context.Context, *grpc.ClientConn, ...grpc.CallOption) (any, error)) (any, error) {
		return grpc.RoundRobin(ctx, r.client.GRPCClient(), f)
	}
	switch plan.Type {
	case config.OpEmbedderInsertWithMetadata:
		return checkLocation(call(func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			return embedderpb.NewEmbedderClient(conn).InsertWithMetadata(ctx, &embedderpb.InsertWithMetadataRequest{
				Request:  &embedderpb.InsertRequest{Document: &embedderpb.Document{Id: id, Text: plan.Text, Timestamp: time.Now().UnixNano()}},
				Metadata: meta,
			}, copts...)
		}))
	case config.OpEmbedderUpdate:
		return checkLocation(call(func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			return embedderpb.NewEmbedderClient(conn).Update(ctx, &embedderpb.UpdateRequest{
				Document: &embedderpb.Document{Id: id, Text: plan.Text, Timestamp: time.Now().UnixNano()},
				Config:   &payload.Update_Config{SkipStrictExistCheck: true},
			}, copts...)
		}))
	case config.OpEmbedderUpdateWithMetadata:
		return checkLocation(call(func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			return embedderpb.NewEmbedderClient(conn).UpdateWithMetadata(ctx, &embedderpb.UpdateWithMetadataRequest{
				Request: &embedderpb.UpdateRequest{
					Document: &embedderpb.Document{Id: id, Text: plan.Text, Timestamp: time.Now().UnixNano()},
					Config:   &payload.Update_Config{SkipStrictExistCheck: true},
				},
				Metadata: meta,
			}, copts...)
		}))
	case config.OpEmbedderUpsert:
		return checkLocation(call(func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			return embedderpb.NewEmbedderClient(conn).Upsert(ctx, &embedderpb.UpsertRequest{
				Document: &embedderpb.Document{Id: id, Text: plan.Text, Timestamp: time.Now().UnixNano()},
			}, copts...)
		}))
	case config.OpEmbedderUpsertWithMetadata:
		return checkLocation(call(func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			return embedderpb.NewEmbedderClient(conn).UpsertWithMetadata(ctx, &embedderpb.UpsertWithMetadataRequest{
				Request:  &embedderpb.UpsertRequest{Document: &embedderpb.Document{Id: id, Text: plan.Text, Timestamp: time.Now().UnixNano()}},
				Metadata: meta,
			}, copts...)
		}))
	case config.OpEmbedderRemove:
		return checkLocation(call(func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			return embedderpb.NewEmbedderClient(conn).Remove(ctx, &embedderpb.RemoveRequest{
				Id:     id,
				Config: &payload.Remove_Config{SkipStrictExistCheck: true},
			}, copts...)
		}))
	case config.OpEmbedderRemoveWithMetadata:
		return checkLocation(call(func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			return embedderpb.NewEmbedderClient(conn).RemoveWithMetadata(ctx, &embedderpb.RemoveRequest{
				Id:     id,
				Config: &payload.Remove_Config{SkipStrictExistCheck: true},
			}, copts...)
		}))
	default:
		return errors.New("unsupported embedder mutation type")
	}
}

func checkLocation(res any, err error) error {
	if err != nil {
		return err
	}
	loc, ok := res.(*payload.Object_Location)
	if !ok || loc == nil || (loc.GetUuid() == "" && loc.GetName() == "") {
		return errors.New("embedder mutation response location is empty")
	}
	return nil
}

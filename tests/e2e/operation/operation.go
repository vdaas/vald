//go:build e2e

//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package operation

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type client struct {
	host string
	port int
}

type Dataset struct {
	Train     [][]float32
	Test      [][]float32
	Neighbors [][]int
}

type Client interface {
	Search(t *testing.T, ctx context.Context, ds Dataset) error
	SearchByID(t *testing.T, ctx context.Context, ds Dataset) error
	SearchWithParameters(
		t *testing.T,
		ctx context.Context,
		ds Dataset,
		num uint32,
		radius float32,
		epsilon float32,
		timeout int64,
		statusValidator StatusValidator,
		errorValidator ErrorValidator,
	) error
	SearchByIDWithParameters(
		t *testing.T,
		ctx context.Context,
		ds Dataset,
		num uint32,
		radius float32,
		epsilon float32,
		timeout int64,
		statusValidator StatusValidator,
		errorValidator ErrorValidator,
	) error
	LinearSearch(t *testing.T, ctx context.Context, ds Dataset) error
	LinearSearchByID(t *testing.T, ctx context.Context, ds Dataset) error
	LinearSearchWithParameters(
		t *testing.T,
		ctx context.Context,
		ds Dataset,
		num uint32,
		timeout int64,
		statusValidator StatusValidator,
		errorValidator ErrorValidator,
	) error
	LinearSearchByIDWithParameters(
		t *testing.T,
		ctx context.Context,
		ds Dataset,
		num uint32,
		timeout int64,
		statusValidator StatusValidator,
		errorValidator ErrorValidator,
	) error
	Insert(t *testing.T, ctx context.Context, ds Dataset) error
	Update(t *testing.T, ctx context.Context, ds Dataset) error
	Upsert(t *testing.T, ctx context.Context, ds Dataset) error
	Remove(t *testing.T, ctx context.Context, ds Dataset) error
	InsertWithParameters(
		t *testing.T,
		ctx context.Context,
		ds Dataset,
		skipStrictExistCheck bool,
		statusValidator StatusValidator,
		errorValidator ErrorValidator,
	) error
	UpdateWithParameters(
		t *testing.T,
		ctx context.Context,
		ds Dataset,
		skipStrictExistCheck bool,
		offset int,
		statusValidator StatusValidator,
		errorValidator ErrorValidator,
	) error
	UpsertWithParameters(
		t *testing.T,
		ctx context.Context,
		ds Dataset,
		skipStrictExistCheck bool,
		offset int,
		statusValidator StatusValidator,
		errorValidator ErrorValidator,
	) error
	RemoveWithParameters(
		t *testing.T,
		ctx context.Context,
		ds Dataset,
		skipStrictExistCheck bool,
		statusValidator StatusValidator,
		errorValidator ErrorValidator,
	) error
	MultiSearch(t *testing.T, ctx context.Context, ds Dataset) error
	MultiSearchByID(t *testing.T, ctx context.Context, ds Dataset) error
	MultiLinearSearch(t *testing.T, ctx context.Context, ds Dataset) error
	MultiLinearSearchByID(t *testing.T, ctx context.Context, ds Dataset) error
	MultiInsert(t *testing.T, ctx context.Context, ds Dataset) error
	MultiUpdate(t *testing.T, ctx context.Context, ds Dataset) error
	MultiUpsert(t *testing.T, ctx context.Context, ds Dataset) error
	MultiRemove(t *testing.T, ctx context.Context, ds Dataset) error
	GetObject(t *testing.T, ctx context.Context, ds Dataset) error
	Exists(t *testing.T, ctx context.Context, id string) error
	CreateIndex(t *testing.T, ctx context.Context) error
	SaveIndex(t *testing.T, ctx context.Context) error
	IndexInfo(t *testing.T, ctx context.Context) (*payload.Info_Index_Count, error)
}

const (
	defaultSearchTimeout = 4 * int64(time.Second)
)

func New(host string, port int) (Client, error) {
	return &client{
		host: host,
		port: port,
	}, nil
}

func (c *client) CreateIndex(t *testing.T, ctx context.Context) error {
	client, err := c.getAgentClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
		PoolSize: 10000,
	})

	return err
}

func (c *client) SaveIndex(t *testing.T, ctx context.Context) error {
	client, err := c.getAgentClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.SaveIndex(ctx, &payload.Empty{})

	return err
}

func (c *client) IndexInfo(t *testing.T, ctx context.Context) (*payload.Info_Index_Count, error) {
	client, err := c.getAgentClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.IndexInfo(ctx, &payload.Empty{})
}

func (c *client) getGRPCConn(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		c.host+":"+strconv.Itoa(c.port),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                time.Second,
				Timeout:             5 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
}

func (c *client) getClient(ctx context.Context) (vald.Client, error) {
	conn, err := c.getGRPCConn(ctx)
	if err != nil {
		return nil, err
	}

	return vald.NewValdClient(conn), nil
}

func (c *client) getAgentClient(ctx context.Context) (core.AgentClient, error) {
	conn, err := c.getGRPCConn(ctx)
	if err != nil {
		return nil, err
	}

	return core.NewAgentClient(conn), nil
}

func (c *client) recall(results []string, neighbors []int) (recall float64) {
	ns := map[string]struct{}{}
	for _, n := range neighbors {
		ns[strconv.Itoa(n)] = struct{}{}
	}

	for _, r := range results {
		if _, ok := ns[r]; ok {
			recall++
		}
	}

	return recall / float64(len(neighbors))
}

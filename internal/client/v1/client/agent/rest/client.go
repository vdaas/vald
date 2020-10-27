//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package rest provides agent ngt REST client functions
package rest

import (
	"context"
	"net/http"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/http/json"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type agentClient struct {
	addr string
}

func New(ctx context.Context, opts ...Option) Client {
	c := new(agentClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	return c
}

func (c *agentClient) Exists(
	ctx context.Context,
	req *client.ObjectID,
	opts ...grpc.CallOption,
) (res *client.ObjectID, err error) {
	res = new(client.ObjectID)
	err = json.Request(ctx, http.MethodGet, c.addr+"/exists/"+req.GetId(), req, res)
	return
}

func (c *agentClient) Search(
	ctx context.Context,
	req *client.SearchRequest,
	opts ...grpc.CallOption,
) (res *client.SearchResponse, err error) {
	res = new(client.SearchResponse)
	err = json.Request(ctx, http.MethodPost, c.addr+"/search", req, res)
	return
}

func (c *agentClient) SearchByID(
	ctx context.Context,
	req *client.SearchIDRequest,
	opts ...grpc.CallOption,
) (res *client.SearchResponse, err error) {
	res = new(client.SearchResponse)
	err = json.Request(ctx, http.MethodPost, c.addr+"/search/id", req, res)
	return
}

func (c *agentClient) StreamSearch(
	ctx context.Context,
	opts ...grpc.CallOption,
) (vald.Search_StreamSearchClient, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) StreamSearchByID(
	ctx context.Context,
	opts ...grpc.CallOption,
) (vald.Search_StreamSearchByIDClient, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) MultiSearch(
	ctx context.Context,
	req *client.SearchMultiRequest,
	opts ...grpc.CallOption,
) (res *client.SearchResponses, err error) {
	res = new(client.SearchResponses)
	err = json.Request(ctx, http.MethodPost, c.addr+"/search", req, res)
	return
}

func (c *agentClient) MultiSearchByID(
	ctx context.Context,
	req *client.SearchIDMultiRequest,
	opts ...grpc.CallOption,
) (res *client.SearchResponses, err error) {
	res = new(client.SearchResponses)
	err = json.Request(ctx, http.MethodPost, c.addr+"/search/id", req, res)
	return
}

func (c *agentClient) Insert(
	ctx context.Context,
	req *client.InsertRequest,
	opts ...grpc.CallOption,
) (*client.ObjectLocation, error) {
	return nil, json.Request(ctx, http.MethodPost, c.addr+"/insert", req, nil)
}

func (c *agentClient) StreamInsert(
	ctx context.Context,
	opts ...grpc.CallOption,
) (vald.Insert_StreamInsertClient, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) MultiInsert(
	ctx context.Context,
	objectVectors *client.InsertMultiRequest,
	opts ...grpc.CallOption,
) (*client.ObjectLocations, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) Update(
	ctx context.Context,
	req *client.UpdateRequest,
	opts ...grpc.CallOption,
) (*client.ObjectLocation, error) {
	return nil, json.Request(ctx, http.MethodPost, c.addr+"/update", req, nil)
}

func (c *agentClient) StreamUpdate(
	ctx context.Context,
	opts ...grpc.CallOption,
) (vald.Update_StreamUpdateClient, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) MultiUpdate(
	ctx context.Context,
	objectVectors *client.UpdateMultiRequest,
	opts ...grpc.CallOption,
) (*client.ObjectLocations, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) Upsert(
	ctx context.Context,
	req *client.UpsertRequest,
	opts ...grpc.CallOption,
) (*client.ObjectLocation, error) {
	return nil, json.Request(ctx, http.MethodPost, c.addr+"/upsert", req, nil)
}

func (c *agentClient) StreamUpsert(
	ctx context.Context,
	opts ...grpc.CallOption,
) (vald.Upsert_StreamUpsertClient, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) MultiUpsert(
	ctx context.Context,
	objectVectors *client.UpsertMultiRequest,
	opts ...grpc.CallOption,
) (*client.ObjectLocations, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) Remove(
	ctx context.Context,
	req *client.RemoveRequest,
	opts ...grpc.CallOption,
) (*client.ObjectLocation, error) {
	return nil, json.Request(ctx, http.MethodDelete, c.addr+"/remove/"+req.GetId().GetId(), req, nil)
}

func (c *agentClient) StreamRemove(
	ctx context.Context,
	opts ...grpc.CallOption,
) (vald.Remove_StreamRemoveClient, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) MultiRemove(
	ctx context.Context,
	req *client.RemoveMultiRequest,
	opts ...grpc.CallOption,
) (*client.ObjectLocations, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) GetObject(
	ctx context.Context,
	req *client.ObjectID,
	opts ...grpc.CallOption,
) (res *client.ObjectVector, err error) {
	res = new(client.ObjectVector)
	err = json.Request(ctx, http.MethodGet, c.addr+"/object/"+req.GetId(), req, res)
	return
}

func (c *agentClient) StreamGetObject(
	ctx context.Context,
	opts ...grpc.CallOption,
) (vald.Object_StreamGetObjectClient, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *agentClient) CreateIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	return nil, json.Request(ctx, http.MethodGet, c.addr+"/index/create", req, nil)
}

func (c *agentClient) SaveIndex(
	ctx context.Context,
	req *client.Empty,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	return nil, json.Request(ctx, http.MethodGet, c.addr+"/index/save", nil, nil)
}

func (c *agentClient) CreateAndSaveIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	return nil, json.Request(ctx, http.MethodGet, c.addr+"/index/createandsave", nil, nil)
}

func (c *agentClient) IndexInfo(
	ctx context.Context,
	req *client.Empty,
	opts ...grpc.CallOption,
) (res *client.InfoIndexCount, err error) {
	res = new(client.InfoIndexCount)
	err = json.Request(ctx, http.MethodGet, c.addr+"/index/info", nil, res)
	return
}

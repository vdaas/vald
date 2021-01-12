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

// Package rest provides rest client functions
package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/http/json"

	"github.com/yahoojapan/ngtd/model"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type ngtdClient struct {
	addr string
}

func New(ctx context.Context, opts ...Option) (Client, error) {
	c := new(ngtdClient)
	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}
	return c, nil
}

func (c *ngtdClient) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (oid *payload.Object_ID, err error) {
	id, err := c.GetObject(ctx, in, opts...)
	if err != nil {
		return nil, err
	}
	return &payload.Object_ID{
		Id: id.GetId(),
	}, nil
}

func (c *ngtdClient) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	vec := make([]float64, 0, len(in.GetVector()))
	for _, v := range in.GetVector() {
		vec = append(vec, float64(v))
	}
	var res model.SearchResponse
	err := json.Request(ctx, http.MethodPost, c.addr+"/search", model.SearchRequest{
		Vector:  vec,
		Size:    int(in.GetConfig().GetNum()),
		Epsilon: in.GetConfig().GetEpsilon(),
	}, &res)
	if err != nil {
		return nil, err
	}
	sr := &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(res.Result)),
	}
	for _, r := range res.Result {
		sr.Results = append(sr.Results, &payload.Object_Distance{
			Id:       r.ID,
			Distance: r.Distance,
		})
	}
	return sr, nil
}

func (c *ngtdClient) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	var res model.SearchResponse
	err := json.Request(ctx, http.MethodPost, c.addr+"/search", model.SearchRequest{
		ID:      in.GetId(),
		Size:    int(in.GetConfig().GetNum()),
		Epsilon: in.GetConfig().GetEpsilon(),
	}, &res)
	if err != nil {
		return nil, err
	}
	sr := &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(res.Result)),
	}
	for _, r := range res.Result {
		sr.Results = append(sr.Results, &payload.Object_Distance{
			Id:       r.ID,
			Distance: r.Distance,
		})
	}
	return sr, nil
}

func (c *ngtdClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchClient, err error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchByIDClient, err error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, 0, len(in.GetRequests())),
	}
	for _, req := range in.GetRequests() {
		r, err := c.Search(ctx, req)
		if err == nil {
			res.Responses = append(res.Responses, r)
		}
	}
	return res, nil
}

func (c *ngtdClient) MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, 0, len(in.GetRequests())),
	}
	for _, req := range in.GetRequests() {
		r, err := c.SearchByID(ctx, req)
		if err == nil {
			res.Responses = append(res.Responses, r)
		}
	}
	return res, nil
}

func (c *ngtdClient) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	vec := make([]float64, 0, len(in.GetVector().GetVector()))
	for _, v := range in.GetVector().GetVector() {
		vec = append(vec, float64(v))
	}
	var res model.InsertResponse
	err := json.Request(ctx, http.MethodPost, c.addr+"/insert", model.InsertRequest{
		ID:     in.GetVector().GetId(),
		Vector: vec,
	}, &res)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *ngtdClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Insert_StreamInsertClient, err error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	req := &model.MultiInsertRequest{
		InsertRequests: make([]model.InsertRequest, 0, len(in.GetRequests())),
	}
	for _, i := range in.GetRequests() {
		vec := make([]float64, 0, len(i.GetVector().GetVector()))
		for _, v := range i.GetVector().GetVector() {
			vec = append(vec, float64(v))
		}
		req.InsertRequests = append(req.InsertRequests, model.InsertRequest{
			ID:     i.GetVector().GetId(),
			Vector: vec,
		})
	}
	var r model.MultiInsertResponse
	return nil, json.Request(ctx, http.MethodPost, c.addr+"/multiinsert", req, &r)
}

func (c *ngtdClient) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	_, err = c.Remove(ctx, &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: in.GetVector().GetId(),
		},
	}, opts...)
	if err != nil {
		return nil, err
	}
	_, err = c.Insert(ctx, &payload.Insert_Request{
		Vector: in.GetVector(),
	}, opts...)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *ngtdClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (res vald.Update_StreamUpdateClient, err error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	for _, req := range in.GetRequests() {
		_, err := c.Update(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (c *ngtdClient) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	id, err := c.Exists(ctx, &payload.Object_ID{
		Id: in.GetVector().GetId(),
	}, opts...)
	if err == nil || len(id.GetId()) != 0 {
		return c.Update(ctx, &payload.Update_Request{
			Vector: in.GetVector(),
		}, opts...)
	}
	return c.Insert(ctx, &payload.Insert_Request{
		Vector: in.GetVector(),
	}, opts...)
}

func (c *ngtdClient) StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Upsert_StreamUpsertClient, err error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiUpsert(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	for _, req := range in.GetRequests() {
		_, err := c.Upsert(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (c *ngtdClient) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	var res model.RemoveResponse
	err := json.Request(ctx, http.MethodGet, c.addr+"/remove/"+in.GetId().GetId(), nil, &res)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *ngtdClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (res vald.Remove_StreamRemoveClient, err error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	req := &model.MultiRemoveRequest{
		IDs: make([]string, 0, len(in.GetRequests())),
	}
	for _, i := range in.GetRequests() {
		req.IDs = append(req.IDs, i.GetId().GetId())
	}
	var r model.MultiInsertResponse
	return nil, json.Request(ctx, http.MethodPost, c.addr+"/multiremove", req, &r)
}

func (c *ngtdClient) GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	var res model.GetObjectsResponse
	err := json.Request(ctx, http.MethodPost, c.addr+"/getobjects", model.GetObjectsRequest{
		IDs: []string{in.GetId()},
	}, &res)
	if err != nil {
		return nil, err
	}
	return &payload.Object_Vector{
		Id:     res.Result[0].ID,
		Vector: res.Result[0].Vector,
	}, nil
}

func (c *ngtdClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (res vald.Object_StreamGetObjectClient, err error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) CreateIndex(
	ctx context.Context,
	in *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	err := json.Request(ctx, http.MethodGet, fmt.Sprintf("%s/index/create/%d", c.addr, in.GetPoolSize()), nil, nil)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *ngtdClient) SaveIndex(
	ctx context.Context,
	req *client.Empty,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	err := json.Request(ctx, http.MethodGet, c.addr+"/index/save", nil, nil)
	return nil, err
}

func (c *ngtdClient) CreateAndSaveIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	_, err := c.CreateIndex(ctx, req)
	if err != nil {
		return nil, err
	}
	_, err = c.SaveIndex(ctx, nil)
	return nil, err
}

func (c *ngtdClient) IndexInfo(
	ctx context.Context,
	req *client.Empty,
	opts ...grpc.CallOption,
) (res *client.InfoIndexCount, err error) {
	return nil, errors.ErrUnsupportedClientMethod
}

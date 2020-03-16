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

// Package rest provides rest client functions
package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/errors"
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

func (c *ngtdClient) Exists(
	ctx context.Context,
	req *client.ObjectID,
) (*client.ObjectID, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) Search(
	ctx context.Context,
	req *client.SearchRequest,
) (*client.SearchResponse, error) {
	res := new(model.SearchResponse)
	err := json.Request(ctx, http.MethodPost, c.addr+"/search", searchRequestToNgtdSearchRequest(req), &res)
	if err != nil {
		return nil, err
	}
	return ngtdSearchResponseToSearchResponse(res), nil
}

func (c *ngtdClient) SearchByID(
	ctx context.Context,
	req *client.SearchIDRequest,
) (*client.SearchResponse, error) {
	res := new(model.SearchResponse)
	err := json.Request(ctx, http.MethodPost, c.addr+"/searchbyid", searchIDRequestToNgtdSearchRequest(req), res)
	if err != nil {
		return nil, err
	}
	return ngtdSearchResponseToSearchResponse(res), nil
}

func (c *ngtdClient) StreamSearch(
	ctx context.Context,
	dataProvider func() *client.SearchRequest,
	f func(*client.SearchResponse, error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) StreamSearchByID(
	ctx context.Context,
	dataProvider func() *client.SearchIDRequest,
	f func(*client.SearchResponse, error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) Insert(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	err := json.Request(ctx, http.MethodPost, c.addr+"/insert", objectVectorToNgtdInsertRequest(req), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ngtdClient) StreamInsert(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiInsert(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	err := json.Request(ctx, http.MethodPost, c.addr+"/multiinsert", objectVectorsToNgtdMultiInsertRequest(req), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ngtdClient) Update(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) StreamUpdate(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiUpdate(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) Remove(
	ctx context.Context,
	req *client.ObjectID,
) error {
	return json.Request(ctx, http.MethodGet, c.addr+"/remove/"+req.GetId(), nil, nil)
}

func (c *ngtdClient) StreamRemove(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiRemove(
	ctx context.Context,
	req *client.ObjectIDs,
) (err error) {
	res := new(model.MultiRemoveResponse)
	err = json.Request(ctx, http.MethodGet, c.addr+"/multiremove/", objectIDsToNgtdMultiRemoveRequest(req), res)
	if err != nil {
		return err
	}

	for _, resErr := range res.Errors {
		if err == nil {
			err = resErr
		} else {
			if resErr != nil {
				err = errors.Wrap(err, resErr.Error())
			}
		}
	}

	return
}

func (c *ngtdClient) GetObject(
	ctx context.Context,
	req *client.ObjectID,
) (*client.ObjectVector, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) StreamGetObject(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(*client.ObjectVector, error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) CreateIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
) error {
	res := new(model.DefaultResponse)
	err := json.Request(ctx, http.MethodGet, c.addr+"/index/create/"+strconv.Itoa(int(req.GetPoolSize())), nil, res)
	if err != nil {
		return err
	}

	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (c *ngtdClient) SaveIndex(ctx context.Context) error {
	return json.Request(ctx, http.MethodGet, c.addr+"/index/save", nil, nil)
}

func (c *ngtdClient) CreateAndSaveIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) IndexInfo(ctx context.Context) (*client.InfoIndex, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func searchRequestToNgtdSearchRequest(in *client.SearchRequest) *model.SearchRequest {
	size, epsilon := getSizeAndEpsilon(in.GetConfig())
	return &model.SearchRequest{
		Vector:  tofloat64(in.GetVector()),
		Size:    size,
		Epsilon: epsilon,
	}
}

func searchIDRequestToNgtdSearchRequest(in *client.SearchIDRequest) *model.SearchRequest {
	size, epsilon := getSizeAndEpsilon(in.GetConfig())
	return &model.SearchRequest{
		ID:      in.GetId(),
		Size:    size,
		Epsilon: epsilon,
	}
}

func objectVectorToNgtdInsertRequest(in *client.ObjectVector) *model.InsertRequest {
	return &model.InsertRequest{
		ID:     in.GetId(),
		Vector: tofloat64(in.GetVector()),
	}
}

func objectVectorsToNgtdMultiInsertRequest(in *client.ObjectVectors) *model.MultiInsertRequest {
	reqs := make([]model.InsertRequest, len(in.GetVectors()))

	for _, v := range in.GetVectors() {
		reqs = append(reqs, model.InsertRequest{
			Vector: tofloat64(v.GetVector()),
			ID:     v.GetId(),
		})
	}

	return &model.MultiInsertRequest{
		InsertRequests: reqs,
	}
}

func objectIDsToNgtdMultiRemoveRequest(in *client.ObjectIDs) *model.MultiRemoveRequest {
	return &model.MultiRemoveRequest{
		IDs: in.GetIds(),
	}
}

func ngtdSearchResponseToSearchResponse(in *model.SearchResponse) *client.SearchResponse {
	results := make([]*client.ObjectDistance, len(in.Result))

	for _, r := range in.Result {
		results = append(results, &client.ObjectDistance{
			Id:       r.ID,
			Distance: r.Distance,
		})
	}

	return &client.SearchResponse{
		Results: results,
	}
}

func getSizeAndEpsilon(cfg *client.SearchConfig) (size int, epsilon float32) {
	if cfg != nil {
		size = int(cfg.GetNum())
		epsilon = float32(cfg.GetEpsilon())
	}
	return
}

func tofloat64(in []float32) (out []float64) {
	out = make([]float64, len(in))
	for i := range in {
		out[i] = float64(in[i])
	}
	return
}

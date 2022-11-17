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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

func (c *client) MultiSearch(t *testing.T, ctx context.Context, ds Dataset) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	cfg := &payload.Search_Config{
		Num:     3,
		Radius:  -1.0,
		Epsilon: 0.1,
	}

	reqs := make([]*payload.Search_Request, 0, len(ds.Test))
	for _, v := range ds.Test {
		reqs = append(reqs, &payload.Search_Request{
			Vector: v,
			Config: cfg,
		})
	}

	req := &payload.Search_MultiRequest{
		Requests: reqs,
	}

	res, err := client.MultiSearch(ctx, req)
	if err != nil {
		return err
	}

	if len(res.GetResponses()) != len(ds.Test) {
		t.Error("number of responses does not match with sent requests")
	}

	return nil
}

func (c *client) MultiSearchByID(t *testing.T, ctx context.Context, ds Dataset) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	cfg := &payload.Search_Config{
		Num:     3,
		Radius:  -1.0,
		Epsilon: 0.1,
	}

	reqs := make([]*payload.Search_IDRequest, 0, len(ds.Test))
	for i := range ds.Test {
		reqs = append(reqs, &payload.Search_IDRequest{
			Id:     strconv.Itoa(i),
			Config: cfg,
		})
	}

	req := &payload.Search_MultiIDRequest{
		Requests: reqs,
	}

	res, err := client.MultiSearchByID(ctx, req)
	if err != nil {
		return err
	}

	if len(res.GetResponses()) != len(ds.Test) {
		t.Error("number of responses does not match with sent requests")
	}

	return nil
}

func (c *client) MultiLinearSearch(t *testing.T, ctx context.Context, ds Dataset) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	cfg := &payload.Search_Config{
		Num: 3,
	}

	reqs := make([]*payload.Search_Request, 0, len(ds.Test))
	for _, v := range ds.Test {
		reqs = append(reqs, &payload.Search_Request{
			Vector: v,
			Config: cfg,
		})
	}

	req := &payload.Search_MultiRequest{
		Requests: reqs,
	}

	res, err := client.MultiLinearSearch(ctx, req)
	if err != nil {
		return err
	}

	if len(res.GetResponses()) != len(ds.Test) {
		t.Error("number of responses does not match with sent requests")
	}

	return nil
}

func (c *client) MultiLinearSearchByID(t *testing.T, ctx context.Context, ds Dataset) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	cfg := &payload.Search_Config{
		Num: 3,
	}

	reqs := make([]*payload.Search_IDRequest, 0, len(ds.Test))
	for i := range ds.Test {
		reqs = append(reqs, &payload.Search_IDRequest{
			Id:     strconv.Itoa(i),
			Config: cfg,
		})
	}

	req := &payload.Search_MultiIDRequest{
		Requests: reqs,
	}

	res, err := client.MultiLinearSearchByID(ctx, req)
	if err != nil {
		return err
	}

	if len(res.GetResponses()) != len(ds.Test) {
		t.Error("number of responses does not match with sent requests")
	}

	return nil
}

func (c *client) MultiInsert(t *testing.T, ctx context.Context, ds Dataset) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	cfg := &payload.Insert_Config{
		SkipStrictExistCheck: true,
	}

	reqs := make([]*payload.Insert_Request, 0, len(ds.Train))
	for i, v := range ds.Train {
		reqs = append(reqs, &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:     strconv.Itoa(i),
				Vector: v,
			},
			Config: cfg,
		})
	}

	req := &payload.Insert_MultiRequest{
		Requests: reqs,
	}

	res, err := client.MultiInsert(ctx, req)
	if err != nil {
		return err
	}

	if len(res.GetLocations()) != len(ds.Train) {
		t.Error("number of responses does not match with sent requests")
	}

	return nil
}

func (c *client) MultiUpdate(t *testing.T, ctx context.Context, ds Dataset) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	cfg := &payload.Update_Config{
		SkipStrictExistCheck: true,
	}

	reqs := make([]*payload.Update_Request, 0, len(ds.Train))
	for i, v := range ds.Train {
		reqs = append(reqs, &payload.Update_Request{
			Vector: &payload.Object_Vector{
				Id:     strconv.Itoa(i),
				Vector: append(v[1:], v[0]),
			},
			Config: cfg,
		})
	}

	req := &payload.Update_MultiRequest{
		Requests: reqs,
	}

	res, err := client.MultiUpdate(ctx, req)
	if err != nil {
		return err
	}

	// Note: The MultiUpdate API internally checks the identity of the vectors to be updated by the LB Gateway,
	// so it is important to remember that the number of responses is not always the same as the number of requested data.
	// The response includes an ID, so the client can check the order of the data based on the requested ID.
	if len(res.GetLocations()) == 0 {
		t.Error("empty response detected")
	}

	return nil
}

func (c *client) MultiUpsert(t *testing.T, ctx context.Context, ds Dataset) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	cfg := &payload.Upsert_Config{
		SkipStrictExistCheck: true,
	}

	reqs := make([]*payload.Upsert_Request, 0, len(ds.Train))
	for i, v := range ds.Train {
		reqs = append(reqs, &payload.Upsert_Request{
			Vector: &payload.Object_Vector{
				Id:     strconv.Itoa(i),
				Vector: v,
			},
			Config: cfg,
		})
	}

	req := &payload.Upsert_MultiRequest{
		Requests: reqs,
	}

	res, err := client.MultiUpsert(ctx, req)
	if err != nil {
		return err
	}

	// Note: The MultiUpsert API internally checks the identity of the vectors to be updated by the LB Gateway,
	// so it is important to remember that the number of responses is not always the same as the number of requested data.
	// The response includes an ID, so the client can check the order of the data based on the requested ID.
	if len(res.GetLocations()) == 0 {
		t.Error("empty response detected")
	}

	return nil
}

func (c *client) MultiRemove(t *testing.T, ctx context.Context, ds Dataset) error {
	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	cfg := &payload.Remove_Config{
		SkipStrictExistCheck: true,
	}

	reqs := make([]*payload.Remove_Request, 0, len(ds.Train))
	for i := range ds.Train {
		reqs = append(reqs, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: strconv.Itoa(i),
			},
			Config: cfg,
		})
	}

	req := &payload.Remove_MultiRequest{
		Requests: reqs,
	}

	res, err := client.MultiRemove(ctx, req)
	if err != nil {
		return err
	}

	if len(res.GetLocations()) != len(ds.Train) {
		t.Error("number of responses does not match with sent requests")
	}

	return nil
}

//go:build e2e
// +build e2e

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
package operation

import (
	"context"
	"io"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func (c *client) Search(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.SearchWithParameters(t, ctx, ds, 100, -1.0, 0.1, 3000000000)
}

func (c *client) SearchWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	num uint32,
	radius float32,
	epsilon float32,
	timeout int64,
) error {
	t.Log("search operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamSearch(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf(
					"error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s",
					err,
					serr,
					st.Code().String(),
					errdetails.Serialize(st.Details()),
					st.Message(),
					st.Err().Error(),
					errdetails.Serialize(st.Proto()),
				)
				continue
			}

			resp := res.GetResponse()
			if resp == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned response is nil")
				continue
			}

			topKIDs := make([]string, 0, len(resp.GetResults()))
			for _, d := range resp.GetResults() {
				topKIDs = append(topKIDs, d.GetId())
			}

			if len(topKIDs) == 0 {
				t.Errorf("empty result is returned for test ID %s: %#v", resp.GetRequestId(), topKIDs)
				continue
			}

			idx, err := strconv.Atoi(resp.GetRequestId())
			if err != nil {
				t.Errorf("an error occurred while converting RequestId into int: %s", err)
				continue
			}

			t.Logf("results: %d, recall: %f", len(topKIDs), c.recall(topKIDs, ds.Neighbors[idx][:len(topKIDs)]))
		}
	}()

	for i := 0; i < len(ds.Test); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Search_Request{
			Vector: ds.Test[i],
			Config: &payload.Search_Config{
				RequestId: id,
				Num:       num,
				Radius:    radius,
				Epsilon:   epsilon,
				Timeout:   timeout,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("search operation finished")

	return nil
}

func (c *client) SearchByID(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.SearchByIDWithParameters(t, ctx, ds, 100, -1.0, 0.1, 3000000000)
}

func (c *client) SearchByIDWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	num uint32,
	radius float32,
	epsilon float32,
	timeout int64,
) error {
	t.Log("searchByID operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamSearchByID(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf(
					"error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s",
					err,
					serr,
					st.Code().String(),
					errdetails.Serialize(st.Details()),
					st.Message(),
					st.Err().Error(),
					errdetails.Serialize(st.Proto()),
				)
				continue
			}

			resp := res.GetResponse()
			if resp == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned response is nil")
			}

			topKIDs := make([]string, 0, len(resp.GetResults()))
			for _, d := range resp.GetResults() {
				topKIDs = append(topKIDs, d.GetId())
			}

			if len(topKIDs) == 0 {
				t.Errorf("empty result is returned: %#v", topKIDs)
			}
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Search_IDRequest{
			Id: id,
			Config: &payload.Search_Config{
				RequestId: id,
				Num:       num,
				Radius:    radius,
				Epsilon:   epsilon,
				Timeout:   timeout,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("searchByID operation finished")

	return nil
}

func (c *client) Insert(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.InsertWithParameters(t, ctx, ds, false)
}

func (c *client) InsertWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	skipStrictExistCheck bool,
) error {
	t.Log("insert operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamInsert(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf(
					"error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s",
					err,
					serr,
					st.Code().String(),
					errdetails.Serialize(st.Details()),
					st.Message(),
					st.Err().Error(),
					errdetails.Serialize(st.Proto()),
				)
				continue
			}

			loc := res.GetLocation()
			if loc == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned loc is nil")
				continue
			}

			t.Logf("returned loc: %s", loc)
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:     id,
				Vector: ds.Train[i],
			},
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: skipStrictExistCheck,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("insert operation finished")

	return nil
}

func (c *client) Update(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.UpdateWithParameters(t, ctx, ds, false, 0)
}

func (c *client) UpdateWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	skipStrictExistCheck bool,
	offset int,
) error {
	t.Log("update operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamUpdate(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf(
					"error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s",
					err,
					serr,
					st.Code().String(),
					errdetails.Serialize(st.Details()),
					st.Message(),
					st.Err().Error(),
					errdetails.Serialize(st.Proto()),
				)
				continue
			}

			loc := res.GetLocation()
			if loc == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned loc is nil")
			}

			t.Logf("returned: %s", loc)
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		v := ds.Train[i]
		err := sc.Send(&payload.Update_Request{
			Vector: &payload.Object_Vector{
				Id:     id,
				Vector: append(v[offset+1:], v[offset]),
			},
			Config: &payload.Update_Config{
				SkipStrictExistCheck: skipStrictExistCheck,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("update operation finished")

	return nil
}

func (c *client) Upsert(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.UpsertWithParameters(t, ctx, ds, false, 1)
}

func (c *client) UpsertWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	skipStrictExistCheck bool,
	offset int,
) error {
	t.Log("upsert operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamUpsert(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf(
					"error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s",
					err,
					serr,
					st.Code().String(),
					errdetails.Serialize(st.Details()),
					st.Message(),
					st.Err().Error(),
					errdetails.Serialize(st.Proto()),
				)
				continue
			}

			loc := res.GetLocation()
			if loc == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned loc is nil")
			}

			t.Logf("returned: %s", loc)
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		v := ds.Train[i]
		err := sc.Send(&payload.Upsert_Request{
			Vector: &payload.Object_Vector{
				Id:     id,
				Vector: append(v[offset+1:], v[offset]),
			},
			Config: &payload.Upsert_Config{
				SkipStrictExistCheck: skipStrictExistCheck,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("upsert operation finished")

	return nil
}

func (c *client) Remove(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.RemoveWithParameters(t, ctx, ds, false)
}

func (c *client) RemoveWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	skipStrictExistCheck bool,
) error {
	t.Log("remove operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamRemove(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf(
					"error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s",
					err,
					serr,
					st.Code().String(),
					errdetails.Serialize(st.Details()),
					st.Message(),
					st.Err().Error(),
					errdetails.Serialize(st.Proto()),
				)
				continue
			}

			loc := res.GetLocation()
			if loc == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned response is nil")
				continue
			}

			t.Logf("returned: %s", loc)
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: id,
			},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: skipStrictExistCheck,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("remove operation finished")

	return nil
}

func (c *client) Exists(t *testing.T, ctx context.Context, id string) error {
	t.Log("exists operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	req := &payload.Object_ID{
		Id: id,
	}

	_, err = client.Exists(ctx, req)
	if err != nil {
		return err
	}

	t.Log("exists operation finished")

	return nil
}

func (c *client) GetObject(t *testing.T, ctx context.Context, ds Dataset) error {
	t.Log("getObject operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamGetObject(ctx)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				st, serr := status.FromError(err)
				t.Errorf(
					"error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s",
					err,
					serr,
					st.Code().String(),
					errdetails.Serialize(st.Details()),
					st.Message(),
					st.Err().Error(),
					errdetails.Serialize(st.Proto()),
				)
				continue
			}

			resp := res.GetVector()
			if resp == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					continue
				}

				t.Error("returned response is nil")
				continue
			}

			idx, err := strconv.Atoi(resp.GetId())
			if err != nil {
				t.Errorf("an error occurred while converting Id into int: %s", err)
				continue
			}

			if !reflect.DeepEqual(res.GetVector().GetVector(), ds.Train[idx]) {
				t.Errorf(
					"got: %#v, expected: %#v",
					res.GetVector().GetVector(),
					ds.Train[idx],
				)
			}
		}
	}()

	for i := 0; i < len(ds.Train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: id,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("getObject operation finished")

	return nil
}

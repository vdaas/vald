//go:build e2e

// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package operation

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

type (
	StatusValidator = func(t *testing.T, status int32, msg string) (err error)
	ErrorValidator  = func(t *testing.T, err error) error
)

func DefaultStatusValidator(t *testing.T, status int32, msg string) (err error) {
	t.Helper()

	// TODO: convert int32->codes.Code
	return errors.Errorf("code: %d, message: %s", status, msg)
}

func ParseAndLogError(t *testing.T, err error) error {
	t.Helper()

	st, _, parsed := status.ParseError(err, codes.Unknown, "nothing")
	if parsed == nil {
		return err
	}

	t.Errorf(
		"error: %v\tserror: %v\tcode: %s\tdetails: %s\tmessage: %s\tstatus-error: %s\tproto: %s",
		err,
		parsed,
		st.Code().String(),
		errdetails.Serialize(st.Details()),
		st.Message(),
		st.Err().Error(),
		errdetails.Serialize(st.Proto()),
	)

	return parsed
}

func (c *client) Search(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.SearchWithParameters(
		t,
		ctx,
		ds,
		100,
		-1.0,
		0.1,
		3000000000,
		DefaultStatusValidator,
		ParseAndLogError,
	)
}

func (c *client) SearchWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	num uint32,
	radius float32,
	epsilon float32,
	timeout int64,
	svalidator StatusValidator,
	evalidator ErrorValidator,
) (rerr error) {
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
	var mu sync.Mutex
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				if err := evalidator(t, err); err != nil {
					mu.Lock()
					rerr = errors.Join(
						rerr,
						errors.Errorf(
							"stream finished by an error: %s",
							err.Error(),
						),
					)
					mu.Unlock()
				}
				return
			}

			resp := res.GetResponse()
			if resp == nil {
				status := res.GetStatus()
				if status != nil {
					if e := svalidator(t, status.GetCode(), status.GetMessage()); e != nil {
						t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s",
							status.GetCode(),
							status.GetMessage(),
							errdetails.Serialize(status.GetDetails()))
						mu.Lock()
						rerr = errors.Join(rerr, e)
						mu.Unlock()
					}
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
			left, right, ok := strings.Cut(resp.GetRequestId(), "-")
			if !ok {
				sid := strings.SplitN(resp.GetRequestId(), "-", 2)
				left, right = sid[0], sid[1]
			}

			idx, err := strconv.Atoi(left)
			if err != nil {
				t.Errorf("an error occurred while converting RequestId into int: %s", err)
				continue
			}

			t.Logf("algo: %s, id: %d, results: %d, recall: %f", right, idx, len(topKIDs), c.recall(topKIDs, ds.Neighbors[idx][:len(topKIDs)]))
		}
	}()

	for i := 0; i < len(ds.Test); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Search_Request{
			Vector: ds.Test[i],
			Config: &payload.Search_Config{
				RequestId: id + "-Unknown",
				Num:       num,
				Radius:    radius,
				Epsilon:   epsilon,
				Timeout:   timeout,
			},
		})
		if err != nil {
			mu.Lock()
			defer mu.Unlock()
			return err
		}
		err = sc.Send(&payload.Search_Request{
			Vector: ds.Test[i],
			Config: &payload.Search_Config{
				RequestId:            id + "-ConcurrentQueue",
				Num:                  num,
				Radius:               radius,
				Epsilon:              epsilon,
				Timeout:              timeout,
				AggregationAlgorithm: payload.Search_ConcurrentQueue,
			},
		})
		if err != nil {
			mu.Lock()
			defer mu.Unlock()
			return err
		}
		err = sc.Send(&payload.Search_Request{
			Vector: ds.Test[i],
			Config: &payload.Search_Config{
				RequestId:            id + "-SortSlice",
				Num:                  num,
				Radius:               radius,
				Epsilon:              epsilon,
				Timeout:              timeout,
				AggregationAlgorithm: payload.Search_SortSlice,
			},
		})
		if err != nil {
			mu.Lock()
			defer mu.Unlock()
			return err
		}
		err = sc.Send(&payload.Search_Request{
			Vector: ds.Test[i],
			Config: &payload.Search_Config{
				RequestId:            id + "-SortPoolSlice",
				Num:                  num,
				Radius:               radius,
				Epsilon:              epsilon,
				Timeout:              timeout,
				AggregationAlgorithm: payload.Search_SortPoolSlice,
			},
		})
		if err != nil {
			mu.Lock()
			defer mu.Unlock()
			return err
		}
		err = sc.Send(&payload.Search_Request{
			Vector: ds.Test[i],
			Config: &payload.Search_Config{
				RequestId:            id + "-PairingHeap",
				Num:                  num,
				Radius:               radius,
				Epsilon:              epsilon,
				Timeout:              timeout,
				AggregationAlgorithm: payload.Search_PairingHeap,
			},
		})
		if err != nil {
			mu.Lock()
			defer mu.Unlock()
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("search operation finished")

	return rerr
}

func (c *client) SearchByID(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.SearchByIDWithParameters(t,
		ctx,
		ds,
		100,
		-1.0,
		0.1,
		3000000000,
		DefaultStatusValidator,
		ParseAndLogError,
	)
}

func (c *client) SearchByIDWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	num uint32,
	radius float32,
	epsilon float32,
	timeout int64,
	svalidator StatusValidator,
	evalidator ErrorValidator,
) (rerr error) {
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
	var mu sync.Mutex
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				if err := evalidator(t, err); err != nil {
					mu.Lock()
					rerr = errors.Join(
						rerr,
						errors.Errorf(
							"stream finished by an error: %s",
							err.Error(),
						),
					)
					mu.Unlock()
				}
				return
			}

			resp := res.GetResponse()
			if resp == nil {
				status := res.GetStatus()
				if status != nil {
					if e := svalidator(t, status.GetCode(), status.GetMessage()); e != nil {
						t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s",
							status.GetCode(),
							status.GetMessage(),
							errdetails.Serialize(status.GetDetails()))
						mu.Lock()
						rerr = errors.Join(rerr, e)
						mu.Unlock()
					}
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
			mu.Lock()
			defer mu.Unlock()
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("searchByID operation finished")

	mu.Lock()
	defer mu.Unlock()
	return rerr
}

func (c *client) LinearSearch(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.LinearSearchWithParameters(
		t,
		ctx,
		ds,
		100,
		3000000000,
		DefaultStatusValidator,
		ParseAndLogError,
	)
}

func (c *client) LinearSearchWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	num uint32,
	timeout int64,
	svalidator StatusValidator,
	evalidator ErrorValidator,
) (rerr error) {
	t.Log("linearsearch operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamLinearSearch(ctx)
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
				if err := evalidator(t, err); err != nil {
					rerr = errors.Join(
						rerr,
						errors.Errorf(
							"stream finished by an error: %s",
							err.Error(),
						),
					)
				}
				return
			}

			resp := res.GetResponse()
			if resp == nil {
				status := res.GetStatus()
				if status != nil {
					if e := svalidator(t, status.GetCode(), status.GetMessage()); e != nil {
						t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s",
							status.GetCode(),
							status.GetMessage(),
							errdetails.Serialize(status.GetDetails()))
						rerr = errors.Join(rerr, e)
					}
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
				Timeout:   timeout,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("linearsearch operation finished")

	return rerr
}

func (c *client) LinearSearchByID(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.LinearSearchByIDWithParameters(t,
		ctx,
		ds,
		100,
		3000000000,
		DefaultStatusValidator,
		ParseAndLogError,
	)
}

func (c *client) LinearSearchByIDWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	num uint32,
	timeout int64,
	svalidator StatusValidator,
	evalidator ErrorValidator,
) (rerr error) {
	t.Log("linearsearchByID operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamLinearSearchByID(ctx)
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
				if err := evalidator(t, err); err != nil {
					rerr = errors.Join(
						rerr,
						errors.Errorf(
							"stream finished by an error: %s",
							err.Error(),
						),
					)
				}
				return
			}

			resp := res.GetResponse()
			if resp == nil {
				status := res.GetStatus()
				if status != nil {
					if e := svalidator(t, status.GetCode(), status.GetMessage()); e != nil {
						t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s",
							status.GetCode(),
							status.GetMessage(),
							errdetails.Serialize(status.GetDetails()))
						rerr = errors.Join(rerr, e)
					}
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
				Timeout:   timeout,
			},
		})
		if err != nil {
			return err
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("linearsearchByID operation finished")

	return rerr
}

func (c *client) Insert(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.InsertWithParameters(t,
		ctx,
		ds,
		false,
		DefaultStatusValidator,
		ParseAndLogError,
	)
}

func (c *client) InsertWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	skipStrictExistCheck bool,
	svalidator StatusValidator,
	evalidator ErrorValidator,
) (rerr error) {
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
				if err := evalidator(t, err); err != nil {
					rerr = errors.Join(
						rerr,
						errors.Errorf(
							"stream finished by an error: %s",
							err.Error(),
						),
					)
				}
				return
			}

			loc := res.GetLocation()
			if loc == nil {
				status := res.GetStatus()
				if status != nil {
					if e := svalidator(t, status.GetCode(), status.GetMessage()); e != nil {
						t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s",
							status.GetCode(),
							status.GetMessage(),
							errdetails.Serialize(status.GetDetails()))
						rerr = errors.Join(rerr, e)
					}
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

	return rerr
}

func (c *client) Update(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.UpdateWithParameters(t,
		ctx,
		ds,
		false,
		1,
		DefaultStatusValidator,
		ParseAndLogError,
	)
}

func (c *client) UpdateWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	skipStrictExistCheck bool,
	offset int,
	svalidator StatusValidator,
	evalidator ErrorValidator,
) (rerr error) {
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
				if err := evalidator(t, err); err != nil {
					rerr = errors.Join(
						rerr,
						errors.Errorf(
							"stream finished by an error: %s",
							err.Error(),
						),
					)
				}
				return
			}

			loc := res.GetLocation()
			if loc == nil {
				status := res.GetStatus()
				if status != nil {
					if e := svalidator(t, status.GetCode(), status.GetMessage()); e != nil {
						t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s",
							status.GetCode(),
							status.GetMessage(),
							errdetails.Serialize(status.GetDetails()))
						rerr = errors.Join(rerr, e)
					}
					continue
				}

				t.Error("returned loc is nil")
				continue
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
				Vector: append(v[offset:], v[:offset]...),
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

	return rerr
}

func (c *client) Upsert(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.UpsertWithParameters(t,
		ctx,
		ds,
		false,
		2,
		DefaultStatusValidator,
		ParseAndLogError,
	)
}

func (c *client) UpsertWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	skipStrictExistCheck bool,
	offset int,
	svalidator StatusValidator,
	evalidator ErrorValidator,
) (rerr error) {
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
				if err := evalidator(t, err); err != nil {
					rerr = errors.Join(
						rerr,
						errors.Errorf(
							"stream finished by an error: %s",
							err.Error(),
						),
					)
				}
				return
			}

			loc := res.GetLocation()
			if loc == nil {
				status := res.GetStatus()
				if status != nil {
					if e := svalidator(t, status.GetCode(), status.GetMessage()); e != nil {
						t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s",
							status.GetCode(),
							status.GetMessage(),
							errdetails.Serialize(status.GetDetails()))
						rerr = errors.Join(rerr, e)
					}
					continue
				}

				t.Error("returned loc is nil")
				continue
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
				Vector: append(v[offset:], v[:offset]...),
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

	return rerr
}

func (c *client) Remove(t *testing.T, ctx context.Context, ds Dataset) error {
	return c.RemoveWithParameters(t,
		ctx,
		ds,
		false,
		DefaultStatusValidator,
		ParseAndLogError,
	)
}

func (c *client) RemoveWithParameters(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
	skipStrictExistCheck bool,
	svalidator StatusValidator,
	evalidator ErrorValidator,
) (rerr error) {
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
				if err := evalidator(t, err); err != nil {
					rerr = errors.Join(
						rerr,
						errors.Errorf(
							"stream finished by an error: %s",
							err.Error(),
						),
					)
				}
				return
			}

			loc := res.GetLocation()
			if loc == nil {
				status := res.GetStatus()
				if status != nil {
					if e := svalidator(t, status.GetCode(), status.GetMessage()); e != nil {
						t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s",
							status.GetCode(),
							status.GetMessage(),
							errdetails.Serialize(status.GetDetails()))
						rerr = errors.Join(rerr, e)
					}
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

	return rerr
}

func (c *client) Flush(t *testing.T, ctx context.Context) error {
	t.Log("flush operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.Flush(ctx, &payload.Flush_Request{})
	if err != nil {
		return err
	}

	t.Log("flush operation finished")

	return nil
}

func (c *client) RemoveByTimestamp(t *testing.T, ctx context.Context, timestamp int64) error {
	t.Log("removeByTimestamp operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	req := &payload.Remove_TimestampRequest{
		Timestamps: []*payload.Remove_Timestamp{
			{
				Timestamp: timestamp,
				Operator:  payload.Remove_Timestamp_Gt,
			},
		},
	}

	_, err = client.RemoveByTimestamp(ctx, req)
	if err != nil {
		return err
	}

	t.Log("removeByTimestamp operation finished")

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

func (c *client) GetObject(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
) (rerr error) {
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
				err = ParseAndLogError(t, err)
				rerr = errors.Join(
					rerr,
					errors.Errorf(
						"stream finished by an error: %s",
						err.Error(),
					),
				)
				return
			}

			resp := res.GetVector()
			if resp == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned:\tcode: %d\tmessage: %s\tdetails: %s", err.GetCode(), err.GetMessage(), errdetails.Serialize(err.GetDetails()))
					rerr = errors.Wrap(rerr, err.String())
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

			if ts := resp.GetTimestamp(); ts <= 0 {
				t.Error("timestamp is not set properly")
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

	return rerr
}

func (c *client) StreamListObject(
	t *testing.T,
	ctx context.Context,
	ds Dataset,
) error {
	t.Log("StreamListObject operation started")

	client, err := c.getClient(ctx)
	if err != nil {
		return err
	}

	sc, err := client.StreamListObject(ctx, &payload.Object_List_Request{})
	if err != nil {
		return err
	}

	// kv : [indexId]count
	indexCnt := make(map[string]int)
exit_loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			res, err := sc.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break exit_loop
				}
				return err
			}
			vec := res.GetVector()
			if vec == nil {
				st := res.GetStatus()
				return fmt.Errorf("returned vector is empty: code: %v, msg: %v, details: %v", st.GetCode(), st.GetMessage(), st.GetDetails())
			}
			if len(vec.GetVector()) == 0 {
				return fmt.Errorf("returned vector is empty: id: %v", vec.GetId())
			}
			indexCnt[vec.GetId()]++
		}
	}

	if len(indexCnt) != len(ds.Train) {
		return fmt.Errorf("the number of vectors returned is different: got %v, want %v", len(indexCnt), len(ds.Train))
	}

	replica := -1
	for k, v := range indexCnt {
		if replica == -1 {
			replica = v
			continue
		}
		if v != replica {
			return fmt.Errorf("the number of vectors returned is different at index id %v: got %v, want %v", k, v, replica)
		}
	}

	t.Log("StreamListObject operation finished successfully and all vectors are returned with correct replica number")
	return nil
}

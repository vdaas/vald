//go:build e2e

// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func (c *client) Search(t *testing.T, ctx context.Context, ds Dataset) (rerr error) {
	t.Log("Search operation started")

	client, err := c.getClient()
	if err != nil {
		return err
	}

	var (
		num     = uint32(100)
		radius  = -1
		epsilon = 0.01
		timeout = time.Second * 3
	)

	eg, egctx := errgroup.New(ctx)
	eg.SetLimit(100)
	mu := sync.Mutex{}
	for i := 0; i < len(ds.Test); i++ {
		eg.Go(func() error {
			id := strconv.Itoa(i)
			res, err := client.Search(egctx, &payload.Search_Request{
				Vector: ds.Test[i],
				Config: &payload.Search_Config{
					RequestId: id,
					Num:       num,
					Radius:    float32(radius),
					Epsilon:   float32(epsilon),
					Timeout:   timeout.Nanoseconds(),
				},
			})
			if err != nil {
				select {
				case <-egctx.Done():
					return nil
				default:
					return err
				}
			}
			resp := res.GetResults()
			topKIDs := make([]string, 0, len(resp))
			for _, d := range resp {
				topKIDs = append(topKIDs, d.Id)
			}
			if len(topKIDs) == 0 {
				t.Errorf("empty result is returned for test ID %s: %#v", res.GetRequestId(), topKIDs)
			}
			t.Logf("id: %s, results: %d, recall: %f", id, len(topKIDs), c.recall(topKIDs, ds.Neighbors[i][:len(topKIDs)]))
			// t.Logf("algo: %s, id: %d, results: %d, recall: %f", right, idx, len(topKIDs), c.recall(topKIDs, ds.Neighbors[idx][:len(topKIDs)]))
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		if err = ParseAndLogError(t, err); err != nil {
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
	}
	if rerr != nil {
		t.Fatalf("an error occured: %s", rerr.Error())
	}

	t.Log("search operation finished")
	return
}

func (c *client) SearchByID(t *testing.T, ctx context.Context, ds Dataset) error {
	return nil
}

func (c *client) LinearSearch(t *testing.T, ctx context.Context, ds Dataset) error {
	return nil
}

func (c *client) LinearSearchByID(t *testing.T, ctx context.Context, ds Dataset) error {
	return nil
}

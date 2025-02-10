//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// package crud provides e2e tests using ann-benchmarks datasets
package crud

import (
	"context"
	"flag"
	"os"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/params"
	"github.com/vdaas/vald/tests/e2e/hdf5"
	k8sclient "github.com/vdaas/vald/tests/e2e/kubernetes/client"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

var (
	cfg config.Data

	ctx    context.Context
	client vald.Client

	kclient k8sclient.Client

	ds *hdf5.Dataset
)

func TestMain(m *testing.M) {
	testing.Init()
	p, fail, err := params.New().Parse()
	flag.Parse()
	log.Init()
	if fail || err != nil || p.ConfigFilePath() == "" {
		log.Fatalf("failed to parse the parameters: %v", err)
	}

	if testing.Short() {
		log.Info("skipping this pkg test when -short because e2e test takes a long time")
		os.Exit(0)
	}

	cfg, err := config.Load(p.ConfigFilePath())
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if cfg.Kubernetes.PortForward.Enabled {
		kclient, err = k8sclient.New(cfg.Kubernetes.KubeConfig)
		if err != nil {
			log.Fatalf("failed to create kubernetes client: %v", err)
		}

		forwarder := kclient.Portforward(
			cfg.Kubernetes.PortForward.Namespace,
			cfg.Kubernetes.PortForward.PodName,
			int(cfg.Kubernetes.PortForward.LocalPort),
			int(cfg.Kubernetes.PortForward.PodPort),
		)
		err = forwarder.Start()
		if err != nil {
			log.Fatalf("failed to start portforward: %v, %#v", err, cfg.Kubernetes)
		}
		defer forwarder.Close()
	}

	ds, err = hdf5.HDF5ToDataset(cfg.Dataset.Name)
	if err != nil {
		log.Fatalf("failed to load dataset: %v", err)
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	client, ctx, err = newClient(ctx, cfg.Metadata)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ech, err := client.Start(ctx)
	if err != nil {
		log.Fatalf("failed to start client: %v", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		case err := <-ech:
			if err != nil {
				log.Fatalf("client daemon returned error: %v", err)
			}
		}
	}()
	defer func() {
		err = client.Stop(ctx)
		if err != nil {
			log.Fatalf("failed to stop client: %v", err)
		}
	}()

	m.Run()
}

func newClient(ctx context.Context, meta map[string]string) (client vald.Client, mctx context.Context, err error) {
	gopts, err := cfg.Target.Opts()
	if err != nil {
		return nil, nil, err
	}
	client, err = vald.New(
		vald.WithClient(
			grpc.New(gopts...),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	if meta != nil {
		mctx = metadata.NewOutgoingContext(ctx, metadata.New(meta))
	}
	return client, mctx, nil
}

func sleep(t *testing.T, dur time.Duration) {
	t.Logf("%v sleep for %s.", time.Now(), dur)
	time.Sleep(dur)
	t.Logf("%v sleep finished.", time.Now())
}

func recall(t *testing.T, resultIDs []string, neighbors []int) (recall float64) {
	t.Helper()
	ns := map[string]struct{}{}
	for _, n := range neighbors {
		ns[strconv.Itoa(n)] = struct{}{}
	}

	for _, r := range resultIDs {
		if _, ok := ns[r]; ok {
			recall++
		}
	}

	return recall / float64(len(neighbors))
}

func calculateRecall(t *testing.T, res *payload.Search_Response, idx int) (rc float64) {
	t.Helper()
	topKIDs := make([]string, 0, len(res.GetResults()))
	for _, d := range res.GetResults() {
		topKIDs = append(topKIDs, d.GetId())
	}

	if len(topKIDs) == 0 {
		t.Errorf("empty result is returned for test ID %s: %#v", res.GetRequestId(), topKIDs)
		return
	}
	rc = recall(t, topKIDs, ds.Neighbors[idx][:len(topKIDs)])
	return rc
}

func indexStatus(t *testing.T, ctx context.Context) {
	t.Helper()
	{
		res, err := client.IndexInfo(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexInfo %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexInfo %v", err)
			}
		}
		t.Logf("IndexInfo: %v", res.String())
	}
	{
		res, err := client.IndexDetail(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexDetail %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexDetail %v", err)
			}
		}
		t.Logf("IndexDetail: %v", res.String())
	}
	{
		res, err := client.IndexStatistics(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexStatistics %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexStatistics %v", err)
			}
		}
		t.Logf("IndexStatistics: %v", res.String())
	}
	{
		res, err := client.IndexStatisticsDetail(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexStatisticsDetail %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexStatisticsDetail %v", err)
			}
		}
		t.Logf("IndexStatisticsDetail: %v", res.String())
	}
}

func TestE2EUnaryCRUD(t *testing.T) {
	timestamp := time.Now().UnixNano()

	{
		res, err := client.IndexProperty(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexProperty %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexProperty %v", err)
			}
		}
		t.Logf("IndexProperty: %v", res.String())
	}

	eg, _ := errgroup.New(ctx)
	eg.SetLimit(int(cfg.Insert.Concurrency))
	for i, vec := range ds.Train[cfg.Insert.Offset : cfg.Insert.Offset+cfg.Insert.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Insert.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Insert(ctx, &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:        id,
					Vector:    vec,
					Timestamp: ts,
				},
				Config: &payload.Insert_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: cfg.Insert.SkipStrictExistCheck,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to insert vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to insert vector: %v", err)
				}
			}
			t.Logf("vector %v id %s inserted to %s", vec, id, res.String())
			return nil
		}))
	}
	eg.Wait()

	sleep(t, cfg.Index.WaitAfterInsert)

	indexStatus(t, ctx)

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Search.Concurrency))
	for i, vec := range ds.Test[cfg.Search.Offset : cfg.Search.Offset+cfg.Search.Num] {
		for _, query := range cfg.Search.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			eg.Go(safety.RecoverFunc(func() error {
				res, err := client.Search(ctx, &payload.Search_Request{
					Vector: vec,
					Config: &payload.Search_Config{
						RequestId:            rid,
						Num:                  query.K,
						Radius:               query.Radius,
						Epsilon:              query.Epsilon,
						Timeout:              query.Timeout.Nanoseconds(),
						AggregationAlgorithm: query.Algorithm,
						MinNum:               query.MinNum,
						Ratio:                wrapperspb.Float(query.Ratio),
						Nprobe:               query.Nprobe,
					},
				})
				if err != nil {
					t.Errorf("failed to search vector: %v", err)
				}
				t.Logf("vector %v id %s searched recall: %f, payload %s", vec, rid, calculateRecall(t, res, i), res.String())
				return nil
			}))
		}
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.SearchByID.Concurrency))
	for i, vec := range ds.Train[cfg.SearchByID.Offset : cfg.SearchByID.Offset+cfg.SearchByID.Num] {
		for _, query := range cfg.SearchByID.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			eg.Go(safety.RecoverFunc(func() error {
				res, err := client.SearchByID(ctx, &payload.Search_IDRequest{
					Id: id,
					Config: &payload.Search_Config{
						RequestId:            rid,
						Num:                  query.K,
						Radius:               query.Radius,
						Epsilon:              query.Epsilon,
						Timeout:              query.Timeout.Nanoseconds(),
						AggregationAlgorithm: query.Algorithm,
						MinNum:               query.MinNum,
						Ratio:                wrapperspb.Float(query.Ratio),
						Nprobe:               query.Nprobe,
					},
				})
				if err != nil {
					t.Errorf("failed to search vector: %v", err)
				}
				t.Logf("vector %v id %s searched recall: %f, payload %s", vec, rid, calculateRecall(t, res, i), res.String())
				return nil
			}))
		}
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.LinearSearch.Concurrency))
	for i, vec := range ds.Test[cfg.LinearSearch.Offset : cfg.LinearSearch.Offset+cfg.LinearSearch.Num] {
		for _, query := range cfg.LinearSearch.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			eg.Go(safety.RecoverFunc(func() error {
				res, err := client.LinearSearch(ctx, &payload.Search_Request{
					Vector: vec,
					Config: &payload.Search_Config{
						RequestId:            rid,
						Num:                  query.K,
						Radius:               query.Radius,
						Epsilon:              query.Epsilon,
						Timeout:              query.Timeout.Nanoseconds(),
						AggregationAlgorithm: query.Algorithm,
						MinNum:               query.MinNum,
						Ratio:                wrapperspb.Float(query.Ratio),
						Nprobe:               query.Nprobe,
					},
				})
				if err != nil {
					t.Errorf("failed to search vector: %v", err)
				}
				t.Logf("vector %v id %s searched recall: %f, payload %s", vec, rid, calculateRecall(t, res, i), res.String())
				return nil
			}))
		}
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.LinearSearchByID.Concurrency))
	for i, vec := range ds.Train[cfg.LinearSearchByID.Offset : cfg.LinearSearchByID.Offset+cfg.LinearSearchByID.Num] {
		for _, query := range cfg.LinearSearchByID.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			eg.Go(safety.RecoverFunc(func() error {
				res, err := client.LinearSearchByID(ctx, &payload.Search_IDRequest{
					Id: id,
					Config: &payload.Search_Config{
						RequestId:            rid,
						Num:                  query.K,
						Radius:               query.Radius,
						Epsilon:              query.Epsilon,
						Timeout:              query.Timeout.Nanoseconds(),
						AggregationAlgorithm: query.Algorithm,
						MinNum:               query.MinNum,
						Ratio:                wrapperspb.Float(query.Ratio),
						Nprobe:               query.Nprobe,
					},
				})
				if err != nil {
					t.Errorf("failed to search vector: %v", err)
				}
				t.Logf("vector %v id %s searched recall: %f, payload %s", vec, rid, calculateRecall(t, res, i), res.String())
				return nil
			}))
		}
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Object.Concurrency))
	for i := range ds.Train[cfg.Object.Offset : cfg.Object.Offset+cfg.Object.Num] {
		id := strconv.Itoa(i)
		eg.Go(safety.RecoverFunc(func() error {
			obj, err := client.GetObject(ctx, &payload.Object_VectorRequest{
				Id: &payload.Object_ID{Id: id},
			})
			if err != nil {
				t.Errorf("failed to get object: %v", err)
			}
			t.Logf("id %s got object: %v", id, obj.String())

			exists, err := client.Exists(ctx, &payload.Object_ID{Id: id})
			if err != nil {
				t.Errorf("failed to check object exists: %v", err)
			}
			t.Logf("id %s exists: %v", id, exists.String())

			res, err := client.GetTimestamp(ctx, &payload.Object_TimestampRequest{
				Id: &payload.Object_ID{Id: id},
			})
			if err != nil {
				t.Errorf("failed to get timestamp: %v", err)
			}
			t.Logf("id %s got timestamp: %v", id, res.String())
			return nil
		}))
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Update.Concurrency))
	for i, vec := range ds.Train[cfg.Update.Offset : cfg.Update.Offset+cfg.Update.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Update.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Update(ctx, &payload.Update_Request{
				Vector: &payload.Object_Vector{
					Id:        id,
					Vector:    vec,
					Timestamp: ts,
				},
				Config: &payload.Update_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: cfg.Update.SkipStrictExistCheck,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to update vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to update vector: %v", err)
				}
			}
			t.Logf("vector %v id %s updated to %s", vec, id, res.String())
			return nil
		}))
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Remove.Concurrency))
	for i := range ds.Train[cfg.Remove.Offset : cfg.Remove.Offset+cfg.Remove.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Remove.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Remove(ctx, &payload.Remove_Request{
				Id: &payload.Object_ID{Id: id},
				Config: &payload.Remove_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: cfg.Remove.SkipStrictExistCheck,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to remove vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to remove vector: %v", err)
				}
			}
			t.Logf("id %s'd vector removed to %s", id, res.String())
			return nil
		}))
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Upsert.Concurrency))
	for i, vec := range ds.Train[cfg.Upsert.Offset : cfg.Upsert.Offset+cfg.Upsert.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Upsert.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Upsert(ctx, &payload.Upsert_Request{
				Vector: &payload.Object_Vector{
					Id:        id,
					Vector:    vec,
					Timestamp: ts,
				},
				Config: &payload.Upsert_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: cfg.Upsert.SkipStrictExistCheck,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to upsert vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to upsert vector: %v", err)
				}
			}
			t.Logf("vector %v id %s upserted to %s", vec, id, res.String())
			return nil
		}))
	}
	eg.Wait()

	{
		rts := time.Now().Add(-time.Hour).UnixNano()
		res, err := client.RemoveByTimestamp(ctx, &payload.Remove_TimestampRequest{
			Timestamps: []*payload.Remove_Timestamp{
				{
					Timestamp: rts,
					Operator:  payload.Remove_Timestamp_Le,
				},
			},
		})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to remove by timestamp vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to remove by timestamp vector: %v", err)
			}
		}
		t.Logf("removed by timestamp %s to %s", time.Unix(0, rts).String(), res.String())
	}

	{
		res, err := client.Flush(ctx, &payload.Flush_Request{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to flush %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to flush %v", err)
			}
		}
		t.Logf("flushed %s", res.String())
	}

	indexStatus(t, ctx)
}
>>>>>>> d400d6723 (Refactor add V2 E2E testing for more maintainability)

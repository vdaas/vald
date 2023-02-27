//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package service manages the main logic of server.
package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/config"
	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/kvs"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service/vqueue"
	"google.golang.org/grpc"
)

type index struct {
	uuid string
	vec  []float32
}

func Test_ngt_InsertUpsert(t *testing.T) {
	if testing.Short() {
		t.Skip("The execution of this test takes a lot of time, so it is not performed during the short test\ttest: Test_ngt_InsertUpsert")
		return
	}
	type args struct {
		idxes    []index
		poolSize uint32
		bulkSize int
	}
	type fields struct {
		svcCfg  *config.NGT
		svcOpts []Option

		core              core.NGT
		eg                errgroup.Group
		kvs               kvs.BidiMap
		fmu               sync.Mutex
		fmap              map[string]uint32
		vq                vqueue.Queue
		indexing          atomic.Value
		saving            atomic.Value
		cimu              sync.Mutex
		lastNocie         uint64
		nocie             uint64
		nogce             uint64
		inMem             bool
		dim               int
		alen              int
		lim               time.Duration
		dur               time.Duration
		sdur              time.Duration
		minLit            time.Duration
		maxLit            time.Duration
		litFactor         time.Duration
		enableProactiveGC bool
		enableCopyOnWrite bool
		path              string
		smu               sync.Mutex
		tmpPath           atomic.Value
		oldPath           string
		basePath          string
		cowmu             sync.Mutex
		backupGen         uint64
		poolSize          uint32
		radius            float32
		epsilon           float32
		idelay            time.Duration
		dcd               bool
		kvsdbConcurrency  int
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	var (
		// default NGT configuration for test
		kvsdbCfg  = &config.KVSDB{}
		vqueueCfg = &config.VQueue{}
	)
	tests := []test{
		func() test {
			count := 10000000
			return test{
				name: fmt.Sprintf("insert & upsert %d random and 11 digits added to each vector element", count),
				args: args{
					idxes: createRandomData(count, &createRandomDataConfig{
						additionaldigits: 11,
					}),
					poolSize: uint32(count / 10),
					bulkSize: count / 10,
				},
				fields: fields{
					svcCfg: &config.NGT{
						Dimension:    128,
						DistanceType: core.Cosine.String(),
						ObjectType:   core.Uint8.String(),
						KVSDB:        kvsdbCfg,
						VQueue:       vqueueCfg,
					},
					svcOpts: []Option{
						WithEnableInMemoryMode(true),
					},
				},
			}
		}(),
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			eg, _ := errgroup.New(ctx)
			n, err := New(test.fields.svcCfg, append(test.fields.svcOpts, WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init ngt service, error = %v", err)
			}
			for _, idx := range test.args.idxes {
				err = n.Insert(idx.uuid, idx.vec)
				if err := checkFunc(test.want, err); err != nil {
					tt.Errorf("error = %v", err)
				}

			}

			log.Warn("start create index operation")
			err = n.CreateAndSaveIndex(ctx, test.args.poolSize)
			if err != nil {
				tt.Errorf("error creating index: %v", err)
			}
			log.Warn("start update operation")
			for i := 0; i < 100; i++ {
				idx := i
				eg.Go(func() error {
					log.Warnf("started %d-1", idx)
					for _, idx := range test.args.idxes[:len(test.args.idxes)/3] {
						_ = n.Delete(idx.uuid)
						_ = n.Insert(idx.uuid, idx.vec)
					}
					log.Warnf("finished %d-1", idx)
					return nil
				})

				eg.Go(func() error {
					log.Warnf("started %d-2", idx)
					for _, idx := range test.args.idxes[len(test.args.idxes)/3 : 2*len(test.args.idxes)/3] {
						_ = n.Delete(idx.uuid)
						_ = n.Insert(idx.uuid, idx.vec)
					}
					log.Warnf("finished %d-2", idx)
					return nil
				})

				eg.Go(func() error {
					log.Warnf("started %d-3", idx)
					for _, idx := range test.args.idxes[2*len(test.args.idxes)/3:] {
						_ = n.Delete(idx.uuid)
						_ = n.Insert(idx.uuid, idx.vec)
					}
					log.Warnf("finished %d-3", idx)
					return nil
				})
			}
			eg.Wait()

			log.Warn("start final create index operation")
			err = n.CreateAndSaveIndex(ctx, test.args.poolSize)
			if err != nil {
				tt.Errorf("error creating index: %v", err)
			}
		})
	}
}

// NOTE: After moving this implementation to the e2e package, remove this test function.
func Test_ngt_E2E(t *testing.T) {
	if testing.Short() {
		t.Skip("The execution of this test takes a lot of time, so it is not performed during the short test\ttest: Test_ngt_E2E")
		return
	}
	type args struct {
		requests []*payload.Upsert_MultiRequest

		addr     string
		dialOpts []grpc.DialOption
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	multiUpsertRequestGenFunc := func(idxes []index, chunk int) (res []*payload.Upsert_MultiRequest) {
		reqs := make([]*payload.Upsert_Request, 0, chunk)
		for i := 0; i < len(idxes); i++ {
			if len(reqs) == chunk-1 {
				res = append(res, &payload.Upsert_MultiRequest{
					Requests: reqs,
				})
				reqs = make([]*payload.Upsert_Request, 0, chunk)
			} else {
				reqs = append(reqs, &payload.Upsert_Request{
					Vector: &payload.Object_Vector{
						Id:     idxes[i].uuid,
						Vector: idxes[i].vec,
					},
					Config: &payload.Upsert_Config{
						SkipStrictExistCheck: true,
					},
				})
			}
		}
		if len(reqs) > 0 {
			res = append(res, &payload.Upsert_MultiRequest{
				Requests: reqs,
			})
		}
		return res
	}

	tests := []test{
		{
			name: "insert & upsert 100 random",
			args: args{
				requests: multiUpsertRequestGenFunc(
					createRandomData(500000, new(createRandomDataConfig)),
					50,
				),
				addr: "127.0.0.1:8080",
				dialOpts: []grpc.DialOption{
					grpc.WithInsecure(),
				},
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			conn, err := grpc.DialContext(ctx, test.args.addr, test.args.dialOpts...)
			if err := checkFunc(test.want, err); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := conn.Close(); err != nil {
					t.Error(err)
				}
			}()
			client := vald.NewValdClient(conn)

			for i := 0; i < 2; i++ {
				for _, req := range test.args.requests {
					_, err := client.MultiUpsert(ctx, req)
					if err != nil {
						t.Error(err)
					}
				}
				log.Info("%d step: finished all requests", i+1)
				time.Sleep(3 * time.Second)
			}
		})
	}
}

type createRandomDataConfig struct {
	additionaldigits int
}

func (cfg *createRandomDataConfig) verify() *createRandomDataConfig {
	if cfg == nil {
		cfg = new(createRandomDataConfig)
	}
	if cfg.additionaldigits < 0 {
		cfg.additionaldigits = 0
	}
	return cfg
}

func createRandomData(num int, cfg *createRandomDataConfig) []index {
	cfg = cfg.verify()

	var ad float32 = 1.0
	for i := 0; i < cfg.additionaldigits; i++ {
		ad = ad * 0.1
	}

	result := make([]index, 0)
	f32s, _ := vector.GenF32Vec(vector.NegativeUniform, num, 128)

	for idx, vec := range f32s {
		for i := range vec {
			if f := vec[i] * ad; f == 0.0 {
				if vec[i] > 0.0 {
					vec[i] = math.MaxFloat32
				} else if vec[i] < 0.0 {
					vec[i] = math.SmallestNonzeroFloat32
				}
				continue
			}
			vec[i] = vec[i] * ad
		}
		result = append(result, index{
			uuid: fmt.Sprintf("%s_%s-%s:%d:%d,%d", uuid.New().String(), uuid.New().String(), uuid.New().String(), idx, idx/100, idx%100),
			vec:  vec,
		})
	}

	return result
}

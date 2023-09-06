// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package grpc

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_server_CreateIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		c *payload.Control_CreateIndexRequest
	}
	type fields struct {
		srvOpts []Option
		svcCfg  *config.NGT
		svcOpts []service.Option
	}
	type want struct {
		wantRes *payload.Empty
		errCode status.Code
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Empty, error) error
		beforeFunc func(*testing.T, context.Context, args, Server)
		afterFunc  func(*testing.T, args)
	}

	// common variables for test
	const (
		name = "vald-agent-ngt-1" // agent name
		dim  = 3                  // vector dimension
		id   = "uuid-1"           // id for getObject request
	)
	var (
		ip = net.LoadLocalIP() // agent ip address

		// default NGT configuration for test
		defaultSvcCfg = &config.NGT{
			Dimension:    dim,
			DistanceType: ngt.Angle.String(),
			ObjectType:   ngt.Float.String(),
			KVSDB:        &config.KVSDB{},
			VQueue:       &config.VQueue{},
		}
		defaultSvcOpts = []service.Option{
			service.WithEnableInMemoryMode(true),
		}

		defaultInsertConfig = &payload.Insert_Config{}
		defaultRemoveConfig = &payload.Remove_Config{}
	)

	insert := func(ctx context.Context, s Server, cnt int) error {
		req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, cnt, dim, defaultInsertConfig)
		if err != nil {
			return err
		}
		if _, err := s.MultiInsert(ctx, req); err != nil {
			return err
		}

		return nil
	}
	remove := func(ctx context.Context, s Server, cnt int) error {
		_, err := s.MultiRemove(ctx, request.GenMultiRemoveReq(cnt, defaultRemoveConfig))
		return err
	}
	defaultCheckFunc := func(w want, gotRes *payload.Empty, err error) error {
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				return errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.errCode {
				return errors.Errorf("got code: \"%#v\",\n\t\t\t\twant code: \"%#v\"", st.Code(), w.errCode)
			}
		}

		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}

	/*
		- Equivalence Class Testing
			- case 1.1: success to create index with 1 uncommitted insert index
			- case 1.2: success to create index with 100 uncommitted insert index
			- case 2.1: success to create index with 1 uncommitted delete index
			- case 2.2: success to create index with 100 uncommitted delete index
			- case 3.1: success to create index with 1 uncommitted insert & delete index
			- case 3.2: success to create index with 100 uncommitted insert & delete index

		- Boundary Value Testing
			- case 1.1: fail to create index with 0 uncommitted index
			- case 2.1: success to create index with invalid dimension

		- Decision Table Testing
			// with uncommitted index count 100
			- case 1.1: success to create index with poolSize > uncommitted index count
			- case 1.2: success to create index with poolSize < uncommitted index count
			- case 1.3: success to create index with poolSize = uncommitted index count
			- case 1.4: success to create index with poolSize = 0
	*/
	tests := []test{
		{
			name: "Equivalence Class Testing case 1.1: success to create index with 1 uncommitted insert index",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 1,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				if err := insert(ctx, s, 1); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Equivalence Class Testing case 1.2: success to create index with 100 uncommitted insert index",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 100,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				if err := insert(ctx, s, 100); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Equivalence Class Testing case 2.1: success to create index with 1 uncommitted delete index",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 1,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				if err := insert(ctx, s, 1); err != nil {
					t.Fatal(err)
				}
				if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
					PoolSize: 1,
				}); err != nil {
					t.Fatal(err)
				}
				if err := remove(ctx, s, 1); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Equivalence Class Testing case 2.2: success to create index with 100 uncommitted delete index",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 1,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				cnt := 100
				if err := insert(ctx, s, cnt); err != nil {
					t.Fatal(err)
				}
				if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
					PoolSize: 1,
				}); err != nil {
					t.Fatal(err)
				}
				if err := remove(ctx, s, cnt); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Equivalence Class Testing case 3.1: success to create index with 1 uncommitted insert & delete index",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 1,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				insertCnt := 1
				removeCnt := 1
				if err := insert(ctx, s, insertCnt); err != nil {
					t.Fatal(err)
				}
				if err := remove(ctx, s, removeCnt); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Equivalence Class Testing case 3.2: success to create index with 100 uncommitted insert & delete index",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 100,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				insertCnt := 100
				removeCnt := 100
				if err := insert(ctx, s, insertCnt); err != nil {
					t.Fatal(err)
				}
				if err := remove(ctx, s, removeCnt); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Boundary Value Testing case 1.1: fail to create index with 0 uncommitted index",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 0,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			want: want{
				errCode: codes.FailedPrecondition,
			},
		},
		{
			name: "Boundary Value Testing case 2.1: success to create index with invalid dimension",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 100,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				srv, ok := s.(*server)
				if !ok {
					t.Error("Server cannot convert to *server")
				}

				insertCnt := 100
				invalidDim := dim + 1
				req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, invalidDim, defaultInsertConfig)
				if err != nil {
					t.Fatal(err)
				}

				for _, r := range req.Requests {
					if err := srv.ngt.Insert(r.Vector.Id, r.Vector.Vector); err != nil {
						t.Fatal(err)
					}
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Decision Table Testing case 1.1: success to create index with poolSize > uncommitted index count",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 10000,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				insertCnt := 100
				if err := insert(ctx, s, insertCnt); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Decision Table Testing case 1.2: success to create index with poolSize < uncommitted index count",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 1,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				insertCnt := 100
				if err := insert(ctx, s, insertCnt); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Decision Table Testing case 1.3: success to create index with poolSize = uncommitted index count",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 100,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				insertCnt := 100
				if err := insert(ctx, s, insertCnt); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
		{
			name: "Decision Table Testing case 1.4: success to create index with poolSize = 0",
			args: args{
				c: &payload.Control_CreateIndexRequest{
					PoolSize: 0,
				},
			},
			fields: fields{
				srvOpts: []Option{
					WithName(name),
					WithIP(ip),
				},
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				insertCnt := 100
				if err := insert(ctx, s, insertCnt); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Empty{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			eg, _ := errgroup.New(ctx)
			ngt, err := service.New(test.fields.svcCfg, append(test.fields.svcOpts, service.WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init ngt service, error = %v", err)
			}

			s, err := New(append(test.fields.srvOpts, WithNGT(ngt), WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init server, error= %v", err)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, ctx, test.args, s)
			}

			gotRes, err := s.CreateIndex(ctx, test.args.c)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_SaveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		in *payload.Empty
	}
	type fields struct {
		srvOpts   []Option
		svcCfg    *config.NGT
		svcOpts   []service.Option
		indexPath string // index path for svcOpts
	}
	type want struct {
		wantRes *payload.Empty
		errCode codes.Code
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(test, context.Context, Server, service.NGT, want, *payload.Empty, error) error
		beforeFunc func(*testing.T, context.Context, Server, service.NGT)
		afterFunc  func(test)
	}
	defaultCheckFunc := func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error {
		if (err == nil && w.errCode != 0) || (err != nil && w.errCode == 0) {
			return errors.Errorf("got error is %v, but want error code is %v", err, w.errCode)
		}
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				return errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.errCode {
				return errors.Errorf("got code: \"%#v\",\n\t\t\t\twant code: \"%#v\"", st.Code(), w.errCode)
			}
		}

		if diff := comparator.Diff(gotRes, w.wantRes, comparator.IgnoreUnexported(payload.Empty{})); diff != "" {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}

	// common variables for test
	const (
		name = "vald-agent-ngt-1" // agent name
		dim  = 3                  // vector dimension
		id   = "uuid-1"           // id for getObject request
	)
	var (
		// agent ip address
		ip = net.LoadLocalIP()

		// default NGT configuration for test
		defaultSvcCfg = &config.NGT{
			Dimension:    dim,
			DistanceType: ngt.Angle.String(),
			ObjectType:   ngt.Float.String(),
			KVSDB:        &config.KVSDB{},
			VQueue:       &config.VQueue{},
		}
		defaultSrvOpts = []Option{
			WithName(name),
			WithIP(ip),
		}
		defaultSvcOpts = []service.Option{
			service.WithEnableInMemoryMode(false),
		}

		emptyPayload = &payload.Empty{}
	)

	mkdirTemp := func() string {
		d, err := os.MkdirTemp("", "")
		if err != nil {
			t.Error(err)
		}
		return d
	}
	defaultAfterFunc := func(test test) {
		os.RemoveAll(test.fields.indexPath)
	}
	// this function checks the backup file can be loaded and check if it contain the wantVecs indexes.
	// it creates a new ngt and server instance with the backup file, and checks if we can retrieve all of wantVecs indexes
	// and check the total index count matches with wantVecs count.
	checkBackupFolder := func(fields fields, ctx context.Context, wantVecs []*payload.Insert_Request) error {
		// create another server instance to check if any vector is inserted and saved to the backup dir
		eg, _ := errgroup.New(ctx)
		ngt, err := service.New(fields.svcCfg, append(fields.svcOpts,
			service.WithErrGroup(eg),
			service.WithIndexPath(fields.indexPath))...)
		if err != nil {
			return errors.Errorf("failed to init ngt service, error = %v", err)
		}
		srv, err := New(append(fields.srvOpts, WithNGT(ngt), WithErrGroup(eg))...)
		if err != nil {
			return errors.Errorf("failed to init server, error= %v", err)
		}

		// get object and check if the vector is equals to inserted one
		for _, ir := range wantVecs {
			obj, err := srv.GetObject(ctx, &payload.Object_VectorRequest{
				Id: &payload.Object_ID{
					Id: ir.GetVector().GetId(),
				},
			})
			if err != nil {
				return err
			}

			// FIXME: remove these 2 lines after migrating Config.Timestamp to Vector.Timestamp
			wantVec := ir.GetVector()
			wantVec.Timestamp = obj.Timestamp

			if !reflect.DeepEqual(obj, wantVec) {
				return errors.Errorf("vector is not match, got: %v, want: %v", obj, ir)
			}
		}

		// check total index count is same
		ii, err := srv.IndexInfo(ctx, emptyPayload)
		if err != nil {
			return err
		}

		wantIndexInfo := &payload.Info_Index_Count{
			Stored: uint32(len(wantVecs)),
		}
		if !reflect.DeepEqual(ii, wantIndexInfo) {
			return errors.Errorf("stored index count not correct, got: %v, want: %v", ii, wantIndexInfo)
		}

		return nil
	}
	/*
		- Equivalence Class Testing (with copy on write disable)
			- case 1.1: success to save 1 inserted index
			- case 1.2: success to save 100 inserted index
			- case 2.1: fail to save index with no write access on backup folder
				- this test case will be check in service layer as it has exteneral dependencies (file permission)
			- case 3.1: success to save index when other save index process is running
		- Boundary Value Testing
			- case 1.1: success to save index with no index
			- case 2.1: success to save index with invalid dimension
				- the invalid index will be removed from NGT and the index file
		- Decision Table Testing
			- case 1.1: success to save index with in-memory mode
				- do nothing and no file will be created
			- case 2.1: success to save 1 inserted index with copy-on-write enabled
			- case 2.2: success to save 100 inserted index with copy-on-write enabled
	*/
	tests := []test{
		func() test {
			irs, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, 1, dim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 1.1: success to save 1 inserted index",
				args: args{
					in: emptyPayload,
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT) {
					t.Helper()
					if _, err := s.Insert(ctx, irs.Requests[0]); err != nil {
						t.Error(err)
					}
					// we need to create index before saving to store the indexed vector
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: 1,
					}); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					if err := checkBackupFolder(test.fields, ctx, irs.Requests); err != nil {
						return err
					}

					return nil
				},
			}
		}(),
		func() test {
			irs, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, 100, dim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 1.2: success to save 100 inserted index",
				args: args{
					in: emptyPayload,
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT) {
					t.Helper()
					if _, err := s.MultiInsert(ctx, irs); err != nil {
						t.Error(err)
					}
					// we need to create index before saving to store the indexed vector
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: 1,
					}); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					if err := checkBackupFolder(test.fields, ctx, irs.Requests); err != nil {
						return err
					}

					return nil
				},
			}
		}(),
		func() test {
			// bulk insert request to make saveIndex cost time
			insertNum := 100000
			irs, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertNum, dim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Equivalence Class Testing case 3.1: success to save index when other save index process is running",
				args: args{
					in: emptyPayload,
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT) {
					t.Helper()
					if _, err := s.MultiInsert(ctx, irs); err != nil {
						t.Error(err)
					}
					// we need to create index before saving to store the indexed vector
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertNum),
					}); err != nil {
						t.Error(err)
					}

					// ensure the goroutine call is scheduled and the first saveIndex() will be executed in the goroutine
					wg := sync.WaitGroup{}
					wg.Add(1)
					go func() {
						wg.Done()
						// execute in goroutine to ensure it is executed twice
						s.SaveIndex(ctx, emptyPayload)
					}()
					wg.Wait() // wait until the goroutine is scheduled and wg.Done() is executed
				},
				want: want{
					wantRes: emptyPayload,
				},
			}
		}(),
		func() test {
			return test{
				name: "Boundary Value Testing case 1.1: success to save index with no index",
				args: args{
					in: emptyPayload,
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					if err := checkBackupFolder(test.fields, ctx, nil); err != nil {
						return err
					}

					return nil
				},
			}
		}(),
		func() test {
			return test{
				name: "Boundary Value Testing case 2.1: success to save index with invalid dimension",
				args: args{
					in: emptyPayload,
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT) {
					t.Helper()
					invalidDim := dim + 1
					vecs, err := vector.GenF32Vec(vector.Gaussian, 1, invalidDim)
					if err != nil {
						t.Error(err)
					}

					// insert invalid vector to ngt directly
					if err := n.Insert("uuid-1", vecs[0]); err != nil {
						t.Error(err)
					}
					// we need to create index before saving to store the indexed vector
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: 1,
					}); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					// the invalid index will be removed inside ngt
					if err := checkBackupFolder(test.fields, ctx, nil); err != nil {
						return err
					}

					return nil
				},
			}
		}(),
		func() test {
			return test{
				name: "Decision Table Testing case 1.1: success to save index with in-memory mode",
				args: args{
					in: emptyPayload,
				},
				fields: fields{
					srvOpts: defaultSrvOpts,
					svcCfg:  defaultSvcCfg,
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
					indexPath: mkdirTemp(),
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}

					files, err := file.ListInDir(test.fields.indexPath)
					if err != nil {
						return err
					}

					// check any file is generated in backup directory
					if len(files) > 0 {
						return errors.New("no file should be created when in memory mode is enabled")
					}

					return nil
				},
			}
		}(),
		func() test {
			irs, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, 1, dim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 2.1: success to save 1 inserted index with copy-on-write enabled",
				args: args{
					in: emptyPayload,
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   append(defaultSvcOpts, service.WithCopyOnWrite(true)),
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT) {
					t.Helper()
					if _, err := s.Insert(ctx, irs.Requests[0]); err != nil {
						t.Error(err)
					}
					// we need to create index before saving to store the indexed vector
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: 1,
					}); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					if err := checkBackupFolder(test.fields, ctx, irs.Requests); err != nil {
						return err
					}

					return nil
				},
			}
		}(),
		func() test {
			irs, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, 100, dim, nil)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "Decision Table Testing case 2.1: success to save 100 inserted index with copy-on-write enabled",
				args: args{
					in: emptyPayload,
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   append(defaultSvcOpts, service.WithCopyOnWrite(true)),
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT) {
					t.Helper()
					if _, err := s.MultiInsert(ctx, irs); err != nil {
						t.Error(err)
					}
					// we need to create index before saving to store the indexed vector
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: 1,
					}); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					if err := checkBackupFolder(test.fields, ctx, irs.Requests); err != nil {
						return err
					}

					return nil
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

			afterFunc := test.afterFunc
			if test.afterFunc == nil {
				afterFunc = defaultAfterFunc
			}
			defer afterFunc(test)

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			eg, _ := errgroup.New(ctx)
			ngt, err := service.New(test.fields.svcCfg, append(test.fields.svcOpts,
				service.WithErrGroup(eg),
				service.WithIndexPath(test.fields.indexPath))...)
			if err != nil {
				tt.Errorf("failed to init ngt service, error = %v", err)
			}

			s, err := New(append(test.fields.srvOpts, WithNGT(ngt), WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init server, error= %v", err)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, ctx, s, ngt)
			}

			gotRes, err := s.SaveIndex(ctx, test.args.in)
			if err := checkFunc(test, ctx, s, ngt, test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_CreateAndSaveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		c *payload.Control_CreateIndexRequest
	}
	type fields struct {
		srvOpts   []Option
		svcCfg    *config.NGT
		svcOpts   []service.Option
		indexPath string // index path for svcOpts
	}
	type want struct {
		wantRes *payload.Empty
		errCode codes.Code
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error
		beforeFunc func(*testing.T, context.Context, Server, service.NGT, test)
		afterFunc  func(*testing.T, test)
	}

	// common variables for test
	const (
		name = "vald-agent-ngt-1" // agent name
		dim  = 3                  // vector dimension
		id   = "uuid-1"           // id for getObject request
	)
	var (
		// agent ip address
		ip = net.LoadLocalIP()

		// default NGT configuration for test
		defaultSvcCfg = &config.NGT{
			Dimension:    dim,
			DistanceType: ngt.Angle.String(),
			ObjectType:   ngt.Float.String(),
			KVSDB:        &config.KVSDB{},
			VQueue:       &config.VQueue{},
		}
		defaultSrvOpts = []Option{
			WithName(name),
			WithIP(ip),
		}
		defaultSvcOpts = []service.Option{
			service.WithEnableInMemoryMode(false),
		}
		defaultInsertConfig = &payload.Insert_Config{}
		emptyPayload        = &payload.Empty{}
	)

	defaultCheckFunc := func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error {
		if (err == nil && w.errCode != 0) || (err != nil && w.errCode == 0) {
			return errors.Errorf("got error is %v, but want error code is %v", err, w.errCode)
		}
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				return errors.Errorf("got error cannot convert to Status: \"%#v\"", err)
			}
			if st.Code() != w.errCode {
				return errors.Errorf("got code: \"%#v\",\n\t\t\t\twant code: \"%#v\"", st.Code(), w.errCode)
			}
		}

		if diff := comparator.Diff(gotRes, w.wantRes, comparator.IgnoreUnexported(payload.Empty{})); diff != "" {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}

	mkdirTemp := func() string {
		d, err := os.MkdirTemp("", "")
		if err != nil {
			t.Error(err)
		}
		return d
	}
	defaultAfterFunc := func(t *testing.T, test test) {
		t.Helper()
		os.RemoveAll(test.fields.indexPath)
	}

	// this function checks the backup file can be loaded and check if it contains the wantVecs indexes.
	// it creates a new ngt and server instance with the backup file, and checks if we can retrieve all of wantVecs indexes
	// and check the total index count matches with wantVecs count.
	checkBackupFolder := func(fields fields, ctx context.Context, wantVecs []*payload.Insert_Request) error {
		// create another server instance to check if any vector is inserted and saved to the backup dir
		eg, _ := errgroup.New(ctx)
		ngt, err := service.New(fields.svcCfg, append(fields.svcOpts,
			service.WithErrGroup(eg),
			service.WithIndexPath(fields.indexPath))...)
		if err != nil {
			return errors.Errorf("failed to init ngt service, error = %v", err)
		}
		srv, err := New(append(fields.srvOpts, WithNGT(ngt), WithErrGroup(eg))...)
		if err != nil {
			return errors.Errorf("failed to init server, error= %v", err)
		}

		// get object and check if the vector is equals to inserted one
		for _, ir := range wantVecs {
			obj, err := srv.GetObject(ctx, &payload.Object_VectorRequest{
				Id: &payload.Object_ID{
					Id: ir.GetVector().GetId(),
				},
			})
			if err != nil {
				return err
			}

			// FIXME: remove these 2 lines after migrating Config.Timestamp to Vector.Timestamp
			wantVec := ir.GetVector()
			wantVec.Timestamp = obj.Timestamp

			if !reflect.DeepEqual(obj, ir.GetVector()) {
				return errors.Errorf("vector is not match, got: %v, want: %v", obj, ir)
			}
		}

		// check total index count is same
		ii, err := srv.IndexInfo(ctx, emptyPayload)
		if err != nil {
			return err
		}

		wantIndexInfo := &payload.Info_Index_Count{
			Stored: uint32(len(wantVecs)),
		}
		if !reflect.DeepEqual(ii, wantIndexInfo) {
			return errors.Errorf("stored index count not correct, got: %v, want: %v", ii, wantIndexInfo)
		}

		return nil
	}

	/*
		- Equivalence Class Testing (with copy on write disable)
			- case 1.1: success to create and save 1 uncommitted insert index
			- case 1.2: success to create and save 100 uncommitted insert index
			- case 2.1: success to create and save 1 uncommitted delete index
			- case 2.2: success to create and save 100 uncommitted delete index
			- case 3.1: success to create and save 1 uncommitted update index
			- case 3.2: success to create and save 100 uncommitted update index
		- Boundary Value Testing
			- case 1.1: fail to create and save 0 index
			- case 2.1: success to create and save index with invalid dimension
				- the invalid index will be removed from NGT and the index file
		- Decision Table Testing
			- case 1.1: success to create and save 100 index with in-memory mode
				- do nothing and no file will be created
			- case 2.1: success to create and save 1 inserted index with copy-on-write enabled
			- case 2.2: success to create and save 100 inserted index with copy-on-write enabled

			// with uncommitted index count 100
			- case 3.1: success to create and save index with poolSize > uncommitted index count
			- case 3.2: success to create and save index with poolSize < uncommitted index count
			- case 3.3: success to create and save index with poolSize = uncommitted index count
			- case 3.4: success to create and save index with poolSize = 0
	*/
	tests := []test{
		func() test {
			insertCnt := 1
			var ir *payload.Insert_MultiRequest

			return test{
				name: "Equivalence Class Testing case 1.1: success to create and save 1 uncommitted insert index",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if ir, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}
					if _, err := s.MultiInsert(ctx, ir); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: &payload.Empty{},
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, gotRes, err); err != nil {
						return err
					}
					return checkBackupFolder(test.fields, ctx, ir.GetRequests())
				},
			}
		}(),
		func() test {
			insertCnt := 100
			var ir *payload.Insert_MultiRequest

			return test{
				name: "Equivalence Class Testing case 1.2: success to create and save 100 uncommitted insert index",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if ir, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}
					if _, err := s.MultiInsert(ctx, ir); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: &payload.Empty{},
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, gotRes, err); err != nil {
						return err
					}
					return checkBackupFolder(test.fields, ctx, ir.GetRequests())
				},
			}
		}(),
		func() test {
			insertCnt := 100
			var ir *payload.Insert_MultiRequest

			return test{
				name: "Equivalence Class Testing case 2.1: success to create and save 1 uncommitted delete index",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if ir, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}

					// insert 100 request
					if _, err := s.MultiInsert(ctx, ir); err != nil {
						t.Error(err)
					}
					if _, err := s.CreateAndSaveIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					}); err != nil {
						t.Error(err)
					}

					// delete 1 request
					if _, err := s.Remove(ctx, &payload.Remove_Request{
						Id: &payload.Object_ID{
							Id: ir.GetRequests()[0].GetVector().GetId(),
						},
					}); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: &payload.Empty{},
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, gotRes, err); err != nil {
						return err
					}
					// we expect the request[0] is removed and shouldn't be exists in backup files
					expectedVecs := ir.GetRequests()[1:]
					return checkBackupFolder(test.fields, ctx, expectedVecs)
				},
			}
		}(),
		func() test {
			insertCnt := 200
			var ir *payload.Insert_MultiRequest

			return test{
				name: "Equivalence Class Testing case 2.2: success to create and save 100 uncommitted delete index",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if ir, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}

					// remove requests
					rr := make([]*payload.Remove_Request, 100)
					for i := 0; i < 100; i++ {
						rr[i] = &payload.Remove_Request{
							Id: &payload.Object_ID{
								Id: ir.GetRequests()[i].GetVector().GetId(),
							},
						}
					}

					// insert 200 request
					if _, err := s.MultiInsert(ctx, ir); err != nil {
						t.Error(err)
					}
					if _, err := s.CreateAndSaveIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					}); err != nil {
						t.Error(err)
					}

					// delete 100 request
					if _, err := s.MultiRemove(ctx, &payload.Remove_MultiRequest{
						Requests: rr,
					}); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: &payload.Empty{},
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, gotRes, err); err != nil {
						return err
					}
					// we expect the request[0-99] is removed and shouldn't be exists in backup files
					expectedVecs := ir.GetRequests()[100:]
					return checkBackupFolder(test.fields, ctx, expectedVecs)
				},
			}
		}(),
		func() test {
			insertCnt := 1
			var ir *payload.Insert_MultiRequest
			var updateVec []float32 // the updated vector

			return test{
				name: "Equivalence Class Testing case 3.1: success to create and save 1 uncommitted update index",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if ir, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}

					updateID := ir.GetRequests()[0].GetVector().GetId()
					updateVecs, err := vector.GenF32Vec(vector.Gaussian, 1, dim)
					if err != nil {
						t.Error(err)
					}
					updateVec = updateVecs[0]

					// insert 1 request
					if _, err := s.MultiInsert(ctx, ir); err != nil {
						t.Error(err)
					}
					if _, err := s.CreateAndSaveIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					}); err != nil {
						t.Error(err)
					}

					// update vector request
					if _, err := s.Update(ctx, &payload.Update_Request{
						Vector: &payload.Object_Vector{
							Id:     updateID,
							Vector: updateVec,
						},
					}); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: &payload.Empty{},
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, gotRes, err); err != nil {
						return err
					}
					// we expect vector is the update vec
					expectedVecs := ir.GetRequests()
					expectedVecs[0].Vector.Vector = updateVec
					return checkBackupFolder(test.fields, ctx, expectedVecs)
				},
			}
		}(),
		func() test {
			insertCnt := 100
			var ir *payload.Insert_MultiRequest
			var updateVecs [][]float32
			updateReqs := make([]*payload.Update_Request, insertCnt)

			return test{
				name: "Equivalence Class Testing case 3.2: success to create and save 100 uncommitted update index",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if ir, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}

					// generate another set of vectors for update
					updateVecs, err = vector.GenF32Vec(vector.Gaussian, insertCnt, dim)
					if err != nil {
						t.Error(err)
					}

					// generate update requests for insert
					for i := range ir.GetRequests() {
						updateReqs[i] = &payload.Update_Request{
							Vector: &payload.Object_Vector{
								Id:     ir.GetRequests()[i].GetVector().GetId(),
								Vector: updateVecs[i],
							},
						}
					}

					// insert 100 request
					if _, err := s.MultiInsert(ctx, ir); err != nil {
						t.Error(err)
					}
					if _, err := s.CreateAndSaveIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					}); err != nil {
						t.Error(err)
					}

					// update vector request
					if _, err := s.MultiUpdate(ctx, &payload.Update_MultiRequest{
						Requests: updateReqs,
					}); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: &payload.Empty{},
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, gotRes, err); err != nil {
						return err
					}
					// we expect vector is the update vec
					expectedVecs := ir.GetRequests()
					for i := range expectedVecs {
						expectedVecs[i].GetVector().Vector = updateVecs[i]
					}

					return checkBackupFolder(test.fields, ctx, expectedVecs)
				},
			}
		}(),
		func() test {
			return test{
				name: "Boundary Value Testing case 1.1: fail to create and save 0 index",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: 0,
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				want: want{
					wantRes: &payload.Empty{},
				},
			}
		}(),
		func() test {
			return test{
				name: "Boundary Value Testing case 2.1: success to create and save index with invalid dimension",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: 1,
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   defaultSvcOpts,
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					invalidDim := dim + 1
					vecs, err := vector.GenF32Vec(vector.Gaussian, 1, invalidDim)
					if err != nil {
						t.Error(err)
					}

					// insert invalid vector to ngt directly
					if err := n.Insert("uuid-1", vecs[0]); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					// the invalid index will be removed inside ngt
					return checkBackupFolder(test.fields, ctx, nil)
				},
			}
		}(),
		func() test {
			insertCnt := 100
			var ir *payload.Insert_MultiRequest

			return test{
				name: "Decision Table Testing case 1.1: success to create and save 100 index with in-memory mode",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts: defaultSrvOpts,
					svcCfg:  defaultSvcCfg,
					svcOpts: []service.Option{
						service.WithEnableInMemoryMode(true),
					},
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if ir, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}
					if _, err := s.MultiInsert(ctx, ir); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: &payload.Empty{},
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, gotRes *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, gotRes, err); err != nil {
						return err
					}
					files, err := file.ListInDir(test.fields.indexPath)
					if err != nil {
						return err
					}

					// check any file is generated in backup directory
					if len(files) > 0 {
						return errors.New("no file should be created when in memory mode is enabled")
					}
					return nil
				},
			}
		}(),
		func() test {
			insertCnt := 1
			var irs *payload.Insert_MultiRequest

			return test{
				name: "Decision Table Testing case 2.1: success to create and save 1 inserted index with copy-on-write enabled",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   append(defaultSvcOpts, service.WithCopyOnWrite(true)),
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if irs, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}
					if _, err := s.Insert(ctx, irs.Requests[0]); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					return checkBackupFolder(test.fields, ctx, irs.Requests)
				},
			}
		}(),
		func() test {
			insertCnt := 100
			var irs *payload.Insert_MultiRequest

			return test{
				name: "Decision Table Testing case 2.2: success to create and save 100 inserted index with copy-on-write enabled",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   append(defaultSvcOpts, service.WithCopyOnWrite(true)),
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if irs, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}
					if _, err := s.MultiInsert(ctx, irs); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					return checkBackupFolder(test.fields, ctx, irs.Requests)
				},
			}
		}(),
		func() test {
			insertCnt := 100
			var irs *payload.Insert_MultiRequest

			return test{
				name: "Decision Table Testing case 3.1: success to create and save index with poolSize > uncommitted index count",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt + 1),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   append(defaultSvcOpts, service.WithCopyOnWrite(true)),
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if irs, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}
					if _, err := s.MultiInsert(ctx, irs); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					return checkBackupFolder(test.fields, ctx, irs.Requests)
				},
			}
		}(),
		func() test {
			insertCnt := 100
			var irs *payload.Insert_MultiRequest

			return test{
				name: "Decision Table Testing case 3.2: success to create and save index with poolSize < uncommitted index count",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt - 1),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   append(defaultSvcOpts, service.WithCopyOnWrite(true)),
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if irs, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}
					if _, err := s.MultiInsert(ctx, irs); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					return checkBackupFolder(test.fields, ctx, irs.Requests)
				},
			}
		}(),
		func() test {
			insertCnt := 100
			var irs *payload.Insert_MultiRequest

			return test{
				name: "Decision Table Testing case 3.3: success to create and save index with poolSize = uncommitted index count",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   append(defaultSvcOpts, service.WithCopyOnWrite(true)),
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if irs, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}
					if _, err := s.MultiInsert(ctx, irs); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					return checkBackupFolder(test.fields, ctx, irs.Requests)
				},
			}
		}(),
		func() test {
			insertCnt := 100
			var irs *payload.Insert_MultiRequest

			return test{
				name: "Decision Table Testing case 3.4: success to create and save index with poolSize = 0",
				args: args{
					c: &payload.Control_CreateIndexRequest{
						PoolSize: 0,
					},
				},
				fields: fields{
					srvOpts:   defaultSrvOpts,
					svcCfg:    defaultSvcCfg,
					svcOpts:   append(defaultSvcOpts, service.WithCopyOnWrite(true)),
					indexPath: mkdirTemp(),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, s Server, n service.NGT, test test) {
					t.Helper()
					var err error
					if irs, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig); err != nil {
						t.Error(err)
					}
					if _, err := s.MultiInsert(ctx, irs); err != nil {
						t.Error(err)
					}
				},
				want: want{
					wantRes: emptyPayload,
				},
				checkFunc: func(test test, ctx context.Context, s Server, n service.NGT, w want, e *payload.Empty, err error) error {
					if err := defaultCheckFunc(test, ctx, s, n, w, e, err); err != nil {
						return err
					}
					return checkBackupFolder(test.fields, ctx, irs.Requests)
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

			afterFunc := test.afterFunc
			if test.afterFunc == nil {
				afterFunc = defaultAfterFunc
			}
			defer afterFunc(tt, test)

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			eg, _ := errgroup.New(ctx)
			ngt, err := service.New(
				test.fields.svcCfg,
				append(test.fields.svcOpts,
					service.WithErrGroup(eg),
					service.WithIndexPath(test.fields.indexPath),
				)...)
			if err != nil {
				tt.Errorf("failed to init ngt service, error = %v", err)
			}

			s, err := New(append(test.fields.srvOpts, WithNGT(ngt), WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init server, error= %v", err)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, ctx, s, ngt, test)
			}

			gotRes, err := s.CreateAndSaveIndex(ctx, test.args.c)
			if err := checkFunc(test, ctx, s, ngt, test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_IndexInfo(t *testing.T) {
	t.Parallel()
	type args struct {
		in1 *payload.Empty
	}
	type fields struct {
		name              string
		ip                string
		streamConcurrency int
		svcCfg            *config.NGT
		svcOpts           []service.Option
	}
	type want struct {
		wantRes *payload.Info_Index_Count
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(Server, context.Context, args, want, *payload.Info_Index_Count, error) error
		beforeFunc func(*testing.T, context.Context, args, Server)
		afterFunc  func(*testing.T, args)
	}

	// common variables for test
	const (
		name      = "vald-agent-ngt-1" // agent name
		dim       = 3                  //  vector dimension
		id        = "uuid-1"           // id for getObject request
		insertCnt = 100                // default insert count
		removeCnt = 100                // default remove count
	)
	var (
		ip = net.LoadLocalIP() // agent ip address

		// default NGT configuration for test
		kvsdbCfg  = &config.KVSDB{}
		vqueueCfg = &config.VQueue{}

		defaultSvcCfg = &config.NGT{
			Dimension:    dim,
			DistanceType: ngt.Angle.String(),
			ObjectType:   ngt.Float.String(),
			KVSDB:        kvsdbCfg,
			VQueue:       vqueueCfg,
		}
		defaultSvcOpts = []service.Option{
			service.WithEnableInMemoryMode(true),
		}

		defaultInsertConfig = &payload.Insert_Config{}
		defaultRemoveConfig = &payload.Remove_Config{}
	)
	insertAndCreateIndex := func(s Server, ctx context.Context, cnt int, createIdx bool) (*payload.Insert_MultiRequest, error) {
		req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, cnt, dim, defaultInsertConfig)
		if err != nil {
			return nil, err
		}

		if _, err := s.MultiInsert(ctx, req); err != nil {
			return nil, err
		}

		if createIdx {
			if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
				PoolSize: uint32(len(req.Requests)),
			}); err != nil {
				return nil, err
			}
		}
		return req, nil
	}
	defaultCheckFunc := func(s Server, ctx context.Context, args args, w want, gotRes *payload.Info_Index_Count, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(gotRes, w.wantRes, comparator.IgnoreUnexported(payload.Info_Index_Count{})); diff != "" {
			return errors.New(diff)
		}
		return nil
	}
	periodicallyCheckIndexInfoFunc := func(s Server, ctx context.Context, args args, w want, chkFunc func(want, *payload.Info_Index_Count, error) error) error {
		timeout := time.After(5 * time.Second)
		ticker := time.Tick(10 * time.Millisecond)
		for {
			select {
			case <-timeout:
				gotRes, err := s.IndexInfo(ctx, args.in1)
				if err := chkFunc(w, gotRes, err); err != nil {
					return err
				}
			case <-ticker:
				gotRes, err := s.IndexInfo(ctx, args.in1)
				if err := chkFunc(w, gotRes, err); err == nil {
					return nil
				}
			}
		}
	}

	/*
		- Equivalence Class Testing
			- case 1.1: return stored count when NGT is empty
			- case 1.2: return stored count with 100 number of indexes
			- case 2.1: return uncommitted count 0 when NGT is empty
			- case 2.2: return uncommitted count with 100 uncommitted insert index
			- case 2.3: return uncommitted count with 100 uncommitted delete index
			- case 2.4: return uncommitted count with 100 uncommitted insert+delete index
			- case 3.1: return when NGT is indexing
			- case 3.2: return when NGT is not indexing
			- case 4.1: return when NGT is saving index
			- case 4.2: return when NGT is not saving index
		- Boundary Value Testing
			- NONE
		- Decision Table Testing
			- NONE
	*/
	tests := []test{
		{
			name: "Equivalence Class Testing case 1.1: return stored count when NGT is empty",
			args: args{
				in1: &payload.Empty{},
			},
			fields: fields{
				name:    name,
				ip:      ip,
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      0,
					Uncommitted: 0,
					Indexing:    false,
					Saving:      false,
				},
			},
		},
		func() test {
			tmpDir, err := os.MkdirTemp("", "")
			if err != nil {
				t.Error(err)
			}

			svcCfg := &config.NGT{
				Dimension:    dim,
				DistanceType: ngt.Angle.String(),
				ObjectType:   ngt.Float.String(),
				KVSDB:        kvsdbCfg,
				VQueue:       vqueueCfg,
			}
			svcOpts := append(defaultSvcOpts,
				service.WithIndexPath(tmpDir),
				service.WithEnableInMemoryMode(false),
			)

			// create server to insert index before test
			ctx, cancel := context.WithCancel(context.Background())
			eg, _ := errgroup.New(ctx)
			ngt, err := service.New(svcCfg, append(svcOpts, service.WithErrGroup(eg))...)
			if err != nil {
				t.Errorf("failed to init ngt service, error = %v", err)
			}
			s, err := New(WithErrGroup(eg),
				WithNGT(ngt),
				WithName(name),
				WithIP(ip),
			)
			if err != nil {
				t.Error(err)
			}

			// insert 100 indexes and create index
			if _, err := insertAndCreateIndex(s, ctx, 100, true); err != nil {
				t.Error(err)
			}

			// save index to file
			if _, err := s.SaveIndex(ctx, &payload.Empty{}); err != nil {
				t.Error(err)
			}

			// close this ngt instance
			ngt.Close(ctx)
			cancel()

			return test{
				name: "Equivalence Class Testing case 1.2: return stored count with 100 number of indexes",
				args: args{
					in1: &payload.Empty{},
				},
				fields: fields{
					name:    name,
					ip:      ip,
					svcCfg:  svcCfg,
					svcOpts: svcOpts,
				},
				want: want{
					wantRes: &payload.Info_Index_Count{
						Stored:      insertCnt,
						Uncommitted: 0,
						Indexing:    false,
						Saving:      false,
					},
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					cancel()
					os.RemoveAll(tmpDir)
				},
			}
		}(),
		{
			name: "Equivalence Class Testing case 2.1: return uncommitted count 0 when NGT is empty",
			args: args{
				in1: &payload.Empty{},
			},
			fields: fields{
				name:    name,
				ip:      ip,
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      0,
					Uncommitted: 0,
					Indexing:    false,
					Saving:      false,
				},
			},
		},
		{
			name: "Equivalence Class Testing case 2.2: return uncommitted count with 100 uncommitted insert index",
			args: args{
				in1: &payload.Empty{},
			},
			fields: fields{
				name:    name,
				ip:      ip,
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				if _, err := insertAndCreateIndex(s, ctx, insertCnt, false); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      0,
					Uncommitted: insertCnt,
					Indexing:    false,
					Saving:      false,
				},
			},
		},
		{
			name: "Equivalence Class Testing case 2.3: return uncommitted count with 100 uncommitted delete index",
			args: args{
				in1: &payload.Empty{},
			},
			fields: fields{
				name:    name,
				ip:      ip,
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				// we need to insert request first before remove
				req, err := insertAndCreateIndex(s, ctx, removeCnt, true)
				if err != nil {
					t.Fatal(err)
				}

				// remove the inserted indexes above for the uncommitted count
				for _, r := range req.Requests {
					if _, err := s.Remove(ctx, &payload.Remove_Request{
						Id: &payload.Object_ID{
							Id: r.GetVector().GetId(),
						},
						Config: defaultRemoveConfig,
					}); err != nil {
						t.Fatal(err)
					}
				}
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      removeCnt,
					Uncommitted: removeCnt,
					Indexing:    false,
					Saving:      false,
				},
			},
		},
		{
			name: "Equivalence Class Testing case 2.4: return uncommitted count with 100 uncommitted insert+delete index",
			args: args{
				in1: &payload.Empty{},
			},
			fields: fields{
				name:    name,
				ip:      ip,
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				// we need vectors inserted before removal
				rreq, err := insertAndCreateIndex(s, ctx, removeCnt, true)
				if err != nil {
					t.Fatal(err)
				}

				// remove requests inserted above for the uncommitted remove count
				for _, r := range rreq.Requests {
					if _, err := s.Remove(ctx, &payload.Remove_Request{
						Id: &payload.Object_ID{
							Id: r.GetVector().GetId(),
						},
						Config: defaultRemoveConfig,
					}); err != nil {
						t.Fatal(err)
					}
				}

				// insert requests for the uncommitted insert count
				if _, err := insertAndCreateIndex(s, ctx, insertCnt, false); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      insertCnt,
					Uncommitted: insertCnt + removeCnt,
					Indexing:    false,
					Saving:      false,
				},
			},
		},
		{
			name: "Equivalence Class Testing case 3.1: return when NGT is indexing",
			args: args{
				in1: &payload.Empty{},
			},
			fields: fields{
				name:   name,
				ip:     ip,
				svcCfg: defaultSvcCfg,
				svcOpts: append(defaultSvcOpts,
					service.WithDefaultPoolSize(100),
					service.WithEnableInMemoryMode(true),
				),
			},
			beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
				t.Helper()
				insertCnt := 10000
				if _, err := insertAndCreateIndex(s, ctx, insertCnt, false); err != nil {
					t.Fatal(err)
				}
				go func() {
					if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
						PoolSize: uint32(insertCnt),
					}); err != nil {
						t.Error(err)
					}
				}()
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Indexing: true,
				},
			},
			checkFunc: func(s Server, ctx context.Context, args args, w want, i *payload.Info_Index_Count, err error) error {
				chk := func(w want, i *payload.Info_Index_Count, err error) error {
					if err != nil {
						return errors.Errorf("expected no error, got error: %v", err)
					}
					if i == nil || i.GetIndexing() != w.wantRes.GetIndexing() {
						return errors.Errorf("expected indexing, got: %#v", i)
					}
					return nil
				}

				return periodicallyCheckIndexInfoFunc(s, ctx, args, w, chk)
			},
		},
		{
			name: "Equivalence Class Testing case 3.2: return when NGT is not indexing",
			args: args{
				in1: &payload.Empty{},
			},
			fields: fields{
				name:    name,
				ip:      ip,
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      0,
					Uncommitted: 0,
					Indexing:    false,
					Saving:      false,
				},
			},
		},
		func() test {
			tmpDir, err := file.MkdirTemp("")
			if err != nil {
				t.Fatal(err)
			}

			return test{
				name: "Equivalence Class Testing case 4.1: return when NGT is saving index",
				args: args{
					in1: &payload.Empty{},
				},
				fields: fields{
					name:   name,
					ip:     ip,
					svcCfg: defaultSvcCfg,
					svcOpts: append(defaultSvcOpts,
						service.WithDefaultPoolSize(100),
						service.WithEnableInMemoryMode(false),
						service.WithIndexPath(tmpDir),
					),
				},
				beforeFunc: func(t *testing.T, ctx context.Context, a args, s Server) {
					t.Helper()
					insertCnt := 10000
					if _, err := insertAndCreateIndex(s, ctx, insertCnt, true); err != nil {
						t.Fatal(err)
					}

					go func() {
						if _, err := s.SaveIndex(ctx, &payload.Empty{}); err != nil {
							// since the context closed error will be returned after checkFunc,
							// we ignore the error here to avoid go test error
							// t.Log(err)
						}
					}()
				},
				want: want{
					wantRes: &payload.Info_Index_Count{
						Saving: true,
					},
				},
				checkFunc: func(s Server, ctx context.Context, args args, w want, _ *payload.Info_Index_Count, _ error) error {
					chk := func(w want, i *payload.Info_Index_Count, err error) error {
						if err != nil {
							return errors.Errorf("expected no error, got error: %v", err)
						}
						if i == nil || i.GetSaving() != w.wantRes.GetSaving() {
							return errors.Errorf("expected indexing, got: %#v", i)
						}
						return nil
					}

					return periodicallyCheckIndexInfoFunc(s, ctx, args, w, chk)
				},
				afterFunc: func(t *testing.T, _ args) {
					t.Helper()
					os.RemoveAll(tmpDir)
				},
			}
		}(),
		{
			name: "Equivalence Class Testing case 4.2: return when NGT is not saving index",
			args: args{
				in1: &payload.Empty{},
			},
			fields: fields{
				name:    name,
				ip:      ip,
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      0,
					Uncommitted: 0,
					Indexing:    false,
					Saving:      false,
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

			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			eg, _ := errgroup.New(ctx)
			ngt, err := service.New(test.fields.svcCfg, append(test.fields.svcOpts, service.WithErrGroup(eg))...)
			if err != nil {
				tt.Errorf("failed to init ngt service, error = %v", err)
			}

			s, err := New(WithErrGroup(eg),
				WithNGT(ngt),
				WithName(test.fields.name),
				WithIP(test.fields.ip),
				WithStreamConcurrency(test.fields.streamConcurrency),
			)
			if err != nil {
				tt.Error(err)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, ctx, test.args, s)
			}

			gotRes, err := s.IndexInfo(ctx, test.args.in1)

			if err := checkFunc(s, ctx, test.args, test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW

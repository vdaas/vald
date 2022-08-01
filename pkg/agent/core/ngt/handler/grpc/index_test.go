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
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	itest "github.com/vdaas/vald/internal/test"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
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
		afterFunc  func(args)
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
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
		ctx context.Context
		in1 *payload.Empty
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Empty
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Empty, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Empty, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           in1: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx: nil,
		           in1: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.SaveIndex(test.args.ctx, test.args.in1)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_CreateAndSaveIndex(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		c   *payload.Control_CreateIndexRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Empty
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Empty, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Empty, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           c: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx: nil,
		           c: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.CreateAndSaveIndex(test.args.ctx, test.args.c)
			if err := checkFunc(test.want, gotRes, err); err != nil {
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
		afterFunc  func(args)
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
		{
			name: "Equivalence Class Testing case 1.2: return stored count with 100 number of indexes",
			args: args{
				in1: &payload.Empty{},
			},
			fields: fields{
				name: name,
				ip:   ip,
				svcCfg: &config.NGT{
					Dimension:    784,
					DistanceType: ngt.Angle.String(),
					ObjectType:   ngt.Float.String(),
					KVSDB:        kvsdbCfg,
					VQueue:       vqueueCfg,
				},
				svcOpts: append(defaultSvcOpts,
					service.WithIndexPath(itest.GetTestdataPath("backup/100index")),
					service.WithEnableInMemoryMode(false),
				),
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      insertCnt,
					Uncommitted: 0,
					Indexing:    false,
					Saving:      false,
				},
			},
		},
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
					insertCnt := 10000
					if _, err := insertAndCreateIndex(s, ctx, insertCnt, true); err != nil {
						t.Fatal(err)
					}

					go func() {
						if _, err := s.SaveIndex(ctx, &payload.Empty{}); err != nil {
							t.Error(err)
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
				afterFunc: func(a args) {
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
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

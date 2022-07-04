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
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_server_CreateIndex(t *testing.T) {
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

			gotRes, err := s.CreateIndex(test.args.ctx, test.args.c)
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
		checkFunc  func(want, *payload.Info_Index_Count, error) error
		beforeFunc func(*testing.T, args, *server)
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
	defaultCheckFunc := func(w want, gotRes *payload.Info_Index_Count, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(gotRes, w.wantRes, comparator.IgnoreUnexported(payload.Info_Index_Count{})); diff != "" {
			return errors.New(diff)
		}
		// if !reflect.DeepEqual(gotRes, w.wantRes) {
		// 	return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		// }
		return nil
	}

	/*
		- Equivalence Class Testing
			- case 1.1: return stored count when NGT is empty
			- case 1.2: return stored count with 100 number of indexes
			- case 2.1: return uncommitted count when NGT is empty
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
				name:    name,
				ip:      ip,
				svcCfg:  defaultSvcCfg,
				svcOpts: defaultSvcOpts,
			},
			beforeFunc: func(t *testing.T, a args, s *server) {
				ctx := context.Background()

				req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
				if err != nil {
					t.Fatal(err)
				}

				if _, err := s.MultiInsert(ctx, req); err != nil {
					t.Fatal(err)
				}
				if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
					PoolSize: uint32(len(req.Requests)),
				}); err != nil {
					t.Fatal(err)
				}
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
			name: "Equivalence Class Testing case 2.1: return uncommitted count when NGT is empty",
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
			beforeFunc: func(t *testing.T, a args, s *server) {
				ctx := context.Background()

				req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
				if err != nil {
					t.Fatal(err)
				}

				if _, err := s.MultiInsert(ctx, req); err != nil {
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
			beforeFunc: func(t *testing.T, a args, s *server) {
				ctx := context.Background()

				// we need to insert request first before remove
				req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
				if err != nil {
					t.Fatal(err)
				}

				if _, err := s.MultiInsert(ctx, req); err != nil {
					t.Fatal(err)
				}
				if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
					PoolSize: uint32(len(req.Requests)),
				}); err != nil {
					t.Fatal(err)
				}

				// remove
				if _, err := s.MultiRemove(ctx, request.GenMultiRemoveReq(removeCnt, defaultRemoveConfig)); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      insertCnt,
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
			beforeFunc: func(t *testing.T, a args, s *server) {
				ctx := context.Background()

				// we need to insert request first before remove
				totalInsertCnt := insertCnt + removeCnt
				req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, totalInsertCnt, dim, defaultInsertConfig)
				if err != nil {
					t.Fatal(err)
				}

				if _, err := s.MultiInsert(ctx, req); err != nil {
					t.Fatal(err)
				}
				if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
					PoolSize: uint32(len(req.Requests)),
				}); err != nil {
					t.Fatal(err)
				}

				// remove
				if _, err := s.MultiRemove(ctx, request.GenMultiRemoveReq(removeCnt, defaultRemoveConfig)); err != nil {
					t.Fatal(err)
				}
				// insert again
				req, err = request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
				if err != nil {
					t.Fatal(err)
				}

				if _, err := s.MultiInsert(ctx, req); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      insertCnt + removeCnt,
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
				name: name,
				ip:   ip,
				svcCfg: &config.NGT{
					Dimension:    dim,
					DistanceType: ngt.Angle.String(),
					ObjectType:   ngt.Float.String(),
					KVSDB: &config.KVSDB{
						Concurrency: 1,
					},
					VQueue: &config.VQueue{
						InsertBufferPoolSize: 1000000000,
						DeleteBufferPoolSize: 1000000000,
					},
				},
				svcOpts: append(defaultSvcOpts,
					service.WithAutoIndexLength(1000000000),
					service.WithAutoIndexCheckDuration("10m"),
					service.WithAutoSaveIndexDuration("10m"),
					service.WithAutoIndexDurationLimit("10m"),
					service.WithDefaultPoolSize(1000000000),
				),
			},
			beforeFunc: func(t *testing.T, a args, s *server) {
				ctx := context.Background()

				insertCnt := 50
				req, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, dim, defaultInsertConfig)
				if err != nil {
					t.Fatal(err)
				}

				if _, err := s.MultiInsert(ctx, req); err != nil {
					t.Fatal(err)
				}

				t.Log("inserted")
				// go func() {
				if _, err := s.CreateIndex(ctx, &payload.Control_CreateIndexRequest{
					PoolSize: uint32(insertCnt),
				}); err != nil {
					// t.Fatal(err)
				}
				// }()
			},
			want: want{
				wantRes: &payload.Info_Index_Count{
					Stored:      0,
					Uncommitted: 0,
					Indexing:    true,
					Saving:      false,
				},
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
		{
			name: "Equivalence Class Testing case 4.1: return when NGT is saving index",
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

			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               ngt,
				eg:                eg,
				streamConcurrency: test.fields.streamConcurrency,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, s)
			}

			gotRes, err := s.IndexInfo(ctx, test.args.in1)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

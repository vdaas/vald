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
package service

import (
	"context"
	"testing"

	tmock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	iconfig "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/test/mock"
	"github.com/vdaas/vald/pkg/index/job/correction/config"
)

type mockDiscovererClient struct {
	client mock.ClientInternal
}

func (*mockDiscovererClient) Start(context.Context) (<-chan error, error) {
	return nil, nil
}

func (*mockDiscovererClient) GetAddrs(context.Context) []string {
	return nil
}

func (m *mockDiscovererClient) GetClient() grpc.Client {
	return &m.client
}

func Test_correct_correctTimestamp(t *testing.T) {
	t.Parallel()

	// This mock just returns nil and record args inside
	m := mockDiscovererClient{}
	m.client.On("Do", tmock.Anything, tmock.Anything, tmock.Anything).Return(nil, nil)
	c := &correct{
		discoverer: &m,
	}

	type args struct {
		target *vectorReplica
		found  []*vectorReplica
	}

	type want struct {
		addrs []string
		err   error
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "nothing happens when no replica is found",
			args: args{
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id:        "target",
						Timestamp: 100,
					},
				},
				found: []*vectorReplica{},
			},
			want: want{
				addrs: nil,
				err:   nil,
			},
		},
		{
			name: "updates one found vec when found vecs are older than target",
			args: args{
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id:        "target",
						Timestamp: 100,
					},
				},
				found: []*vectorReplica{
					{
						addr: "found",
						vec: &payload.Object_Vector{
							Id:        "found",
							Timestamp: 99,
						},
					},
				},
			},
			want: want{
				addrs: []string{"found"},
				err:   nil,
			},
		},
		{
			name: "updates multiple found vecs when found vecs are older than target",
			args: args{
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id:        "target",
						Timestamp: 100,
					},
				},
				found: []*vectorReplica{
					{
						addr: "found1",
						vec: &payload.Object_Vector{
							Id:        "found",
							Timestamp: 99,
						},
					},
					{
						addr: "found2",
						vec: &payload.Object_Vector{
							Id:        "found",
							Timestamp: 98,
						},
					},
				},
			},
			want: want{
				addrs: []string{"found1", "found2"},
				err:   nil,
			},
		},
		{
			name: "updates target vec when found vecs are newer than target",
			args: args{
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id:        "target",
						Timestamp: 0,
					},
				},
				found: []*vectorReplica{
					{
						addr: "found1",
						vec: &payload.Object_Vector{
							Id:        "found",
							Timestamp: 99,
						},
					},
				},
			},
			want: want{
				addrs: []string{"target"},
				err:   nil,
			},
		},
		{
			name: "updates target vec and one of found vecs with the latest found vec",
			args: args{
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id:        "target",
						Timestamp: 0,
					},
				},
				found: []*vectorReplica{
					{
						addr: "found1",
						vec: &payload.Object_Vector{
							Id:        "found",
							Timestamp: 99,
						},
					},
					{
						addr: "latest",
						vec: &payload.Object_Vector{
							Id:        "found",
							Timestamp: 100,
						},
					},
				},
			},
			want: want{
				addrs: []string{"target", "found1"},
				err:   nil,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			err := c.correctTimestamp(context.Background(), test.args.target, test.args.found)
			require.Equal(tt, test.want.err, err)

			for _, addr := range test.want.addrs {
				// check if the agents which need to be corrected are called
				// checking calling parameter, like timestamp, is impossible because its inside of the function arg
				m.client.AssertCalled(tt, "Do", tmock.Anything, addr, tmock.Anything)
			}
		})
	}
}

func Test_correct_correctReplica(t *testing.T) {
	t.Parallel()

	// This mock just returns nil and record args inside
	m := mockDiscovererClient{}
	m.client.On("Do", tmock.Anything, tmock.Anything, tmock.Anything).Return(nil, nil)

	type args struct {
		indexReplica   int
		target         *vectorReplica
		found          []*vectorReplica
		availableAddrs []string
	}

	type addrMethod struct {
		addr   string
		method string
	}

	type want struct {
		addrMethods []addrMethod
		err         error
	}

	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "nothing happens when replica number sutisfies",
			args: args{
				indexReplica: 2,
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id: "target",
					},
				},
				found: []*vectorReplica{
					{
						addr: "found",
						vec: &payload.Object_Vector{
							Id: "found",
						},
					},
				},
				availableAddrs: []string{},
			},
			want: want{
				addrMethods: nil,
				err:         nil,
			},
		},
		{
			name: "insert replica when replica number is not enough",
			args: args{
				indexReplica: 2,
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id: "target",
					},
				},
				found:          []*vectorReplica{},
				availableAddrs: []string{"available"},
			},
			want: want{
				addrMethods: []addrMethod{
					{
						addr:   "available",
						method: insertMethod,
					},
				},
				err: nil,
			},
		},
		{
			name: "insert replica to the agent with most memory available",
			args: args{
				indexReplica: 2,
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id: "target",
					},
				},
				found: []*vectorReplica{},
				// this is supposed to be sorted by memory usage with descending order
				availableAddrs: []string{"most memory used", "second memory used"},
			},
			want: want{
				addrMethods: []addrMethod{
					{
						addr:   "second memory used",
						method: insertMethod,
					},
				},
				err: nil,
			},
		},
		{
			name: "delete replica from myself when replica number is too much by one",
			args: args{
				indexReplica: 2,
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id: "target",
					},
				},
				found: []*vectorReplica{
					{
						addr: "found1",
					},
					{
						addr: "found2",
					},
				},
				availableAddrs: []string{},
			},
			want: want{
				addrMethods: []addrMethod{
					{
						addr:   "target",
						method: deleteMethod,
					},
				},
				err: nil,
			},
		},
		{
			name: "delete replica from myself and most memory used agent when replica number is too much by more than one",
			args: args{
				indexReplica: 2,
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id: "target",
					},
				},
				found: []*vectorReplica{
					{
						addr: "found1",
					},
					{
						addr: "found2",
					},
					{
						addr: "found3",
					},
				},
				availableAddrs: []string{},
			},
			want: want{
				addrMethods: []addrMethod{
					{
						addr:   "target",
						method: deleteMethod,
					},
					{
						addr:   "found1",
						method: deleteMethod,
					},
				},
				err: nil,
			},
		},
		{
			name: "return ErrNoAvailableAgentToInsert when availableAddrs is empty when insertion required",
			args: args{
				indexReplica: 2,
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id: "target",
					},
				},
				found:          []*vectorReplica{},
				availableAddrs: []string{},
			},
			want: want{
				addrMethods: nil,
				err:         errors.ErrNoAvailableAgentToInsert,
			},
		},
		{
			name: "return ErrFailedToCorrectReplicaNum when there is not enough number of availableAddrs",
			args: args{
				indexReplica: 3,
				target: &vectorReplica{
					addr: "target",
					vec: &payload.Object_Vector{
						Id: "target",
					},
				},
				found:          []*vectorReplica{},
				availableAddrs: []string{"available"},
			},
			want: want{
				addrMethods: nil,
				err:         errors.ErrFailedToCorrectReplicaNum,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		c := &correct{
			discoverer: &m,
			cfg: &config.Data{
				Corrector: &iconfig.Corrector{
					IndexReplica: test.args.indexReplica,
				},
			},
		}

		// agentAddrs = availableAddrs + target.addr + found.addr
		// skipcq: CRT-D0001
		c.agentAddrs = append(test.args.availableAddrs, test.args.target.addr)
		for _, found := range test.args.found {
			c.agentAddrs = append(c.agentAddrs, found.addr)
		}

		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			err := c.correctReplica(context.Background(), test.args.target, test.args.found)
			if test.want.err != nil {
				require.ErrorIs(t, test.want.err, err)
			}

			for _, am := range test.want.addrMethods {
				// check if the agents which need to be corrected are called with the required method
				// checking calling parameter, like timestamp, is impossible because its inside of the function arg
				m.client.AssertCalled(tt, "Do", tmock.MatchedBy(func(ctx context.Context) bool {
					method := ctx.Value(grpc.GrpcMethodContextKey)
					val, ok := method.(string)
					if !ok {
						return false
					}
					return val == am.method
				}), am.addr, tmock.Anything)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		cfg        *config.Data
// 		discoverer discoverer.Client
// 	}
// 	type want struct {
// 		want Corrector
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Corrector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Corrector, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cfg:nil,
// 		           discoverer:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           cfg:nil,
// 		           discoverer:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := New(test.args.cfg, test.args.discoverer)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_correct_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		cfg                   *config.Data
// 		discoverer            discoverer.Client
// 		agentAddrs            []string
// 		indexInfos            sync.Map[string, *payload.Info_Index_Count]
// 		uuidsCount            uint32
// 		uncommittedUUIDsCount uint32
// 		checkedID             bbolt.Bbolt
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           cfg:nil,
// 		           discoverer:nil,
// 		           agentAddrs:nil,
// 		           indexInfos:nil,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           checkedID:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           cfg:nil,
// 		           discoverer:nil,
// 		           agentAddrs:nil,
// 		           indexInfos:nil,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           checkedID:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			c := &correct{
// 				cfg:                   test.fields.cfg,
// 				discoverer:            test.fields.discoverer,
// 				agentAddrs:            test.fields.agentAddrs,
// 				indexInfos:            test.fields.indexInfos,
// 				uuidsCount:            test.fields.uuidsCount,
// 				uncommittedUUIDsCount: test.fields.uncommittedUUIDsCount,
// 				checkedID:             test.fields.checkedID,
// 			}
//
// 			got, err := c.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_correct_PreStop(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		cfg                   *config.Data
// 		discoverer            discoverer.Client
// 		agentAddrs            []string
// 		indexInfos            sync.Map[string, *payload.Info_Index_Count]
// 		uuidsCount            uint32
// 		uncommittedUUIDsCount uint32
// 		checkedID             bbolt.Bbolt
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           cfg:nil,
// 		           discoverer:nil,
// 		           agentAddrs:nil,
// 		           indexInfos:nil,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           checkedID:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           cfg:nil,
// 		           discoverer:nil,
// 		           agentAddrs:nil,
// 		           indexInfos:nil,
// 		           uuidsCount:0,
// 		           uncommittedUUIDsCount:0,
// 		           checkedID:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			c := &correct{
// 				cfg:                   test.fields.cfg,
// 				discoverer:            test.fields.discoverer,
// 				agentAddrs:            test.fields.agentAddrs,
// 				indexInfos:            test.fields.indexInfos,
// 				uuidsCount:            test.fields.uuidsCount,
// 				uncommittedUUIDsCount: test.fields.uncommittedUUIDsCount,
// 				checkedID:             test.fields.checkedID,
// 			}
//
// 			err := c.PreStop(test.args.in0)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }

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
package request

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

var defaultObjectLocationComparators = []cmp.Option{
	comparator.IgnoreUnexported(payload.Object_Locations{}),
	comparator.IgnoreUnexported(payload.Object_Location{}),
}

func TestGenObjectLocations(t *testing.T) {
	type args struct {
		num    int
		name   string
		ipAddr string
	}
	type want struct {
		want *payload.Object_Locations
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *payload.Object_Locations) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *payload.Object_Locations) error {
		if diff := comparator.Diff(got, w.want, defaultObjectLocationComparators...); diff != "" {
			return errors.Errorf("diff: %v", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "success to generate 1 object location",
			args: args{
				num:    1,
				name:   "vald-agent-01",
				ipAddr: "127.0.0.1",
			},
			want: want{
				want: &payload.Object_Locations{
					Locations: []*payload.Object_Location{
						{
							Name: "vald-agent-01",
							Uuid: "uuid-1",
							Ips:  []string{"127.0.0.1"},
						},
					},
				},
			},
		},
		{
			name: "success to generate 5 object location",
			args: args{
				num:    5,
				name:   "vald-agent-01",
				ipAddr: "127.0.0.1",
			},
			want: want{
				want: &payload.Object_Locations{
					Locations: []*payload.Object_Location{
						{
							Name: "vald-agent-01",
							Uuid: "uuid-1",
							Ips:  []string{"127.0.0.1"},
						},
						{
							Name: "vald-agent-01",
							Uuid: "uuid-2",
							Ips:  []string{"127.0.0.1"},
						},
						{
							Name: "vald-agent-01",
							Uuid: "uuid-3",
							Ips:  []string{"127.0.0.1"},
						},
						{
							Name: "vald-agent-01",
							Uuid: "uuid-4",
							Ips:  []string{"127.0.0.1"},
						},
						{
							Name: "vald-agent-01",
							Uuid: "uuid-5",
							Ips:  []string{"127.0.0.1"},
						},
					},
				},
			},
		},
		{
			name: "success to generate 0 object location",
			args: args{
				num:    0,
				name:   "vald-agent-01",
				ipAddr: "127.0.0.1",
			},
			want: want{
				want: &payload.Object_Locations{
					Locations: []*payload.Object_Location{},
				},
			},
		},
		// max num test is ignored due to test timeout
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

			got := GenObjectLocations(test.args.num, test.args.name, test.args.ipAddr)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestGenObjectStreamLocation(t *testing.T) {
	type args struct {
		num    int
		name   string
		ipAddr string
	}
	type want struct {
		want []*payload.Object_StreamLocation
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []*payload.Object_StreamLocation) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []*payload.Object_StreamLocation) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "success to generate 1 object stream location",
			args: args{
				num:    1,
				name:   "vald-agent-01",
				ipAddr: "127.0.0.1",
			},
			want: want{
				want: []*payload.Object_StreamLocation{
					{
						Payload: &payload.Object_StreamLocation_Location{
							Location: &payload.Object_Location{
								Name: "vald-agent-01",
								Uuid: "uuid-1",
								Ips:  []string{"127.0.0.1"},
							},
						},
					},
				},
			},
		},
		{
			name: "success to generate 5 object stream location",
			args: args{
				num:    5,
				name:   "vald-agent-01",
				ipAddr: "127.0.0.1",
			},
			want: want{
				want: []*payload.Object_StreamLocation{
					{
						Payload: &payload.Object_StreamLocation_Location{
							Location: &payload.Object_Location{
								Name: "vald-agent-01",
								Uuid: "uuid-1",
								Ips:  []string{"127.0.0.1"},
							},
						},
					},
					{
						Payload: &payload.Object_StreamLocation_Location{
							Location: &payload.Object_Location{
								Name: "vald-agent-01",
								Uuid: "uuid-2",
								Ips:  []string{"127.0.0.1"},
							},
						},
					},
					{
						Payload: &payload.Object_StreamLocation_Location{
							Location: &payload.Object_Location{
								Name: "vald-agent-01",
								Uuid: "uuid-3",
								Ips:  []string{"127.0.0.1"},
							},
						},
					},
					{
						Payload: &payload.Object_StreamLocation_Location{
							Location: &payload.Object_Location{
								Name: "vald-agent-01",
								Uuid: "uuid-4",
								Ips:  []string{"127.0.0.1"},
							},
						},
					},
					{
						Payload: &payload.Object_StreamLocation_Location{
							Location: &payload.Object_Location{
								Name: "vald-agent-01",
								Uuid: "uuid-5",
								Ips:  []string{"127.0.0.1"},
							},
						},
					},
				},
			},
		},
		{
			name: "success to generate 0 object stream location",
			args: args{
				num:    0,
				name:   "vald-agent-01",
				ipAddr: "127.0.0.1",
			},
			want: want{
				want: []*payload.Object_StreamLocation{},
			},
		},
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

			got := GenObjectStreamLocation(test.args.num, test.args.name, test.args.ipAddr)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

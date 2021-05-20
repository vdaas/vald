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

// Package info provides general info metrics functions
package info

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		name        string
		fullname    string
		description string
		i           interface{}
	}
	type want struct {
		want metrics.Metric
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, metrics.Metric, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got metrics.Metric, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			type x struct {
				A string `info:"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
			}
			return test{
				name: "returns error when the passed struct's field is invalid",
				args: args{
					name:        "x/info",
					fullname:    "x info",
					description: "description",
					i: x{
						A: "a",
					},
				},
				want: want{
					want: nil,
					err:  errors.New("invalid key name: only ASCII characters accepted; max length must be 255 characters"),
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
		func() test {
			type x struct {
				A string `info:"data1"`
				B string `info:"data2"`
			}
			return test{
				name: "returns new metric when the passed struct is valid",
				args: args{
					name:        "x/info",
					fullname:    "x info",
					description: "description",
					i: x{
						A: "a",
						B: "b",
					},
				},
				want: want{},
				checkFunc: func(w want, got metrics.Metric, err error) error {
					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					if got == nil {
						return errors.New("got is nil")
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.name, test.args.fullname, test.args.description, test.args.i)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_labelKVs(t *testing.T) {
	t.Parallel()
	type args struct {
		i interface{}
	}
	type want struct {
		want map[metrics.Key]string
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, map[metrics.Key]string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got map[metrics.Key]string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			type x struct {
				A string `info:"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
			}
			return test{
				name: "returns error when the passed struct's field is too long",
				args: args{
					i: x{},
				},
				want: want{
					want: nil,
					err:  errors.New("invalid key name: only ASCII characters accepted; max length must be 255 characters"),
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
		func() test {
			type y struct {
				A string `info:"a"`
			}
			type x struct {
				A string   `info:"a"`
				B []string `info:"b"`
				C uint     `info:"c"`
				D uint8    `info:"d"`
				E uint16   `info:"e"`
				F uint32   `info:"f"`
				G uint64   `info:"g"`
				H int      `info:"h"`
				I int8     `info:"i"`
				J int16    `info:"j"`
				K int32    `info:"k"`
				L int64    `info:"l"`
				M float32  `info:"m"`
				N float64  `info:"n"`
				O bool     `info:"o"`
				P y        `info:"p"`
				Q string
				R []rune `info:"r"`
				S []y    `info:"s"`
			}
			return test{
				name: "returns kvs when the passed struct is valid",
				args: args{
					i: x{},
				},
				want: want{
					want: map[metrics.Key]string{
						metrics.MustNewKey("a"): "",
						metrics.MustNewKey("b"): "[]",
						metrics.MustNewKey("c"): "0",
						metrics.MustNewKey("d"): "0",
						metrics.MustNewKey("e"): "0",
						metrics.MustNewKey("f"): "0",
						metrics.MustNewKey("g"): "0",
						metrics.MustNewKey("h"): "0",
						metrics.MustNewKey("i"): "0",
						metrics.MustNewKey("j"): "0",
						metrics.MustNewKey("k"): "0",
						metrics.MustNewKey("l"): "0",
						metrics.MustNewKey("m"): "0E+00",
						metrics.MustNewKey("n"): "0E+00",
						metrics.MustNewKey("o"): "false",
						metrics.MustNewKey("r"): "",
					},
					err: nil,
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := labelKVs(test.args.i)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_info_Measurement(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		name     string
		fullname string
		info     metrics.Int64Measure
		kvs      map[metrics.Key]string
	}
	type want struct {
		want []metrics.Measurement
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []metrics.Measurement, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []metrics.Measurement, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "always returns empty measurement",
			args: args{
				ctx: nil,
			},
			fields: fields{
				name:     "",
				fullname: "",
				info:     *metrics.Int64(metrics.ValdOrg+"/test", "test", metrics.UnitDimensionless),
				kvs:      nil,
			},
			want: want{
				want: []metrics.Measurement{},
				err:  nil,
			},
			checkFunc: defaultCheckFunc,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			i := &info{
				name:     test.fields.name,
				fullname: test.fields.fullname,
				info:     test.fields.info,
				kvs:      test.fields.kvs,
			}

			got, err := i.Measurement(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_info_MeasurementWithTags(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	type fields struct {
		name     string
		fullname string
		info     metrics.Int64Measure
		kvs      map[metrics.Key]string
	}
	type want struct {
		want []metrics.MeasurementWithTags
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []metrics.MeasurementWithTags, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []metrics.MeasurementWithTags, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			m := metrics.Int64(metrics.ValdOrg+"/test", "test", metrics.UnitDimensionless)
			kvs := map[metrics.Key]string{
				metrics.MustNewKey("a"): "",
			}
			return test{
				name: "always returns const measurement",
				args: args{
					ctx: nil,
				},
				fields: fields{
					name:     "",
					fullname: "",
					info:     *m,
					kvs:      kvs,
				},
				want: want{
					want: []metrics.MeasurementWithTags{
						{
							Measurement: m.M(int64(1)),
							Tags:        kvs,
						},
					},
					err: nil,
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			i := &info{
				name:     test.fields.name,
				fullname: test.fields.fullname,
				info:     test.fields.info,
				kvs:      test.fields.kvs,
			}

			got, err := i.MeasurementWithTags(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_info_View(t *testing.T) {
	t.Parallel()
	type fields struct {
		name     string
		fullname string
		info     metrics.Int64Measure
		kvs      map[metrics.Key]string
	}
	type want struct {
		want []*metrics.View
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []*metrics.View) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []*metrics.View) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			m := metrics.Int64(metrics.ValdOrg+"/test", "test", metrics.UnitDimensionless)
			kvs := map[metrics.Key]string{
				metrics.MustNewKey("a"): "",
			}
			return test{
				name: "always returns view",
				fields: fields{
					name:     "name",
					fullname: "fullname",
					info:     *m,
					kvs:      kvs,
				},
				want: want{
					want: []*metrics.View{
						{
							Name:        "fullname",
							Description: m.Description(),
							TagKeys: []metrics.Key{
								metrics.MustNewKey("a"),
							},
							Measure:     m,
							Aggregation: metrics.LastValue(),
						},
					},
				},
				checkFunc: func(w want, got []*metrics.View) error {
					if !reflect.DeepEqual(got[0].Name, w.want[0].Name) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
					}

					if !reflect.DeepEqual(got[0].Description, w.want[0].Description) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
					}

					if !reflect.DeepEqual(got[0].TagKeys, w.want[0].TagKeys) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
					}

					if !reflect.DeepEqual(got[0].Measure, w.want[0].Measure) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
					}

					if !reflect.DeepEqual(got[0].Aggregation.Type.String(), w.want[0].Aggregation.Type.String()) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			i := &info{
				name:     test.fields.name,
				fullname: test.fields.fullname,
				info:     test.fields.info,
				kvs:      test.fields.kvs,
			}

			got := i.View()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

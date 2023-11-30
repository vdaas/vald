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

package kvs

import (
	"context"
	"math"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type want struct {
		want BidiMap
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, BidiMap) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got BidiMap) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			var (
				wantOu [slen]*sync.Map[uint32, valueStructOu]
				wantUo [slen]*sync.Map[string, ValueStructUo]
			)
			for i := 0; i < slen; i++ {
				wantOu[i] = new(sync.Map[uint32, valueStructOu])
				wantUo[i] = new(sync.Map[string, ValueStructUo])
			}
			return test{
				name: "return the bidi struct",
				want: want{
					want: &bidi{
						concurrency: runtime.GOMAXPROCS(-1) * 10,
						l:           0,
						ou:          wantOu,
						uo:          wantUo,
						eg:          errgroup.Get(),
					},
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := New(WithErrGroup(errgroup.Get()))
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_bidi_Get(t *testing.T) {
	t.Parallel()
	type args struct {
		key string
	}
	type fields struct {
		ou [slen]*sync.Map[uint32, valueStructOu]
		uo [slen]*sync.Map[string, ValueStructUo]
		l  uint64
	}
	type want struct {
		want    uint32
		want1   int64
		want2   bool
		wantLen uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint32, int64, bool, *bidi) error
		beforeFunc func(*testing.T, args, BidiMap)
		afterFunc  func(args, BidiMap)
	}
	defaultCheckFunc := func(w want, got uint32, got1 int64, got2 bool, bm *bidi) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		if !reflect.DeepEqual(got2, w.want2) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got2, w.want2)
		}
		if want, got := w.wantLen, atomic.LoadUint64(&bm.l); want != got {
			return errors.Errorf("l got: \"%#v\",\n\t\t\t\tl want: \"%#v\"", got, want)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return the value when there is a value for the key",
				args: args{
					key: key,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(a.key, val, ts)
				},
				want: want{
					want:    val,
					want1:   ts,
					want2:   true,
					wantLen: 1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 100,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return the value when there is a value for the key and l of fields is 100",
				args: args{
					key: key,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(a.key, val, ts)
				},
				want: want{
					want:    val,
					want1:   ts,
					want2:   true,
					wantLen: 101,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: math.MaxUint64,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return the value when there is a value for the key and l of fields is maximun value of uint64",
				args: args{
					key: key,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(a.key, val, ts)
				},
				want: want{
					want:    val,
					want1:   ts,
					want2:   true,
					wantLen: 0,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return (0, 0, false) when there is no value for the key",
				args: args{
					key: "84a333-59633fd4-4553-414a",
				},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					want:    0,
					want1:   0,
					want2:   false,
					wantLen: 1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name:   "return (0, 0, false) when there is no value for the key and the key is empty string",
				args:   args{},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					want:    0,
					want1:   0,
					want2:   false,
					wantLen: 1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = ""
				val uint32 = 0
				ts  int64  = 0
			)

			return test{
				name:   "return (0, 0, true) when the default value is set for the key and the key is empty string",
				args:   args{},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					want:    0,
					want1:   0,
					want2:   true,
					wantLen: 1,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, b)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args, b)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, got1, got2 := b.Get(test.args.key)
			if err := checkFunc(test.want, got, got1, got2, b); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_bidi_GetInverse(t *testing.T) {
	t.Parallel()
	type args struct {
		val uint32
	}
	type fields struct {
		ou [slen]*sync.Map[uint32, valueStructOu]
		uo [slen]*sync.Map[string, ValueStructUo]
		l  uint64
	}
	type want struct {
		want    string
		want1   int64
		want2   bool
		wantLen uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, string, int64, bool, *bidi) error
		beforeFunc func(*testing.T, args, BidiMap)
		afterFunc  func(args, BidiMap)
	}
	defaultCheckFunc := func(w want, got string, got1 int64, got2 bool, bm *bidi) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		if !reflect.DeepEqual(got2, w.want2) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got2, w.want2)
		}
		if want, got := w.wantLen, atomic.LoadUint64(&bm.l); want != got {
			return errors.Errorf("l got: \"%#v\",\n\t\t\t\tl want: \"%#v\"", got, want)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return key and timestamp and true when there is a key for the value",
				args: args{
					val: val,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, _ args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					want:    key,
					want1:   ts,
					want2:   true,
					wantLen: 1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 100,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return key and timestamp and true when there is a key for the value and l of fields is 100",
				args: args{
					val: val,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, _ args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					want:    key,
					want1:   ts,
					want2:   true,
					wantLen: 101,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: math.MaxUint64,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return key and timestamp and true when there is a key for the value and l of fields is maximun value of uint64",
				args: args{
					val: val,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, _ args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					want:    key,
					want1:   ts,
					want2:   true,
					wantLen: 0,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return false when there is a no key for the value",
				args: args{
					val: 10000,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, _ args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					want:    "",
					want1:   0,
					want2:   false,
					wantLen: 1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name:   "return false when there is a no key for the value and the val is 0",
				args:   args{},
				fields: fields,
				beforeFunc: func(t *testing.T, _ args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					want:    "",
					want1:   0,
					want2:   false,
					wantLen: 1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = ""
				val uint32 = 0
				ts  int64  = 0
			)

			return test{
				name:   "return (0, 0, true) when the default value is set for the key and the val is 0",
				args:   args{},
				fields: fields,
				beforeFunc: func(t *testing.T, _ args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					want:    "",
					want1:   0,
					want2:   true,
					wantLen: 1,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, b)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args, b)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, got1, got2 := b.GetInverse(test.args.val)
			if err := checkFunc(test.want, got, got1, got2, b); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_bidi_Set(t *testing.T) {
	t.Parallel()
	type args struct {
		key string
		val uint32
		ts  int64
	}
	type fields struct {
		ou [slen]*sync.Map[uint32, valueStructOu]
		uo [slen]*sync.Map[string, ValueStructUo]
		l  uint64
	}
	type want struct {
		key string
		val uint32
		ts  int64
		l   uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(w want, args args, b *bidi) error
		beforeFunc func(*testing.T, args, *bidi)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, args args, b *bidi) error {
		val, ts, ok := b.Get(args.key)
		if !ok {
			return errors.New("val not found")
		}
		key, tsi, ok := b.GetInverse(args.val)
		if !ok {
			return errors.New("key not found")
		}
		if val != w.val {
			return errors.Errorf("val is not equals.\twant: %v, but got: %v", w.val, val)
		}
		if key != w.key {
			return errors.Errorf("key is not equals.\twant: %v, but got: %v", w.key, key)
		}
		if ts != w.ts || tsi != w.ts {
			return errors.Errorf("timestamp is not equals.\twant: %v, but got: %v, %v", w.ts, ts, tsi)
		}
		if l := atomic.LoadUint64(&b.l); l != w.l {
			return errors.Errorf("l is not equals.\twant: %v, but got: %v", l, w.l)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "set success when the key is not empty string and val is not 0",
				args: args{
					key: key,
					val: val,
					ts:  ts,
				},
				fields: fields,
				want: want{
					key: key,
					val: val,
					ts:  ts,
					l:   1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 100,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "set success when the key is not empty string and val is not 0 and l of fields is 100",
				args: args{
					key: key,
					val: val,
					ts:  ts,
				},
				fields: fields,
				want: want{
					key: key,
					val: val,
					ts:  ts,
					l:   101,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: math.MaxUint64,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "set success when the key is not empty string and val is not 0 and l of fields is maximun value of uint64",
				args: args{
					key: key,
					val: val,
					ts:  ts,
				},
				fields: fields,
				want: want{
					key: key,
					val: val,
					ts:  ts,
					l:   0,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				oldVal uint32 = 10000
				oldTs  int64  = 20000
			)

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "set success when the key is already set and the same key is set twie",
				args: args{
					key: key,
					val: val,
					ts:  ts,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, b *bidi) {
					t.Helper()
					b.Set(a.key, oldVal, oldTs)
				},
				want: want{
					key: key,
					val: val,
					ts:  ts,
					l:   1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			key := "45637ec4-c85f-11ea-87d0"

			return test{
				name: "set success when the val is 0",
				args: args{
					key: key,
				},
				fields: fields,
				want: want{
					key: key,
					val: 0,
					ts:  0,
					l:   1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var val uint32 = 14438

			return test{
				name: "set success when the key is empty string",
				args: args{
					val: val,
				},
				fields: fields,
				want: want{
					val: val,
					ts:  0,
					l:   1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			return test{
				name:   "set success when the key and empty and the val is 0",
				args:   args{},
				fields: fields,
				want: want{
					val: 0,
					ts:  0,
					l:   1,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, b)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			b.Set(test.args.key, test.args.val, test.args.ts)
			if err := checkFunc(test.want, test.args, b); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_bidi_Delete(t *testing.T) {
	t.Parallel()
	type args struct {
		key string
	}
	type fields struct {
		ou [slen]*sync.Map[uint32, valueStructOu]
		uo [slen]*sync.Map[string, ValueStructUo]
		l  uint64
	}
	type want struct {
		wantVal uint32
		wantOk  bool
		wantLen uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, args, *bidi, uint32, bool) error
		beforeFunc func(*testing.T, args, BidiMap)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, _ args, _ *bidi, gotVal uint32, gotOk bool) error {
		if !reflect.DeepEqual(gotVal, w.wantVal) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVal, w.wantVal)
		}
		if !reflect.DeepEqual(gotOk, w.wantOk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return val and true when the delete successes",
				args: args{
					key: key,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(a.key, val, ts)
				},
				checkFunc: func(w want, a args, b *bidi, gotVal uint32, gotOk bool) error {
					if err := defaultCheckFunc(w, a, b, gotVal, gotOk); err != nil {
						return err
					}
					if want, got := w.wantLen, atomic.LoadUint64(&b.l); want != got {
						return errors.Errorf("l is not equals.\twant: %v, but got: %v", want, got)
					}
					if _, _, ok := b.Get(key); ok {
						return errors.New("the value for the key exists")
					}
					if _, _, ok := b.GetInverse(val); ok {
						return errors.New("the key for the val has not disappeared")
					}
					return nil
				},
				want: want{
					wantVal: val,
					wantOk:  true,
					wantLen: 0,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 100,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return val and true when the delete successes and l of fields is 100",
				args: args{
					key: key,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(a.key, val, ts)
				},
				checkFunc: func(w want, a args, b *bidi, gotVal uint32, gotOk bool) error {
					if err := defaultCheckFunc(w, a, b, gotVal, gotOk); err != nil {
						return err
					}
					if want, got := w.wantLen, atomic.LoadUint64(&b.l); want != got {
						return errors.Errorf("l is not equals.\twant: %v, but got: %v", want, got)
					}
					if _, _, ok := b.Get(key); ok {
						return errors.New("the value for the key exists")
					}
					if _, _, ok := b.GetInverse(val); ok {
						return errors.New("the key for the val has not disappeared")
					}
					return nil
				},
				want: want{
					wantVal: val,
					wantOk:  true,
					wantLen: 100,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: math.MaxUint64,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return val and true when the delete successes and l of fields is maximun value of uint64",
				args: args{
					key: key,
				},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(a.key, val, ts)
				},
				checkFunc: func(w want, a args, b *bidi, gotVal uint32, gotOk bool) error {
					if err := defaultCheckFunc(w, a, b, gotVal, gotOk); err != nil {
						return err
					}
					if want, got := w.wantLen, atomic.LoadUint64(&b.l); want != got {
						return errors.Errorf("l is not equals.\twant: %v, but got: %v", want, got)
					}
					if _, _, ok := b.Get(key); ok {
						return errors.New("the value for the key exists")
					}
					if _, _, ok := b.GetInverse(val); ok {
						return errors.New("the key for the val has not disappeared")
					}
					return nil
				},
				want: want{
					wantVal: val,
					wantOk:  true,
					wantLen: math.MaxUint64,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name:   "return val and true when the delete successes and the key is empty string",
				args:   args{},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(a.key, val, ts)
				},
				checkFunc: func(w want, a args, b *bidi, gotVal uint32, gotOk bool) error {
					if err := defaultCheckFunc(w, a, b, gotVal, gotOk); err != nil {
						return err
					}
					if want, got := w.wantLen, atomic.LoadUint64(&b.l); want != got {
						return errors.Errorf("l is not equals.\twant: %v, but got: %v", want, got)
					}
					if _, _, ok := b.Get(a.key); ok {
						return errors.New("the value for the key exists")
					}
					if _, _, ok := b.GetInverse(val); ok {
						return errors.New("the key for the val has not disappeared")
					}
					return nil
				},
				want: want{
					wantVal: val,
					wantOk:  true,
					wantLen: 0,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return (0, false) when the delete fails",
				args: args{
					key: "95314ec4-d95f-14ea-19d0",
				},
				fields: fields,
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				want: want{
					wantVal: 0,
					wantOk:  false,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, b)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotVal, gotOk := b.Delete(test.args.key)
			if err := checkFunc(test.want, test.args, b, gotVal, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_bidi_DeleteInverse(t *testing.T) {
	t.Parallel()
	type args struct {
		val uint32
	}
	type fields struct {
		ou [slen]*sync.Map[uint32, valueStructOu]
		uo [slen]*sync.Map[string, ValueStructUo]
		l  uint64
	}
	type want struct {
		wantKey string
		wantOk  bool
		wantLen uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, args, *bidi, string, bool) error
		beforeFunc func(*testing.T, args, BidiMap)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, _ args, _ *bidi, gotKey string, gotOk bool) error {
		if !reflect.DeepEqual(gotKey, w.wantKey) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotKey, w.wantKey)
		}
		if !reflect.DeepEqual(gotOk, w.wantOk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return key and true when the delete successes",
				args: args{
					val: val,
				},
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				checkFunc: func(w want, a args, b *bidi, gotKey string, gotOk bool) error {
					if err := defaultCheckFunc(w, a, b, gotKey, gotOk); err != nil {
						return err
					}
					if want, got := w.wantLen, atomic.LoadUint64(&b.l); want != got {
						return errors.Errorf("l is not equals.\twant: %v, but got: %v", want, got)
					}
					if _, _, ok := b.Get(key); ok {
						return errors.New("the value for the key exists")
					}
					if _, _, ok := b.GetInverse(val); ok {
						return errors.New("the key for the val has not disappeared")
					}
					return nil
				},
				fields: fields,
				want: want{
					wantKey: key,
					wantOk:  true,
					wantLen: 0,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 100,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return key and true when the delete successes and l of fields is 100",
				args: args{
					val: val,
				},
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				checkFunc: func(w want, a args, b *bidi, gotKey string, gotOk bool) error {
					if err := defaultCheckFunc(w, a, b, gotKey, gotOk); err != nil {
						return err
					}
					if want, got := w.wantLen, atomic.LoadUint64(&b.l); want != got {
						return errors.Errorf("l is not equals.\twant: %v, but got: %v", want, got)
					}
					if _, _, ok := b.Get(key); ok {
						return errors.New("the value for the key exists")
					}
					if _, _, ok := b.GetInverse(val); ok {
						return errors.New("the key for the val has not disappeared")
					}
					return nil
				},
				fields: fields,
				want: want{
					wantKey: key,
					wantOk:  true,
					wantLen: 100,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: math.MaxUint64,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return key and true when the delete successes and l of fields is maximun value of uint64",
				args: args{
					val: val,
				},
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				checkFunc: func(w want, a args, b *bidi, gotKey string, gotOk bool) error {
					if err := defaultCheckFunc(w, a, b, gotKey, gotOk); err != nil {
						return err
					}
					if want, got := w.wantLen, atomic.LoadUint64(&b.l); want != got {
						return errors.Errorf("l is not equals.\twant: %v, but got: %v", want, got)
					}
					if _, _, ok := b.Get(key); ok {
						return errors.New("the value for the key exists")
					}
					if _, _, ok := b.GetInverse(val); ok {
						return errors.New("the key for the val has not disappeared")
					}
					return nil
				},
				fields: fields,
				want: want{
					wantKey: key,
					wantOk:  true,
					wantLen: math.MaxUint64,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key       = "45637ec4-c85f-11ea-87d0"
				ts  int64 = 24438
			)

			return test{
				name: "return key and true when the delete successes and the val is 0",
				args: args{},
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(key, a.val, ts)
				},
				checkFunc: func(w want, a args, b *bidi, gotKey string, gotOk bool) error {
					if err := defaultCheckFunc(w, a, b, gotKey, gotOk); err != nil {
						return err
					}
					if want, got := w.wantLen, atomic.LoadUint64(&b.l); want != got {
						return errors.Errorf("l is not equals.\twant: %v, but got: %v", want, got)
					}
					if _, _, ok := b.Get(key); ok {
						return errors.New("the value for the key exists")
					}
					if _, _, ok := b.GetInverse(a.val); ok {
						return errors.New("the key for the val has not disappeared")
					}
					return nil
				},
				fields: fields,
				want: want{
					wantKey: key,
					wantOk:  true,
					wantLen: 0,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
				ts  int64  = 24438
			)

			return test{
				name: "return false when the delete fails",
				args: args{
					val: 10000,
				},
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					bm.Set(key, val, ts)
				},
				fields: fields,
				want: want{
					wantKey: "",
					wantOk:  false,
					wantLen: 0,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, b)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotKey, gotOk := b.DeleteInverse(test.args.val)
			if err := checkFunc(test.want, test.args, b, gotKey, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_bidi_Range(t *testing.T) {
	t.Parallel()
	type args struct {
		f func(string, uint32, int64) bool
	}
	type fields struct {
		ou [slen]*sync.Map[uint32, valueStructOu]
		uo [slen]*sync.Map[string, ValueStructUo]
		l  uint64
	}
	type want struct {
		want    map[string]ValueStructUo
		wantLen uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *bidi) error
		beforeFunc func(*testing.T, args, BidiMap)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, bm *bidi) error {
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				got     = make(map[string]ValueStructUo)
				wantMap = map[string]ValueStructUo{
					"1ec4-c85f-11ea-87d0": {10000, 20000},
					"2ec4-c85f-11ea-87d0": {10001, 20001},
					"3ec4-c85f-11ea-87d0": {10002, 20002},
					"4ec4-c85f-11ea-87d0": {10003, 20003},
				}
			)
			var mu sync.Mutex

			return test{
				name: "rage get successes",
				args: args{
					f: func(s string, u uint32, t int64) bool {
						mu.Lock()
						got[s] = ValueStructUo{u, t}
						mu.Unlock()
						return true
					},
				},
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					for key, val := range wantMap {
						bm.Set(key, val.value, val.timestamp)
					}
				},
				checkFunc: func(w want, bm *bidi) error {
					if !reflect.DeepEqual(got, w.want) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
					}
					if want, got := w.wantLen, atomic.LoadUint64(&bm.l); want != got {
						return errors.Errorf("l got: \"%d\",\n\t\t\t\tl want: \"%d\"", got, want)
					}
					return nil
				},
				fields: fields,
				want: want{
					want:    wantMap,
					wantLen: 4,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 100,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				got     = make(map[string]ValueStructUo)
				wantMap = map[string]ValueStructUo{
					"1ec4-c85f-11ea-87d0": {10000, 20000},
					"2ec4-c85f-11ea-87d0": {10001, 20001},
					"3ec4-c85f-11ea-87d0": {10002, 20002},
					"4ec4-c85f-11ea-87d0": {10003, 20003},
				}
			)
			var mu sync.Mutex

			return test{
				name: "rage get successes when l of fields is 100",
				args: args{
					f: func(s string, u uint32, t int64) bool {
						mu.Lock()
						got[s] = ValueStructUo{u, t}
						mu.Unlock()
						return true
					},
				},
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					for key, val := range wantMap {
						bm.Set(key, val.value, val.timestamp)
					}
				},
				checkFunc: func(w want, bm *bidi) error {
					if !reflect.DeepEqual(got, w.want) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
					}
					if want, got := w.wantLen, atomic.LoadUint64(&bm.l); want != got {
						return errors.Errorf("l got: \"%d\",\n\t\t\t\tl want: \"%d\"", got, want)
					}
					return nil
				},
				fields: fields,
				want: want{
					want:    wantMap,
					wantLen: 104,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: math.MaxUint64,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(sync.Map[uint32, valueStructOu])
				fields.uo[i] = new(sync.Map[string, ValueStructUo])
			}

			var (
				got     = make(map[string]ValueStructUo)
				wantMap = map[string]ValueStructUo{
					"1ec4-c85f-11ea-87d0": {10000, 20000},
					"2ec4-c85f-11ea-87d0": {10001, 20001},
					"3ec4-c85f-11ea-87d0": {10002, 20002},
					"4ec4-c85f-11ea-87d0": {10003, 20003},
				}
			)
			var mu sync.Mutex

			return test{
				name: "rage get successes when l of fields is maximun value of uint64",
				args: args{
					f: func(s string, u uint32, t int64) bool {
						mu.Lock()
						got[s] = ValueStructUo{u, t}
						mu.Unlock()
						return true
					},
				},
				beforeFunc: func(t *testing.T, a args, bm BidiMap) {
					t.Helper()
					for key, val := range wantMap {
						bm.Set(key, val.value, val.timestamp)
					}
				},
				checkFunc: func(w want, bm *bidi) error {
					if !reflect.DeepEqual(got, w.want) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
					}
					if want, got := w.wantLen, atomic.LoadUint64(&bm.l); want != got {
						return errors.Errorf("l got: \"%d\",\n\t\t\t\tl want: \"%d\"", got, want)
					}
					return nil
				},
				fields: fields,
				want: want{
					want:    wantMap,
					wantLen: 3,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			eg, egctx := errgroup.New(ctx)
			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
				eg: eg,
			}
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args, b)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			b.Range(egctx, test.args.f)
			if err := checkFunc(test.want, b); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_bidi_Len(t *testing.T) {
	t.Parallel()
	type fields struct {
		ou [slen]*sync.Map[uint32, valueStructOu]
		uo [slen]*sync.Map[string, ValueStructUo]
		l  uint64
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return 100 when l of field is 100",
			fields: fields{
				l: 100,
			},
			want: want{
				want: 100,
			},
		},
		{
			name: "return maximun value when l of field is maximun value of uint64",
			fields: fields{
				l: math.MaxUint64,
			},
			want: want{
				want: math.MaxUint64,
			},
		},
		{
			name:   "return 0 when l of field is default value",
			fields: fields{},
			want: want{
				want: 0,
			},
		},
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}

			got := b.Len()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_bidi_Close(t *testing.T) {
// 	type fields struct {
// 		concurrency int
// 		l           uint64
// 		ou          [slen]*sync.Map[uint32, valueStructOu]
// 		uo          [slen]*sync.Map[string, ValueStructUo]
// 		eg          errgroup.Group
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
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
// 		       fields: fields {
// 		           concurrency:0,
// 		           l:0,
// 		           ou:nil,
// 		           uo:nil,
// 		           eg:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           concurrency:0,
// 		           l:0,
// 		           ou:nil,
// 		           uo:nil,
// 		           eg:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			b := &bidi{
// 				concurrency: test.fields.concurrency,
// 				l:           test.fields.l,
// 				ou:          test.fields.ou,
// 				uo:          test.fields.uo,
// 				eg:          test.fields.eg,
// 			}
//
// 			err := b.Close()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }

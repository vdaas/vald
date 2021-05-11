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

package kvs

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
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
				wantOu [slen]*ou
				wantUo [slen]*uo
			)
			for i := 0; i < slen; i++ {
				wantOu[i] = new(ou)
				wantUo[i] = new(uo)
			}
			return test{
				name: "return BindiMap",
				want: want{
					want: &bidi{
						l:  0,
						ou: wantOu,
						uo: wantUo,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := New()
			if err := test.checkFunc(test.want, got); err != nil {
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
		ou [slen]*ou
		uo [slen]*uo
		l  uint64
	}
	type want struct {
		want  uint32
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint32, bool) error
		beforeFunc func(args, BidiMap)
		afterFunc  func(args, BidiMap)
	}
	defaultCheckFunc := func(w want, got uint32, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
			)

			return test{
				name: "return (14438, true) when there is a value for the key",
				args: args{
					key: key,
				},
				fields: fields,
				beforeFunc: func(a args, bm BidiMap) {
					bm.Set(a.key, val)
				},
				want: want{
					want:  val,
					want1: true,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				key1        = "45637ec4-c85f-11ea-87d0"
				val1 uint32 = 14438
			)
			var key2 = "84a333-59633fd4-4553-414a"

			return test{
				name: "return (0, false) when there is no value for the key",
				args: args{
					key: key2,
				},
				fields: fields,
				beforeFunc: func(_ args, bm BidiMap) {
					bm.Set(key1, val1)
				},
				want: want{
					want:  0,
					want1: false,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
			)

			return test{
				name:   "return (0, false) when there is no value for the key and the key is default value",
				args:   args{},
				fields: fields,
				beforeFunc: func(_ args, bm BidiMap) {
					bm.Set(key, val)
				},
				want: want{
					want:  0,
					want1: false,
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
				test.beforeFunc(test.args, b)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args, b)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, got1 := b.Get(test.args.key)
			if err := test.checkFunc(test.want, got, got1); err != nil {
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
		ou [slen]*ou
		uo [slen]*uo
		l  uint64
	}
	type want struct {
		want  string
		want1 bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, string, bool) error
		beforeFunc func(args, BidiMap)
		afterFunc  func(args, BidiMap)
	}
	defaultCheckFunc := func(w want, got string, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
			)

			return test{
				name: "return (45637ec4-c85f-11ea-87d0, true) when there is a key for the value",
				args: args{
					val: val,
				},
				fields: fields,
				beforeFunc: func(_ args, bm BidiMap) {
					bm.Set(key, val)
				},
				want: want{
					want:  key,
					want1: true,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				key         = "45637ec4-c85f-11ea-87d0"
				val1 uint32 = 14438
			)
			var val2 uint32 = 10000

			return test{
				name: "return false when there is a no key for the value",
				args: args{
					val: val2,
				},
				fields: fields,
				beforeFunc: func(_ args, bm BidiMap) {
					bm.Set(key, val1)
				},
				want: want{
					want:  "",
					want1: false,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
			)

			return test{
				name:   "return false when there is a no key for the value and the val is default value",
				args:   args{},
				fields: fields,
				beforeFunc: func(_ args, bm BidiMap) {
					bm.Set(key, val)
				},
				want: want{
					want:  "",
					want1: false,
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
				test.beforeFunc(test.args, b)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args, b)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, got1 := b.GetInverse(test.args.val)
			if err := test.checkFunc(test.want, got, got1); err != nil {
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
	}
	type fields struct {
		ou [slen]*ou
		uo [slen]*uo
		l  uint64
	}
	type want struct {
		key string
		val uint32
		l   uint64
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(w want, args args, b *bidi) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, args args, b *bidi) error {
		val, ok := b.Get(args.key)
		if !ok {
			return errors.New("uuid not found")
		}
		key, ok := b.GetInverse(args.val)
		if !ok {
			return errors.New("object id not found")
		}
		if val != w.val {
			return errors.Errorf("val is not equals. want: %v, but got: %v", w.val, val)
		}
		if key != w.key {
			return errors.Errorf("key is not equals. want: %v, but got: %v", w.key, key)
		}
		if l := atomic.LoadUint64(&b.l); l != w.l {
			return errors.Errorf("l is not equals. want: %v, but got: %v", l, w.l)
		}
		return nil
	}
	tests := []test{
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
			)

			return test{
				name: "set success",
				args: args{
					key: key,
					val: val,
				},
				fields: fields,
				want: want{
					key: key,
					val: val,
					l:   1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				key = "45637ec4-c85f-11ea-87d0"
			)

			return test{
				name: "set success when the val is default value",
				args: args{
					key: key,
				},
				fields: fields,
				want: want{
					key: key,
					val: 0,
					l:   1,
				},
			}
		}(),
		func() test {
			fields := fields{
				l: 0,
			}
			for i := 0; i < slen; i++ {
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				val uint32 = 14438
			)

			return test{
				name:   "set success when the key is default value",
				args:   args{},
				fields: fields,
				want: want{
					val: val,
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
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}

			b.Set(test.args.key, test.args.val)
			if err := test.checkFunc(test.want, test.args, b); err != nil {
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
		ou [slen]*ou
		uo [slen]*uo
		l  uint64
	}
	type want struct {
		wantVal uint32
		wantOk  bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, args, *bidi, uint32, bool) error
		beforeFunc func(args, BidiMap)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, a args, b *bidi, gotVal uint32, gotOk bool) error {
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
				fields.ou[i] = new(ou)
				fields.uo[i] = new(uo)
			}

			var (
				key        = "45637ec4-c85f-11ea-87d0"
				val uint32 = 14438
			)
			var wantl = 0

			return test{
				name: "return (14438, true) when the delete successes",
				args: args{
					key: key,
				},
				fields: fields,
				beforeFunc: func(a args, bm BidiMap) {
					bm.Set(a.key, val)
				},
				checkFunc: func(w want, a args, b *bidi, gotVal uint32, gotOk bool) error {
					if err := defaultCheckFunc(w, a, b, gotVal, gotOk); err != nil {
						return err
					}
					if l := atomic.LoadUint64(&b.l); wantl != 0 {
						return errors.Errorf("l is not equals. want: %v, but got: %v", wantl, l)
					}
					if _, ok := b.Get(key); ok {
						return errors.New("the value for the key exists")
					}
					if _, ok := b.GetInverse(val); ok {
						return errors.New("the key for the val has not disappeared")
					}
					return nil
				},
				want: want{
					wantVal: val,
					wantOk:  true,
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
				test.beforeFunc(test.args, b)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			gotVal, gotOk := b.Delete(test.args.key)
			if err := test.checkFunc(test.want, test.args, b, gotVal, gotOk); err != nil {
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
		ou [slen]*ou
		uo [slen]*uo
		l  uint64
	}
	type want struct {
		wantKey string
		wantOk  bool
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, string, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotKey string, gotOk bool) error {
		if !reflect.DeepEqual(gotKey, w.wantKey) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotKey, w.wantKey)
		}
		if !reflect.DeepEqual(gotOk, w.wantOk) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           val: 0,
		       },
		       fields: fields {
		           ou: nil,
		           uo: nil,
		           l: 0,
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
		           val: 0,
		           },
		           fields: fields {
		           ou: nil,
		           uo: nil,
		           l: 0,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}

			gotKey, gotOk := b.DeleteInverse(test.args.val)
			if err := test.checkFunc(test.want, gotKey, gotOk); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_bidi_Range(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		f   func(string, uint32) bool
	}
	type fields struct {
		ou [slen]*ou
		uo [slen]*uo
		l  uint64
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           f: nil,
		       },
		       fields: fields {
		           ou: nil,
		           uo: nil,
		           l: 0,
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
		           f: nil,
		           },
		           fields: fields {
		           ou: nil,
		           uo: nil,
		           l: 0,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}

			b.Range(test.args.ctx, test.args.f)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_bidi_Len(t *testing.T) {
	t.Parallel()
	type fields struct {
		ou [slen]*ou
		uo [slen]*uo
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			b := &bidi{
				ou: test.fields.ou,
				uo: test.fields.uo,
				l:  test.fields.l,
			}

			got := b.Len()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_stringToBytes(t *testing.T) {
	t.Parallel()
	type args struct {
		s string
	}
	type want struct {
		wantB []byte
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []byte) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotB []byte) error {
		if !reflect.DeepEqual(gotB, w.wantB) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotB, w.wantB)
		}
		return nil
	}
	tests := []test{
		func() test {
			s := "vdaas/vald"
			return test{
				name: "return bytes when s is vdaas/vald",
				args: args{
					s: s,
				},
				want: want{
					wantB: []byte(s),
				},
			}
		}(),
		func() test {
			return test{
				name: "return nil bytes when s is default value",
				args: args{},
				want: want{
					wantB: []byte(nil),
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

			gotB := stringToBytes(test.args.s)
			if err := test.checkFunc(test.want, gotB); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

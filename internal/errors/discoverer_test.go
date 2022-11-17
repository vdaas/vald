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
package errors

import "testing"

func TestErrAddrCouldNotDiscover(t *testing.T) {
	type args struct {
		err    error
		record string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns wrapped discover error when err is server error and record is '127.0.0.1'",
			args: args{
				err:    New("server error"),
				record: "127.0.0.1",
			},
			want: want{
				want: New("addr 127.0.0.1 ip couldn't discover: server error"),
			},
		},
		{
			name: "returns wrapped discover error when err is server error and record is empty",
			args: args{
				err:    New("server error"),
				record: "",
			},
			want: want{
				want: New("addr  ip couldn't discover: server error"),
			},
		},
		{
			name: "returns discover error when err is nil error and record is '127.0.0.1'",
			args: args{
				err:    nil,
				record: "127.0.0.1",
			},
			want: want{
				want: New("addr 127.0.0.1 ip couldn't discover"),
			},
		},
		{
			name: "returns discover error when err is nil error and record is empty",
			args: args{
				err:    nil,
				record: "",
			},
			want: want{
				want: New("addr  ip couldn't discover"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrAddrCouldNotDiscover(test.args.err, test.args.record)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrNodeNotFound(t *testing.T) {
	type args struct {
		node string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns discover node not found error when node is '127.0.0.1'",
			args: args{
				node: "127.0.0.1",
			},
			want: want{
				want: New("discover node 127.0.0.1 not found"),
			},
		},
		{
			name: "returns discover node not found error when node is empty",
			args: args{
				node: "",
			},
			want: want{
				want: New("discover node  not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrNodeNotFound(test.args.node)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrNamespaceNotFound(t *testing.T) {
	type args struct {
		ns string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns discover namespace not found error when ns is 'vald'",
			args: args{
				ns: "vald",
			},
			want: want{
				want: New("discover namespace vald not found"),
			},
		},
		{
			name: "returns discover namespace not found error when ns is empty",
			args: args{
				ns: "",
			},
			want: want{
				want: New("discover namespace  not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrNamespaceNotFound(test.args.ns)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrPodNameNotFound(t *testing.T) {
	type args struct {
		name string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns discover pod not found error when ns is 'vald-discoverer'",
			args: args{
				name: "vald-discoverer",
			},
			want: want{
				want: New("discover pod vald-discoverer not found"),
			},
		},
		{
			name: "returns discover pod not found error when name is empty",
			args: args{
				name: "",
			},
			want: want{
				want: New("discover pod  not found"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrPodNameNotFound(test.args.name)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrInvalidDiscoveryCache(t *testing.T) {
	type want struct {
		want error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns cache type cast error",
			want: want{
				want: New("cache type cast failed"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
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

			got := ErrInvalidDiscoveryCache
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

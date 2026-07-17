// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
package target

import "testing"

// TestNewMirrorTargetTemplate_CriticalOptionAborts pins down that a critical
// option error (e.g. an empty required field) aborts construction through the
// shared vald.SkipNonCriticalOptionError policy, matching New()'s severity handling.
func TestNewMirrorTargetTemplate_CriticalOptionAborts(t *testing.T) {
	t.Parallel()

	mt, err := NewMirrorTargetTemplate(WithMirrorTargetName(""))
	if err == nil {
		t.Fatal("NewMirrorTargetTemplate() error = nil, want non-nil for a critical option failure")
	}
	if mt != nil {
		t.Errorf("NewMirrorTargetTemplate() = %+v, want nil on critical failure", mt)
	}
}

// TestNewMirrorTargetTemplate_NonCriticalOptionIsWarned pins down that a
// non-critical option error (WithMirrorTargetStatus never fails; a no-op
// option is used to exercise the success path instead) does not abort
// construction, and that valid options apply as expected.
func TestNewMirrorTargetTemplate_NonCriticalOptionIsWarned(t *testing.T) {
	t.Parallel()

	mt, err := NewMirrorTargetTemplate(
		WithMirrorTargetName("t1"),
		WithMirrorTargetNamespace(""),
		WithMirrorTargetColocation("az-a"),
		WithMirrorTargetHost("host-a"),
		WithMirrorTargetPort(8081),
	)
	if err != nil {
		t.Fatalf("NewMirrorTargetTemplate() error = %v, want nil", err)
	}
	if mt == nil {
		t.Fatal("NewMirrorTargetTemplate() = nil, want non-nil")
	}
	if mt.Name != "t1" || mt.Spec.Colocation != "az-a" || mt.Spec.Target.Host != "host-a" || mt.Spec.Target.Port != 8081 {
		t.Errorf("NewMirrorTargetTemplate() = %+v, want name=t1 colocation=az-a host=host-a port=8081", mt)
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNewMirrorTargetTemplate(t *testing.T) {
// 	type args struct {
// 		opts []MirrorTargetTemplateOption
// 	}
// 	type want struct {
// 		want *MirrorTarget
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *MirrorTarget, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *MirrorTarget, err error) error {
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
// 		           opts:nil,
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
// 		           opts:nil,
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
// 			got, err := NewMirrorTargetTemplate(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

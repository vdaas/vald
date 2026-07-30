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
package v1

import (
	"testing"

	"github.com/vdaas/vald/internal/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	testJobAppLabelValue = "vald-benchmark"
	testJobMutatedValue  = "mutated"
)

func newValdBenchmarkJobFixture() *ValdBenchmarkJob {
	return &ValdBenchmarkJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job-fixture",
			Namespace: "default",
			Labels:    map[string]string{"app": testJobAppLabelValue},
		},
		Spec: BenchmarkJobSpec{
			JobType: "search",
			Target:  &BenchmarkTarget{Host: "vald-lb-gateway", Port: 8081},
			Dataset: &BenchmarkDataset{
				Name:  "fashion-mnist",
				Range: &config.BenchmarkDatasetRange{Start: 1, End: 1000},
			},
			Rules: []*config.BenchmarkJobRule{{Name: "rule-a"}},
		},
		Status: BenchmarkJobAvailable,
	}
}

func TestValdBenchmarkJob_DeepCopy(t *testing.T) {
	t.Parallel()

	orig := newValdBenchmarkJobFixture()
	cp := orig.DeepCopy()
	if cp == nil {
		t.Fatal("DeepCopy() = nil, want copy")
	}

	cp.Labels["app"] = testJobMutatedValue
	cp.Spec.Target.Host = testJobMutatedValue
	cp.Spec.Dataset.Range.End = 9
	cp.Spec.Rules[0].Name = testJobMutatedValue

	if orig.Labels["app"] != testJobAppLabelValue {
		t.Errorf("original labels mutated: %v", orig.Labels)
	}
	if orig.Spec.Target.Host != "vald-lb-gateway" {
		t.Errorf("original target mutated: %v", orig.Spec.Target)
	}
	if orig.Spec.Dataset.Range.End != 1000 {
		t.Errorf("original dataset range mutated: %v", orig.Spec.Dataset.Range)
	}
	if orig.Spec.Rules[0].Name != "rule-a" {
		t.Errorf("original rules mutated: %v", orig.Spec.Rules[0])
	}
}

func TestValdBenchmarkJob_DeepCopyObject(t *testing.T) {
	t.Parallel()

	orig := newValdBenchmarkJobFixture()
	obj := orig.DeepCopyObject()
	cp, ok := obj.(*ValdBenchmarkJob)
	if !ok {
		t.Fatalf("DeepCopyObject() = %T, want *ValdBenchmarkJob", obj)
	}
	if cp.GetName() != orig.GetName() {
		t.Errorf("copied name = %q, want %q", cp.GetName(), orig.GetName())
	}

	list := &ValdBenchmarkJobList{Items: []ValdBenchmarkJob{*orig}}
	lobj := list.DeepCopyObject()
	lcp, ok := lobj.(*ValdBenchmarkJobList)
	if !ok {
		t.Fatalf("DeepCopyObject() = %T, want *ValdBenchmarkJobList", lobj)
	}
	lcp.Items[0].Labels["app"] = testJobMutatedValue
	if orig.Labels["app"] != testJobAppLabelValue {
		t.Errorf("original mutated through list copy: %v", orig.Labels)
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestBenchmarkDataset_DeepCopyInto(t *testing.T) {
// 	type args struct {
// 		out *BenchmarkDataset
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		in         *BenchmarkDataset
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           out:nil,
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
// 		           out:nil,
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
// 			test.in.DeepCopyInto(test.args.out)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestBenchmarkJobSpec_DeepCopyInto(t *testing.T) {
// 	type args struct {
// 		out *BenchmarkJobSpec
// 	}
// 	type fields struct {
// 		ObjectConfig            *config.ObjectConfig
// 		ClientConfig            *config.GRPCClient
// 		Target                  *BenchmarkTarget
// 		Dataset                 *BenchmarkDataset
// 		UpdateConfig            *config.UpdateConfig
// 		GlobalConfig            *config.GlobalConfig
// 		RemoveConfig            *config.RemoveConfig
// 		InsertConfig            *config.InsertConfig
// 		ServerConfig            *config.Servers
// 		SearchConfig            *config.SearchConfig
// 		UpsertConfig            *config.UpsertConfig
// 		JobType                 string
// 		Rules                   []*config.BenchmarkJobRule
// 		Repetition              int
// 		Replica                 int
// 		RPS                     int
// 		ConcurrencyLimit        int
// 		TTLSecondsAfterFinished int
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           out:BenchmarkJobSpec{},
// 		       },
// 		       fields: fields {
// 		           ObjectConfig:nil,
// 		           ClientConfig:nil,
// 		           Target:nil,
// 		           Dataset:nil,
// 		           UpdateConfig:nil,
// 		           GlobalConfig:nil,
// 		           RemoveConfig:nil,
// 		           InsertConfig:nil,
// 		           ServerConfig:nil,
// 		           SearchConfig:nil,
// 		           UpsertConfig:nil,
// 		           JobType:"",
// 		           Rules:nil,
// 		           Repetition:0,
// 		           Replica:0,
// 		           RPS:0,
// 		           ConcurrencyLimit:0,
// 		           TTLSecondsAfterFinished:0,
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
// 		           out:BenchmarkJobSpec{},
// 		           },
// 		           fields: fields {
// 		           ObjectConfig:nil,
// 		           ClientConfig:nil,
// 		           Target:nil,
// 		           Dataset:nil,
// 		           UpdateConfig:nil,
// 		           GlobalConfig:nil,
// 		           RemoveConfig:nil,
// 		           InsertConfig:nil,
// 		           ServerConfig:nil,
// 		           SearchConfig:nil,
// 		           UpsertConfig:nil,
// 		           JobType:"",
// 		           Rules:nil,
// 		           Repetition:0,
// 		           Replica:0,
// 		           RPS:0,
// 		           ConcurrencyLimit:0,
// 		           TTLSecondsAfterFinished:0,
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
// 			in := &BenchmarkJobSpec{
// 				ObjectConfig:            test.fields.ObjectConfig,
// 				ClientConfig:            test.fields.ClientConfig,
// 				Target:                  test.fields.Target,
// 				Dataset:                 test.fields.Dataset,
// 				UpdateConfig:            test.fields.UpdateConfig,
// 				GlobalConfig:            test.fields.GlobalConfig,
// 				RemoveConfig:            test.fields.RemoveConfig,
// 				InsertConfig:            test.fields.InsertConfig,
// 				ServerConfig:            test.fields.ServerConfig,
// 				SearchConfig:            test.fields.SearchConfig,
// 				UpsertConfig:            test.fields.UpsertConfig,
// 				JobType:                 test.fields.JobType,
// 				Rules:                   test.fields.Rules,
// 				Repetition:              test.fields.Repetition,
// 				Replica:                 test.fields.Replica,
// 				RPS:                     test.fields.RPS,
// 				ConcurrencyLimit:        test.fields.ConcurrencyLimit,
// 				TTLSecondsAfterFinished: test.fields.TTLSecondsAfterFinished,
// 			}
//
// 			in.DeepCopyInto(test.args.out)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestBenchmarkJobSpec_DeepCopy(t *testing.T) {
// 	type fields struct {
// 		ObjectConfig            *config.ObjectConfig
// 		ClientConfig            *config.GRPCClient
// 		Target                  *BenchmarkTarget
// 		Dataset                 *BenchmarkDataset
// 		UpdateConfig            *config.UpdateConfig
// 		GlobalConfig            *config.GlobalConfig
// 		RemoveConfig            *config.RemoveConfig
// 		InsertConfig            *config.InsertConfig
// 		ServerConfig            *config.Servers
// 		SearchConfig            *config.SearchConfig
// 		UpsertConfig            *config.UpsertConfig
// 		JobType                 string
// 		Rules                   []*config.BenchmarkJobRule
// 		Repetition              int
// 		Replica                 int
// 		RPS                     int
// 		ConcurrencyLimit        int
// 		TTLSecondsAfterFinished int
// 	}
// 	type want struct {
// 		want *BenchmarkJobSpec
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *BenchmarkJobSpec) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *BenchmarkJobSpec) error {
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
// 		       fields: fields {
// 		           ObjectConfig:nil,
// 		           ClientConfig:nil,
// 		           Target:nil,
// 		           Dataset:nil,
// 		           UpdateConfig:nil,
// 		           GlobalConfig:nil,
// 		           RemoveConfig:nil,
// 		           InsertConfig:nil,
// 		           ServerConfig:nil,
// 		           SearchConfig:nil,
// 		           UpsertConfig:nil,
// 		           JobType:"",
// 		           Rules:nil,
// 		           Repetition:0,
// 		           Replica:0,
// 		           RPS:0,
// 		           ConcurrencyLimit:0,
// 		           TTLSecondsAfterFinished:0,
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
// 		           ObjectConfig:nil,
// 		           ClientConfig:nil,
// 		           Target:nil,
// 		           Dataset:nil,
// 		           UpdateConfig:nil,
// 		           GlobalConfig:nil,
// 		           RemoveConfig:nil,
// 		           InsertConfig:nil,
// 		           ServerConfig:nil,
// 		           SearchConfig:nil,
// 		           UpsertConfig:nil,
// 		           JobType:"",
// 		           Rules:nil,
// 		           Repetition:0,
// 		           Replica:0,
// 		           RPS:0,
// 		           ConcurrencyLimit:0,
// 		           TTLSecondsAfterFinished:0,
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
// 			in := &BenchmarkJobSpec{
// 				ObjectConfig:            test.fields.ObjectConfig,
// 				ClientConfig:            test.fields.ClientConfig,
// 				Target:                  test.fields.Target,
// 				Dataset:                 test.fields.Dataset,
// 				UpdateConfig:            test.fields.UpdateConfig,
// 				GlobalConfig:            test.fields.GlobalConfig,
// 				RemoveConfig:            test.fields.RemoveConfig,
// 				InsertConfig:            test.fields.InsertConfig,
// 				ServerConfig:            test.fields.ServerConfig,
// 				SearchConfig:            test.fields.SearchConfig,
// 				UpsertConfig:            test.fields.UpsertConfig,
// 				JobType:                 test.fields.JobType,
// 				Rules:                   test.fields.Rules,
// 				Repetition:              test.fields.Repetition,
// 				Replica:                 test.fields.Replica,
// 				RPS:                     test.fields.RPS,
// 				ConcurrencyLimit:        test.fields.ConcurrencyLimit,
// 				TTLSecondsAfterFinished: test.fields.TTLSecondsAfterFinished,
// 			}
//
// 			got := in.DeepCopy()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestValdBenchmarkJob_DeepCopyInto(t *testing.T) {
// 	type args struct {
// 		out *ValdBenchmarkJob
// 	}
// 	type fields struct {
// 		Base       resource.Base[ValdBenchmarkJob, *ValdBenchmarkJob]
// 		TypeMeta   metav1.TypeMeta
// 		Status     BenchmarkJobStatus
// 		ObjectMeta metav1.ObjectMeta
// 		Spec       BenchmarkJobSpec
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           out:ValdBenchmarkJob{},
// 		       },
// 		       fields: fields {
// 		           Base:nil,
// 		           TypeMeta:nil,
// 		           Status:nil,
// 		           ObjectMeta:nil,
// 		           Spec:BenchmarkJobSpec{},
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
// 		           out:ValdBenchmarkJob{},
// 		           },
// 		           fields: fields {
// 		           Base:nil,
// 		           TypeMeta:nil,
// 		           Status:nil,
// 		           ObjectMeta:nil,
// 		           Spec:BenchmarkJobSpec{},
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
// 			in := &ValdBenchmarkJob{
// 				Base:       test.fields.Base,
// 				TypeMeta:   test.fields.TypeMeta,
// 				Status:     test.fields.Status,
// 				ObjectMeta: test.fields.ObjectMeta,
// 				Spec:       test.fields.Spec,
// 			}
//
// 			in.DeepCopyInto(test.args.out)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

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

package job

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	corev1 "k8s.io/api/core/v1"
)

// TestNewBenchmarkJob_OptionErrorSeverity pins down the severity
// classification of the template options, aligned with the mirror target
// template: an empty required field (container name/image, operator config
// map) fails with errors.ErrCriticalOption and aborts construction, while any
// other invalid value is warned and skipped so the defaults survive.
func TestNewBenchmarkJob_OptionErrorSeverity(t *testing.T) {
	t.Parallel()

	t.Run("empty required fields abort construction", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			opt  BenchmarkTemplateOption
			name string
		}{
			{name: "container name", opt: WithContainerName("")},
			{name: "container image", opt: WithContainerImage("")},
			{name: "operator config map", opt: WithOperatorConfigMap("")},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got, err := NewBenchmarkJob(tc.opt)
				if err == nil {
					t.Fatal("NewBenchmarkJob() error = nil, want a critical option failure")
				}
				e := &errors.ErrCriticalOption{}
				if !errors.As(err, &e) {
					t.Errorf("NewBenchmarkJob() error = %v, want errors.ErrCriticalOption", err)
				}
				if got != nil {
					t.Errorf("NewBenchmarkJob() = %v, want nil on critical failure", got)
				}
			})
		}
	})

	t.Run("empty image pull policy is warned and keeps the default", func(t *testing.T) {
		t.Parallel()

		got, err := NewBenchmarkJob(WithImagePullPolicy(""))
		if err != nil {
			t.Fatalf("NewBenchmarkJob() error = %v, want nil (non-critical invalid options must be skipped)", err)
		}
		tpl, ok := got.(*benchmarkJobTemplate)
		if !ok {
			t.Fatalf("NewBenchmarkJob() = %T, want *benchmarkJobTemplate", got)
		}
		if tpl.imagePullPolicy != PullAlways {
			t.Errorf("imagePullPolicy = %q, want default %q", tpl.imagePullPolicy, PullAlways)
		}
	})
}

// TestCreateJobTpl_OptionErrorSeverity pins down the severity classification
// of the k8s Job construction options: an empty job name is a critical
// failure, every other invalid value is warned and skipped.
func TestCreateJobTpl_OptionErrorSeverity(t *testing.T) {
	t.Parallel()

	t.Run("empty job name aborts", func(t *testing.T) {
		t.Parallel()

		bj, err := NewBenchmarkJob()
		if err != nil {
			t.Fatalf("NewBenchmarkJob() error = %v, want nil", err)
		}
		if _, err := bj.CreateJobTpl(WithName("")); err == nil {
			t.Fatal("CreateJobTpl() error = nil, want a critical option failure")
		} else {
			e := &errors.ErrCriticalOption{}
			if !errors.As(err, &e) {
				t.Errorf("CreateJobTpl() error = %v, want errors.ErrCriticalOption", err)
			}
		}
	})

	t.Run("invalid non-critical options are warned and skipped", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			opt  BenchmarkOption
			name string
		}{
			{name: "empty namespace", opt: WithNamespace("")},
			{name: "empty service account name", opt: WithSvcAccountName("")},
			{name: "empty restart policy", opt: WithRestartPolicy("")},
			{name: "empty owner refs", opt: WithOwnerRef(nil)},
			{name: "empty label", opt: WithLabel(nil)},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				bj, err := NewBenchmarkJob()
				if err != nil {
					t.Fatalf("NewBenchmarkJob() error = %v, want nil", err)
				}
				job, err := bj.CreateJobTpl(tc.opt)
				if err != nil {
					t.Fatalf("CreateJobTpl() error = %v, want nil (non-critical invalid options must be skipped)", err)
				}
				if job.Namespace != "" {
					t.Errorf("Namespace = %q, want empty (invalid value must be skipped, not applied)", job.Namespace)
				}
				// The defaults must survive the skipped invalid option.
				if got := job.Spec.Template.Spec.RestartPolicy; got != corev1.RestartPolicyNever {
					t.Errorf("RestartPolicy = %q, want default %q", got, corev1.RestartPolicyNever)
				}
				if got := job.Spec.Template.Spec.ServiceAccountName; got != svcAccount {
					t.Errorf("ServiceAccountName = %q, want default %q", got, svcAccount)
				}
			})
		}
	})
}

// TestCreateJobTpl_ProductionCallPattern mirrors the exact option usage of
// pkg/operator/benchmark/service createJob — including the zero Repetition /
// Replica / TTL values an unset CRD spec produces — to guarantee the severity
// classification keeps the production path succeeding unchanged.
func TestCreateJobTpl_ProductionCallPattern(t *testing.T) {
	t.Parallel()

	const jobName = "sample-benchmark-job"

	bj, err := NewBenchmarkJob(
		WithContainerName(jobName),
		WithContainerImage("vdaas/vald-benchmark-job:latest"),
		WithImagePullPolicy(PullAlways),
		WithOperatorConfigMap("vald-benchmark-operator-config"),
	)
	if err != nil {
		t.Fatalf("NewBenchmarkJob() error = %v, want nil for the production option pattern", err)
	}
	job, err := bj.CreateJobTpl(
		WithName(jobName),
		WithNamespace("default"),
		WithLabel(map[string]string{"app": "vald-benchmark-operator", "benchmark-name": jobName + "1"}),
		WithCompletions(0),
		WithParallelism(0),
		WithOwnerRef([]k8s.OwnerReference{{
			APIVersion: "vald.vdaas.org/v1",
			Kind:       "ValdBenchmarkJob",
			Name:       jobName,
			UID:        "uid-1",
		}}),
		WithTTLSecondsAfterFinished(0),
	)
	if err != nil {
		t.Fatalf("CreateJobTpl() error = %v, want nil for the production option pattern", err)
	}
	if job.Name != jobName {
		t.Errorf("Name = %q, want %q", job.Name, jobName)
	}
	if job.Namespace != "default" {
		t.Errorf("Namespace = %q, want %q", job.Namespace, "default")
	}
	if len(job.OwnerReferences) != 1 || job.OwnerReferences[0].Name != jobName {
		t.Errorf("OwnerReferences = %+v, want a single ref for %q", job.OwnerReferences, jobName)
	}
	if job.Spec.Completions != nil {
		t.Errorf("Completions = %v, want nil (zero repetition keeps the k8s default)", *job.Spec.Completions)
	}
	if job.Spec.Parallelism != nil {
		t.Errorf("Parallelism = %v, want nil (zero replica keeps the k8s default)", *job.Spec.Parallelism)
	}
	if got := job.Spec.TTLSecondsAfterFinished; got == nil || *got != defaultTTLSeconds {
		t.Errorf("TTLSecondsAfterFinished = %v, want default %d", got, defaultTTLSeconds)
	}
	if got := job.Spec.Template.Spec.ServiceAccountName; got != svcAccount {
		t.Errorf("ServiceAccountName = %q, want %q", got, svcAccount)
	}
	if got := job.Spec.Template.Spec.RestartPolicy; got != corev1.RestartPolicyNever {
		t.Errorf("RestartPolicy = %q, want %q", got, corev1.RestartPolicyNever)
	}
	containers := job.Spec.Template.Spec.Containers
	if len(containers) != 1 {
		t.Fatalf("Containers = %+v, want exactly one container", containers)
	}
	if containers[0].Name != jobName {
		t.Errorf("container Name = %q, want %q", containers[0].Name, jobName)
	}
	if want := "vdaas/vald-benchmark-job:latest"; containers[0].Image != want {
		t.Errorf("container Image = %q, want %q", containers[0].Image, want)
	}
	if got := containers[0].ImagePullPolicy; got != corev1.PullPolicy(PullAlways) {
		t.Errorf("container ImagePullPolicy = %q, want %q", got, PullAlways)
	}
	volumes := job.Spec.Template.Spec.Volumes
	if len(volumes) != 1 || volumes[0].ConfigMap == nil {
		t.Fatalf("Volumes = %+v, want exactly one config map volume", volumes)
	}
	if want := "vald-benchmark-operator-config"; volumes[0].ConfigMap.Name != want {
		t.Errorf("config map volume Name = %q, want %q", volumes[0].ConfigMap.Name, want)
	}
}

// TestCreateJobTpl_FieldRefEnvVars pins down the exact downward-API EnvVar
// list (names, field paths and order) so the fieldRefEnvVar helper extraction
// cannot drift from the previous hand-written entries.
func TestCreateJobTpl_FieldRefEnvVars(t *testing.T) {
	t.Parallel()

	fieldRef := func(name, fieldPath string) corev1.EnvVar {
		return corev1.EnvVar{
			Name: name,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: fieldPath,
				},
			},
		}
	}
	want := []corev1.EnvVar{
		fieldRef("CRD_NAMESPACE", "metadata.namespace"),
		fieldRef("CRD_NAME", "metadata.labels['job-name']"),
		fieldRef("MY_NODE_NAME", "spec.nodeName"),
		fieldRef("MY_POD_NAMESPACE", "metadata.namespace"),
		fieldRef("MY_POD_NAME", "metadata.name"),
	}

	bj, err := NewBenchmarkJob()
	if err != nil {
		t.Fatalf("NewBenchmarkJob() error = %v, want nil", err)
	}
	job, err := bj.CreateJobTpl()
	if err != nil {
		t.Fatalf("CreateJobTpl() error = %v, want nil", err)
	}
	got := job.Spec.Template.Spec.Containers[0].Env
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Env = %+v, want %+v", got, want)
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_fieldRefEnvVar(t *testing.T) {
// 	type args struct {
// 		name      string
// 		fieldPath string
// 	}
// 	type want struct {
// 		want corev1.EnvVar
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, corev1.EnvVar) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got corev1.EnvVar) error {
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
// 		           name:"",
// 		           fieldPath:"",
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
// 		           name:"",
// 		           fieldPath:"",
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
// 			got := fieldRefEnvVar(test.args.name, test.args.fieldPath)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestNewBenchmarkJob(t *testing.T) {
// 	type args struct {
// 		opts []BenchmarkTemplateOption
// 	}
// 	type want struct {
// 		want BenchmarkTpl
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, BenchmarkTpl, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got BenchmarkTpl, err error) error {
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
// 			got, err := NewBenchmarkJob(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_benchmarkJobTemplate_CreateJobTpl(t *testing.T) {
// 	type args struct {
// 		opts []BenchmarkOption
// 	}
// 	type fields struct {
// 		jobTpl             k8s.Job
// 		containerName      string
// 		containerImageName string
// 		configMapName      string
// 		imagePullPolicy    ImagePullPolicy
// 	}
// 	type want struct {
// 		want k8s.Job
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, k8s.Job, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.Job, err error) error {
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
// 		       fields: fields {
// 		           jobTpl:nil,
// 		           containerName:"",
// 		           containerImageName:"",
// 		           configMapName:"",
// 		           imagePullPolicy:nil,
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
// 		           fields: fields {
// 		           jobTpl:nil,
// 		           containerName:"",
// 		           containerImageName:"",
// 		           configMapName:"",
// 		           imagePullPolicy:nil,
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
// 			b := &benchmarkJobTemplate{
// 				jobTpl:             test.fields.jobTpl,
// 				containerName:      test.fields.containerName,
// 				containerImageName: test.fields.containerImageName,
// 				configMapName:      test.fields.configMapName,
// 				imagePullPolicy:    test.fields.imagePullPolicy,
// 			}
//
// 			got, err := b.CreateJobTpl(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

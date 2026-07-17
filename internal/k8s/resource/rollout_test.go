//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package resource

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/client-go/util/retry"
)

// rolloutTestName is the workload name shared by the RolloutRestart
// characterization tests below; testNamespaceDefault is declared in
// syncer_test.go.
const rolloutTestName = "vald-agent"

// newRolloutDeployment returns a Deployment fixture whose pod template
// carries annotations, so RolloutRestart's stamping behavior can be observed
// against both a nil and a pre-populated annotation map.
func newRolloutDeployment(name string, annotations map[string]string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespaceDefault},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Annotations: annotations},
			},
		},
	}
}

// TestRolloutRestart_StampsAnnotation characterizes RolloutRestart's core
// contract: it stamps rolloutAnnotationKey with an RFC3339 timestamp onto the
// pod template, initializing a nil annotation map when necessary and leaving
// any pre-existing annotations untouched.
func TestRolloutRestart_StampsAnnotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		existing map[string]string
		name     string
	}{
		{
			name:     "nil annotation map is initialized",
			existing: nil,
		},
		{
			name:     "existing annotations are preserved",
			existing: map[string]string{"custom": "keep-me"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			deploy := newRolloutDeployment(rolloutTestName, tc.existing)
			dc := Deployment(&fakeClientSet{cs: clientfake.NewClientset(deploy)}, testNamespaceDefault)

			before := time.Now().UTC()
			if err := RolloutRestart(ctx, dc, rolloutTestName); err != nil {
				t.Fatalf("RolloutRestart() error = %v", err)
			}
			after := time.Now().UTC()

			got, err := dc.Get(ctx, rolloutTestName, EmptyGetOptions)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}

			ann := got.Spec.Template.Annotations
			if ann == nil {
				t.Fatal("Spec.Template.Annotations = nil, want an initialized map")
			}
			for k, want := range tc.existing {
				if got := ann[k]; got != want {
					t.Errorf("existing annotation %q = %q, want %q", k, got, want)
				}
			}

			stamp, ok := ann[rolloutAnnotationKey]
			if !ok {
				t.Fatalf("annotation %q missing, got %v", rolloutAnnotationKey, ann)
			}
			ts, err := time.Parse(time.RFC3339, stamp)
			if err != nil {
				t.Fatalf("annotation value %q is not RFC3339: %v", stamp, err)
			}
			if ts.Before(before.Add(-time.Second)) || ts.After(after.Add(time.Second)) {
				t.Errorf("stamped time %v outside expected window [%v, %v]", ts, before, after)
			}
		})
	}
}

// TestRolloutRestart_TargetNotFound characterizes the error path: a name
// with no backing workload surfaces the apierrors NotFound directly, since
// retry.RetryOnConflict only retries "Conflict" errors.
func TestRolloutRestart_TargetNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dc := Deployment(&fakeClientSet{cs: clientfake.NewClientset()}, testNamespaceDefault)

	err := RolloutRestart(ctx, dc, "missing")
	if err == nil {
		t.Fatal("RolloutRestart() error = nil, want a NotFound error")
	}
	if !apierrors.IsNotFound(err) {
		t.Errorf("RolloutRestart() error = %v, want apierrors NotFound", err)
	}
}

// conflictReactor returns a reactor that fails the first n update attempts
// with a Conflict error and lets every attempt after that fall through to
// the fake clientset's default handling.
func conflictReactor(n int) (k8stesting.ReactionFunc, *int) {
	attempts := 0
	return func(action k8stesting.Action) (bool, runtime.Object, error) {
		attempts++
		if attempts <= n {
			return true, nil, apierrors.NewConflict(
				schema.GroupResource{Group: "apps", Resource: "deployments"},
				rolloutTestName,
				errors.New("simulated update conflict"),
			)
		}
		return false, nil, nil
	}, &attempts
}

// TestRolloutRestart_RetriesOnConflict characterizes RolloutRestart's use of
// retry.RetryOnConflict: transient Conflict errors from the update call are
// retried in-place, with the pod template refetched on every attempt, until
// the update finally succeeds.
func TestRolloutRestart_RetriesOnConflict(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	deploy := newRolloutDeployment(rolloutTestName, nil)
	cs := clientfake.NewClientset(deploy)

	const conflictsBeforeSuccess = 2
	reactor, attempts := conflictReactor(conflictsBeforeSuccess)
	cs.PrependReactor("update", "deployments", reactor)

	dc := Deployment(&fakeClientSet{cs: cs}, testNamespaceDefault)
	if err := RolloutRestart(ctx, dc, rolloutTestName); err != nil {
		t.Fatalf("RolloutRestart() error = %v, want nil after retrying past transient conflicts", err)
	}
	if want := conflictsBeforeSuccess + 1; *attempts != want {
		t.Errorf("update attempts = %d, want %d (retries + final success)", *attempts, want)
	}

	got, err := dc.Get(ctx, rolloutTestName, EmptyGetOptions)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if _, ok := got.Spec.Template.Annotations[rolloutAnnotationKey]; !ok {
		t.Errorf("annotation %q missing after retry succeeded", rolloutAnnotationKey)
	}
}

// TestRolloutRestart_ConflictBudgetExhausted characterizes the give-up path:
// once retry.DefaultRetry's step budget is exhausted, RolloutRestart returns
// the last Conflict error instead of retrying indefinitely.
func TestRolloutRestart_ConflictBudgetExhausted(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	deploy := newRolloutDeployment(rolloutTestName, nil)
	cs := clientfake.NewClientset(deploy)

	// retry.DefaultRetry.Steps caps the total attempts; conflicting on every
	// call exercises that cap directly instead of an arbitrary large count.
	reactor, attempts := conflictReactor(retry.DefaultRetry.Steps)
	cs.PrependReactor("update", "deployments", reactor)

	dc := Deployment(&fakeClientSet{cs: cs}, testNamespaceDefault)
	err := RolloutRestart(ctx, dc, rolloutTestName)
	if err == nil {
		t.Fatal("RolloutRestart() error = nil, want a Conflict error after exhausting the retry budget")
	}
	if !apierrors.IsConflict(err) {
		t.Errorf("RolloutRestart() error = %v, want apierrors Conflict", err)
	}
	if want := retry.DefaultRetry.Steps; *attempts != want {
		t.Errorf("update attempts = %d, want retry.DefaultRetry.Steps = %d", *attempts, want)
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestRolloutRestart(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		client podAnnotationInterface[T]
// 		name   string
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
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
// 		           ctx:nil,
// 		           client:nil,
// 		           name:"",
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
// 		           client:nil,
// 		           name:"",
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
// 			err := RolloutRestart(test.args.ctx, test.args.client, test.args.name)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

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
	"strings"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

// fakeClientSet wires the client-go fake clientset into the client.ClientSet
// interface consumed by the resource client factories.
type fakeClientSet struct {
	cs kubernetes.Interface
}

func (f *fakeClientSet) GetClientSet() kubernetes.Interface {
	return f.cs
}

func (*fakeClientSet) GetRESTConfig() *rest.Config {
	return &rest.Config{}
}

// TestBaseClient_NilGuards verifies that calling methods on a nil or
// uninitialized baseClient returns ErrKubernetesClientNotFound instead of
// panicking: the previous nil check sat after mu.RLock and was dead code for
// nil receivers.
func TestBaseClient_NilGuards(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("nil receiver returns client-not-found", func(t *testing.T) {
		t.Parallel()

		var b *pod
		if _, err := b.Get(ctx, "name", EmptyGetOptions); !errors.Is(err, errors.ErrKubernetesClientNotFound) {
			t.Errorf("Get() error = %v, want %v", err, errors.ErrKubernetesClientNotFound)
		}
		if _, err := b.List(ctx, metav1.ListOptions{}); !errors.Is(err, errors.ErrKubernetesClientNotFound) {
			t.Errorf("List() error = %v, want %v", err, errors.ErrKubernetesClientNotFound)
		}
		if err := b.Delete(ctx, "name", EmptyDeleteOptions); !errors.Is(err, errors.ErrKubernetesClientNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, errors.ErrKubernetesClientNotFound)
		}
		if _, err := b.Watch(ctx, metav1.ListOptions{}); !errors.Is(err, errors.ErrKubernetesClientNotFound) {
			t.Errorf("Watch() error = %v, want %v", err, errors.ErrKubernetesClientNotFound)
		}
		if _, err := b.GetScale(ctx, "name", EmptyGetOptions); !errors.Is(err, errors.ErrKubernetesClientNotFound) {
			t.Errorf("GetScale() error = %v, want %v", err, errors.ErrKubernetesClientNotFound)
		}
	})

	t.Run("nil interface returns client-not-found", func(t *testing.T) {
		t.Parallel()

		b := new(pod)
		if _, err := b.Get(ctx, "name", EmptyGetOptions); !errors.Is(err, errors.ErrKubernetesClientNotFound) {
			t.Errorf("Get() error = %v, want %v", err, errors.ErrKubernetesClientNotFound)
		}
	})
}

// TestBaseClient_CreateJob_UniqueNames guards the job name generation against
// the previous second-precision timestamp: two jobs created from the same
// source within the same second must not collide.
func TestBaseClient_CreateJob_UniqueNames(t *testing.T) {
	t.Parallel()

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "vald-index-creation", Namespace: "default"},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "job", Image: "vdaas/vald-index-creation"}},
				},
			},
		},
	}
	dc := Deployment(&fakeClientSet{cs: clientfake.NewClientset(deploy)}, "default")

	ctx := context.Background()
	seen := make(map[string]bool, 3)
	for i := range 3 {
		job, err := dc.CreateJob(ctx, "vald-index-creation", EmptyGetOptions, EmptyCreateOptions)
		if err != nil {
			t.Fatalf("CreateJob() #%d error = %v", i, err)
		}
		name := job.GetName()
		if !strings.HasPrefix(name, "vald-index-creation-") {
			t.Errorf("CreateJob() #%d name = %q, want prefix %q", i, name, "vald-index-creation-")
		}
		if seen[name] {
			t.Errorf("CreateJob() #%d name %q collides with an earlier job", i, name)
		}
		seen[name] = true
	}
}

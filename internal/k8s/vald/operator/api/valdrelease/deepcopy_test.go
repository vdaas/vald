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

package valdrelease

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

// securityContexts returns fresh, non-nil pod/container security contexts.
func securityContexts() (*corev1.PodSecurityContext, *corev1.SecurityContext) {
	return &corev1.PodSecurityContext{RunAsUser: new(int64(1000))},
		&corev1.SecurityContext{RunAsNonRoot: new(true)}
}

// TestDeepCopyInto_SecurityContexts guards that DeepCopyInto does not share the
// corev1 security-context pointers (nor their nested pointers) with the source.
func TestDeepCopyInto_SecurityContexts(t *testing.T) {
	t.Parallel()
	type pair struct {
		pod *corev1.PodSecurityContext
		sc  *corev1.SecurityContext
	}
	tests := []struct {
		// copyFn builds a source with both contexts set, deep-copies it, and
		// returns the source and copy pointers.
		copyFn func() (src, dst pair)
		name   string
	}{
		{name: "*Agent", copyFn: func() (pair, pair) {
			in, out := new(Agent), new(Agent)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*Discoverer", copyFn: func() (pair, pair) {
			in, out := new(Discoverer), new(Discoverer)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*GatewayFilter", copyFn: func() (pair, pair) {
			in, out := new(GatewayFilter), new(GatewayFilter)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*GatewayLb", copyFn: func() (pair, pair) {
			in, out := new(GatewayLb), new(GatewayLb)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*GatewayMirror", copyFn: func() (pair, pair) {
			in, out := new(GatewayMirror), new(GatewayMirror)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*ManagerIndex", copyFn: func() (pair, pair) {
			in, out := new(ManagerIndex), new(ManagerIndex)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*ManagerIndexOperator", copyFn: func() (pair, pair) {
			in, out := new(ManagerIndexOperator), new(ManagerIndexOperator)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
		{name: "*ManagerIndexReadreplicaRotator", copyFn: func() (pair, pair) {
			in, out := new(ManagerIndexReadreplicaRotator), new(ManagerIndexReadreplicaRotator)
			in.PodSecurityContext, in.SecurityContext = securityContexts()
			in.DeepCopyInto(out)
			return pair{in.PodSecurityContext, in.SecurityContext}, pair{out.PodSecurityContext, out.SecurityContext}
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			src, dst := tc.copyFn()
			if dst.pod == nil || dst.sc == nil {
				t.Fatalf("copy dropped security contexts: pod=%v sc=%v", dst.pod, dst.sc)
			}
			if dst.pod == src.pod || dst.pod.RunAsUser == src.pod.RunAsUser {
				t.Error("PodSecurityContext is shared with the source")
			}
			if dst.sc == src.sc || dst.sc.RunAsNonRoot == src.sc.RunAsNonRoot {
				t.Error("SecurityContext is shared with the source")
			}
			if *dst.pod.RunAsUser != 1000 || !*dst.sc.RunAsNonRoot {
				t.Errorf("copied values differ: RunAsUser=%d RunAsNonRoot=%t", *dst.pod.RunAsUser, *dst.sc.RunAsNonRoot)
			}
		})
	}
}

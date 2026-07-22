//go:build e2e

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

// package crud provides end-to-end tests using ann-benchmarks datasets.
package crud

import (
	"context"
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

func (r *runner) processHTTP(t testing.TB, ctx context.Context, plan *config.Execution) {
	t.Helper()
	if plan == nil || plan.HTTP == nil {
		t.Fatal("http plan is nil")
		return
	}
	method := plan.HTTP.Method
	if method == "" {
		method = http.MethodGet
	}
	var body io.Reader
	if plan.HTTP.Body != "" {
		body = strings.NewReader(plan.HTTP.Body)
	}
	req, err := http.NewRequestWithContext(ctx, method, plan.HTTP.URL, body)
	if err != nil {
		t.Errorf("failed to create http request for %s: %v", plan.HTTP.URL, err)
		return
	}
	for key, val := range plan.HTTP.Headers {
		req.Header.Set(key, val)
	}
	if plan.HTTP.Auth != nil && plan.HTTP.Auth.ServiceAccount != nil {
		if r.k8s == nil {
			t.Error("kubernetes client is nil, cannot create service account token")
			return
		}
		sa := plan.HTTP.Auth.ServiceAccount
		token, err := resource.CreateServiceAccountToken(ctx, r.k8s, sa.Namespace, sa.Name, sa.ExpirationSeconds)
		if err != nil {
			t.Errorf("failed to create service account token for %s/%s: %v", sa.Namespace, sa.Name, err)
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := http.DefaultClient
	if plan.HTTP.TLS != nil && plan.HTTP.TLS.InsecureSkipVerify {
		tr, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			t.Error("http.DefaultTransport is not *http.Transport")
			return
		}
		tr = tr.Clone()
		tr.TLSClientConfig = &tls.Config{
			// #nosec G402 -- skipcq: GSC-G402 configured intentionally for e2e verification against self-signed endpoints
			InsecureSkipVerify: true,
		}
		client = &http.Client{Transport: tr}
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("failed to send http request to %s: %v", plan.HTTP.URL, err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body for %s: %v", plan.HTTP.URL, err)
		}
	}()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read response body from %s: %v", plan.HTTP.URL, err)
		return
	}
	log.Infof("http %s %s returned status: %s, body: %d bytes", method, plan.HTTP.URL, resp.Status, len(res))
	for _, expect := range plan.Expect {
		value, ok := expect.Value.(string)
		if !ok {
			t.Errorf("Expect.Value must be a string for http execution %s, got: %v", plan.Name, expect.Value)
			continue
		}
		switch expect.Op {
		case config.Contains:
			if !strings.Contains(string(res), value) {
				t.Errorf("response body from %s does not contain %q, body: %s", plan.HTTP.URL, value, string(res))
			}
		case config.NotContains:
			if strings.Contains(string(res), value) {
				t.Errorf("response body from %s contains %q, body: %s", plan.HTTP.URL, value, string(res))
			}
		default:
			t.Errorf("unsupported operator %s for http execution %s", expect.Op, plan.Name)
		}
	}
}

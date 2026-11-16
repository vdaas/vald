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
package json

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

func benchSearchResponse(n int) *payload.Search_Response {
	res := &payload.Search_Response{RequestId: "request-id", Results: make([]*payload.Object_Distance, n)}
	for i := range res.Results {
		res.Results[i] = &payload.Object_Distance{Id: "uuid-0000000001", Distance: float32(i) * 0.01}
	}
	return res
}

// BenchmarkEncodeResponse measures the REST response path for a top-10 search result.
func BenchmarkEncodeResponse(b *testing.B) {
	res := benchSearchResponse(10)
	b.ReportAllocs()
	for b.Loop() {
		w := httptest.NewRecorder()
		if err := EncodeResponse(w, res, http.StatusOK, "application/json"); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDecodeRequest measures the REST request path for a 768-dimensional search request.
func BenchmarkDecodeRequest(b *testing.B) {
	vec := make([]float32, 768)
	for i := range vec {
		vec[i] = float32(i) * 0.001
	}
	var buf bytes.Buffer
	w := httptest.NewRecorder()
	if err := EncodeResponse(w, &payload.Search_Request{Vector: vec, Config: &payload.Search_Config{Num: 10}}, http.StatusOK); err != nil {
		b.Fatal(err)
	}
	buf.Write(w.Body.Bytes())
	body := buf.Bytes()
	b.ReportAllocs()
	b.SetBytes(int64(len(body)))
	for b.Loop() {
		r := httptest.NewRequest(http.MethodPost, "/search", io.NopCloser(bytes.NewReader(body)))
		r.ContentLength = int64(len(body))
		var req payload.Search_Request
		if err := DecodeRequest(r, &req); err != nil {
			b.Fatal(err)
		}
		if len(req.GetVector()) != len(vec) {
			b.Fatalf("decoded %d dims, want %d", len(req.GetVector()), len(vec))
		}
	}
}

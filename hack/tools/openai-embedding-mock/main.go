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

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type embeddingRequest struct {
	Input any    `json:"input"`
	Model string `json:"model"`
}

type embeddingObject struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

type embeddingResponse struct {
	Object string            `json:"object"`
	Data   []embeddingObject `json:"data"`
	Model  string            `json:"model"`
}

func main() {
	addr := os.Getenv("MOCK_OPENAI_ADDR")
	if addr == "" {
		addr = "127.0.0.1:18000"
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	http.HandleFunc("/v1/embeddings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req embeddingRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		model := req.Model
		if model == "" {
			model = "text-embedding-3-small"
		}

		res := embeddingResponse{
			Object: "list",
			Data: []embeddingObject{
				{
					Object:    "embedding",
					Index:     0,
					Embedding: []float32{1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			Model: model,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Printf("openai embedding mock listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

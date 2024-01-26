// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package rest

import (
	"net/http"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/json"
)

// Handler represents an interface for rest handler.
type Handler interface {
	Register(w http.ResponseWriter, r *http.Request) (int, error)
	Index(w http.ResponseWriter, r *http.Request) (int, error)
	Exists(w http.ResponseWriter, r *http.Request) (int, error)
	Search(w http.ResponseWriter, r *http.Request) (int, error)
	SearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	MultiSearch(w http.ResponseWriter, r *http.Request) (int, error)
	MultiSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearch(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	MultiLinearSearch(w http.ResponseWriter, r *http.Request) (int, error)
	MultiLinearSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	Insert(w http.ResponseWriter, r *http.Request) (int, error)
	MultiInsert(w http.ResponseWriter, r *http.Request) (int, error)
	Update(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpdate(w http.ResponseWriter, r *http.Request) (int, error)
	Upsert(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpsert(w http.ResponseWriter, r *http.Request) (int, error)
	Remove(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveByTimestamp(w http.ResponseWriter, r *http.Request) (int, error)
	MultiRemove(w http.ResponseWriter, r *http.Request) (int, error)
	GetObject(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	vald vald.ServerWithMirror
}

// New returns a Vald server as rest handler with mirror using the provided options.
func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(h)
	}
	return h
}

// Register is an HTTP handler function that processes registration requests.
// It decodes the incoming JSON payload into a payload.Mirror_Targets struct,
// then invokes the vald.Register method to handle the registration logic.
// The response is written to the http.ResponseWriter.
func (h *handler) Register(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Mirror_Targets
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Register(r.Context(), req)
	})
}

// Index is an HTTP handler function that handles requests to the index endpoint.
// It returns an HTTP status code and an error. It creates a map to store data,
// then uses json.Handler to process the request, extract data, and log the request using dump.Request.
func (*handler) Index(w http.ResponseWriter, r *http.Request) (int, error) {
	data := make(map[string]interface{})
	return json.Handler(w, r, &data, func() (interface{}, error) {
		return dump.Request(nil, data, r)
	})
}

// Search is an HTTP handler function that processes search requests.
// It decodes the incoming JSON payload into a payload.Search_Request struct,
// then invokes the vald.Search method to handle the search logic.
func (h *handler) Search(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Search(r.Context(), req)
	})
}

// SearchByID is an HTTP handler function that processes search by ID requests.
// It decodes the incoming JSON payload into a payload.Search_IDRequest struct,
// then invokes the vald.SearchByID method to handle the search by ID logic.
func (h *handler) SearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_IDRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.SearchByID(r.Context(), req)
	})
}

// MultiSearch is an HTTP handler function that processes multi-search requests.
// It decodes the incoming JSON payload into a payload.Search_MultiRequest struct,
// then invokes the vald.MultiSearch method to handle the multi-search logic.
func (h *handler) MultiSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiSearch(r.Context(), req)
	})
}

// MultiSearchByID is an HTTP handler function that processes multi-search by ID requests.
// It decodes the incoming JSON payload into a payload.Search_MultiIDRequest struct,
// then invokes the vald.MultiSearchByID method to handle the multi-search by ID logic.
func (h *handler) MultiSearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiIDRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiSearchByID(r.Context(), req)
	})
}

// LinearSearch is an HTTP handler function that processes linear search requests.
// It decodes the incoming JSON payload into a payload.Search_Request struct,
// then invokes the vald.LinearSearch method to handle the linear search logic.
func (h *handler) LinearSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.LinearSearch(r.Context(), req)
	})
}

// LinearSearchByID is an HTTP handler function that processes linear search by ID requests.
// It decodes the incoming JSON payload into a payload.Search_IDRequest struct,
// then invokes the vald.LinearSearchByID method to handle the linear search by ID logic.
func (h *handler) LinearSearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_IDRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.LinearSearchByID(r.Context(), req)
	})
}

// MultiLinearSearch is an HTTP handler function that processes multi-linear search requests.
// It decodes the incoming JSON payload into a payload.Search_MultiRequest struct,
// then invokes the vald.MultiLinearSearch method to handle the multi-linear search logic.
func (h *handler) MultiLinearSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiLinearSearch(r.Context(), req)
	})
}

// MultiLinearSearchByID is an HTTP handler function that processes multi-linear search by ID requests.
// It decodes the incoming JSON payload into a payload.Search_MultiIDRequest struct,
// then invokes the vald.MultiLinearSearchByID method to handle the multi-linear search by ID logic.
func (h *handler) MultiLinearSearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiIDRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiLinearSearchByID(r.Context(), req)
	})
}

// Insert is an HTTP handler function that processes insert requests.
// It decodes the incoming JSON payload into a payload.Insert_Request struct,
// then invokes the vald.Insert method to handle the insert logic.
func (h *handler) Insert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Insert_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Insert(r.Context(), req)
	})
}

// MultiInsert is an HTTP handler function that processes multi-insert requests.
// It decodes the incoming JSON payload into a payload.Insert_MultiRequest struct,
// then invokes the vald.MultiInsert method to handle the multi-insert logic.
func (h *handler) MultiInsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Insert_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiInsert(r.Context(), req)
	})
}

// Update is an HTTP handler function that processes update requests.
// It decodes the incoming JSON payload into a payload.Update_Request struct,
// then invokes the vald.Update method to handle the update logic.
func (h *handler) Update(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Update_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Update(r.Context(), req)
	})
}

// MultiUpdate is an HTTP handler function that processes multi-update requests.
// It decodes the incoming JSON payload into a payload.Update_MultiRequest struct,
// then invokes the vald.MultiUpdate method to handle the multi-update logic.
func (h *handler) MultiUpdate(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Update_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiUpdate(r.Context(), req)
	})
}

// Upsert is an HTTP handler function that processes upsert requests.
// It decodes the incoming JSON payload into a payload.Upsert_Request struct,
// then invokes the vald.Upsert method to handle the upsert logic.
func (h *handler) Upsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Upsert_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Upsert(r.Context(), req)
	})
}

// MultiUpsert is an HTTP handler function that processes multi-upsert requests.
// It decodes the incoming JSON payload into a payload.Upsert_MultiRequest struct,
// then invokes the vald.MultiUpsert method to handle the multi-upsert logic.
func (h *handler) MultiUpsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Upsert_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiUpsert(r.Context(), req)
	})
}

// Remove is an HTTP handler function that processes remove requests.
// It decodes the incoming JSON payload into a payload.Remove_Request struct,
// then invokes the vald.Remove method to handle the remove logic.
func (h *handler) Remove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Remove_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Remove(r.Context(), req)
	})
}

// RemoveByTimestamp is an HTTP handler function that processes remove-by-timestamp requests.
// It decodes the incoming JSON payload into a payload.Remove_TimestampRequest struct,
// then invokes the vald.RemoveByTimestamp method to handle the remove-by-timestamp logic.
func (h *handler) RemoveByTimestamp(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Remove_TimestampRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.RemoveByTimestamp(r.Context(), req)
	})
}

// MultiRemove is an HTTP handler function that processes multi-remove requests.
// It decodes the incoming JSON payload into a payload.Remove_MultiRequest struct,
// then invokes the vald.MultiRemove method to handle the multi-remove logic.
func (h *handler) MultiRemove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Remove_MultiRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.MultiRemove(r.Context(), req)
	})
}

// GetObject is an HTTP handler function that processes get-object requests.
// It decodes the incoming JSON payload into a payload.Object_VectorRequest struct,
// then invokes the vald.GetObject method to handle the get-object logic.
func (h *handler) GetObject(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_VectorRequest
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.GetObject(r.Context(), req)
	})
}

// Exists is an HTTP handler function that processes exists requests.
// It decodes the incoming JSON payload into a payload.Object_ID struct,
// then invokes the vald.Exists method to handle the exists logic.
func (h *handler) Exists(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_ID
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.vald.Exists(r.Context(), req)
	})
}

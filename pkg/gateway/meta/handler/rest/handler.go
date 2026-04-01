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

package rest

import (
	"context"
	"net/http"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/json"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request) (int, error)
	Exists(w http.ResponseWriter, r *http.Request) (int, error)
	SearchWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	Search(w http.ResponseWriter, r *http.Request) (int, error)
	SearchByIDWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	SearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	MultiSearchWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	MultiSearch(w http.ResponseWriter, r *http.Request) (int, error)
	MultiSearchByIDWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	MultiSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearchWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearch(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearchByIDWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	LinearSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	MultiLinearSearchWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	MultiLinearSearch(w http.ResponseWriter, r *http.Request) (int, error)
	MultiLinearSearchByIDWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	MultiLinearSearchByID(w http.ResponseWriter, r *http.Request) (int, error)
	InsertWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	Insert(w http.ResponseWriter, r *http.Request) (int, error)
	MultiInsertWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	MultiInsert(w http.ResponseWriter, r *http.Request) (int, error)
	UpdateWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	Update(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpdateWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpdate(w http.ResponseWriter, r *http.Request) (int, error)
	UpdateTimestampWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	UpsertWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	Upsert(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpsertWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	MultiUpsert(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	Remove(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveByTimestampWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	MultiRemoveWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	MultiRemove(w http.ResponseWriter, r *http.Request) (int, error)
	Flush(w http.ResponseWriter, r *http.Request) (int, error)
	GetObjectWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
	GetObject(w http.ResponseWriter, r *http.Request) (int, error)
	StreamListObjectWithMetadata(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	vald vald.ServerWithMetadata
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(h)
	}
	return h
}

func (*handler) Index(w http.ResponseWriter, r *http.Request) (int, error) {
	data := make(map[string]any)
	return json.Handler(w, r, &data, func() (any, error) {
		return dump.Request(nil, data, r)
	})
}

func (h *handler) SearchWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.SearchWithMetadata(r.Context(), req)
	})
}

func (h *handler) Search(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.SearchWithMetadata(w, r)
}

func (h *handler) SearchByIDWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_IDRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.SearchByIDWithMetadata(r.Context(), req)
	})
}

func (h *handler) SearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.SearchByIDWithMetadata(w, r)
}

func (h *handler) MultiSearchWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.MultiSearchWithMetadata(r.Context(), req)
	})
}

func (h *handler) MultiSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.MultiSearchWithMetadata(w, r)
}

func (h *handler) MultiSearchByIDWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiIDRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.MultiSearchByIDWithMetadata(r.Context(), req)
	})
}

func (h *handler) MultiSearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.MultiSearchByIDWithMetadata(w, r)
}

func (h *handler) LinearSearchWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.LinearSearchWithMetadata(r.Context(), req)
	})
}

func (h *handler) LinearSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.LinearSearchWithMetadata(w, r)
}

func (h *handler) LinearSearchByIDWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_IDRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.LinearSearchByIDWithMetadata(r.Context(), req)
	})
}

func (h *handler) LinearSearchByID(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.LinearSearchByIDWithMetadata(w, r)
}

func (h *handler) MultiLinearSearchWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Search_MultiRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.MultiLinearSearchWithMetadata(r.Context(), req)
	})
}

func (h *handler) MultiLinearSearch(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.MultiLinearSearchWithMetadata(w, r)
}

func (h *handler) MultiLinearSearchByIDWithMetadata(
	w http.ResponseWriter, r *http.Request,
) (code int, err error) {
	var req *payload.Search_MultiIDRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.MultiLinearSearchByIDWithMetadata(r.Context(), req)
	})
}

func (h *handler) MultiLinearSearchByID(
	w http.ResponseWriter, r *http.Request,
) (code int, err error) {
	return h.MultiLinearSearchByIDWithMetadata(w, r)
}

func (h *handler) InsertWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Insert_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.InsertWithMetadata(r.Context(), req)
	})
}

func (h *handler) Insert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.InsertWithMetadata(w, r)
}

func (h *handler) MultiInsertWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Insert_MultiRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.MultiInsertWithMetadata(r.Context(), req)
	})
}

func (h *handler) MultiInsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.MultiInsertWithMetadata(w, r)
}

func (h *handler) UpdateWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Update_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.UpdateWithMetadata(r.Context(), req)
	})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.UpdateWithMetadata(w, r)
}

func (h *handler) MultiUpdateWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Update_MultiRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.MultiUpdateWithMetadata(r.Context(), req)
	})
}

func (h *handler) MultiUpdate(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.MultiUpdateWithMetadata(w, r)
}

func (h *handler) UpdateTimestampWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Update_TimestampRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.UpdateTimestampWithMetadata(r.Context(), req)
	})
}

func (h *handler) UpsertWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Upsert_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.UpsertWithMetadata(r.Context(), req)
	})
}

func (h *handler) Upsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.UpsertWithMetadata(w, r)
}

func (h *handler) MultiUpsertWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Upsert_MultiRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.MultiUpsertWithMetadata(r.Context(), req)
	})
}

func (h *handler) MultiUpsert(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.MultiUpsertWithMetadata(w, r)
}

func (h *handler) RemoveWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Remove_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.RemoveWithMetadata(r.Context(), req)
	})
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.RemoveWithMetadata(w, r)
}

func (h *handler) RemoveByTimestampWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Remove_TimestampRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.RemoveByTimestampWithMetadata(r.Context(), req)
	})
}

func (h *handler) MultiRemoveWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Remove_MultiRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.MultiRemoveWithMetadata(r.Context(), req)
	})
}

func (h *handler) MultiRemove(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.MultiRemoveWithMetadata(w, r)
}

func (h *handler) Flush(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Flush_Request
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.Flush(r.Context(), req)
	})
}

func (h *handler) GetObjectWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_VectorRequest
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.GetObjectWithMetadata(r.Context(), req)
	})
}

func (h *handler) GetObject(w http.ResponseWriter, r *http.Request) (code int, err error) {
	return h.GetObjectWithMetadata(w, r)
}

func (h *handler) StreamListObjectWithMetadata(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_List_Request
	return json.Handler(w, r, &req, func() (any, error) {
		srv := &listObjectWithMetadataHTTPServer{
			ctx:       r.Context(),
			responses: make([]*payload.Object_List_Response, 0, 8),
		}
		if err := h.vald.StreamListObjectWithMetadata(req, srv); err != nil {
			return nil, err
		}
		return srv.responses, nil
	})
}

func (h *handler) Exists(w http.ResponseWriter, r *http.Request) (code int, err error) {
	var req *payload.Object_ID
	return json.Handler(w, r, &req, func() (any, error) {
		return h.vald.Exists(r.Context(), req)
	})
}

type listObjectWithMetadataHTTPServer struct {
	ctx       context.Context
	header    metadata.MD
	trailer   metadata.MD
	responses []*payload.Object_List_Response
}

func (s *listObjectWithMetadataHTTPServer) SetHeader(md metadata.MD) error {
	s.header = md
	return nil
}

func (s *listObjectWithMetadataHTTPServer) SendHeader(md metadata.MD) error {
	s.header = md
	return nil
}

func (s *listObjectWithMetadataHTTPServer) SetTrailer(md metadata.MD) {
	s.trailer = md
}

func (s *listObjectWithMetadataHTTPServer) Context() context.Context {
	return s.ctx
}

func (s *listObjectWithMetadataHTTPServer) Send(res *payload.Object_List_Response) error {
	s.responses = append(s.responses, res)
	return nil
}

func (s *listObjectWithMetadataHTTPServer) SendMsg(m any) error {
	res, ok := m.(*payload.Object_List_Response)
	if !ok {
		return nil
	}
	return s.Send(res)
}

func (*listObjectWithMetadataHTTPServer) RecvMsg(any) error {
	return nil
}

var _ vald.ObjectWithMetadata_StreamListObjectWithMetadataServer = (*listObjectWithMetadataHTTPServer)(nil)
var _ ggrpc.ServerStream = (*listObjectWithMetadataHTTPServer)(nil)

//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package rest provides rest api logic
package rest

import (
	"net/http"

	"github.com/vdaas/vald/apis/grpc/v1/manager/backup"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/http/json"
)

type Handler interface {
	GetVector(w http.ResponseWriter, r *http.Request) (int, error)
	Locations(w http.ResponseWriter, r *http.Request) (int, error)
	Register(w http.ResponseWriter, r *http.Request) (int, error)
	RegisterMulti(w http.ResponseWriter, r *http.Request) (int, error)
	Remove(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveMulti(w http.ResponseWriter, r *http.Request) (int, error)
	RegisterIPs(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveIPs(w http.ResponseWriter, r *http.Request) (int, error)
}

type handler struct {
	backup backup.BackupServer
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOptions, opts...) {
		opt(h)
	}
	return h
}

func (h *handler) GetVector(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Backup_GetVector_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.backup.GetVector(r.Context(), req)
	})
}

func (h *handler) Locations(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Backup_Locations_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.backup.Locations(r.Context(), req)
	})
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Backup_Compressed_Vector
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.backup.Register(r.Context(), req)
	})
}

func (h *handler) RegisterMulti(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Backup_Compressed_Vectors
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.backup.RegisterMulti(r.Context(), req)
	})
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Backup_Remove_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.backup.Remove(r.Context(), req)
	})
}

func (h *handler) RemoveMulti(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Backup_Remove_RequestMulti
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.backup.RemoveMulti(r.Context(), req)
	})
}

func (h *handler) RegisterIPs(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Backup_IP_Register_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.backup.RegisterIPs(r.Context(), req)
	})
}

func (h *handler) RemoveIPs(w http.ResponseWriter, r *http.Request) (int, error) {
	var req *payload.Backup_IP_Remove_Request
	return json.Handler(w, r, &req, func() (interface{}, error) {
		return h.backup.RemoveIPs(r.Context(), req)
	})
}

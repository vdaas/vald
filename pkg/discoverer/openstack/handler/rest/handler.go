//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//



// Package setting stores all server application settings
package rest

import (
	"fmt"
	"net/http"

	"github.com/vdaas/vald/pkg/discoverer/openstack/service"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request) error
	Search(w http.ResponseWriter, r *http.Request) error
	SearchByID(w http.ResponseWriter, r *http.Request) error
	Insert(w http.ResponseWriter, r *http.Request) error
	MultiInsert(w http.ResponseWriter, r *http.Request) error
	Update(w http.ResponseWriter, r *http.Request) error
	MultiUpdate(w http.ResponseWriter, r *http.Request) error
	Remove(w http.ResponseWriter, r *http.Request) error
	MultiRemove(w http.ResponseWriter, r *http.Request) error
	CreateIndex(w http.ResponseWriter, r *http.Request) error
	SaveIndex(w http.ResponseWriter, r *http.Request) error
	GetObject(w http.ResponseWriter, r *http.Request) error
}

type handler struct {
	ngt service.NGT
}

func New() Handler {
	return &handler{}
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, r.URL.String())
	return nil
}

func (h *handler) Search(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) SearchByID(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) Insert(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) MultiInsert(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) MultiUpdate(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) MultiRemove(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) CreateIndex(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) SaveIndex(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) GetObject(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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

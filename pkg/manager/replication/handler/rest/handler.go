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

// Package rest provides rest api logic
package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request) error
	Exists(w http.ResponseWriter, r *http.Request) error
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
	agent agent.AgentServer
}

func New(opts ...Option) Handler {
	h := new(handler)

	for _, opt := range append(defaultOpts, opts...) {
		opt(h)
	}
	return h
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, r.URL.String())
	return nil
}

func (h *handler) Search(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Search_Request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.Search(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) SearchByID(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Search_IDRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.SearchByID(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) Insert(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Object_Vector
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.Insert(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) MultiInsert(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Object_Vectors
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.MultiInsert(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Object_Vector
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.Update(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) MultiUpdate(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Object_Vectors
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.MultiUpdate(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Object_ID
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.Remove(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) MultiRemove(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Object_IDs
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.MultiRemove(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) CreateIndex(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Controll_CreateIndexRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.CreateIndex(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) SaveIndex(w http.ResponseWriter, r *http.Request) (err error) {
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	_, err = h.agent.SaveIndex(r.Context(), nil)
	return
}

func (h *handler) GetObject(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Object_ID
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.GetObject(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) Exists(w http.ResponseWriter, r *http.Request) (err error) {
	var req *payload.Object_ID
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	res, err := h.agent.Exists(r.Context(), req)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return err
	}
	return nil
}

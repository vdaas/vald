//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
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

package json

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/rest"
)

type RFC7807Error struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Detail   interface{} `json:"detail"`
	Instance string      `json:"instance"`
	Status   int         `json:"status"`
	Error    string      `json:"error"`
}

func Encode(w io.Writer, data interface{}) (err error) {
	return jsoniter.NewEncoder(w).Encode(data)
}

func Decode(r io.Reader, data interface{}) (err error) {
	return jsoniter.NewDecoder(r).Decode(data)
}

func MarshalIndent(data interface{}, pref, ind string) ([]byte, error) {
	return jsoniter.MarshalIndent(data, pref, ind)
}

func EncodeResponse(w http.ResponseWriter,
	data interface{}, status int, contentTypes ...string) error {
	for _, ct := range contentTypes {
		w.Header().Add(rest.ContentType, ct)
	}
	w.WriteHeader(status)
	return Encode(w, data)
}

func DecodeRequest(r *http.Request, data interface{}) (err error) {
	if r != nil && r.Body != nil && r.ContentLength != 0 {
		err = Decode(r.Body, data)
		if err != nil {
			return err
		}
		_, err := io.Copy(ioutil.Discard, r.Body)
		if err != nil {
			return errors.ErrRequestBodyFlush(err)
		}
		// close
		err = r.Body.Close()
		if err != nil {
			return errors.ErrRequestBodyClose(err)
		}
	}
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request,
	data interface{}, logic func() (interface{},
		error)) (code int, err error) {
	err = DecodeRequest(r, &data)
	if err != nil {
		return http.StatusBadRequest, err
	}
	res, err := logic()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = EncodeResponse(w, res, http.StatusOK,
		rest.ApplicationJSON, rest.CharsetUTF8)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	return http.StatusOK, nil
}

func ErrorHandler(w http.ResponseWriter, r *http.Request,
	msg string, code int, err error) error {
	data := RFC7807Error{
		Type:   r.RequestURI,
		Title:  msg,
		Status: code,
		Error:  err.Error(),
	}
	data.Instance, err = os.Hostname()
	if err != nil {
		log.Error(err)
	}
	body := make(map[string]interface{})
	err = Decode(r.Body, &body)
	if err != nil {
		log.Error(err)
	}
	data.Detail, err = dump.Request(nil, body, r)
	if err != nil {
		log.Error(err)
	}

	res, err := MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	w.Header().Add(rest.ContentType, rest.ProblemJSON)
	w.Header().Add(rest.ContentType, rest.CharsetUTF8)
	w.WriteHeader(code)
	w.Write(res)
	return nil
}

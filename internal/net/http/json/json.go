//
// Copyright (C) 2019 kpango (Yusuke Kato)
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

package json

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/rest"
)

type RFC7807Error struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
	Status   int    `json:"status"`
	Error    string `json:"error"`
}

func Encode(w io.Writer, data interface{}) (err error) {
	return jsoniter.NewEncoder(w).Encode(data)
}

func Decode(r io.Reader, data interface{}) (err error) {
	return jsoniter.NewDecoder(r).Decode(data)
}

func EncodeResponse(w http.ResponseWriter, data interface{}, status int, contentTypes ...string) error {
	for _, ct := range contentTypes {
		w.Header().Add(rest.ContentType, ct)
	}
	w.WriteHeader(status)
	return Encode(w, data)
}

func DecodeRequest(r *http.Request, data interface{}) (err error) {
	if r != nil && r.Body != nil {
		err = Decode(r.Body, data)
		if err != nil {
			return err
		}
		io.Copy(ioutil.Discard, r.Body)
		return r.Body.Close()
	}
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request,
	data interface{}, logic func() (interface{}, error)) (code int, err error) {
	err = DecodeRequest(r, &data)
	if err != nil {
		return http.StatusBadRequest, err
	}
	res, err := logic()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = EncodeResponse(w, res, http.StatusOK, rest.ApplicationJSON, rest.CharsetUTF8)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func ErrorHandler(w http.ResponseWriter, data RFC7807Error) (err error) {
	data.Instance, err = os.Hostname()
	if err != nil {
		log.Error(err)
	}
	return EncodeResponse(w, data, data.Status, rest.ProblemJSON, rest.CharsetUTF8)
}

//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"bytes"
	"context"
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

func Unmarshal(data []byte, i interface{}) error {
	return jsoniter.Unmarshal(data, i)
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

func DecodeResponse(res *http.Response, data interface{}) (err error) {
	if res != nil && res.Body != nil && res.ContentLength != 0 {
		err = Decode(res.Body, data)
		if err != nil {
			return err
		}
		_, err := io.Copy(ioutil.Discard, res.Body)
		if err != nil {
			return errors.ErrRequestBodyFlush(err)
		}
		// close
		err = res.Body.Close()
		if err != nil {
			return errors.ErrRequestBodyClose(err)
		}
	}
	return nil
}

func EncodeRequest(req *http.Request,
	data interface{}, contentTypes ...string) error {
	for _, ct := range contentTypes {
		req.Header.Add(rest.ContentType, ct)
	}
	buf := new(bytes.Buffer)
	if err := Encode(buf, data); err != nil {
		return err
	}
	return req.Write(buf)
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

func Request(ctx context.Context, method string, url string, payloyd interface{}, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}

	if payloyd != nil && method != http.MethodGet {
		if err := EncodeRequest(req, payloyd, rest.ApplicationJSON, rest.CharsetUTF8); err != nil {
			return err
		}
	}

	// TODO replace vald original client.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return DecodeResponse(resp, data)
}

// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package json

import (
	"bytes"
	"net/http"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/net/http/dump"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/os"
)

// statusClientClosedRequest is the non-standard 499 status that grpc-gateway uses for codes.Canceled.
const statusClientClosedRequest = 499

// RFC7807Error represents RFC 7807 error.
type RFC7807Error struct {
	Detail   any    `json:"detail"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Instance string `json:"instance"`
	Error    string `json:"error"`
	Status   int    `json:"status"`
}

// EncodeResponse encodes http response body.
// data is encoded before the header is written, so an encode failure leaves the response untouched.
func EncodeResponse(w http.ResponseWriter, data any, status int, contentTypes ...string) error {
	var buf bytes.Buffer
	err := json.Encode(&buf, data)
	if err != nil {
		return err
	}
	return writeResponse(w, buf.Bytes(), status, contentTypes...)
}

func writeResponse(w http.ResponseWriter, body []byte, status int, contentTypes ...string) error {
	for _, ct := range contentTypes {
		w.Header().Add(rest.ContentType, ct)
	}
	w.WriteHeader(status)
	_, err := w.Write(body)
	return err
}

// DecodeRequest decodes http request body.
func DecodeRequest(r *http.Request, data any) (err error) {
	if r != nil && r.Body != nil && r.ContentLength != 0 {
		err = json.Decode(r.Body, data)
		if err != nil {
			return err
		}
		_, err := io.Copy(io.Discard, r.Body)
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

// Handler responds to an HTTP request to perform a logic function.
func Handler(
	w http.ResponseWriter, r *http.Request, data any, logic func() (any,
		error),
) (code int, err error) {
	err = DecodeRequest(r, &data)
	if err != nil {
		return http.StatusBadRequest, err
	}
	res, err := logic()
	if err != nil {
		code = httpStatusFromError(err)
		if code == http.StatusServiceUnavailable {
			// The router sends no error body for 503 because 503 means the response
			// has already been written, so write the RFC 7807 error here.
			if eerr := ErrorHandler(w, r, err.Error(), code, err); eerr != nil {
				log.Error(eerr)
			}
		}
		return code, err
	}
	var buf bytes.Buffer
	err = json.Encode(&buf, res)
	if err != nil {
		// Nothing has been written yet and the router sends no error body for 503,
		// so send the status here instead of leaving the client an implicit 200.
		w.WriteHeader(http.StatusServiceUnavailable)
		return http.StatusServiceUnavailable, err
	}
	err = writeResponse(w, buf.Bytes(), http.StatusOK,
		rest.ApplicationJSON, rest.CharsetUTF8)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	return http.StatusOK, nil
}

// ErrorHandler responds to an HTTP request to send RFC 7807 error.
func ErrorHandler(w http.ResponseWriter, r *http.Request, msg string, code int, err error) error {
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
	body := make(map[string]any)
	err = json.Decode(r.Body, &body)
	if err != nil && !isConsumedBody(err) {
		log.Error(err)
	}
	data.Detail, err = dump.Request(nil, body, r)
	if err != nil {
		log.Error(err)
	}

	res, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	w.Header().Add(rest.ContentType, rest.ProblemJSON)
	w.Header().Add(rest.ContentType, rest.CharsetUTF8)
	w.WriteHeader(code)
	_, err = w.Write(res)
	return err
}

// httpStatusFromError returns the HTTP status for the gRPC code carried by err,
// using the grpc-gateway mapping. Errors without a known gRPC code map to 500.
func httpStatusFromError(err error) int {
	st, _ := status.FromError(err)
	if st == nil {
		return http.StatusInternalServerError
	}
	//nolint:exhaustive // only codes with a dedicated grpc-gateway HTTP mapping are listed; every other code intentionally falls through to 500
	switch st.Code() {
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Canceled:
		return statusClientClosedRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

// isConsumedBody reports whether err comes from decoding a request body that has already been
// read to the end or closed (e.g. by DecodeRequest in Handler), which is not a decode failure.
func isConsumedBody(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, http.ErrBodyReadAfterClose)
}

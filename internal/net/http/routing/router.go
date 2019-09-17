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

// Package routing provides implementation of Go API for routing http Handler wrapped by rest.Func
package routing

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/json"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/safety"
)

type router struct {
	eg      errgroup.Group
	timeout time.Duration
	routes  []Route
}

//New returns Routed http.Handler
func New(opts ...Option) http.Handler {
	r := new(router)
	for _, opt := range append(defaultOpts, opts...) {
		opt(r)
	}

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 32

	rt := mux.NewRouter().StrictSlash(true)
	for _, route := range r.routes {
		rt.Handle(route.Pattern, r.routing(route.Name, route.Pattern, route.Methods, route.HandlerFunc))
	}

	return rt
}

// routing wraps the handler.Func and returns a new http.Handler.
// routing helps to handle unsupported HTTP method, timeout,
// and the error returned from the handler.Func.
func (rt *router) routing(name, path string, m []string, h rest.Func) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		for _, method := range m {
			if strings.EqualFold(r.Method, method) {
				// execute only if the request method is inside the method list

				// context for timeout
				ctx, cancel := context.WithTimeout(r.Context(), rt.timeout)
				defer cancel()
				start := fastime.UnixNanoNow()

				// run the custom handler logic in go routine,
				// report error to error channel
				ech := make(chan error)
				sch := make(chan int)
				defer close(ech)
				defer close(sch)
				rt.eg.Go(safety.RecoverFunc(func() (err error) {
					// it is the responsibility for handler to close the request
					var code int
					code, err = h(w, r.WithContext(ctx))
					if err != nil {
						sch <- code
					}
					ech <- err
					return nil
				}))

				select {
				case err = <-ech:
					// handler finished first, may have error returned
					if err != nil {
						code := <-sch
						err = errors.ErrHandler(err)

						log.Error(err)

						err = json.ErrorHandler(w, json.RFC7807Error{
							Type:  path,
							Title: fmt.Sprintf("Error on Handler: %s\tPath: %s\tMethod: %s\t", name, path, method),
							Detail: fmt.Sprintf("Handler: %s\tPath: %s\tMethod: %s\tStatus: %s\nError: %s",
								name, path, method, http.StatusText(code), err.Error()),
							Status: code,
							Error:  err.Error(),
						})
						if err != nil {
							log.Error(err)
						}
					}
					return
				case <-ctx.Done():
					// timeout passed or parent context canceled first,
					// it is the responsibility for handler to response to the user
					err = errors.ErrHandlerTimeout(ctx.Err(), time.Unix(0, fastime.UnixNanoNow()-start))
					if err != nil {
						log.Error(err)
					}

					err = json.ErrorHandler(w, json.RFC7807Error{
						Type:  path,
						Title: fmt.Sprintf("Timeout Error on Handler: %s\tPath: %s\tMethod: %s\t", name, path, method),
						Detail: fmt.Sprintf("Handler: %s\tPath: %s\tMethod: %s\tStatus: %s\nError: %s",
							name, path, method, http.StatusText(http.StatusRequestTimeout), err.Error()),
						Status: http.StatusRequestTimeout,
						Error:  err.Error(),
					})
					if err != nil {
						log.Error(err)
					}
					return
				}
			}
		}

		// flush and close the request body; for GET method, r.Body may be nil
		err = errors.Wrap(err, flushAndClose(r.Body).Error())
		if err != nil {
			err = errors.ErrRequestBodyCloseAndFlush(err)
			log.Error(err)
		}

		err = json.ErrorHandler(w, json.RFC7807Error{
			Type:  path,
			Title: fmt.Sprintf("Invalid Request Method Handler: %s\tPath: %s\tMethod: %s\t", name, path, r.Method),
			Detail: fmt.Sprintf("Handler: %s\tPath: %s\tMethod: %s\tStatus: %s\nError: %s",
				name, path, r.Method, http.StatusText(http.StatusMethodNotAllowed), err.Error()),
			Status: http.StatusMethodNotAllowed,
			Error:  err.Error(),
		})
		if err != nil {
			log.Error(err)
		}

	})
}

// flushAndClose helps to flush and close a ReadCloser. Used for request body internal.
// Returns if there is any errors.
func flushAndClose(rc io.ReadCloser) error {
	if rc != nil {
		// flush
		_, err := io.Copy(ioutil.Discard, rc)
		if err != nil {
			return errors.ErrRequestBodyFlush(err)
		}
		// close
		err = rc.Close()
		if err != nil {
			return errors.ErrRequestBodyClose(err)
		}
	}
	return nil
}

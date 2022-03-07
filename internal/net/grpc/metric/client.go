//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package metric provides metrics functions for grpc
package metric

import (
	"go.opencensus.io/plugin/ocgrpc"
)

// ClientHandler is a type alias of ocgrpc.ClientHandler to record OpenCensus stats and traces.
type ClientHandler = ocgrpc.ClientHandler

// NewClientHandler returns the client handler of OpenCensus stats and traces.
func NewClientHandler(opts ...ClientOption) *ClientHandler {
	handler := new(ClientHandler)

	for _, opt := range append(clientDefaultOpts, opts...) {
		opt(handler)
	}

	return handler
}

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

// Package ngtd provides ngtd starter  functionality
package ngtd

type Option func(*server)

var defaultOptions = []Option{
	WithDimension(128),
	WithIndexDir("/tmp/ngtd/"),
	WithServerType(HTTP),
	WithPort(8200),
}

func WithDimension(dim int) Option {
	return func(n *server) {
		if dim > 0 {
			n.dim = dim
		}
	}
}

func WithServerType(t ServerType) Option {
	return func(n *server) {
		n.srvType = t
	}
}

func WithIndexDir(path string) Option {
	return func(n *server) {
		if len(path) != 0 {
			n.indexDir = path
		}
	}
}

func WithPort(port int) Option {
	return func(n *server) {
		if port > 0 {
			n.port = port
		}
	}
}

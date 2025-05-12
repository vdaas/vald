//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package grpc provides grpc server logic
package grpc

import (
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
)

type Option func(*server) error

var defaultOptions = []Option{}

// WithName returns the functional option to the server name.
func WithName(name string) Option {
}

// TODO: oai 切り替え
// TODO: errors 切り出し
// TODO: handler → usecase で yaml 決まれば動くはず
// WithValdClient returns the functional option to set the vald client.
func WithValdClient(c vald.Client) Option {
	return func(n *embedder) error {
		if c == nil {
			return errors.New("vald client is nil")
		}
		n.client = c
		return nil
	}
}

//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package service

import (
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
)

type Option func(*embedder) error

var defaultOptions = []Option{}

func WithValdClient(c vald.Client) Option {
	return func(e *embedder) error {
		if c == nil {
			return errors.NewErrInvalidOption("ValdClient", c, errors.New("vald client is nil"))
		}
		e.client = c
		return nil
	}
}

func WithMetaClient(c MetaClient) Option {
	return func(e *embedder) error {
		if c == nil {
			return errors.NewErrInvalidOption("MetaClient", c, errors.New("meta client is nil"))
		}
		e.mclient = c
		return nil
	}
}

func WithLLM(c LLM) Option {
	return func(e *embedder) error {
		if c == nil {
			return errors.NewErrInvalidOption("LLM", c, errors.New("llm is nil"))
		}
		e.llm = c
		return nil
	}
}

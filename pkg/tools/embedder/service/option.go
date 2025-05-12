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

package service

import (
	agent "github.com/vdaas/vald/internal/client/v1/client/agent/core"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
)

// Option represent the functional option for ngt.
type Option func(n *embedder) error

var defaultOptions = []Option{}

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

// WithAgentClient returns the functional option to set the agent client.
func WithAgentClient(c agent.Client) Option {
	return func(n *embedder) error {
		if c == nil {
			return errors.New("agent client is nil")
		}
		n.aclient = c
		return nil
	}
}

// WithLLM returns the functional option to set the LLM client.
func WithLLM(c LLM) Option {
	return func(n *embedder) error {
		if c == nil {
			return errors.New("LLM client is nil")
		}
		n.llm = c
		return nil
	}
}

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

// Package tensorflow provides implementation of Go API for extract data to vector
package tensorflow

type Option func(*tensorflow)

var (
	defaultOpts = []Option{
		WithOperations(),        // set to default
		WithSessionOptions(nil), // set to default
		WithNdim(0),             // set to default
	}
)

func WithSessionOptions(opts *SessionOptions) Option {
	return func(t *tensorflow) {
		t.options = opts
	}
}

func WithSessionTarget(tgt string) Option {
	return func(t *tensorflow) {
		if len(tgt) != 0 {
			t.sessionTarget = tgt
		}
	}
}

func WithSessionConfig(cfg []byte) Option {
	return func(t *tensorflow) {
		if cfg != nil {
			t.sessionConfig = cfg
		}
	}
}

func WithOperations(opes ...*Operation) Option {
	return func(t *tensorflow) {
		if opes != nil {
			if t.operations != nil {
				t.operations = append(t.operations, opes...)
			} else {
				t.operations = opes
			}
		}
	}
}

func WithExportPath(path string) Option {
	return func(t *tensorflow) {
		if len(path) != 0 {
			t.exportDir = path
		}
	}
}

func WithTags(tags ...string) Option {
	return func(t *tensorflow) {
		if tags != nil {
			if t.tags != nil {
				t.tags = append(t.tags, tags...)
			} else {
				t.tags = tags
			}
		}
	}
}

func WithNdim(ndim int8) Option {
	return func(t *tensorflow) {
		t.ndim = ndim
	}
}

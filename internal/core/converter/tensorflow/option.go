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

// Package tensorflow provides implementation of Go API for extract data to vector
package tensorflow

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// Option is tensorflow configure.
type Option func(*tensorflow)

var defaultOptions = []Option{
	withLoadFunc(tf.LoadSavedModel), // set to default
	WithOperations(),                // set to default
	WithSessionOptions(nil),         // set to default
	WithNdim(0),                     // set to default
}

// WithSessionOptions returns Option that sets options.
func WithSessionOptions(opts *SessionOptions) Option {
	return func(t *tensorflow) {
		if opts != nil {
			t.options = opts
		}
	}
}

// WithSessionTarget returns Option that sets target.
func WithSessionTarget(tgt string) Option {
	return func(t *tensorflow) {
		if tgt != "" {
			if t.options == nil {
				t.options = &SessionOptions{
					Target: tgt,
				}
			} else {
				t.options.Target = tgt
			}
		}
	}
}

// WithSessionConfig returns Option that sets config.
func WithSessionConfig(cfg []byte) Option {
	return func(t *tensorflow) {
		if cfg != nil {
			if t.options == nil {
				t.options = &SessionOptions{
					Config: cfg,
				}
			} else {
				t.options.Config = cfg
			}
		}
	}
}

// WithOperations returns Option that sets operations.
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

// WithExportPath returns Option that sets exportDir.
func WithExportPath(path string) Option {
	return func(t *tensorflow) {
		if path != "" {
			t.exportDir = path
		}
	}
}

// WithTags returns Option that sets tags.
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

func withLoadFunc(
	loadFunc func(exportDir string, tags []string, options *SessionOptions) (*tf.SavedModel, error)) Option {
	return func(t *tensorflow) {
		if loadFunc != nil {
			t.loadFunc = loadFunc
		}
	}
}

// WithFeed returns Option that sets feeds.
func WithFeed(operationName string, outputIndex int) Option {
	return func(t *tensorflow) {
		t.feeds = append(t.feeds, OutputSpec{operationName, outputIndex})
	}
}

// WithFeeds returns Option that sets feeds.
func WithFeeds(operationNames []string, outputIndexes []int) Option {
	return func(t *tensorflow) {
		if operationNames != nil && outputIndexes != nil && len(operationNames) == len(outputIndexes) {
			for i := range operationNames {
				t.feeds = append(t.feeds, OutputSpec{operationNames[i], outputIndexes[i]})
			}
		}
	}
}

// WithFetch returns Option that sets fetches.
func WithFetch(operationName string, outputIndex int) Option {
	return func(t *tensorflow) {
		t.fetches = append(t.fetches, OutputSpec{operationName, outputIndex})
	}
}

// WithFetches returns Option that sets fetches.
func WithFetches(operationNames []string, outputIndexes []int) Option {
	return func(t *tensorflow) {
		if operationNames != nil && outputIndexes != nil && len(operationNames) == len(outputIndexes) {
			for i := range operationNames {
				t.fetches = append(t.fetches, OutputSpec{operationNames[i], outputIndexes[i]})
			}
		}
	}
}

// WithWarmupInputs returns Option that sets warmupInputs.
func WithWarmupInputs(warmupInputs ...string) Option {
	return func(t *tensorflow) {
		if warmupInputs != nil {
			if t.warmupInputs != nil {
				t.warmupInputs = append(t.warmupInputs, warmupInputs...)
			} else {
				t.warmupInputs = warmupInputs
			}
		}
	}
}

// WithNdim returns Option that sets ndim.
func WithNdim(ndim uint8) Option {
	return func(t *tensorflow) {
		t.ndim = ndim
	}
}

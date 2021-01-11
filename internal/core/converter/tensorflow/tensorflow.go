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
	"io"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/vdaas/vald/internal/errors"
)

// SessionOptions is a type alias for tensorflow.SessionOptions.
type SessionOptions = tf.SessionOptions

// Operation is a type alias for tensorflow.Operation.
type Operation = tf.Operation

// Closer is a type alias io.Closer.
type Closer = io.Closer

// TF represents a tensorflow interface.
type TF interface {
	GetVector(inputs ...string) ([]float64, error)
	GetValue(inputs ...string) (interface{}, error)
	GetValues(inputs ...string) (values []interface{}, err error)
	Closer
}

type session interface {
	Run(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*Operation) ([]*tf.Tensor, error)
	Closer
}

type tensorflow struct {
	exportDir    string
	tags         []string
	loadFunc     func(exportDir string, tags []string, options *SessionOptions) (*tf.SavedModel, error)
	feeds        []OutputSpec
	fetches      []OutputSpec
	operations   []*Operation
	options      *SessionOptions
	graph        *tf.Graph
	session      session
	warmupInputs []string
	ndim         uint8
}

// OutputSpec is the specification of an feed/fetch.
type OutputSpec struct {
	operationName string
	outputIndex   int
}

const (
	twoDim uint8 = iota + 2
	threeDim
)

// New load a tensorlfow model and returns a new tensorflow struct.
func New(opts ...Option) (TF, error) {
	t := new(tensorflow)

	for _, opt := range append(defaultOptions, opts...) {
		opt(t)
	}

	model, err := t.loadFunc(t.exportDir, t.tags, t.options)
	if err != nil {
		return nil, err
	}

	t.graph = model.Graph
	t.session = model.Session

	err = t.warmup()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *tensorflow) warmup() error {
	if t.warmupInputs != nil {
		_, err := t.run(t.warmupInputs...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *tensorflow) Close() error {
	return t.session.Close()
}

func (t *tensorflow) run(inputs ...string) ([]*tf.Tensor, error) {
	if len(inputs) != len(t.feeds) {
		return nil, errors.ErrInputLength(len(inputs), len(t.feeds))
	}

	feeds := make(map[tf.Output]*tf.Tensor, len(inputs))

	for i, val := range inputs {
		inputTensor, err := tf.NewTensor(val)
		if err != nil {
			return nil, err
		}

		feeds[t.graph.Operation(t.feeds[i].operationName).Output(t.feeds[i].outputIndex)] = inputTensor
	}

	fetches := make([]tf.Output, 0, len(t.fetches))
	for _, fetch := range t.fetches {
		fetches = append(fetches, t.graph.Operation(fetch.operationName).Output(fetch.outputIndex))
	}

	return t.session.Run(feeds, fetches, t.operations)
}

func (t *tensorflow) GetVector(inputs ...string) ([]float64, error) {
	tensors, err := t.run(inputs...)
	if err != nil {
		return nil, err
	}

	if len(tensors) == 0 || tensors[0] == nil || tensors[0].Value() == nil {
		return nil, errors.ErrNilTensorTF(tensors)
	}

	switch t.ndim {
	case twoDim:
		value, ok := tensors[0].Value().([][]float64)
		if ok {
			if value == nil {
				return nil, errors.ErrNilTensorValueTF(value)
			}

			return value[0], nil
		}

		return nil, errors.ErrFailedToCastTF(tensors[0].Value())
	case threeDim:
		value, ok := tensors[0].Value().([][][]float64)
		if ok {
			if len(value) == 0 || value[0] == nil {
				return nil, errors.ErrNilTensorValueTF(value)
			}

			return value[0][0], nil
		}

		return nil, errors.ErrFailedToCastTF(tensors[0].Value())
	default:
		value, ok := tensors[0].Value().([]float64)
		if ok {
			return value, nil
		}

		return nil, errors.ErrFailedToCastTF(tensors[0].Value())
	}
}

func (t *tensorflow) GetValue(inputs ...string) (interface{}, error) {
	tensors, err := t.run(inputs...)
	if err != nil {
		return nil, err
	}

	if len(tensors) == 0 || tensors[0] == nil {
		return nil, errors.ErrNilTensorTF(tensors)
	}

	return tensors[0].Value(), nil
}

func (t *tensorflow) GetValues(inputs ...string) (values []interface{}, err error) {
	tensors, err := t.run(inputs...)
	if err != nil {
		return nil, err
	}

	values = make([]interface{}, 0, len(tensors))
	for _, tensor := range tensors {
		values = append(values, tensor.Value())
	}

	return values, nil
}

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

package tensorflow

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/vdaas/vald/internal/errors"
)

type SessionOptions = tf.SessionOptions
type Operation = tf.Operation

type TF interface {
	GetVector(feeds []Feed, fetches []Fetch, targets ...*Operation) (values [][][]float64, err error)
	Close() error
}

type tensorflow struct {
	exportDir     string
	tags          []string
	operations    []*Operation
	sessionTarget string
	sessionConfig []byte
	options       *SessionOptions
	graph         *tf.Graph
	session       *tf.Session
}

type Feed struct {
	InputBytes    []byte
	OperationName string
	OutputIndex   int
}

type Fetch struct {
	OperationName string
	OutputIndex   int
}

func New(opts ...Option) (TF, error) {
	t := new(tensorflow)
	for _, opt := range append(defaultOpts, opts...) {
		opt(t)
	}

	if t.options == nil && (len(t.sessionTarget) != 0 || t.sessionConfig != nil) {
		t.options = &tf.SessionOptions{
			Target: t.sessionTarget,
			Config: t.sessionConfig,
		}
	}

	model, err := tf.LoadSavedModel(t.exportDir, t.tags, t.options)
	if err != nil {
		return nil, err
	}
	t.graph = model.Graph
	t.session = model.Session
	return t, nil
}

func (t *tensorflow) Close() error {
	return t.session.Close()
}

func (t *tensorflow) GetVector(feeds []Feed, fetches []Fetch, targets ...*Operation) (values [][][]float64, err error) {
	input := make(map[tf.Output]*tf.Tensor, len(feeds))
	for _, feed := range feeds {
		inputTensor, err := tf.NewTensor([]string{string(feed.InputBytes)})
		if err != nil {
			return nil, err
		}
		input[t.graph.Operation(feed.OperationName).Output(feed.OutputIndex)] = inputTensor
	}

	output := make([]tf.Output, 0, len(fetches))
	for _, fetch := range fetches {
		output = append(output, t.graph.Operation(fetch.OperationName).Output(fetch.OutputIndex))
	}

	if targets == nil {
		targets = t.operations
	}

	results, err := t.session.Run(input, output, targets)
	if err != nil {
		return nil, err
	}

	values = make([][][]float64, 0, len(results))
	for _, result := range results {
		value, ok := result.Value().([][]float64)
		if ok {
			values = append(values, value)
		} else {
			return nil, errors.ErrFailedToCastTF(result.Value())
		}
	}
	return values, nil
}

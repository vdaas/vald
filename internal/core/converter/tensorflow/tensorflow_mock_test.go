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

type mockSession struct {
	RunFunc   func(map[tf.Output]*tf.Tensor, []tf.Output, []*Operation) ([]*tf.Tensor, error)
	CloseFunc func() error
}

func (m *mockSession) Run(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*Operation) ([]*tf.Tensor, error) {
	return m.RunFunc(feeds, fetches, operations)
}

func (m *mockSession) Close() error {
	return m.CloseFunc()
}

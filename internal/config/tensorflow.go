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

// Package config providers configuration type and load configuration logic
package config

import "encoding/base64"

// Tensorflow represent the Tensorflow converter core configuration for server.
type Tensorflow struct {
	SessiontOption        *SessionOption `json:"sessiont_option,omitempty" yaml:"sessiont_option"`
	ExportPath            string         `json:"export_path,omitempty" yaml:"export_path"`
	Tags                  []string       `json:"tags,omitempty" yaml:"tags"`
	Feeds                 []*OutputSpec  `json:"feeds,omitempty" yaml:"feeds"`
	FeedsMap              map[string]int `json:"-" yaml:"-"`
	Fetches               []*OutputSpec  `json:"fetches,omitempty" yaml:"fetches"`
	FetchesMap            map[string]int `json:"-" yaml:"-"`
	WarmupInputs          []string       `json:"warmup_inputs,omitempty" yaml:"warmup_inputs"`
	ResultNestedDimension uint8          `json:"result_nested_dimension,omitempty" yaml:"result_nested_dimension"`
}

type SessionOption struct {
	Target       string `json:"target,omitempty" yaml:"target"`
	Base64Config string `json:"base64_config,omitempty" yaml:"base64_config"`
	Config       []byte `json:"-" yaml:"-"`
}

type OutputSpec struct {
	OperationName string `json:"operation_name,omitempty" yaml:"operation_name"`
	OutputIndex   int    `json:"output_index,omitempty" yaml:"output_index"`
}

// Bind returns Tensorflow object whose some string value is filed value or environment value.
func (tf *Tensorflow) Bind() *Tensorflow {
	tf.SessiontOption = tf.SessiontOption.Bind()
	for i, spec := range tf.Feeds {
		tf.Feeds[i] = spec.Bind()
		tf.FeedsMap[tf.Feeds[i].OperationName] = tf.Feeds[i].OutputIndex
	}
	for i, spec := range tf.Fetches {
		tf.Fetches[i] = spec.Bind()
		tf.FetchesMap[tf.Fetches[i].OperationName] = tf.Fetches[i].OutputIndex
	}
	tf.ExportPath = GetActualValue(tf.ExportPath)
	tf.Tags = GetActualValues(tf.Tags)
	tf.WarmupInputs = GetActualValues(tf.WarmupInputs)
	return tf
}

func (s *SessionOption) Bind() *SessionOption {
	s.Target = GetActualValue(s.Target)
	s.Base64Config = GetActualValue(s.Base64Config)
	b, err := base64.StdEncoding.DecodeString(s.Base64Config)
	if err == nil {
		s.Config = b
	}
	return s
}

func (o *OutputSpec) Bind() *OutputSpec {
	o.OperationName = GetActualValue(o.OperationName)
	return o
}

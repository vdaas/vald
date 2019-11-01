//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package grpc provides grpc server logic
package model

import (
	"fmt"
	"strings"
)

const (
	comma = ","
)

type MetaVector struct {
	UUID   string
	Vector []float64
	Meta   string
	IPs    []string
}

func (m *MetaVector) GetUUID() string               { return m.UUID }
func (m *MetaVector) GetVector() ([]float64, error) { return m.Vector, nil }
func (m *MetaVector) GetVectorString() string {
	ss := make([]string, 0, len(m.Vector))
	for _, f := range m.Vector {
		ss = append(ss, fmt.Sprint(f))
	}
	return strings.Join(ss, comma)
}
func (m *MetaVector) GetMeta() string  { return m.Meta }
func (m *MetaVector) GetIPs() []string { return m.IPs }

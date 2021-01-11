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

type EgressFilter struct {
	Client *GRPCClient `json:"client" yaml:"client"`
}

type IngressFilter struct {
	Client *GRPCClient `json:"client,omitempty" yaml:"client"`
	Search []string    `json:"search,omitempty" yaml:"search"`
	Insert []string    `json:"insert,omitempty" yaml:"insert"`
	Update []string    `json:"update,omitempty" yaml:"update"`
	Upsert []string    `json:"upsert,omitempty" yaml:"upsert"`
}

func (e *EgressFilter) Bind() *EgressFilter {
	if e.Client != nil {
		e.Client.Bind()
	}
	return e
}

func (i *IngressFilter) Bind() *IngressFilter {
	if i.Client != nil {
		i.Client.Bind()
	}
	if i.Search != nil {
		i.Search = GetActualValues(i.Search)
	}
	if i.Insert != nil {
		i.Insert = GetActualValues(i.Insert)
	}
	if i.Update != nil {
		i.Update = GetActualValues(i.Update)
	}
	if i.Upsert != nil {
		i.Upsert = GetActualValues(i.Upsert)
	}
	return i
}

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

package config

// EgressFilter represents the EgressFilter configuration.
type EgressFilter struct {
	// Client represents the client configuration.
	Client *GRPCClient `json:"client,omitempty" yaml:"client"`
	// DistanceFilters represents the distance filters.
	DistanceFilters []string `json:"distance_filters,omitempty" yaml:"distance_filters"`
	// ObjectFilters represents the object filters.
	ObjectFilters []string `json:"object_filters,omitempty" yaml:"object_filters"`
}

// IngressFilter represents the IngressFilter configuration.
type IngressFilter struct {
	// Client represents the client configuration.
	Client *GRPCClient `json:"client,omitempty" yaml:"client"`
	// Vectorizer represents the vectorizer.
	Vectorizer string `json:"vectorizer,omitempty" yaml:"vectorizer"`
	// SearchFilters represents the search filters.
	SearchFilters []string `json:"search_filters,omitempty" yaml:"search_filters"`
	// InsertFilters represents the insert filters.
	InsertFilters []string `json:"insert_filters,omitempty" yaml:"insert_filters"`
	// UpdateFilters represents the update filters.
	UpdateFilters []string `json:"update_filters,omitempty" yaml:"update_filters"`
	// UpsertFilters represents the upsert filters.
	UpsertFilters []string `json:"upsert_filters,omitempty" yaml:"upsert_filters"`
}

// Bind binds the actual data from the EgressFilter receiver field.
func (e *EgressFilter) Bind() *EgressFilter {
	if e.Client != nil {
		e.Client.Bind()
	} else {
		e.Client = newGRPCClientConfig() // newGRPCClientConfig calls Bind internally
	}
	if e.DistanceFilters != nil {
		e.DistanceFilters = GetActualValues(e.DistanceFilters)
	}
	if e.ObjectFilters != nil {
		e.ObjectFilters = GetActualValues(e.ObjectFilters)
	}
	return e
}

// Bind binds the actual data from the IngressFilter receiver field.
func (i *IngressFilter) Bind() *IngressFilter {
	if i.Client != nil {
		i.Client.Bind()
	} else {
		i.Client = newGRPCClientConfig() // newGRPCClientConfig calls Bind internally
	}

	i.Vectorizer = GetActualValue(i.Vectorizer)

	if i.SearchFilters != nil {
		i.SearchFilters = GetActualValues(i.SearchFilters)
	}
	if i.InsertFilters != nil {
		i.InsertFilters = GetActualValues(i.InsertFilters)
	}
	if i.UpdateFilters != nil {
		i.UpdateFilters = GetActualValues(i.UpdateFilters)
	}
	if i.UpsertFilters != nil {
		i.UpsertFilters = GetActualValues(i.UpsertFilters)
	}
	return i
}

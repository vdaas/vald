//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package cacher provides implementation of cache type definition
package cacher

import (
	"context"

	"github.com/vdaas/vald/internal/strings"
)

// Cache represent the cache interface to store cache.
type Cache interface {
	Start(context.Context)
	Get(string) (interface{}, bool)
	Set(string, interface{})
	Delete(string)
	GetAndDelete(string) (interface{}, bool)
}

// Type represents the cacher type. Currently it support GACHE only.
type Type uint8

const (
	Unknown Type = iota
	GACHE
)

// String returns the type name.
func (m Type) String() string {
	switch m {
	case GACHE:
		return "gache"
	default:
		return "unknown"
	}
}

// ToType returns the type based on the string.
func ToType(str string) Type {
	switch strings.ToLower(str) {
	case "gache":
		return GACHE
	default:
		return Unknown
	}
}

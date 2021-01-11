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

// Package ngt provides ngt
package ngt

type Option func(*core)

var defaultOptions = []Option{
	WithIndexPath("tmpdir"),
	WithObjectType("float"),
	WithDimension(128),
}

func WithIndexPath(path string) Option {
	return func(c *core) {
		if len(path) != 0 {
			c.idxPath = path
		}
	}
}

func WithObjectType(typ string) Option {
	return func(c *core) {
		switch typ {
		case "uint8":
			c.objectType = Uint8
		case "float":
			c.objectType = Float
		default:
			c.objectType = ObjectNone
		}
	}
}

func WithDimension(dimension int) Option {
	return func(c *core) {
		if dimension > 0 {
			c.dimension = dimension
		}
	}
}

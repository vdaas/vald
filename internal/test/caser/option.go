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

import "testing"

type Option func(*caser)

var (
	defaultOptions = []Option{}
)

func WithName(str string) Option {
	return func(c *caser) {
		if len(str) != 0 {
			c.name = str
		}
	}
}

func WithArg(args ...interface{}) Option {
	return func(c *caser) {
		if len(args) != 0 {
			c.args = args
		}
	}
}

func WithField(fields ...interface{}) Option {
	return func(c *caser) {
		if len(fields) != 0 {
			c.fields = fields
		}
	}
}

func WithFieldFunc(fn func(*testing.T) []interface{}) Option {
	return func(c *caser) {
		if fn != nil {
			c.fieldFunc = fn
		}
	}
}

func WithWant(wants ...interface{}) Option {
	return func(c *caser) {
		if len(wants) != 0 {
			c.wants = wants
		}
	}
}

func WithAssertFunc(fn func(gots, wants []interface{}) error) Option {
	return func(c *caser) {
		if fn != nil {
			c.assertFunc = fn
		}
	}
}

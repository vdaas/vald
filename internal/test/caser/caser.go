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
package caser

import (
	"testing"

	"github.com/vdaas/vald/internal/test"
)

type caser struct {
	name       string
	args       []interface{}
	fields     []interface{}
	fieldFunc  func(*testing.T) []interface{}
	wants      []interface{}
	assertFunc func(gots, wants []interface{}) error
}

func New(opts ...Option) test.Caser {
	c := new(caser)
	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}
	return c
}

func (c *caser) Name() string {
	return c.name
}

func (c *caser) Args() []interface{} {
	return c.args
}

func (c *caser) Fields() []interface{} {
	return c.fields
}

func (c *caser) SetField(fields ...interface{}) {
	if len(fields) != 0 {
		c.fields = fields
	}
}

func (c *caser) FieldFunc() func(*testing.T) []interface{} {
	return c.fieldFunc
}

func (c *caser) Wants() []interface{} {
	return c.wants
}

func (c *caser) AssertFunc() func(gots, wants []interface{}) error {
	return c.assertFunc
}

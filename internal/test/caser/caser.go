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

func (c *caser) SetFields(fields ...interface{}) {
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

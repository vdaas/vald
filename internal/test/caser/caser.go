package caser

import (
	"github.com/vdaas/vald/internal/test"
)

type caser struct {
	name      string
	args      []interface{}
	fields    []interface{}
	wants     []interface{}
	checkFunc func() error
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

func (c *caser) Wants() []interface{} {
	return c.wants
}

func (c *caser) CheckFunc() func() error {
	return c.checkFunc
}

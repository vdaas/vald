package caser

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

package test

import "testing"

type Caser interface {
	Name() string
	Args() []interface{}
	Fields() []interface{}
	CheckFunc() func(*testing.T, ...interface{})
}

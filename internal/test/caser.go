package test

import "testing"

type Caser interface {
	Name() string
	DataProvider
	SetField([]interface{})
	FieldFunc() func(*testing.T) []interface{}
	AssertFunc() func(gots, want []interface{}) error
}

type DataProvider interface {
	Args() []interface{}
	Fields() []interface{}
	Wants() []interface{}
}

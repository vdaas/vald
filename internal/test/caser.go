package test

type Caser interface {
	Name() string
	Args() []interface{}
	Fields() []interface{}
	CheckFunc() func(...interface{}) error
}

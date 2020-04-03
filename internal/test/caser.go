package test

type Caser interface {
	Name() string
	Args() []interface{}
	Fields() []interface{}
	Check() func(...interface{}) error
}

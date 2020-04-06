package test

type Caser interface {
	Name() string
	Args() []interface{}
	Fields() []interface{}
	Wants() []interface{}
	CheckFunc() func(gots, want []interface{}) error
}

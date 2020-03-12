// Package strategy provides strategy for e2e testing functions
package strategy

type InsertOption func(*insert)

var (
	defaultInsertOption = []InsertOption{}
)

func WithParallelInsert() InsertOption {
	return func(e *insert) {
		e.parallel = true
	}
}

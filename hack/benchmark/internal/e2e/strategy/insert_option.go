// Package strategy provides strategy for e2e testing functions
package strategy

type InsertOption func(*insert)

var (
	defaultInsertOption = []InsertOption{
		WithParallelInsert(false),
	}
)

func WithParallelInsert(flag bool) InsertOption {
	return func(e *insert) {
		e.parallel = flag
	}
}

// Package strategy provides strategy for e2e testing functions
package strategy

type RemoveOption func(*remove)

var (
	defaultRemoveOptions = []RemoveOption{
		WithParallelRemove(false),
	}
)

func WithParallelRemove(flag bool) RemoveOption {
	return func(e *remove) {
		e.parallel = flag
	}
}

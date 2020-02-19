package strategy

type InsertOption func(*insert)

var (
	defaultInsertOption = []InsertOption{}
)

func WithParallel() InsertOption {
	return func(e *insert) {
		e.parallel = true
	}
}

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

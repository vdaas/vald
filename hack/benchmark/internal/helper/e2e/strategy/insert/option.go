package insert

type Option func(*insert)

var (
	defaultOption = []Option{}
)

func WithParallel() Option {
	return func(e *insert) {
		e.parallel = true
	}
}

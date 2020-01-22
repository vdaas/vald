package log

type Option func(*option)

var (
	defaultOptions = []Option{}
)

type option struct {
	logger Logger
}

func (o *option) apply(opts ...Option) *option {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithLogger(logger Logger) Option {
	return func(o *option) {
		o.logger = logger
	}
}

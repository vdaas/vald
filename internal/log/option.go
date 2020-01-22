package log

type Option func(*option)

type option struct {
	logger Logger
}

func (o *option) evaluateOption(opts ...Option) *option {
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

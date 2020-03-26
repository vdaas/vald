package benchmark

type Option func(*benchmark)

var (
	defaultOptions = []Option{}
)

func WithName(name string) Option {
	return func(b *benchmark) {
		if len(name) != 0 {
			b.name = name
		}
	}
}

func WithStrategy(strategies ...Strategy) Option {
	return func(b *benchmark) {
		if len(strategies) != 0 {
			b.strategies = strategies
		}
	}
}

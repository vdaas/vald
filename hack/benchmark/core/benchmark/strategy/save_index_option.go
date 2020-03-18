package strategy

type SaveIndexOption func(*saveIndex)

var (
	defaultSaveIndexOptions = []SaveIndexOption{
		WithSaveIndexPoolSize(10000),
		WithSaveIndexPreStart(
			(new(preStart)).Func,
		),
	}
)

func WithSaveIndexPoolSize(size int) SaveIndexOption {
	return func(s *saveIndex) {
		if size > 0 {
			s.poolSize = uint32(size)
		}
	}
}

func WithSaveIndexPreStart(fn PreStart) SaveIndexOption {
	return func(s *saveIndex) {
		if fn != nil {
			s.preStart = fn
		}
	}
}

package strategy

type SearchOption func(*search)

var (
	defaultSearchOptions = []SearchOption{
		WithSearchSize(10),
		WithSearchEpsilon(0.01),
		WithSearchRadius(-1),
		WithSearchPreStart(
			(new(defaultPreStart)).PreStart,
		),
	}
)

func WithSearchSize(size int) SearchOption {
	return func(s *search) {
		if size > 0 {
			s.size = size
		}
	}
}

func WithSearchEpsilon(epsilon float32) SearchOption {
	return func(s *search) {
		s.epsilon = epsilon
	}
}

func WithSearchRadius(radius float32) SearchOption {
	return func(s *search) {
		s.radius = radius
	}
}

func WithSearchPreStart(fn PreStart) SearchOption {
	return func(s *search) {
		if s.preStart != nil {
			s.preStart = fn
		}
	}
}

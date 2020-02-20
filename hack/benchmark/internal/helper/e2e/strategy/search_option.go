package strategy

type SearchOption func(*search)

var (
	defaultSearchOptions = []SearchOption{}
)

func WithSearchParallel() SearchOption {
	return func(s *search) {
		s.parallel = true
	}
}

func WithSearchSize(size int) SearchOption {
	return func(s *search) {
		if size > 0 {
			s.size = uint32(size)
		}
	}
}

func WithSearchEpsilon(e float32) SearchOption {
	return func(s *search) {
		if e > 0.0 {
			s.epsilon = e
		}
	}
}

func WithSearchRadius(r float32) SearchOption {
	return func(s *search) {
		if r > 0.0 {
			s.radius = r
		}
	}
}

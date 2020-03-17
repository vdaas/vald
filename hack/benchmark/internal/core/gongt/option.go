package gongt

// Option is NGT confiure.
type Option func(n *NGT)

var (
	defaultOptions = []Option{
		WithIndexPath("/tmp/gongt"),
	}
)

// WithIndexPath returns Option that set indexPath.
func WithIndexPath(path string) Option {
	return func(n *NGT) {
		if len(path) != 0 {
			n.indexPath = path
		}
	}
}

package gongt

type Option func(n *ngt)

var (
	defaultOptions = []Option{
		WithIndexPath("/tmp/gongt"),
		WithObjectType(Float),
		WithDimension(128),
	}
)

func WithIndexPath(path string) Option {
	return func(n *ngt) {
		if len(path) != 0 {
			n.indexPath = path
		}
	}
}

func WithObjectType(typ ObjectType) Option {
	return func(n *ngt) {
		switch typ {
		case Uint8, Float:
			n.objectType = typ
		default:
			n.objectType = ObjectNone
		}
	}
}

func WithDimension(dimension int) Option {
	return func(n *ngt) {
		if dimension > 0 {
			n.dimension = dimension
		}
	}
}

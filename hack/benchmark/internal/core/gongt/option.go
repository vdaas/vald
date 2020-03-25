package gongt

type Option func(c *core)

var (
	defaultOptions = []Option{
		WithIndexPath("/tmp/gongt"),
		WithObjectType(Float),
		WithDimension(128),
	}
)

func WithIndexPath(path string) Option {
	return func(c *core) {
		if len(path) != 0 {
			c.indexPath = path
		}
	}
}

func WithObjectType(typ ObjectType) Option {
	return func(c *core) {
		switch typ {
		case Uint8, Float:
			c.objectType = typ
		default:
			c.objectType = ObjectNone
		}
	}
}

func WithDimension(dimension int) Option {
	return func(c *core) {
		if dimension > 0 {
			c.dimension = dimension
		}
	}
}

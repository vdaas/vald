package ngt

type Option func(*core)

var (
	defaultOptions = []Option{
		WithIndexPath("/tmp/ngt"),
		WithObjectType(Float),
		WithDimension(128),
	}
)

func WithIndexPath(path string) Option {
	return func(c *core) {
		if len(path) != 0 {
			c.idxPath = path
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

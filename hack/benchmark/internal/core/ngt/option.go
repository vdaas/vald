package ngt

type Option func(*core)

var (
	defaultOptions = []Option{
		WithIndexPath("tmpdir"),
		WithObjectType("float"),
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

func WithObjectType(typ string) Option {
	return func(c *core) {
		switch typ {
		case "uint8":
			c.objectType = Uint8
		case "float":
			c.objectType = Float
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

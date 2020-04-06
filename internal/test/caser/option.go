package caser

type Option func(*caser)

var (
	defaultOptions = []Option{
		WithCheckFunc(func() error {
			return nil
		}),
	}
)

func WithName(str string) Option {
	return func(c *caser) {
		if len(str) != 0 {
			c.name = str
		}
	}
}

func WithArg(args ...interface{}) Option {
	return func(c *caser) {
		if len(args) != 0 {
			c.args = args
		}
	}
}

func WithField(fields ...interface{}) Option {
	return func(c *caser) {
		if len(fields) != 0 {
			c.fields = fields
		}
	}
}

func WithWant(wants []interface{}) Option {
	return func(c *caser) {
		if len(wants) != 0 {
			c.wants = wants
		}
	}
}

func WithCheckFunc(fn func() error) Option {
	return func(c *caser) {
		if fn != nil {
			c.checkFunc = fn
		}
	}
}

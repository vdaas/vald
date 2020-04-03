package caser

type Option func(*caser)

var (
	defaultOptions = []Option{
		WithCheck(func(...interface{}) error {
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

func WithCheck(fn func(gots ...interface{}) error) Option {
	return func(c *caser) {
		if fn != nil {
			c.checkFunc = fn
		}
	}
}

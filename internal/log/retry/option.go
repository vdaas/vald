package retry

type Option func(r *retry)

var (
	defaultOption = []Option{
		WithError(nopFunc),
		WithWarn(nopFunc),
	}

	nopFunc = func(vals ...interface{}) {}
)

func WithError(fn func(vals ...interface{})) Option {
	return func(r *retry) {
		if fn == nil {
			return
		}
		r.errorfn = fn
	}
}

func WithWarn(fn func(vals ...interface{})) Option {
	return func(r *retry) {
		if fn == nil {
			return
		}
		r.errorfn = fn
	}
}

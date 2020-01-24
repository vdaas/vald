package retry

type nopRetry struct{}

func NewNop() Retry {
	return new(nopRetry)
}

func (nr *nopRetry) Out(
	fn func(vals ...interface{}) error,
	vals ...interface{},
) {
	fn(vals...)
}

func (nr *nopRetry) Outf(
	fn func(format string, vals ...interface{}) error,
	format string, vals ...interface{},
) {
	fn(format, vals...)
}

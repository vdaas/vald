package mock

type Retry struct {
	OutFunc func(
		fn func(vals ...interface{}) error,
		vals ...interface{},
	)

	OutfFunc func(
		fn func(format string, vals ...interface{}) error,
		format string,
		vals ...interface{},
	)
}

func (r *Retry) Out(
	fn func(vals ...interface{}) error,
	vals ...interface{},
) {
	r.OutFunc(fn, vals...)
}

func (r *Retry) Outf(
	fn func(format string, vals ...interface{}) error,
	format string, vals ...interface{},
) {
	r.OutfFunc(fn, format, vals...)
}

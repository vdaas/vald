package retry

type (
	Out  func(fn func(vals ...interface{}) error, vals ...interface{})
	Outf func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{})
)

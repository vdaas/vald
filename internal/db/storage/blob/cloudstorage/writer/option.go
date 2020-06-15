package writer

// Option configures *writer.
type Option func(w *writer) error

var (
	defaultOpts = []Option{}
)

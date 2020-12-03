package info

type Option func(*Detail) error

var (
	defaultOpts = []Option{}
)

func WithServerName(s string) Option {
	return func(d *Detail) error {
		d.ServerName = s
		return nil
	}
}

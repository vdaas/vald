package strategy

type RemoveOption func(*remove)

var (
	defaultRemoveOptions = []RemoveOption{
		WithPreStart(
			(new(defaultPreStart)).PreStart,
		),
	}
)

func WithPreStart(fn PreStart) RemoveOption {
	return func(d *remove) {
		if fn != nil {
			d.preStart = fn
		}
	}
}

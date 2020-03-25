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
	return func(r *remove) {
		if fn != nil {
			r.preStart = fn
		}
	}
}

package strategy

type RemoveOption func(*remove)

var (
	defaultRemoveOptions = []RemoveOption{
		WithRemovePreStart(
			(new(preStart)).Func,
		),
	}
)

func WithRemovePreStart(fn PreStart) RemoveOption {
	return func(r *remove) {
		if fn != nil {
			r.preStart = fn
		}
	}
}

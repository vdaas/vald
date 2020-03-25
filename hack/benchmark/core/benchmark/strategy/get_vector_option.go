package strategy

type GetVectorOption func(*getVector)

var (
	defaultGetVectorOptions = []GetVectorOption{
		WithGetVectorPreStart(
			(new(defaultPreStart)).PreStart,
		),
	}
)

func WithGetVectorPreStart(fn PreStart) GetVectorOption {
	return func(g *getVector) {
		if g.preStart != nil {
			g.preStart = fn
		}
	}
}

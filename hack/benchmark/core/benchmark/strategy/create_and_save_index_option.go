package strategy

type CreateAndSaveIndexOption func(*createAndSaveIndex)

var (
	defaultCreateAndSaveIndexOptions = []CreateAndSaveIndexOption{
		WithCreateAndSaveIndexPoolSize(10000),
		WithCreateAndSaveIndexPreStart(
			(new(preStart)).Func,
		),
	}
)

func WithCreateAndSaveIndexPoolSize(size int) CreateAndSaveIndexOption {
	return func(c *createAndSaveIndex) {
		if size > 0 {
			c.poolSize = uint32(size)
		}
	}
}

func WithCreateAndSaveIndexPreStart(fn PreStart) CreateAndSaveIndexOption {
	return func(c *createAndSaveIndex) {
		if fn != nil {
			c.preStart = fn
		}
	}
}

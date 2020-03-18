package strategy

type CreateIndexOption func(*createIndex)

var (
	defaultCreateIndexOptions = []CreateIndexOption{
		WithCreateIndexPoolSize(10000),
		WithCreateIndexPreStart(
			(new(preStart)).Func,
		),
	}
)

func WithCreateIndexPoolSize(poolSize int) CreateIndexOption {
	return func(c *createIndex) {
		if poolSize > 0 {
			c.poolSize = uint32(poolSize)
		}
	}
}

func WithCreateIndexPreStart(fn PreStart) CreateIndexOption {
	return func(c *createIndex) {
		if fn != nil {
			c.preStart = fn
		}
	}
}

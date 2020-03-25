package strategy

type CreateIndexOption func(*createIndex)

var (
	defaultCreateIndexOptions = []CreateIndexOption{
		WithCreateIndexPoolSize(10000),
		WithCreateIndexPreStart(
			(new(defaultInsert)).PreStart,
		),
	}
)

func WithCreateIndexPreStart(fn PreStart) CreateIndexOption {
	return func(ci *createIndex) {
		if ci.preStart != nil {
			ci.preStart = fn
		}
	}
}

func WithCreateIndexPoolSize(size int) CreateIndexOption {
	return func(ci *createIndex) {
		if size > 0 {
			ci.poolSize = uint32(size)
		}
	}
}

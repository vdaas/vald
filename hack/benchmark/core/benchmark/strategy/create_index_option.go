package strategy

type CreateIndexOption func(*createIndex)

var (
	defaultCreateIndexOptions = []CreateIndexOption{}
)

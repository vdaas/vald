package strategy

type SaveIndexOption func(*saveIndex)

var (
	defaultSaveIndexOptions = []SaveIndexOption{}
)

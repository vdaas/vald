package strategy

type CreateAndSaveIndexOption func(*createAndSaveIndex)

var (
	defaultCreateAndSaveIndexOptions = []CreateAndSaveIndexOption{}
)

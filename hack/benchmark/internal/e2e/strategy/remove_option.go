package strategy

type RemoveOption func(*remove)

var (
	defaultRemoveOptions = []RemoveOption{}
)

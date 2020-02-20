package strategy

type StreamRemoveOption func(*streamRemove)

var (
	defaultStreamRemoveOptions = []StreamRemoveOption{}
)

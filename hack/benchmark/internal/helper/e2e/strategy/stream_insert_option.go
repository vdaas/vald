package strategy

type StreamInsertOption func(*streamInsert)

var (
	defaultStreamInsertOptions = []StreamInsertOption{}
)

package strategy

type InsertOption func(*insert)

var (
	defaultInsertOptions = []InsertOption{}
)

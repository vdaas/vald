package stratedy

type InsertOption func(*insert)

var (
	defaultInsertOptions = []InsertOption{}
)

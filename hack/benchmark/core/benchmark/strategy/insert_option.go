package strategy

// InsertOption is insert strategy configure.
type InsertOption func(*insert)

var (
	defaultInsertOption = []InsertOption{}
)

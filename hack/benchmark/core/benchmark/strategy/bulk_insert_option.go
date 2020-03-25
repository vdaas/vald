package strategy

type BulkInsertOption func(*bulkInsert)

var (
	defaultBulkInsertOptions = []BulkInsertOption{}
)

package strategy

type BulkRemoveOption func(*bulkRemove)

var (
	defaultBulkRemoveOptions = []BulkRemoveOption{}
)

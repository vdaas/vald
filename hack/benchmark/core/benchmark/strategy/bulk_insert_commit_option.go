package strategy

type BulkInsertCommitOption func(*bulkInsertCommit)

var (
	defaultBulkInsertCommitOptions = []BulkInsertCommitOption{
		WithBulkInsertCommitPoolSize(10000),
	}
)

func WithBulkInsertCommitPoolSize(size int) BulkInsertCommitOption {
	return func(bi *bulkInsertCommit) {
		if size > 0 {
			bi.poolSize = uint32(size)
		}
	}
}

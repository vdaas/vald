package strategy

type BulkInsertCommitOption func(*bulkInsertCommit)

var (
	defaultBulkInsertCommitOptions = []BulkInsertCommitOption{
		WithBulkInsertCommitPoolSize(10000),
	}
)

func WithBulkInsertCommitPoolSize(size int) BulkInsertCommitOption {
	return func(bic *bulkInsertCommit) {
		if size > 0 {
			bic.poolSize = uint32(size)
		}
	}
}

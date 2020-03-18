package strategy

type BulkInsertCommitOption func(*bulkInsertCommit)

var (
	defaultBulkInsertCommitOptions = []BulkInsertCommitOption{
		WithBulkInsertCommitPoolSize(10000),
	}
)

func WithBulkInsertCommitPoolSize(poolSize int) BulkInsertCommitOption {
	return func(bc *bulkInsertCommit) {
		if poolSize > 0 {
			bc.poolSize = uint32(poolSize)
		}
	}
}

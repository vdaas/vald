package strategy

type InsertCommitOption func(*insertCommit)

var (
	defaultInsertCommitOptions = []InsertCommitOption{
		WithInsertCommitPoolSize(10000),
	}
)

func WithInsertCommitPoolSize(size int) InsertCommitOption {
	return func(ic *insertCommit) {
		if size > 0 {
			ic.poolSize = uint32(size)
		}
	}
}

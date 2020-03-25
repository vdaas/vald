package strategy

type InsertCommitOption func(*insertCommit)

var (
	defaultInsertCommitOptions = []InsertCommitOption{
		WithInsertCommitPoolSize(10000),
	}
)

func WithInsertCommitPoolSize(size int) InsertCommitOption {
	return func(isrt *insertCommit) {
		if size > 0 {
			isrt.poolSize = uint32(size)
		}
	}
}

package strategy

type BulkRemoveOption func(*bulkRemove)

var (
	defaultBulkRemoveOptions = []BulkRemoveOption{
		WithBulkRemoveChunkSize(1000),
		WithBulkRemovePrestart(
			(new(preStart)).Func,
		),
	}
)

func WithBulkRemoveChunkSize(chunk int) BulkRemoveOption {
	return func(br *bulkRemove) {
		if chunk > 0 {
			br.chunkSize = chunk
		}
	}
}

func WithBulkRemovePrestart(
	fn PreStart,
) BulkRemoveOption {
	return func(br *bulkRemove) {
		if fn != nil {
			br.fn = fn
		}
	}
}

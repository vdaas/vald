package deleter

import "github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"

type Deleter interface {
	Delete(key string)
}

type deleter struct {
	service s3iface.S3API
	bucket  string
}

func New(opts ...Option) (Deleter, error) {
	return nil, nil
}

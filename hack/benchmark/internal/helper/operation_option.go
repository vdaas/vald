package helper

import "github.com/vdaas/vald/hack/benchmark/internal/assets"

type OperationHelperOption func(*operationHelper)

var (
	defaultOperationHelperOption = []OperationHelperOption{}
)

// func WithNGT() OperationHelperOption {
// 	return func(o *operationHelper) {}
// }

func WithNGTInitializer(fn func()) OperationHelperOption {
	return func(o *operationHelper) {}
}

// func WithGoNGT() OperationHelperOption {
// 	return func(o *operationHelper) {}
// }

func WithGoNGTInitializer(fn func()) OperationHelperOption {
	return func(o *operationHelper) {}
}

func WithDataset(dataSet assets.Dataset) OperationHelperOption {
	return func(o *operationHelper) {}
}

package strategy

type GetVectorOption func(*getVector)

var (
	defaultGetVectorOptions = []GetVectorOption{}
)

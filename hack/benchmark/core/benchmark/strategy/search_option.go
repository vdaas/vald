package strategy

type SearchOption func(*search)

var (
	defaultSearchOptions = []SearchOption{}
)

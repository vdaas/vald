package info

type Option func(*Detail) error

var (
	defaultOpts = []Option{}
)

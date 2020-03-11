package vald

type Option func(*server)

var (
	defaultOptions = []Option{}
)

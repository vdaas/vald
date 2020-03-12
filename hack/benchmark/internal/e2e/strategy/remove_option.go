// Package strategy provides strategy for e2e testing functions
package strategy

type RemoveOption func(*remove)

var (
	defaultRemoveOptions = []RemoveOption{}
)

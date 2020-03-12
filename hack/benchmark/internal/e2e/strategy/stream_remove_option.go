// Package strategy provides strategy for e2e testing functions
package strategy

type StreamRemoveOption func(*streamRemove)

var (
	defaultStreamRemoveOptions = []StreamRemoveOption{}
)

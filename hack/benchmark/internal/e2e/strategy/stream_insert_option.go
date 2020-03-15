// Package strategy provides strategy for e2e testing functions
package strategy

type StreamInsertOption func(*streamInsert)

var (
	defaultStreamInsertOptions = []StreamInsertOption{}
)

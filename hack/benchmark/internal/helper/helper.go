package helper

import "testing"

type Helper interface {
	Run(b *testing.B) error
}

type helper struct {
	parallel bool
}

func NewHelper(opts ...HelperOption) Helper {
	h := new(helper)

	for _, opt := range append(defaultHelperOption, opts...) {
		opt(h)
	}

	return h
}

func (h *helper) Run(b *testing.B) error {
	if h.parallel {
		return h.runParallel(b)
	}
	return h.run(b)
}

func (h *helper) runParallel(b *testing.B) error {
	return nil
}

func (h *helper) run(b *testing.B) error {
	return nil
}

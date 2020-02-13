package helper

import (
	"io/ioutil"
	"os"
	"testing"
)

type Helper interface {
	Run(b *testing.B) error
}

type helper struct {
	parallel  bool
	targets   []string
	operation OperationHelper
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

func (h *helper) run(b *testing.B) error {
	for _, target := range h.targets {
		tmpdir, err := ioutil.TempDir("", "tmpdir")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpdir)

		h.operation.Insert()(b)
		h.operation.CreateIndex()(b)
		h.operation.Search()(b)

		// TODO
		_ = target
	}
	return nil
}

func (h *helper) runParallel(b *testing.B) error {
	for _, target := range h.targets {
		tmpdir, err := ioutil.TempDir("", "tmpdir")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpdir)

		h.operation.InsertParallel()(b)
		h.operation.CreateIndexParallel()(b)
		h.operation.SearchParallel()(b)

		// TODO
		_ = target
	}
	return nil
}

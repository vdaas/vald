package ngtd

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/e2e/internal/starter"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd"
	"github.com/yahoojapan/ngtd/kvs"
)

type ServerType = ngtd.ServerType

type server struct {
	dim     int
	srvType ServerType
	baseDir string
	port    int
}

func New(opts ...Option) starter.Starter {
	srv := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(srv)
	}

	return srv
}

func (ns *server) Run(ctx context.Context, tb testing.TB) func() {
	tb.Helper()

	if err := ns.createIndexDir(); err != nil {
		tb.Error(err)
	}

	gongt.SetDimension(ns.dim)

	db, err := kvs.NewGoLevel(ns.baseDir + "meta")
	if err != nil {
		tb.Error(err)
	}

	n, err := ngtd.NewNGTD(ns.baseDir+"ngt", db, ns.port)
	if err != nil {
		tb.Error(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		wg.Done()

		if err := n.ListenAndServe(ns.srvType); err != nil {
			tb.Errorf("ngtd returned error: %s", err.Error())
		}
	}()

	wg.Wait()

	return func() {
		n.Stop()

		if err := ns.clearIndexDir(); err != nil {
			tb.Error(err)
		}
	}
}

func (ns *server) createIndexDir() error {
	if err := os.RemoveAll(ns.baseDir); err != nil {
		return err
	}

	if err := os.MkdirAll(ns.baseDir, 0755); err != nil {
		return err
	}

	return nil
}

func (ns *server) clearIndexDir() error {
	if err := os.RemoveAll(ns.baseDir + "meta"); err != nil {
		return err
	}

	if err := os.RemoveAll(ns.baseDir + "ngt"); err != nil {
		return err
	}

	return nil
}

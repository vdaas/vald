package ngtd

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/vdaas/vald/hack/benchmark/internal/starter"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd"
	"github.com/yahoojapan/ngtd/kvs"
)

type ServerType = ngtd.ServerType

const (
	HTTP ServerType = 1
	GRPC ServerType = 2
)

type server struct {
	dim      int
	srvType  ServerType
	indexDir string
	port     int
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

	db, err := kvs.NewGoLevel(ns.indexDir + "meta")
	if err != nil {
		tb.Error(err)
	}

	n, err := ngtd.NewNGTD(ns.indexDir+"ngt", db, ns.port)
	if err != nil {
		tb.Error(err)
	}

	go func() {
		if err := n.ListenAndServe(ns.srvType); err != nil {
			tb.Errorf("ngtd returned error: %s", err.Error())
		}
	}()

	time.Sleep(4 * time.Second)

	return func() {
		n.Stop()

		if err := ns.clearIndexDir(); err != nil {
			tb.Error(err)
		}
	}
}

func (ns *server) createIndexDir() error {
	if err := os.RemoveAll(ns.indexDir); err != nil {
		return err
	}

	if err := os.MkdirAll(ns.indexDir, 0755); err != nil {
		return err
	}

	return nil
}

func (ns *server) clearIndexDir() error {
	if err := os.RemoveAll(ns.indexDir + "meta"); err != nil {
		return err
	}

	if err := os.RemoveAll(ns.indexDir + "ngt"); err != nil {
		return err
	}

	return nil
}

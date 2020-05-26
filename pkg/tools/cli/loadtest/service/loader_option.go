package service

import (
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"time"
)

type Option func(*loader) error

var (
	defaultOpts = []Option{
		WithConcurrency(100),
		WithErrGroup(errgroup.Get()),
		WithProgressDuration(5 * time.Second),
	}
)

func WithAddr(a string) Option {
	return func(l *loader) error {
		l.addr = a
		return nil
	}
}

func WithClient(c grpc.Client) Option {
	return func(l *loader) error {
		if c != nil {
			l.client = c
			return nil
		}
		return errors.Errorf("client must not be nil")
	}
}

func WithConcurrency(c int) Option {
	return func(l *loader) error {
		if c > 0 {
			l.concurrency = c
		}
		return nil
	}
}

func WithDataset(n string) Option {
	return func(l *loader) error {
		l.dataset = n
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(l *loader) error {
		if eg != nil {
			l.eg = eg
		}
		return nil
	}
}

func WithProgressDuration(pd time.Duration) Option {
	return func(l *loader) error {
		l.progressDuration = pd
		return nil
	}
}

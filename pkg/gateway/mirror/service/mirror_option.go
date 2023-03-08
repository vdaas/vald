package service

import (
	"time"

	"github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

type MirrorOption func(m *mirr) error

var defaultMirrOpts = []MirrorOption{
	WithAdvertiseInterval("1s"),
}

func WithErrorGroup(eg errgroup.Group) MirrorOption {
	return func(m *mirr) error {
		if eg != nil {
			m.eg = eg
		}
		return nil
	}
}

func WithValdAddrs(addrs ...string) MirrorOption {
	return func(m *mirr) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("lbAddrs", addrs)
		}
		if m.gwAddrs == nil {
			m.gwAddrs = make([]string, 0, len(addrs))
		}
		m.gwAddrs = append(m.gwAddrs, addrs...)
		return nil
	}
}

func WithSelfMirrorAddrs(addrs ...string) MirrorOption {
	return func(m *mirr) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("selfMirrorAddrs", addrs)
		}
		if m.selfMirrAddrs == nil {
			m.selfMirrAddrs = make([]string, 0, len(addrs))
		}
		m.selfMirrAddrs = append(m.selfMirrAddrs, addrs...)
		return nil
	}
}

func WithMirror(c mirror.Client) MirrorOption {
	return func(m *mirr) error {
		if c != nil {
			m.client = c
		}
		return nil
	}
}

func WithAdvertiseInterval(s string) MirrorOption {
	return func(m *mirr) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("advertiseInterval", s)
		}
		dur, err := time.ParseDuration(s)
		if err != nil {
			return errors.NewErrInvalidOption("advertiseInterval", s, err)
		}
		m.advertiseDur = dur
		return nil
	}
}

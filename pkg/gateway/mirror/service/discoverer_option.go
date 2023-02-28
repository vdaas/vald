package service

import (
	"time"

	"github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

type DiscovererOption func(d *discoverer) error

var defaultMirrOpts = []DiscovererOption{
	WithAdvertiseInterval("1s"),
}

func WithErrorGroup(eg errgroup.Group) DiscovererOption {
	return func(d *discoverer) error {
		if eg != nil {
			d.eg = eg
		}
		return nil
	}
}

func WithValdAddrs(addrs ...string) DiscovererOption {
	return func(d *discoverer) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("lbAddrs", addrs)
		}
		if d.gwAddrs == nil {
			d.gwAddrs = make([]string, 0, len(addrs))
		}
		d.gwAddrs = append(d.gwAddrs, addrs...)
		return nil
	}
}

func WithSelfMirrorAddrs(addrs ...string) DiscovererOption {
	return func(d *discoverer) error {
		if len(addrs) == 0 {
			return errors.NewErrCriticalOption("selfMirrorAddrs", addrs)
		}
		if d.selfMirrAddrs == nil {
			d.selfMirrAddrs = make([]string, 0, len(addrs))
		}
		d.selfMirrAddrs = append(d.selfMirrAddrs, addrs...)
		return nil
	}
}

func WithDiscoverer(c mirror.Client) DiscovererOption {
	return func(d *discoverer) error {
		if c != nil {
			d.client = c
		}
		return nil
	}
}

func WithAdvertiseInterval(s string) DiscovererOption {
	return func(d *discoverer) error {
		if len(s) == 0 {
			return errors.NewErrInvalidOption("advertiseInterval", s)
		}
		dur, err := time.ParseDuration(s)
		if err != nil {
			return errors.NewErrInvalidOption("advertiseInterval", s, err)
		}
		d.advertiseDur = dur
		return nil
	}
}

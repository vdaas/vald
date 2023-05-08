package service

import (
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/timeutil"
)

type DiscovererOption func(d *discoverer) error

var (
	defaultDiscovererOpts = []DiscovererOption{
		WithDiscovererDuration("1s"),
		WithDiscovererErrGroup(errgroup.Get()),
		WithDiscovererColocation("dc1"),
	}
)

func WithDiscovererMirror(m Mirror) DiscovererOption {
	return func(d *discoverer) error {
		if m == nil {
			return errors.NewErrCriticalOption("discovererMirror", m)
		}
		d.mirr = m
		return nil
	}
}

func WithDiscovererDialer(der net.Dialer) DiscovererOption {
	return func(d *discoverer) error {
		if der != nil {
			d.der = der
		}
		return nil
	}
}

func WithDiscovererNamespace(ns string) DiscovererOption {
	return func(d *discoverer) error {
		if len(ns) != 0 {
			d.namespace = ns
		}
		return nil
	}
}

func WithDiscovererColocation(loc string) DiscovererOption {
	return func(d *discoverer) error {
		if len(loc) != 0 {
			d.colocation = loc
		}
		return nil
	}
}

func WithDiscovererDuration(s string) DiscovererOption {
	return func(d *discoverer) error {
		if s == "" {
			return nil
		}
		dur, err := timeutil.Parse(s)
		if err != nil {
			return errors.NewErrInvalidOption("discovererDuration", s, err)
		}
		d.dur = dur
		return nil
	}
}

func WithDiscovererErrGroup(eg errgroup.Group) DiscovererOption {
	return func(d *discoverer) error {
		if eg != nil {
			d.eg = eg
		}
		return nil
	}
}

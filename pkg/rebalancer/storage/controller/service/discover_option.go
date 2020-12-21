package service

import (
	"time"

	"github.com/vdaas/vald/internal/errgroup"
)

type DiscovererOption func(d *discoverer) error

var (
	defaultDiscovererOpts = []DiscovererOption{}
)

func WithJobName(name string) DiscovererOption {
	return func(d *discoverer) error {
		d.jobName = name
		return nil
	}
}

func WithJobNamespace(ns string) DiscovererOption {
	return func(d *discoverer) error {
		d.jobNamespace = ns
		return nil
	}
}

func WithJobTemplateKey(k string) DiscovererOption {
	return func(d *discoverer) error {
		d.jobTemplateKey = k
		return nil
	}
}

func WithAgentName(an string) DiscovererOption {
	return func(d *discoverer) error {
		d.agentName = an
		return nil
	}
}

func WithAgentNamespace(ans string) DiscovererOption {
	return func(d *discoverer) error {
		d.agentNamespace = ans
		return nil
	}
}

func WithAgentResourceType(art string) DiscovererOption {
	return func(d *discoverer) error {
		d.agentResourceType = art
		return nil
	}
}

func WithConfigMapName(n string) DiscovererOption {
	return func(d *discoverer) error {
		d.configmapName = n
		return nil
	}
}

func WithConfigMapNamespace(ns string) DiscovererOption {
	return func(d *discoverer) error {
		d.configmapNamespace = ns
		return nil
	}
}

func WithReconcileCheckDuration(t string) DiscovererOption {
	return func(d *discoverer) error {
		rcd, err := time.ParseDuration(t)
		if err != nil {
			return err
		}
		d.rcd = rcd
		return nil
	}
}

func WithTolerance(t float64) DiscovererOption {
	return func(d *discoverer) error {
		d.tolerance = t
		return nil
	}
}

func WithRateThreshold(t float64) DiscovererOption {
	return func(d *discoverer) error {
		d.rateThreshold = t
		return nil
	}
}

func WithErrorGroup(eg errgroup.Group) DiscovererOption {
	return func(d *discoverer) error {
		d.eg = eg
		return nil
	}
}

func WithLeaderElectionID(id string) DiscovererOption {
	return func(d *discoverer) error {
		d.leaderElectionID = id
		return nil
	}
}

package service

import (
	"time"

	"github.com/vdaas/vald/internal/errgroup"
)

type RebalancerOption func(r *rebalancer) error

var (
	defaultRebalancerOpts = []RebalancerOption{}
)

func WithJobName(name string) RebalancerOption {
	return func(r *rebalancer) error {
		r.jobName = name
		return nil
	}
}

func WithJobNamespace(ns string) RebalancerOption {
	return func(r *rebalancer) error {
		r.jobNamespace = ns
		return nil
	}
}

func WithJobTemplateKey(k string) RebalancerOption {
	return func(r *rebalancer) error {
		r.jobTemplateKey = k
		return nil
	}
}

func WithAgentName(an string) RebalancerOption {
	return func(r *rebalancer) error {
		r.agentName = an
		return nil
	}
}

func WithAgentNamespace(ans string) RebalancerOption {
	return func(r *rebalancer) error {
		r.agentNamespace = ans
		return nil
	}
}

func WithAgentResourceType(art string) RebalancerOption {
	return func(r *rebalancer) error {
		r.agentResourceType = art
		return nil
	}
}

func WithConfigMapName(n string) RebalancerOption {
	return func(r *rebalancer) error {
		r.configmapName = n
		return nil
	}
}

func WithConfigMapNamespace(ns string) RebalancerOption {
	return func(r *rebalancer) error {
		r.configmapNamespace = ns
		return nil
	}
}

func WithReconcileCheckDuration(t string) RebalancerOption {
	return func(r *rebalancer) error {
		rcd, err := time.ParseDuration(t)
		if err != nil {
			return err
		}
		r.rcd = rcd
		return nil
	}
}

func WithTolerance(t float64) RebalancerOption {
	return func(r *rebalancer) error {
		r.tolerance = t
		return nil
	}
}

func WithRateThreshold(t float64) RebalancerOption {
	return func(r *rebalancer) error {
		r.rateThreshold = t
		return nil
	}
}

func WithErrorGroup(eg errgroup.Group) RebalancerOption {
	return func(r *rebalancer) error {
		r.eg = eg
		return nil
	}
}

func WithLeaderElectionID(id string) RebalancerOption {
	return func(r *rebalancer) error {
		r.leaderElectionID = id
		return nil
	}
}

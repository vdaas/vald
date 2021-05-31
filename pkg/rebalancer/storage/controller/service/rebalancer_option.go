//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package service

import (
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/rebalancer/storage/controller/config"
)

type RebalancerOption func(r *rebalancer) error

var (
	defaultRebalancerOpts = []RebalancerOption{}
)

func WithPodName(name string) RebalancerOption {
	return func(r *rebalancer) error {
		r.podName = name
		return nil
	}
}

func WithPodNamespace(ns string) RebalancerOption {
	return func(r *rebalancer) error {
		r.podNamespace = ns
		return nil
	}
}

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

func WithJobTemplate(tpl string) RebalancerOption {
	return func(r *rebalancer) error {
		r.jobTemplate = tpl
		return nil
	}
}

func WithAgentName(an string) RebalancerOption {
	return func(r *rebalancer) error {
		r.agentName = an
		return nil
	}
}

func WithAgentPort(port int) RebalancerOption {
	return func(r *rebalancer) error {
		r.agentPort = port
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
		r.agentResourceType = config.AToAgentResourceType(art)
		if r.agentResourceType == config.UNKNOWN_RESOURCE_TYPE {
			return errors.NewErrCriticalOption("agentResourceType", art)
		}
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

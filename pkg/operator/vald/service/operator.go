// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
	"github.com/vdaas/vald/pkg/operator/vald/config"
)

const (
	apiName = "vald/operator/vald"

	// ctrlErrChanBuf is the buffer size for the error channel returned by Start:
	// one slot for ctrl errors and one for context-cancellation drain.
	ctrlErrChanBuf = 2
)

type Operator interface {
	Start(ctx context.Context) (<-chan error, error)
}

type operator struct {
	ctrl k8s.Controller
	eg   errgroup.Group
}

func New(cfg *config.Operator, opts ...Option) (o Operator, err error) {
	if cfg == nil {
		return nil, errors.ErrInvalidConfig
	}
	operator := new(operator)
	for _, opt := range append(defaultOpts, opts...) {
		if oerr := opt(operator); oerr != nil {
			oerr2 := errors.ErrOptionFailed(oerr, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr2, &e) {
				log.Error(oerr)
				return nil, oerr2
			}
			log.Warn(oerr2)
		}
	}

	// Load binds cfg first, so every section (Controller, Vrs, ...) is
	// guaranteed to be non-nil afterwards.
	rcfg, err := cfg.Load()
	if err != nil {
		return nil, err
	}

	ctrlCfg := cfg.Controller
	le := ctrlCfg.LeaderElection
	lease, err := timeutil.Parse(le.LeaseDuration)
	if err != nil {
		return nil, errors.Wrap(err, "invalid controller.leader_election.lease_duration")
	}
	renew, err := timeutil.Parse(le.RenewDeadline)
	if err != nil {
		return nil, errors.Wrap(err, "invalid controller.leader_election.renew_deadline")
	}
	retry, err := timeutil.Parse(le.RetryPeriod)
	if err != nil {
		return nil, errors.Wrap(err, "invalid controller.leader_election.retry_period")
	}
	syncPeriod, err := timeutil.Parse(ctrlCfg.SyncPeriod)
	if err != nil {
		return nil, errors.Wrap(err, "invalid controller.sync_period")
	}

	operator.ctrl, err = k8s.New(
		k8s.WithControllerName(cfg.Name),
		k8s.WithResourceController(newResourceController(rcfg)),
		k8s.WithLeaderElection(le.Enabled, le.ID, le.Namespace),
		k8s.WithLeaderElectionDetail(lease, renew, retry),
		k8s.WithMetricsAddress(ctrlCfg.MetricsAddress),
		k8s.WithSyncPeriod(syncPeriod),
		k8s.WithCacheNamespaces(ctrlCfg.CacheNamespaces...),
	)
	if err != nil {
		return nil, err
	}

	return operator, nil
}

func (o *operator) Start(ctx context.Context) (<-chan error, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/service/operator.Start")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	dech, err := o.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, ctrlErrChanBuf)
	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-dech:
				if err != nil {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ech <- err:
					}
				}
			}
		}
	}))

	return ech, nil
}

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

// Package service manages the main logic of server.
package service

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/manager/replication/agent"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
)

type Replicator interface {
	Start(context.Context) (<-chan error, error)
	GetCurrentPodIPs() ([]string, bool)
}

type replicator struct {
	pods      atomic.Value
	ctrl      k8s.Controller
	namespace string
	name      string
	eg        errgroup.Group
	rdur      time.Duration
	rpods     sync.Map
	client    grpc.Client
}

func New(opts ...Option) (rp Replicator, err error) {
	r := new(replicator)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	r.pods.Store(make([]string, 0))

	r.ctrl, err = k8s.New(
		k8s.WithControllerName("vald k8s replication manager controller"),
		k8s.WithEnableLeaderElection(),
		k8s.WithResourceController(pod.New(
			pod.WithControllerName("pod discoverer"),
			pod.WithOnErrorFunc(func(err error) {
				log.Error("failed to reconcile:", err)
			}),
			pod.WithOnReconcileFunc(func(podList map[string][]pod.Pod) {
				var pods map[string]pod.Pod

				currentPods, ok := r.pods.Load().(map[string]pod.Pod)
				if ok {
					pods = make(map[string]pod.Pod, len(currentPods))
				} else {
					pods = make(map[string]pod.Pod)
				}

				for name, ps := range podList {
					if name == r.name {
						for _, p := range ps {
							if p.Namespace == r.namespace {
								pods[p.Name] = p
							}
						}
					}
				}

				r.pods.Store(pods)

				for name, cpod := range currentPods {
					p, ok := pods[name]
					if !ok ||
						p.Name != cpod.Name ||
						p.Namespace != cpod.Namespace ||
						p.IP != cpod.IP ||
						p.NodeName != cpod.NodeName {
						if _, ok := r.rpods.Load(name); !ok {
							r.rpods.Store(name, cpod)
						}
					}
				}

				log.Debugf("pod resource reconciled\t%#v", podList)
			}),
		)),
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *replicator) Start(ctx context.Context) (<-chan error, error) {
	rech, err := r.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 2)
	r.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		rt := time.NewTicker(r.rdur)
		defer rt.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-rt.C:
				err = r.SendRecoveryRequest(ctx)
				if err != nil {
					log.Error("failed to send recovery request", err)
					ech <- err
				}
			case err = <-rech:
				if err != nil {
					ech <- err
				}
			}
		}
	}))
	return ech, nil
}

func (r *replicator) GetCurrentPodIPs() ([]string, bool) {
	pods, ok := r.pods.Load().(map[string]pod.Pod)
	if !ok {
		return nil, false
	}

	ips := make([]string, 0, len(pods))
	for _, pod := range pods {
		ips = append(ips, pod.IP)
	}

	return ips, true
}

func (r *replicator) SendRecoveryRequest(ctx context.Context) (err error) {
	var mu sync.Mutex
	r.rpods.Range(func(name, rpod interface{}) bool {
		select {
		case <-ctx.Done():
			return false
		default:
			r.rpods.Delete(name)
			r.eg.Go(safety.RecoverFunc(func() error {
				_, cerr := r.client.Do(ctx, rpod.(pod.Pod).IP, func(ctx context.Context,
					conn *grpc.ClientConn,
					copts ...grpc.CallOption) (interface{}, error) {
					return agent.NewReplicationClient(conn).Recover(ctx, &payload.Replication_Recovery{}, copts...)
				})
				if cerr != nil {
					r.rpods.Store(name, rpod)
					mu.Lock()
					err = errors.Wrap(err, cerr.Error())
					mu.Unlock()
				}
				return nil
			}))
		}
		return true
	})
	return err
}

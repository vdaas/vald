//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"sync/atomic"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

type Replicator interface {
	Start(context.Context) (<-chan error, error)
	GetDeletedPods() ([]string, bool)
}

type replicator struct {
	pods      atomic.Value
	ctrl      k8s.Controller
	namespace string
	name      string
	eg        errgroup.Group
}

func New(opts ...Option) (rp Replicator, err error) {
	r := new(replicator)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	r.pods.Store(make([]string, 0, 0))

	r.ctrl, err = k8s.New(
		k8s.WithControllerName("vald k8s replication manager controller"),
		k8s.WithEnableLeaderElection(),
		k8s.WithResourceController(pod.New(
			pod.WithControllerName("pod discoverer"),
			pod.WithOnErrorFunc(func(err error) {
				log.Error(err)
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

                deletedIPs := make([]string, 0, len(currentPods))
				for name, cpod := range currentPods {
					p, ok := pods[name]
					if !ok ||
						p.Name != cpod.Name ||
						p.Namespace != cpod.Namespace ||
						p.IP != cpod.IP ||
						p.NodeName != cpod.NodeName {
                            deletedIPs = append(deletedIPs, p.IP)
                    }
				}

                r.eg.Go(safety.RecoverFunc(func() error{
                    return nil
                }))

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
	return r.ctrl.Start(ctx)
}

func (r *replicator) GetDeletedPods() ([]string, bool) {
	return nil, false
}

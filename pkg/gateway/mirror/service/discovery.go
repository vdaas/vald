// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/vald/mirror/target"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/zeebo/xxh3"
)

const (
	resourcePrefix = "mirror-target"
	groupKey       = "group"
)

// Discovery represents an interface for the main logic of service discovery.
// The primary purpose of the Discovery interface is to reconcile custom resources,
// initiating or terminating gRPC connections based on the state of the custom resources.
type Discovery interface {
	Start(ctx context.Context) (<-chan error, error)
}

type discovery struct {
	namespace  string
	labels     map[string]string
	colocation string
	der        net.Dialer

	targetsByName   atomic.Pointer[map[string]target.Target] // latest reconciliation results.
	ctrl            k8s.Controller
	dur             time.Duration
	selfMirrAddrs   []string
	selfMirrAddrStr string

	mirr Mirror
	eg   errgroup.Group
}

// NewDiscovery creates the Discovery object with optional configuration options.
// It returns the initialized Discovery object and an error if the creation process fails.
func NewDiscovery(opts ...DiscoveryOption) (dsc Discovery, err error) {
	d := new(discovery)
	for _, opt := range append(defaultDiscovererOpts, opts...) {
		if err := opt(d); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(err, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	d.targetsByName.Store(&map[string]target.Target{})
	d.selfMirrAddrStr = strings.Join(d.selfMirrAddrs, ",")

	watcher, err := target.New(
		target.WithControllerName("mirror discovery"),
		target.WithNamespace(d.namespace),
		target.WithLabels(d.labels),
		target.WithOnErrorFunc(func(err error) {
			log.Error("failed to reconcile:", err)
		}),
		target.WithOnReconcileFunc(d.onReconcile),
	)
	if err != nil {
		return nil, err
	}

	if d.ctrl == nil {
		d.ctrl, err = k8s.New(
			k8s.WithDialer(d.der),
			k8s.WithControllerName("vald k8s mirror discovery"),
			k8s.WithDisableLeaderElection(),
			k8s.WithResourceController(watcher),
		)
	}
	return d, err
}

func (d *discovery) onReconcile(ctx context.Context, list map[string]target.Target) {
	log.Debugf("mirror reconciled\t%#v", list)
	d.targetsByName.Store(&list)
}

// Start initiates the service discovery process.
// It returns a channel for receiving errors and an error if the initialization fails.
func (d *discovery) Start(ctx context.Context) (<-chan error, error) {
	dech, err := d.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 2)
	d.eg.Go(func() (err error) {
		defer close(ech)
		tic := time.NewTicker(d.dur)
		defer tic.Stop()

		prev := d.loadTargets()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-tic.C:
				prev, err = d.startSync(ctx, prev)
				if err != nil {
					select {
					case <-ctx.Done():
						return errors.Join(err, ctx.Err())
					case ech <- err:
					}
				}
			case err := <-dech:
				if err != nil {
					select {
					case <-ctx.Done():
						return errors.Join(err, ctx.Err())
					case ech <- err:
					}
				}
			}
		}
	})
	return ech, nil
}

func (d *discovery) loadTargets() map[string]target.Target {
	if v := d.targetsByName.Load(); v != nil {
		return *v
	}
	return map[string]target.Target{}
}

type createdTarget struct {
	name string
	tgt  target.Target
}

type updatedTarget struct {
	name string
	old  target.Target
	new  target.Target
}

type deletedTarget struct {
	name string
	host string
	port uint32
}

func (d *discovery) startSync(ctx context.Context, prev map[string]target.Target) (current map[string]target.Target, errs error) {
	current = d.loadTargets()
	curAddrs := map[string]string{} // map[addr: metadata.name]

	created := map[string]*createdTarget{} // map[addr: target.Target]
	updated := map[string]*updatedTarget{} // map[addr: *updatedTarget]
	for name, ctgt := range current {
		addr := net.JoinHostPort(ctgt.Host, uint16(ctgt.Port))
		curAddrs[addr] = name
		if ptgt, ok := prev[name]; !ok {
			created[addr] = &createdTarget{
				name: name,
				tgt:  ctgt,
			}
		} else {
			if ptgt.Host != ctgt.Host || ptgt.Port != ctgt.Port {
				updated[addr] = &updatedTarget{
					name: name,
					old:  ptgt,
					new:  ctgt,
				}
			}
		}
	}

	deleted := map[string]*deletedTarget{} // map[addr: *deletedTarget]
	for name, ptgt := range prev {
		if _, ok := current[name]; !ok {
			addr := net.JoinHostPort(ptgt.Host, uint16(ptgt.Port))
			deleted[addr] = &deletedTarget{
				name: name,
				host: ptgt.Host,
				port: uint32(ptgt.Port),
			}
		}
	}

	if len(created) != 0 || len(deleted) != 0 || len(updated) != 0 {
		log.Infof("created: %#v\tupdated: %#v\tdeleted: %#v", created, updated, deleted)
		if err := d.connectTarget(ctx, created); err != nil {
			errs = errors.Join(errs, err)
		}
		if err := d.disconnectTarget(ctx, deleted); err != nil {
			errs = errors.Join(errs, err)
		}
		if err := d.updateTarget(ctx, updated); err != nil {
			errs = errors.Join(errs, err)
		}
		return current, errs
	}
	return current, d.syncWithAddr(ctx, current, curAddrs)
}

func (d *discovery) syncWithAddr(ctx context.Context, current map[string]target.Target, curAddrs map[string]string) (errs error) {
	for addr, name := range curAddrs {
		// When the status code of a regularly running Register RPC is Unimplemented, the connection to the target will be disconnected
		// so the status of the resource (CR) may be misaligned. To prevent this, change the status of the resource to Disconnected.
		connected := d.mirr.IsConnected(ctx, addr)
		if !connected && isConnectedPhase(current[name].Phase) {
			errs = errors.Join(errs, d.updateMirrorTargetPhase(ctx, name, target.MirrorTargetPhaseDisconnected))
		} else if connected && !isConnectedPhase(current[name].Phase) {
			errs = errors.Join(errs, d.updateMirrorTargetPhase(ctx, name, target.MirrorTargetPhaseConnected))
		}
	}

	d.mirr.RangeMirrorAddr(func(addr string, _ any) bool {
		connected := d.mirr.IsConnected(ctx, addr)
		if name, ok := curAddrs[addr]; ok {
			if connected && !isConnectedPhase(current[name].Phase) {
				errs = errors.Join(errs,
					d.updateMirrorTargetPhase(ctx, name, target.MirrorTargetPhaseConnected),
				)
			} else if !connected && !isDisconnectedPhase(current[name].Phase) {
				errs = errors.Join(errs,
					d.updateMirrorTargetPhase(ctx, name, target.MirrorTargetPhaseDisconnected),
				)
			}
		} else if !ok && connected {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				log.Error(err)
				return true
			}
			name := resourcePrefix + "-" + strconv.FormatUint(xxh3.HashString(d.selfMirrAddrStr+addr), 10)
			errs = errors.Join(errs, d.createMirrorTargetResource(ctx, name, host, int(port)))
		}
		return true
	})
	return errs
}

func (d *discovery) connectTarget(ctx context.Context, req map[string]*createdTarget) (errs error) {
	for _, created := range req {
		phase := target.MirrorTargetPhaseConnected
		err := d.mirr.Connect(ctx, &payload.Mirror_Target{
			Host: created.tgt.Host,
			Port: uint32(created.tgt.Port),
		})
		if err != nil {
			errs = errors.Join(errs, err)
			phase = target.MirrorTargetPhaseDisconnected
		}
		errs = errors.Join(errs, d.updateMirrorTargetPhase(ctx, created.name, phase))
	}
	return errs
}

func (d *discovery) createMirrorTargetResource(ctx context.Context, name, host string, port int) error {
	mt, err := target.NewMirrorTargetTemplate(
		target.WithMirrorTargetName(name),
		target.WithMirrorTargetNamespace(d.namespace),
		target.WithMirrorTargetStatus(&target.MirrorTargetStatus{
			Phase: target.MirrorTargetPhasePending,
		}),
		target.WithMirrorTargetLabels(d.labels),
		target.WithMirrorTargetColocation(d.colocation),
		target.WithMirrorTargetHost(host),
		target.WithMirrorTargetPort(port),
	)
	if err != nil {
		return err
	}
	return d.ctrl.GetManager().GetClient().Create(ctx, mt)
}

func (d *discovery) disconnectTarget(ctx context.Context, req map[string]*deletedTarget) (errs error) {
	for _, deleted := range req {
		phase := target.MirrorTargetPhaseDisconnected
		err := d.mirr.Disconnect(ctx, &payload.Mirror_Target{
			Host: deleted.host,
			Port: deleted.port,
		})
		if err != nil {
			errs = errors.Join(errs, err)
			phase = target.MirrorTargetPhaseUnknown
		}
		errs = errors.Join(errs, d.updateMirrorTargetPhase(ctx, deleted.name, phase))
	}
	return errs
}

func (d *discovery) updateMirrorTargetPhase(ctx context.Context, name string, phase target.MirrorTargetPhase) error {
	c := d.ctrl.GetManager().GetClient()
	mt := &target.MirrorTarget{}
	err := c.Get(ctx, k8s.ObjectKey{
		Namespace: d.namespace,
		Name:      name,
	}, mt)
	if err != nil {
		return err
	}
	if mt.Status.Phase == phase {
		return nil
	}
	mt.Status.Phase = phase
	mt.Status.LastTransitionTime = k8s.Now()
	return c.Status().Update(ctx, mt)
}

func (d *discovery) updateTarget(ctx context.Context, req map[string]*updatedTarget) (errs error) {
	for _, updated := range req {
		err := d.mirr.Disconnect(ctx, &payload.Mirror_Target{
			Host: updated.old.Host,
			Port: uint32(updated.old.Port),
		})
		if err != nil {
			errs = errors.Join(errs, err, d.updateMirrorTargetPhase(ctx, updated.name, target.MirrorTargetPhaseUnknown))
		} else {
			err := d.mirr.Connect(ctx, &payload.Mirror_Target{
				Host: updated.new.Host,
				Port: uint32(updated.new.Port),
			})
			if err != nil {
				errs = errors.Join(errs, err, d.updateMirrorTargetPhase(ctx, updated.name, target.MirrorTargetPhaseDisconnected))
			}
		}
	}
	return errs
}

func isConnectedPhase(phase target.MirrorTargetPhase) bool {
	return phase == target.MirrorTargetPhaseConnected
}

func isDisconnectedPhase(phase target.MirrorTargetPhase) bool {
	return phase == target.MirrorTargetPhaseDisconnected
}

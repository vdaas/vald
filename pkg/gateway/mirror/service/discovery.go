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
		target.WithOnReconcileFunc(func(ctx context.Context, list map[string]target.Target) {
			log.Debugf("mirror reconciled\t%#v", list)
			d.targetsByName.Store(&list)
		}),
	)
	if err != nil {
		return nil, err
	}
	d.ctrl, err = k8s.New(
		k8s.WithDialer(d.der),
		k8s.WithControllerName("vald k8s mirror discovery"),
		k8s.WithDisableLeaderElection(),
		k8s.WithResourceController(watcher),
	)
	return d, err
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
	return nil
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

func (d *discovery) startSync(ctx context.Context, prev map[string]target.Target) (current map[string]target.Target, err error) {
	current = d.loadTargets()
	curAddrs := make(map[string]string) // map[addr: metadata.name]

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
		err = errors.Join(
			errors.Join(
				d.connectTarget(ctx, created),
				d.disconnectTarget(ctx, deleted)),
			d.updateTarget(ctx, updated))
		return current, err
	}

	for addr, name := range curAddrs {
		// When the status code of a regularly running Register RPC is Unimplemented, the connection to the target will be disconnected
		// so the status of the resource (CR) may be misaligned. To prevent this, change the status of the resource to Disconnected.
		if !d.mirr.IsConnected(ctx, addr) && current[name].Phase == target.MirrorTargetPhaseConnected {
			err = errors.Join(err, d.updateMirrorTargetPhase(ctx, name, target.MirrorTargetPhaseDisconnected))
		}
	}

	d.mirr.RangeMirrorAddr(func(addr string, _ any) bool {
		connected := d.mirr.IsConnected(ctx, addr)
		if name, ok := curAddrs[addr]; ok {
			if st := target.MirrorTargetPhaseConnected; connected && current[name].Phase != st {
				err = errors.Join(err,
					d.updateMirrorTargetPhase(ctx, name, st))
			} else if st := target.MirrorTargetPhaseDisconnected; !connected && current[name].Phase != st {
				err = errors.Join(err,
					d.updateMirrorTargetPhase(ctx, name, st))
			}
		} else if !ok && connected {
			var (
				host string
				port uint16
			)
			host, port, err = net.SplitHostPort(addr)
			if err != nil {
				log.Error(err)
			}
			name := resourcePrefix + "-" + strconv.FormatUint(xxh3.HashString(d.selfMirrAddrStr+addr), 10)
			err = errors.Join(err, d.createMirrorTargetResource(ctx, name, host, int(port)))
		}
		return true
	})
	return current, err
}

func (d *discovery) connectTarget(ctx context.Context, req map[string]*createdTarget) (err error) {
	for _, created := range req {
		phase := target.MirrorTargetPhaseConnected
		cerr := d.mirr.Connect(ctx, &payload.Mirror_Target{
			Host: created.tgt.Host,
			Port: uint32(created.tgt.Port),
		})
		if cerr != nil {
			err = errors.Join(err, cerr)
			phase = target.MirrorTargetPhaseDisconnected
		}
		uerr := d.updateMirrorTargetPhase(ctx, created.name, phase)
		if uerr != nil {
			err = errors.Join(err, uerr)
		}
	}
	return err
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

func (d *discovery) disconnectTarget(ctx context.Context, req map[string]*deletedTarget) error {
	if len(req) == 0 {
		return nil
	}
	tgts := make([]*payload.Mirror_Target, 0, len(req))
	for _, deleted := range req {
		tgts = append(tgts, &payload.Mirror_Target{
			Host: deleted.host,
			Port: deleted.port,
		})
	}
	return d.mirr.Disconnect(ctx, tgts...)
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

func (d *discovery) updateTarget(ctx context.Context, req map[string]*updatedTarget) (err error) {
	for _, updated := range req {
		derr := d.mirr.Disconnect(ctx, &payload.Mirror_Target{
			Host: updated.old.Host,
			Port: uint32(updated.old.Port),
		})
		if derr != nil {
			err = errors.Join(err, derr)
		} else {
			if cerr := d.mirr.Connect(ctx, &payload.Mirror_Target{
				Host: updated.new.Host,
				Port: uint32(updated.new.Port),
			}); cerr != nil {
				err = errors.Join(cerr, err)
				if uerr := d.updateMirrorTargetPhase(ctx, updated.name, target.MirrorTargetPhaseDisconnected); uerr != nil {
					err = errors.Join(err, uerr)
				}
			}
		}
	}
	return err
}

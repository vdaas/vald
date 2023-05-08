package service

import (
	"context"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/vald/mirror/target"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/zeebo/xxh3"
)

const (
	resourcePrefix = "mirror-target"
)

type Discoverer interface {
	Start(ctx context.Context) (<-chan error, error)
}

type discoverer struct {
	namespace  string
	colocation string
	der        net.Dialer

	targetsByName atomic.Pointer[map[string]target.Target] // latest reconciliation results.
	ctrl          k8s.Controller
	dur           time.Duration

	mirr Mirror
	eg   errgroup.Group
}

func NewDiscoverer(opts ...DiscovererOption) (dsc Discoverer, err error) {
	d := new(discoverer)
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

	watcher, err := target.New(
		target.WithControllerName("mirror discoverer"),
		target.WithNamespace(d.namespace),
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
		k8s.WithControllerName("vald k8s mirror discoverer"),
		k8s.WithDisableLeaderElection(),
		k8s.WithResourceController(watcher),
	)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *discoverer) Start(ctx context.Context) (<-chan error, error) {
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

func (d *discoverer) loadTargets() map[string]target.Target {
	if v := d.targetsByName.Load(); v != nil {
		return *v
	}
	return nil
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

func (d *discoverer) startSync(ctx context.Context, prev map[string]target.Target) (cur map[string]target.Target, err error) {
	cur = d.loadTargets()
	curAddrs := make(map[string]string) // map[addr: metadata.name]

	created := map[string]target.Target{}  // map[addr: target.Target]
	updated := map[string]*updatedTarget{} // map[addr: *updatedTarget]
	for name, ctgt := range cur {
		addr := net.JoinHostPort(ctgt.Host, uint16(ctgt.Port))
		curAddrs[addr] = name
		if ptgt, ok := prev[name]; !ok {
			created[addr] = ctgt
		} else {
			if ptgt.Host != ctgt.Host || ptgt.Port != ctgt.Port {
				updated[addr] = &updatedTarget{name: name, old: ptgt, new: ctgt}
			}
		}
	}

	deleted := map[string]*deletedTarget{} // map[addr: *deletedTarget]
	for name, ptgt := range prev {
		if _, ok := cur[name]; !ok {
			addr := net.JoinHostPort(ptgt.Host, uint16(ptgt.Port))
			deleted[addr] = &deletedTarget{name: name, host: ptgt.Host, port: uint32(ptgt.Port)}
		}
	}
	log.Infof("created: %#v\tupdated: %#v\tdeleted: %#v", created, updated, deleted)

	err = errors.Join(
		errors.Join(
			d.createTarget(ctx, created),
			d.deleteTarget(ctx, deleted)),
		d.updateTarget(ctx, updated))
	if err != nil {
		return cur, err
	}

	d.mirr.RangeAllMirrorAddr(func(addr string, _ any) bool {
		connected := d.mirr.IsConnected(ctx, addr)
		if name, ok := curAddrs[addr]; ok && !connected {
			err = errors.Join(err, d.deleteMirrorTargetResource(ctx, name))
		} else if !ok && connected {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				log.Error(err)
			}
			name := resourcePrefix + "-" + strconv.Itoa(int(xxh3.HashString(addr)))
			err = errors.Join(err, d.createMirrorTargetResource(ctx, name, host, int(port)))
		}
		return true
	})
	return cur, err
}

func (d *discoverer) createTarget(ctx context.Context, req map[string]target.Target) error {
	if len(req) == 0 {
		return nil
	}
	conns := make([]*payload.Mirror_Target, 0, len(req))
	for _, created := range req {
		conns = append(conns, &payload.Mirror_Target{
			Ip:   created.Host,
			Port: uint32(created.Port),
		})
	}
	return d.mirr.Connect(ctx, conns...)
}

func (d *discoverer) createMirrorTargetResource(ctx context.Context, name, host string, port int) error {
	mt, err := target.NewMirrorTargetTemplate(
		target.WithMirrorTargetName(name),
		target.WithMirrorTargetNamespace(d.namespace),
		target.WithMirrorTargetColocation(d.colocation),
		target.WithMirrorTargetHost(host),
		target.WithMirrorTargetPort(port),
	)
	if err != nil {
		return err
	}
	return d.ctrl.GetManager().GetClient().Create(ctx, mt)
}

func (d *discoverer) deleteTarget(ctx context.Context, req map[string]*deletedTarget) (err error) {
	if len(req) == 0 {
		return nil
	}
	tgts := make([]*payload.Mirror_Target, 0, len(req))
	for _, deleted := range req {
		tgts = append(tgts, &payload.Mirror_Target{
			Ip:   deleted.host,
			Port: deleted.port,
		})
	}
	return d.mirr.Disconnect(ctx, tgts...)
}

func (d *discoverer) deleteMirrorTargetResource(ctx context.Context, name string) error {
	mt, err := target.NewMirrorTargetTemplate(
		target.WithMirrorTargetName(name),
		target.WithMirrorTargetNamespace(d.namespace),
	)
	if err != nil {
		return err
	}
	return d.ctrl.GetManager().GetClient().Delete(ctx, mt)
}

func (d *discoverer) updateTarget(ctx context.Context, req map[string]*updatedTarget) error {
	if len(req) == 0 {
		return nil
	}
	delTgts := make([]*payload.Mirror_Target, 0, len(req)/2)
	newTgts := make([]*payload.Mirror_Target, 0, len(req)/2)
	for _, updated := range req {
		delTgts = append(delTgts, &payload.Mirror_Target{
			Ip:   updated.old.Host,
			Port: uint32(updated.old.Port),
		})
		newTgts = append(newTgts, &payload.Mirror_Target{
			Ip:   updated.new.Host,
			Port: uint32(updated.new.Port),
		})
	}
	return errors.Join(d.mirr.Disconnect(ctx, delTgts...), d.mirr.Connect(ctx, newTgts...))
}

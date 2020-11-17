package service

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/job"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/rebalancer/storage/controller/model"
)

// Discoverer --
type Discoverer interface {
	// Start --
	Start(ctx context.Context) (<-chan error, error)
}

type discoverer struct {
	jobs         jobsMap
	jobsCache    atomic.Value
	jobName      string
	jobNamespace string

	dcd  time.Duration // discover check duration
	eg   errgroup.Group
	ctrl k8s.Controller
}

func NewDiscoverer(opts ...DiscovererOption) (Discoverer, error) {
	d := new(discoverer)

	for _, opt := range append(defaultDiscovererOpts, opts...) {
		if err := opt(d); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	job, err := job.New(
		job.WithControllerName("job discoverer"),
		job.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		job.WithOnReconcileFunc(func(jobsMap map[string][]job.Job) {
			for name, jobs := range jobsMap {
				if name == d.jobName {
					d.jobs.Store(name, jobs)
				}
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	d.ctrl, err = k8s.New(
		k8s.WithControllerName("rebalance controller"),
		k8s.WithEnableLeaderElection(),
		k8s.WithResourceController(job),
	)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *discoverer) Start(ctx context.Context) (<-chan error, error) {
	cech, err := d.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}

	ech := make(chan error, 1)
	d.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		dt := time.NewTicker(d.dcd)
		defer dt.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-dt.C:
				var wg sync.WaitGroup
				wg.Add(1)
				d.eg.Go(safety.RecoverFunc(func() error {
					defer wg.Done()
					jobs, ok := d.jobs.Load(d.jobName)
					if !ok {
						return nil
					}
					models := make([]*model.Job, 0, len(jobs))
					for _, job := range jobs {
						var t time.Time
						if job.Status.StartTime != nil {
							t = job.Status.StartTime.Time
						}
						models = append(models, &model.Job{
							Name:      job.Name,
							Namespace: job.Namespace,
							Active:    job.Status.Active,
							StartTime: t,
						})
					}
					d.jobsCache.Store(models)
					return nil
				}))
				wg.Wait()
			case err := <-cech:
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

package service

import (
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/job"
	"github.com/vdaas/vald/internal/log"
)

type Discoverer interface {
}

type discoverer struct {
	jobs         jobsMap
	jobName      string
	jobNamespace string

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
					// TODO create job model to only store required information
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

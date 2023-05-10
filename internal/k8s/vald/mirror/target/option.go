package target

import (
	"context"

	"github.com/vdaas/vald/internal/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Option func(r *reconciler) error

var defaultOptions = []Option{}

func WithControllerName(name string) Option {
	return func(r *reconciler) error {
		if len(name) == 0 {
			return errors.NewErrInvalidOption("controllerName", name)
		}
		r.name = name
		return nil
	}
}

func WithManager(mgr manager.Manager) Option {
	return func(r *reconciler) error {
		if mgr == nil {
			return errors.NewErrInvalidOption("manager", mgr)
		}
		r.mgr = mgr
		return nil
	}
}

func WithOnErrorFunc(f func(error)) Option {
	return func(r *reconciler) error {
		if f == nil {
			return errors.NewErrInvalidOption("onErrorFunc", f)
		}
		r.onError = f
		return nil
	}
}

func WithOnReconcileFunc(f func(context.Context, map[string]Target)) Option {
	return func(r *reconciler) error {
		if f == nil {
			return errors.NewErrInvalidOption("onReconcileFunc", f)
		}
		r.onReconcile = f
		return nil
	}
}

func WithNamespace(ns string) Option {
	return func(r *reconciler) error {
		if ns == "" {
			return errors.NewErrInvalidOption("namespace", ns)
		}
		r.addListOpts(client.InNamespace(ns))
		return nil
	}
}

func WithLabels(labels map[string]string) Option {
	return func(r *reconciler) error {
		if len(labels) == 0 {
			return errors.NewErrInvalidOption("labels", labels)
		}
		r.addListOpts(client.MatchingLabels(labels))
		return nil
	}
}

package replicaset

import "sigs.k8s.io/controller-runtime/pkg/manager"

// Option represents functional option for reconciler.
type Option func(*reconciler) error

var defaultOpts = []Option{}

// WithControllerName returns Option to set the controller name.
func WithControllerName(name string) Option {
	return func(r *reconciler) error {
		r.name = name
		return nil
	}
}

// WithNamespace returns Option to set the namespace.
func WithNamespace(ns string) Option {
	return func(r *reconciler) error {
		r.namespace = ns
		return nil
	}
}

// WithManager returns Option to set the resource manager.
func WithManager(mgr manager.Manager) Option {
	return func(r *reconciler) error {
		r.mgr = mgr
		return nil
	}
}

// WithOnErrorFunc returns Option to set the onError hook.
func WithOnErrorFunc(f func(err error)) Option {
	return func(r *reconciler) error {
		r.onError = f
		return nil
	}
}

// WithOnReconcileFunc returns Option to set the onReconcile hook.
func WithOnReconcileFunc(f func(replicasetList map[string][]ReplicaSet)) Option {
	return func(r *reconciler) error {
		r.onReconcile = f
		return nil
	}
}

package job

import "sigs.k8s.io/controller-runtime/pkg/manager"

// Option represents functional option for reconciler.
type Option func(*reconciler) error

var defaultOpts = []Option{}

// WithControllerName returns Option that sets r.name.
func WithControllerName(name string) Option {
	return func(r *reconciler) error {
		r.name = name
		return nil
	}
}

// WithManager returns Option that sets r.mgr.
func WithManager(mgr manager.Manager) Option {
	return func(r *reconciler) error {
		r.mgr = mgr
		return nil
	}
}

// WithOnErrorFunc returns Option that sets r.onError.
func WithOnErrorFunc(f func(err error)) Option {
	return func(r *reconciler) error {
		r.onError = f
		return nil
	}
}

// WithOnReconcileFunc returns Option that sets r.onReconcile.
func WithOnReconcileFunc(f func(jobList map[string][]Job)) Option {
	return func(r *reconciler) error {
		r.onReconcile = f
		return nil
	}
}

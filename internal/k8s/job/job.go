package job

import (
	"context"
	"reflect"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
)

// JobWatcher is a type alias for k8s resource controller.
type JobWatcher k8s.ResourceController

type reconciler struct {
	ctx         context.Context
	mgr         manager.Manager
	name        string
	namespance  string
	onError     func(err error)
	onReconcile func(jobList map[string][]Job)
}

// Job is a type alias for the k8s job definition.
type Job = batchv1.Job

// New returns the JobWatcher that implements reconciliation loop, or any error occurred.
func New(opts ...Option) (JobWatcher, error) {
	r := new(reconciler)

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return nil, nil
}

// Reconcile implements k8s reconciliation loop to retrieve the Job information from k8s.
func (r *reconciler) Reconcile(req reconcile.Request) (res reconcile.Result, err error) {
	js := new(batchv1.JobList)

	err = r.mgr.GetClient().List(r.ctx, js)
	if err != nil {
		if r.onError != nil {
			r.onError(err)
		}
		res = reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Millisecond * 100,
		}
		if k8serrors.IsNotFound(err) {
			log.Error("not found", err)
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, nil
		}
		return
	}

	jobs := make(map[string][]Job, len(js.Items))

	for _, job := range js.Items {
		name, ok := job.GetObjectMeta().GetLabels()["app"]
		if !ok {
			jns := strings.Split(job.GetName(), "-")
			name = strings.Join(jns[:len(jns)-1], "-")
		}

		// TODO: pre-alocate job slice.
		if _, ok := jobs[name]; !ok {
			jobs[name] = make([]Job, 0, len(js.Items))
		}

		jobs[name] = append(jobs[name], job)
	}

	if r.onReconcile != nil {
		r.onReconcile(jobs)
	}
	return
}

// GetName returns the name of resource controller.
func (r *reconciler) GetName() string {
	return r.name
}

// NewReconciler returns the reconciler for the Job.
func (r *reconciler) NewReconciler(ctx context.Context, mgr manager.Manager) reconcile.Reconciler {
	if r.ctx == nil && ctx != nil {
		r.ctx = ctx
	}
	if r.mgr == nil && mgr != nil {
		r.mgr = mgr
	}

	batchv1.AddToScheme(r.mgr.GetScheme())
	return r
}

// For returns the runtime.Object which is job.
func (r *reconciler) For() runtime.Object {
	return new(batchv1.Job)
}

// Owns returns the owner of the job watcher.
// It will always return nil.
func (r *reconciler) Owns() runtime.Object {
	return nil
}

// Watches returns the kind of the job and the event handler.
// It will always return nil.
func (r *reconciler) Watches() (*source.Kind, handler.EventHandler) {
	// return &source.Kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil
}

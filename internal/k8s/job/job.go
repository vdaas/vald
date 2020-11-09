package job

import (
	"context"
	"reflect"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	batchv1 "k8s.io/api/batch/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type reconciler struct {
	ctx         context.Context
	mgr         manager.Manager
	name        string
	namespance  string
	onError     func(err error)
	onReconcile func(jobList map[string][]Job)
}

// Job represents k8s job imformation
// https://godoc.org/k8s.io/api/batch/v1#Job
type Job struct {
	Name      string
	Namespace string
	Active    int32
	Succeeded int32
	Failed    int32
	StartTime *time.Time
}

// New returns k8s.ResourceController(*reconciler) implementation.
func New(opts ...Option) (k8s.ResourceController, error) {
	r := new(reconciler)

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return nil, nil
}

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

	var (
		jobs = make(map[string][]Job, len(js.Items))
	)

	for _, job := range js.Items {
		name := job.GetName()
		// TODO: pre-alocate job slice.
		if _, ok := jobs[name]; !ok {
			jobs[name] = make([]Job, 0, len(js.Items))
		}

		var t *time.Time
		if job.Status.StartTime != nil {
			t = &job.Status.StartTime.Time
		}
		jobs[name] = append(jobs[name], Job{
			Name:      job.GetName(),
			Namespace: job.GetNamespace(),
			Active:    job.Status.Active,
			Succeeded: job.Status.Succeeded,
			Failed:    job.Status.Failed,
			StartTime: t,
		})
	}

	if r.onReconcile != nil {
		r.onReconcile(jobs)
	}
	return
}

func (r *reconciler) GetName() string {
	return r.name
}

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

func (r *reconciler) For() runtime.Object {
	return new(batchv1.Job)
}

func (r *reconciler) Owns() runtime.Object {
	return nil
}

func (r *reconciler) Watches() (*source.Kind, handler.EventHandler) {
	// return &source.Kind{Type: new(corev1.Pod)}, &handler.EnqueueRequestForObject{}
	return nil, nil
}

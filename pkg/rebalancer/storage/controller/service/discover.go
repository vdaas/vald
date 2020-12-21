package service

import (
	"context"
	"math"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/configmap"
	"github.com/vdaas/vald/internal/k8s/job"
	mpod "github.com/vdaas/vald/internal/k8s/metrics/pod"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/k8s/statefulset"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/rebalancer/storage/controller/model"

	// TODO: move to internal

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"
)

// Discoverer represents the discoverer interface.
// TODO: rename discoverer to rebalancer
type Discoverer interface {
	// Start --
	Start(ctx context.Context) (<-chan error, error)
}

type discoverer struct {
	jobs           atomic.Value
	jobName        string
	jobNamespace   string
	jobTemplate    atomic.Value
	jobTemplateKey string // config map key of job template

	agentName         string
	agentNamespace    string
	agentResourceType string // TODO: use custom type insteaf of string
	pods              atomic.Value
	podMetrics        atomic.Value

	// job template config map
	configmapName      string
	configmapNamespace string

	statefulSets atomic.Value

	leaderElectionID string

	rcd       time.Duration // reconcile check duration
	eg        errgroup.Group
	ctrl      k8s.Controller
	tolerance float64
	decoder   *conversion.Decoder
}

// NewDiscoverer --
func NewDiscoverer(opts ...DiscovererOption) (Discoverer, error) {
	d := new(discoverer)

	for _, opt := range append(defaultDiscovererOpts, opts...) {
		if err := opt(d); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	job, err := job.New(
		job.WithControllerName("job discoverer"),
		job.WithNamespaces(d.jobNamespace),
		job.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		job.WithOnReconcileFunc(func(jobList map[string][]job.Job) {
			jobs, ok := jobList[d.jobName]
			if ok {
				d.jobs.Store(jobs)
			} else {
				log.Infof("job not found: %s", d.jobName)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	cm, err := configmap.New(
		configmap.WithControllerName("configmap discoverer"),
		configmap.WithNamespaces(d.configmapNamespace),
		configmap.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		configmap.WithOnReconcileFunc(func(configmapList map[string][]configmap.ConfigMap) {
			configmaps, ok := configmapList[d.configmapNamespace]
			if ok {
				for _, cm := range configmaps {
					if cm.Name == d.configmapName {
						if tmpl, ok := cm.Data[d.jobTemplateKey]; ok {
							d.jobTemplate.Store(tmpl)
						} else {
							log.Infof("job template not found: %s", d.jobTemplateKey)
						}
						break
					}
				}
			} else {
				log.Infof("configmap not found: %s", d.configmapName)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	var rc k8s.ResourceController
	switch d.agentResourceType {
	case "statefulset":
		rc, err = statefulset.New(
			statefulset.WithControllerName("statefulset discoverer"),
			statefulset.WithNamespaces(d.agentNamespace),
			statefulset.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			statefulset.WithOnReconcileFunc(func(statefulSetList map[string][]statefulset.StatefulSet) {
				sss, ok := statefulSetList[d.agentName]
				if ok {
					if len(sss) == 1 {
						d.statefulSets.Store(sss[0])
					} else {
						log.Infof("too many statefulset list: want 1, but %d", len(sss))
					}
				} else {
					log.Infof("statefuleset not found: %s", d.agentName)
				}
			}),
		)
		if err != nil {
			return nil, err
		}
	case "replicaset":
		// TODO: implment get replicaset reconciled result
		return nil, nil
	case "daemonset":
		// TODO: implment get daemonset reconciled result
		return nil, nil
	default:
		return nil, errors.New("invalid agent resource type: " + d.agentResourceType)
	}

	d.ctrl, err = k8s.New(
		k8s.WithControllerName("rebalance storage controller"),
		k8s.WithEnableLeaderElection(),
		k8s.WithLeaderElectionID(d.leaderElectionID),
		k8s.WithResourceController(job),
		k8s.WithResourceController(rc), // statefulset controller
		k8s.WithResourceController(pod.New(
			pod.WithControllerName("pod discover"),
			pod.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			pod.WithOnReconcileFunc(func(podList map[string][]pod.Pod) {
				pods, ok := podList[d.agentName]
				if ok {
					d.pods.Store(pods)
				} else {
					log.Infof("pod not found: %s", d.agentName)
				}
			}),
		)),
		k8s.WithResourceController(mpod.New(
			mpod.WithControllerName("pod metrics discover"),
			mpod.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			mpod.WithOnReconcileFunc(func(podList map[string]mpod.Pod) {
				if len(podList) > 0 {
					d.podMetrics.Store(podList)
				} else {
					log.Info("pod metrics not found")
				}
			}),
		)),
		k8s.WithResourceController(cm),
	)
	if err != nil {
		return nil, err
	}

	d.decoder, err = conversion.NewDecoder(runtime.NewScheme())
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
		dt := time.NewTicker(d.rcd)
		defer dt.Stop()

		var (
			prevSsModel   map[string]*model.StatefulSet
			prevPodModels map[string][]*model.Pod
		)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-dt.C:
				var (
					ss statefulset.StatefulSet
					ok bool

					podModels map[string][]*model.Pod
					jobModels map[string][]*model.Job
					ssModel   map[string]*model.StatefulSet
					jobTpl    *job.Job
				)

				podModels, err := d.genPodModels()
				if err != nil {
					log.Infof("error generating pod models: %s", err.Error())
					continue
				}

				jobModels, err = d.genJobModels()
				if err != nil {
					log.Infof("error generating job models: %s", err.Error())
				}

				jobTpl, err = d.genJobTpl()
				if err != nil {
					log.Infof("error generating job template: %s", err.Error())
					continue
				}

				// TODO: cache specified reconciled result based on agentResourceType.
				switch d.agentResourceType {
				case "statefulset":
					ss, ok = d.statefulSets.Load().(statefulset.StatefulSet)
					if !ok {
						log.Info("statefulset is empty")
						continue
					}
					ssModel = make(map[string]*model.StatefulSet)
					ssModel[ss.Namespace] = &model.StatefulSet{
						Name:            ss.Name,
						Namespace:       ss.Namespace,
						DesiredReplicas: ss.Spec.Replicas,
						Replicas:        ss.Status.Replicas,
					}
				default:
					// TODO: define error for return
					return nil
				}

				if prevSsModel != nil {
					for ns, psm := range prevSsModel {
						if _, ok := ssModel[ns]; !ok {
							continue
						}
						if *ssModel[ns].DesiredReplicas != ssModel[ns].Replicas {
							continue
						}

						decreasedPodNames := d.isSsReplicaDecreased(psm, ssModel[ns], prevPodModels[ns], podModels[ns])
						if len(decreasedPodNames) > 0 {
							for _, name := range decreasedPodNames {
								log.Debugf("[decrease] creating job for pod %s", name)
								if err := d.createJob(ctx, *jobTpl, name); err != nil {
									log.Errorf("failed to create job: %s", err)
									continue
								}
							}
						} else {
							maxPodName, rate := d.getBiasOverDetail(podModels[ns])
							if maxPodName != "" && !d.isJobRunning(jobModels, ns) {
								log.Debugf("[bias] creating job for pod %s, rate: %v", maxPodName, rate)
								if err := d.createJob(ctx, *jobTpl, maxPodName); err != nil {
									log.Errorf("failed to create job: %s", err)
									continue
								}
							}
						}
					}
				}
				// Store reconciled result for next loop.
				prevSsModel = ssModel
				prevPodModels = podModels

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

func (d *discoverer) createJob(ctx context.Context, jobTpl job.Job, agentName string) error {
	if jobTpl.Labels == nil {
		jobTpl.Labels = make(map[string]string)
	}
	if len(d.jobNamespace) != 0 {
		jobTpl.Namespace = d.jobNamespace
	}
	jobTpl.Labels["target_agent_name"] = agentName

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	scheme := runtime.NewScheme()
	if err = batchv1.AddToScheme(scheme); err != nil {
		return err
	}

	c, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return err
	}

	if err := c.Create(ctx, &jobTpl); err != nil {
		return err
	}

	return nil
}

func (d *discoverer) genPodModels() (podModels map[string][]*model.Pod, err error) {
	mpods, ok := d.podMetrics.Load().(map[string]mpod.Pod)
	if !ok {
		return nil, errors.New("pod metrics is empty")
	}

	pods, ok := d.pods.Load().([]pod.Pod)
	if !ok {
		return nil, errors.New("pod is empty")
	}

	podModels = make(map[string][]*model.Pod)
	for _, p := range pods {
		if _, ok := podModels[p.Namespace]; !ok {
			podModels[p.Namespace] = make([]*model.Pod, 0)
		}
		log.Debugf("%s limit: %#v, metrics: %#v", p.Name, p.MemLimit, mpods[p.Name])
		if mpod, ok := mpods[p.Name]; ok {
			podModels[p.Namespace] = append(podModels[p.Namespace], &model.Pod{
				Name:        p.Name,
				Namespace:   p.Namespace,
				MemoryLimit: p.MemLimit,
				MemoryUsage: mpod.Mem,
			})
		}
	}

	return
}

func (d *discoverer) genJobModels() (jobModels map[string][]*model.Job, err error) {
	jobs, ok := d.jobs.Load().([]job.Job)
	if !ok {
		return nil, errors.New("job is empty")
	}

	jobModels = make(map[string][]*model.Job)
	for _, j := range jobs {
		var t time.Time
		if j.Status.StartTime != nil {
			t = j.Status.StartTime.Time
		}
		if _, ok := jobModels[j.Namespace]; !ok {
			jobModels[j.Namespace] = make([]*model.Job, 0)
		}

		jobModels[j.Namespace] = append(jobModels[j.Namespace], &model.Job{
			Name:                 j.Name,
			Namespace:            j.Namespace,
			Active:               j.Status.Active,
			StartTime:            t,
			Type:                 j.Labels["type"],
			TargetAgentNamespace: j.Labels["target_agent_namespace"],
			TargetAgentName:      j.Labels["target_agent_name"],
			ControllerNamespace:  j.Labels["controller_namespace"],
			ControllerName:       j.Labels["controller_name"],
		})
	}

	return
}

func (d *discoverer) genJobTpl() (jobTpl *job.Job, err error) {
	tmpl, ok := d.jobTemplate.Load().(string)
	if !ok {
		return nil, errors.New("job template is empty")
	}
	jobTpl = &job.Job{}
	err = d.decoder.DecodeInto([]byte(tmpl), jobTpl)
	if err != nil {
		return nil, errors.Wrap(err, "fails decoding template")
	}
	return
}

func (d *discoverer) isSsReplicaDecreased(psm, sm *model.StatefulSet, ppm, pm []*model.Pod) (podNames []string) {
	if *psm.DesiredReplicas > *sm.DesiredReplicas {
		podNames = make([]string, 0)
		for _, prevPod := range ppm {
			var ok bool
			for _, pod := range pm {
				if prevPod.Name == pod.Name {
					ok = true
					break
				}
			}
			if !ok {
				podNames = append(podNames, prevPod.Name)
			}
		}
	}
	return
}

func (d *discoverer) getBiasOverDetail(pm []*model.Pod) (string, float64) {
	var unlimited bool
	for _, p := range pm {
		if p.MemoryLimit <= 0 {
			unlimited = true
			break
		}
	}

	if unlimited {
		podName, avgMemUsg, maxMemUsg := calAvgMemUsg(pm)
		sig := calSigMemUsg(pm, avgMemUsg)

		log.Debugf("podName, avgMemUsg, maxMemUsg, sig: %s, %.3f, %f, %.3f ", podName, avgMemUsg, maxMemUsg, sig)

		if maxMemUsg >= (1+d.tolerance)*sig {
			return podName, 1 - (avgMemUsg / maxMemUsg)
		}
		return "", 0
	}
	podName, avgMemUsg, maxMemUsg := calAvgMemUsgWithMemLimit(pm)

	if maxMemUsg >= avgMemUsg+d.tolerance {
		return podName, 1 - (avgMemUsg / maxMemUsg)
	}
	return "", 0
}

func calAvgMemUsgWithMemLimit(pm []*model.Pod) (podName string, avgMemUsg, maxMemUsg float64) {
	for _, p := range pm {
		u := p.MemoryUsage / p.MemoryLimit
		avgMemUsg += u
		if u > maxMemUsg {
			maxMemUsg = u
			podName = p.Name
		}
	}
	avgMemUsg = avgMemUsg / float64(len(pm))
	return
}

func calAvgMemUsg(pm []*model.Pod) (podName string, avgMemUsg, maxMemUsg float64) {
	for _, p := range pm {
		u := p.MemoryUsage
		avgMemUsg += u
		if u > maxMemUsg {
			maxMemUsg = u
			podName = p.Name
		}
	}
	avgMemUsg = avgMemUsg / float64(len(pm))
	return
}

func calSigMemUsg(pm []*model.Pod, avgMemUsg float64) (sig float64) {
	for _, p := range pm {
		sig += math.Pow((p.MemoryUsage - avgMemUsg), 2.0)
	}
	sig /= float64(len(pm))
	sig = math.Sqrt(sig)
	return
}

func (d *discoverer) isJobRunning(jobModels map[string][]*model.Job, ns string) bool {
	for _, jobs := range jobModels {
		for _, job := range jobs {
			if job.Type == "rebalance" && job.Active != 0 && job.TargetAgentNamespace == ns {
				return true
			}
		}
	}
	return false
}

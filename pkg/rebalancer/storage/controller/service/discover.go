package service

import (
	"context"
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
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"
)

// Discoverer --
type Discoverer interface {
	// Start --
	Start(ctx context.Context) (<-chan error, error)
}

type discoverer struct {
	jobs         atomic.Value
	jobName      string
	jobNamespace string
	jobTemplate  atomic.Value
	// TODO: thiknig about this name.
	jobTemplateKey string

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
			log.Infof("[reconcile] reconcile of jobList: %v", jobList)
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
			log.Infof("[reconcile] configmap: %#v", configmaps)
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
	log.Infof("using agent type: %s", d.agentResourceType)
	switch d.agentResourceType {
	case "statefulset":
		log.Info("using statefulset agent type")
		rc, err = statefulset.New(
			statefulset.WithControllerName("statefulset discoverer"),
			statefulset.WithNamespaces(d.agentNamespace),
			statefulset.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			statefulset.WithOnReconcileFunc(func(statefulSetList map[string][]statefulset.StatefulSet) {
				sss, ok := statefulSetList[d.agentName]
				log.Infof("[reconcile] statefulset for agent: %v, list: %#v", d.agentName, sss)
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
		k8s.WithControllerName("rebalance controller"),
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
				// log.Infof("[reconcile] pod list: %#v", len(podList))
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
				// log.Infof("[reconcile] pod metrics: %#v", podList)
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
			prevJobTpl    *job.Job
		)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-dt.C:
				var (
					mpods map[string]mpod.Pod
					pods  []pod.Pod
					jobs  []job.Job
					ss    statefulset.StatefulSet
					ok    bool

					podModels map[string][]*model.Pod
					jobModels map[string][]*model.Job
					ssModel   map[string]*model.StatefulSet
					jobTpl    job.Job
				)

				mpods, ok = d.podMetrics.Load().(map[string]mpod.Pod)
				if !ok {
					log.Info("pod metrics is empty")
					continue
				}

				pods, ok = d.pods.Load().([]pod.Pod)
				if !ok {
					log.Info("pod is empty")
					continue
				}

				podModels = make(map[string][]*model.Pod)
				for _, p := range pods {
					if _, ok := podModels[p.Namespace]; !ok {
						podModels[p.Namespace] = make([]*model.Pod, 0)
					}
					if mpod, ok := mpods[p.Name]; ok {
						podModels[p.Namespace] = append(podModels[p.Namespace], &model.Pod{
							Name:        p.Name,
							Namespace:   p.Namespace,
							MemoryLimit: p.MemLimit,
							MemoryUsage: mpod.Mem,
						})
					}
				}

				jobModels = make(map[string][]*model.Job)
				jobs, ok = d.jobs.Load().([]job.Job)
				if !ok {
					log.Info("job is empty")
				} else {
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
				}

				tmpl, ok := d.jobTemplate.Load().(string)
				if !ok {
					log.Info("job template is empty")
					continue
				}
				err = d.decoder.DecodeInto([]byte(tmpl), &jobTpl)
				if err != nil {
					log.Infof("fails decoding template: %s", err.Error())
					continue
				}
				if prevJobTpl == nil {
					prevJobTpl = &jobTpl
				} else {
					if !equality.Semantic.DeepEqual(*prevJobTpl, jobTpl) {
						prevJobTpl = &jobTpl
					}
				}

				// TODO: cache specified reconciled result based on agentResourceType.
				switch d.agentResourceType {
				case "statefulset":
					ss, ok = d.statefulSets.Load().(statefulset.StatefulSet)
					if !ok {
						log.Info("statefulset is empty")
						continue
					}

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

				// TODO: export below logic to other internal function.
				var (
					mmu        map[string]float64 = make(map[string]float64) // max memory usage
					amu        map[string]float64 = make(map[string]float64) // average memory usage
					rate       map[string]float64 = make(map[string]float64) // rabalance rate
					maxPodName map[string]string  = make(map[string]string)
				)
				if prevSsModel != nil {
					for ns, psm := range prevSsModel {
						if _, ok := ssModel[ns]; !ok {
							continue
						}
						if *ssModel[ns].DesiredReplicas != ssModel[ns].Replicas {
							continue
						}
						if *psm.DesiredReplicas > *ssModel[ns].DesiredReplicas {
							for _, prevPod := range prevPodModels[ns] {
								var ok bool
								for _, pod := range podModels[ns] {
									if prevPod.Name != pod.Name {
										ok = true
										break
									}
								}
								if !ok {
									log.Infof("[decrease] creating job for pod %s", prevPod.Name)
									if err := d.createJob(ctx, jobTpl, prevPod.Name); err != nil {
										log.Errorf("failed to create job: %s", err)
										continue
									}
									prevSsModel[ns] = ssModel[ns]
									prevPodModels[ns] = podModels[ns]
								}
							}
						} else {
							for _, p := range podModels[ns] {
								u := p.MemoryUsage / p.MemoryLimit
								amu[ns] += u
								if u > mmu[ns] {
									mmu[ns] = u
									maxPodName[ns] = p.Name
								}
							}
							amu[ns] = amu[ns] / float64(len(podModels[ns]))
							bias := mmu[ns] - amu[ns]
							log.Infof("bias: %v, tolerance: %v, maxPodName: %v, average memory usage: %v", bias, d.tolerance, maxPodName[ns], amu[ns])
							if bias > d.tolerance {
								rate[ns] = 1 - (amu[ns] / mmu[ns])
								var ok bool
								for _, jobs := range jobModels {
									for _, job := range jobs {
										if job.Type == "rebalance" && job.Active != 0 && job.TargetAgentNamespace == ns {
											ok = true
											break
										}
									}

									if ok {
										break
									}
								}
								if !ok || len(jobModels[ns]) == 0 {
									log.Infof("[bias] creating job for pod %s", maxPodName[ns])
									if err := d.createJob(ctx, jobTpl, maxPodName[ns]); err != nil {
										log.Errorf("failed to create job: %s", err)
										continue
									}
									prevSsModel = ssModel
									prevPodModels = podModels
								}
							}
						}
					}
				} else {
					// Store reconciled result for next loop.
					prevSsModel = ssModel
					prevPodModels = podModels
				}

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

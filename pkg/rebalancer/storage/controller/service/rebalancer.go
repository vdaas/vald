//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package service manages the main logic of server.
package service

import (
	"context"
	"math"
	"reflect"
	"strconv"
	"sync"
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

	// TODO: move to internal after internal/k8s refactoring
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"
)

// Rebalancer represents the rebalancer interface.
type Rebalancer interface {
	Start(ctx context.Context) (<-chan error, error)
}

type reason uint8

const (
	BIAS reason = iota
	RECOVERY
	MANUAL

	jobType = "rebalancer"
)

func (r reason) String() string {
	switch r {
	case BIAS:
		return "bias"
	case RECOVERY:
		return "recovery"
	case MANUAL:
		return "manual"
	default:
		return "unknown"
	}
}

type rebalancer struct {
	podName      string
	podNamespace string

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

	jobConfigmapName      string
	jobConfigmapNamespace string

	statefulSets atomic.Value

	leaderElectionID string

	rcd           time.Duration // reconcile check duration
	eg            errgroup.Group
	ctrl          k8s.Controller
	tolerance     float64
	rateThreshold float64
	decoder       *conversion.Decoder
}

// NewRebalancer initialize job, configmap, pod, podMetrics, statefulset, replicaset, daemonset reconciler.
// And it returns the rebalancer implemenation or any error occurred.
func NewRebalancer(opts ...RebalancerOption) (Rebalancer, error) {
	r := new(rebalancer)

	for _, opt := range append(defaultRebalancerOpts, opts...) {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	job, err := job.New(
		job.WithControllerName("job rebalancer"),
		job.WithNamespaces(r.jobNamespace),
		job.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		job.WithOnReconcileFunc(func(jobList map[string][]job.Job) {
			log.Debugf("[reconcile Job] length Joblist: %d", len(jobList))
			jobs, ok := jobList[r.jobName]
			if ok {
				r.jobs.Store(jobs)
			} else {
				r.jobs.Store(make([]job.Job, 0))
				log.Infof("job not found: %s", r.jobName)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	cm, err := configmap.New(
		configmap.WithControllerName("configmap rebalancer"),
		configmap.WithNamespaces(r.jobConfigmapNamespace),
		configmap.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		configmap.WithOnReconcileFunc(func(configmapList map[string][]configmap.ConfigMap) {
			configmaps, ok := configmapList[r.jobConfigmapNamespace]
			if ok {
				for _, cm := range configmaps {
					if cm.Name == r.jobConfigmapName {
						if tmpl, ok := cm.Data[r.jobTemplateKey]; ok {
							r.jobTemplate.Store(tmpl)
						} else {
							log.Infof("job template not found: %s", r.jobTemplateKey)
						}
						break
					}
				}
			} else {
				log.Infof("configmap not found: %s", r.jobConfigmapName)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	desiredAgentReplicas := make([]int32, 0)
	var rc k8s.ResourceController
	switch r.agentResourceType {
	case "statefulset":
		rc, err = statefulset.New(
			statefulset.WithControllerName("statefulset rebalancer"),
			statefulset.WithNamespaces(r.agentNamespace),
			statefulset.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			statefulset.WithOnReconcileFunc(func(statefulSetList map[string][]statefulset.StatefulSet) {
				log.Debugf("[reconcile StatefulSet] length StatefulSet[%s]: %d", r.agentName, len(statefulSetList))
				sss, ok := statefulSetList[r.agentName]
				if ok {
					log.Debugf("[reconcile StatefulSet] StatefulSet[%s]: desired replica: %d, current replica: %d", r.agentName, *sss[0].Spec.Replicas, sss[0].Status.Replicas)
					if len(sss) == 1 {
						if *sss[0].Spec.Replicas != sss[0].Status.Replicas {
							log.Debug("[test] skip store")
							return
						}

						pss, ok := r.statefulSets.Load().(statefulset.StatefulSet)
						if ok && *sss[0].Spec.Replicas < *pss.Spec.Replicas {
							mu.Lock()
							desiredAgentReplicas = append(desiredAgentReplicas, *sss[0].Spec.Replicas)
							mu.Unlock()
						}
						r.statefulSets.Store(sss[0])
					} else {
						log.Infof("too many statefulset list: want 1, but %r", len(sss))
					}
				} else {
					log.Infof("statefuleset not found: %s", r.agentName)
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
		return nil, errors.New("invalid agent resource type: " + r.agentResourceType)
	}

	r.ctrl, err = k8s.New(
		k8s.WithControllerName("rebalance storage controller"),
		k8s.WithEnableLeaderElection(),
		k8s.WithLeaderElectionID(r.leaderElectionID),
		k8s.WithResourceController(job),
		k8s.WithResourceController(rc), // statefulset controller
		k8s.WithResourceController(pod.New(
			pod.WithControllerName("pod discover"),
			pod.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			pod.WithOnReconcileFunc(func(podList map[string][]pod.Pod) {
				log.Debugf("[reconcile pod] length podList[%s]: %d", r.agentName, len(podList[r.agentName]))
				pods, ok := podList[r.agentName]
				if ok {
					mu.Lock()
					dar := make([]int32, len(desiredAgentReplicas))
					_ = copy(dar, desiredAgentReplicas)
					log.Debugf("[test] before creating job: desiredAgentReplicas: %v", dar)
					mu.Unlock()

					// var muPod sync.Mutex
					// muPod.Lock()
					// defer muPod.Unlock()
					if len(dar) > 0 {
						ppod, ok := r.pods.Load().([]pod.Pod)
						log.Debugf("[test] len(pod): %d, len(prev pod): %d", len(pods), len(ppod))
						if ok && len(pods) < len(ppod) && len(pods) == int(dar[0]) {
							decreasedPodNames := getDecreasedPodNames(ppod, pods, r.agentNamespace)
							jobTpl, err := r.genJobTpl()
							if err != nil {
								log.Infof("[recovery] error generating job template: %s", err.Error())
								return
							}
							// create jobs
							for _, name := range decreasedPodNames {
								log.Debugf("[recovery] creating job for pod %s", name)
								ctx := context.TODO()
								if err := r.createJob(ctx, *jobTpl, RECOVERY, name, r.agentNamespace); err != nil {
									log.Errorf("[recovery] failed to create job: %s", err)
									continue
								}
							}
							mu.Lock()
							desiredAgentReplicas = desiredAgentReplicas[1:]
							log.Debugf("[test] after creating job: desiredAgentReplicas: %v", desiredAgentReplicas)
							mu.Unlock()

							log.Debug("[test] Store pods when len(dar) > 0")
							r.pods.Store(pods)
						}
					} else {
						ss, ok := r.statefulSets.Load().(statefulset.StatefulSet)
						if !ok || *ss.Spec.Replicas != int32(len(pods)) {
							log.Debugf("[test] Skip store pods, ss.Spec.Replica: %d\tlen(pods): %d", *ss.Spec.Replicas != int32(len(pods)))
							return
						}
						log.Debug("[test] Store pods when len(dar) = 0")
						r.pods.Store(pods)
					}
				} else {
					log.Infof("pod not found: %s", r.agentName)
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
					r.podMetrics.Store(podList)
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

	r.decoder, err = conversion.NewDecoder(runtime.NewScheme())
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Start starts the rebalancer controller loop for the Vald agent index rebalancer.
func (r *rebalancer) Start(ctx context.Context) (<-chan error, error) {
	cech, err := r.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}

	ech := make(chan error, 1)
	r.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		dt := time.NewTicker(r.rcd)
		defer dt.Stop()

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

				podModels, err := r.genPodModels()
				if err != nil {
					log.Infof("error generating pod models: %s", err.Error())
					continue
				}

				jobModels, err = r.genJobModels()
				if err != nil {
					log.Infof("error generating job models: %s", err.Error())
				}

				jobTpl, err = r.genJobTpl()
				if err != nil {
					log.Infof("error generating job template: %s", err.Error())
					continue
				}

				// TODO: cache specified reconciled result based on agentResourceType.
				switch r.agentResourceType {
				case "statefulset":
					ss, ok = r.statefulSets.Load().(statefulset.StatefulSet)
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

					for ns, sm := range ssModel {
						if *sm.DesiredReplicas != sm.Replicas {
							log.Debugf("[not desired] desired replica: %d, current replica: %d", *sm.DesiredReplicas, sm.Replicas)
							continue
						}

						maxPodName, rate := r.getBiasOverDetail(podModels[ns])
						if maxPodName == "" || rate < r.rateThreshold {
							log.Debugf("[rate/podname checking] pod name, rate, rateThreshold: %s, %.3f, %f", maxPodName, rate, r.rateThreshold)
							continue
						}
						log.Debugf("[bias/jobcheck] job: %#v", jobModels[r.jobNamespace])

						if !r.isJobRunning(jobModels, ns) {
							log.Debugf("[bias] creating job for pod %s, rate: %v", maxPodName, rate)
							if err := r.createJob(ctx, *jobTpl, BIAS, maxPodName, ns); err != nil {
								log.Errorf("[bias] failed to create job: %s", err)
								continue
							}
						} else {
							log.Debugf("[bias] job is already running")
						}
					}

				default:
					// TODO: define error for return
					return nil
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

func (r *rebalancer) createJob(ctx context.Context, jobTpl job.Job, reason reason, agentName, agentNs string) error {
	jobTpl.Name += "-" + strconv.FormatInt(time.Now().Unix(), 10)

	if len(r.jobNamespace) != 0 {
		jobTpl.Namespace = r.jobNamespace
	}

	if jobTpl.Labels == nil {
		jobTpl.Labels = make(map[string]string)
	}
	jobTpl.Labels["type"] = jobType
	jobTpl.Labels["reason"] = reason.String()
	jobTpl.Labels["target_agent_name"] = agentName
	jobTpl.Labels["target_agent_namespace"] = agentNs
	jobTpl.Labels["controller_name"] = r.podName
	jobTpl.Labels["controller_namespace"] = r.podNamespace

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

func (r *rebalancer) genPodModels() (podModels map[string][]*model.Pod, err error) {
	mpods, ok := r.podMetrics.Load().(map[string]mpod.Pod)
	if !ok {
		return nil, errors.New("pod metrics is empty")
	}

	pods, ok := r.pods.Load().([]pod.Pod)
	if !ok {
		return nil, errors.New("pod is empty")
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

	return
}

func (r *rebalancer) genJobModels() (jobModels map[string][]*model.Job, err error) {
	jobs, ok := r.jobs.Load().([]job.Job)
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
			Reason:               j.Labels["reason"],
			TargetAgentNamespace: j.Labels["target_agent_namespace"],
			TargetAgentName:      j.Labels["target_agent_name"],
			ControllerNamespace:  j.Labels["controller_namespace"],
			ControllerName:       j.Labels["controller_name"],
		})
	}

	return
}

func (r *rebalancer) genJobTpl() (jobTpl *job.Job, err error) {
	tmpl, ok := r.jobTemplate.Load().(string)
	if !ok {
		return nil, errors.New("job template is empty")
	}
	jobTpl = &job.Job{}
	err = r.decoder.DecodeInto([]byte(tmpl), jobTpl)
	if err != nil {
		return nil, errors.Wrap(err, "fails decoding template")
	}
	return
}

func getDecreasedPodNames(prev, cur []pod.Pod, ns string) (podNames []string) {
	podNames = make([]string, 0)
	for _, prevPod := range prev {
		if prevPod.Namespace != ns {
			continue
		}
		var ok bool
		for _, pod := range cur {
			if pod.Namespace != ns {
				continue
			}
			if prevPod.Name == pod.Name {
				ok = true
				break
			}
		}
		if !ok {
			podNames = append(podNames, prevPod.Name)
		}
	}
	return
}

func (r *rebalancer) getBiasOverDetail(pm []*model.Pod) (string, float64) {
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

		if maxMemUsg >= (1+r.tolerance)*sig {
			return podName, 1 - (avgMemUsg / maxMemUsg)
		}
		return "", 0
	}

	podName, avgMemUsg, maxMemUsg := calAvgMemUsgWithMemLimit(pm)

	if maxMemUsg >= avgMemUsg+r.tolerance {
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

func (r *rebalancer) isJobRunning(jobModels map[string][]*model.Job, ns string) bool {
	for _, jobs := range jobModels {
		for _, job := range jobs {
			if job.Type == jobType && job.Active != 0 && job.TargetAgentNamespace == ns {
				return true
			}
		}
	}
	return false
}

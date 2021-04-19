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
	"sync/atomic"
	"time"

	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	agent "github.com/vdaas/vald/internal/client/v1/client/agent/core"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/configmap"
	"github.com/vdaas/vald/internal/k8s/decoder"
	"github.com/vdaas/vald/internal/k8s/job"
	mpod "github.com/vdaas/vald/internal/k8s/metrics/pod"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/k8s/statefulset"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/rebalancer/storage/controller/config"
	"github.com/vdaas/vald/pkg/rebalancer/storage/controller/model"
)

const (
	qualifiedNamePrefix string = "rebalancer.vald.vdaas.org/"
)

// Rebalancer represents the rebalancer interface.
type Rebalancer interface {
	Start(ctx context.Context) (<-chan error, error)
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
	agentResourceType config.AgentResourceType
	pods              atomic.Value
	podMetrics        atomic.Value

	jobConfigmapName string

	statefulSets atomic.Value

	leaderElectionID string

	rcd           time.Duration // reconcile check duration
	eg            errgroup.Group
	ctrl          k8s.Controller
	tolerance     float64
	rateThreshold float64
	decoder       *decoder.Decoder
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

	err := r.initCtrl()
	if err != nil {
		return nil, err
	}

	r.decoder, err = decoder.NewDecoder(r.ctrl.GetManager().GetScheme())
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *rebalancer) initCtrl() (err error) {
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
				return
			}

			r.jobs.Store(make([]job.Job, 0))
			log.Infof("job not found: %s", r.jobName)
		}),
	)
	if err != nil {
		return err
	}

	// TODO: delete this logic
	cm, err := configmap.New(
		configmap.WithControllerName("configmap rebalancer"),
		configmap.WithNamespaces(r.jobNamespace),
		configmap.WithOnErrorFunc(func(err error) {
			log.Error(err)
		}),
		configmap.WithOnReconcileFunc(func(configmapList map[string][]configmap.ConfigMap) {
			configmaps, ok := configmapList[r.jobNamespace]
			if !ok {
				log.Infof("configmap not found: %s", r.jobConfigmapName)
				return
			}
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
		}),
	)
	if err != nil {
		return err
	}

	var rc k8s.ResourceController
	switch r.agentResourceType {
	case config.STATEFULSET:
		rc, err = statefulset.New(
			statefulset.WithControllerName("statefulset rebalancer"),
			statefulset.WithNamespaces(r.agentNamespace),
			statefulset.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			statefulset.WithOnReconcileFunc(func(statefulSetList map[string][]statefulset.StatefulSet) {
				log.Debugf("[reconcile StatefulSet] length StatefulSet[%s]: %d", r.agentName, len(statefulSetList))
				sss, ok := statefulSetList[r.agentName]
				if !ok {
					log.Infof("statefuleset not found: %s", r.agentName)
					return
				}

				// If there are multiple sets of stateful agents in the namespace or none
				// it will return early because it is an unexpected error.
				if len(sss) != 1 {
					log.Infof("too many statefulset list: want 1, but %r", len(sss))
					return
				}

				log.Debugf("[reconcile StatefulSet] StatefulSet[%s]: desired replica: %d, current replica: %d", r.agentName, *sss[0].Spec.Replicas, sss[0].Status.Replicas)

				pss, ok := r.statefulSets.Load().(statefulset.StatefulSet)
				if ok && *sss[0].Spec.Replicas < *pss.Spec.Replicas {
					jobTpl, err := r.genJobTpl()
					if err != nil {
						log.Errorf("[recovery] error generating job template: %s", err.Error())
					} else {
						for i := int(*pss.Spec.Replicas); i > int(*sss[0].Spec.Replicas); i-- {
							name := r.agentName + "-" + strconv.Itoa(i-1)
							log.Debugf("[recovery] creating job for pod %s", name)
							ctx := context.TODO()
							if err := r.createJob(ctx, *jobTpl, config.RECOVERY, name, r.agentNamespace, 1); err != nil {
								log.Errorf("[recovery] failed to create job: %s", err)
							}
						}
					}
				}
				r.statefulSets.Store(sss[0])
			}),
		)
		if err != nil {
			return err
		}
	case config.REPLICASET:
		// TODO: implment get replicaset reconciled result
		return nil
	case config.DAEMONSET:
		// TODO: implment get daemonset reconciled result
		return nil
	default:
		return errors.ErrInvalidAgentResourceType(r.agentResourceType.String())
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
				if !ok {
					log.Infof("pod not found: %s", r.agentName)
					return
				}
				r.pods.Store(pods)
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
		return err
	}

	return nil
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

					podModels       map[string][]*model.Pod
					namespaceByJobs map[string][]job.Job
					ssModel         map[string]*model.StatefulSet
					jobTpl          *job.Job
				)

				podModels, err := r.genPodModels()
				if err != nil {
					log.Infof("error generating pod models: %s", err.Error())
					continue
				}

				namespaceByJobs, err = r.namespaceByJobs()
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
				case config.STATEFULSET:
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

						if !r.isJobRunning(namespaceByJobs, ns) {
							log.Debugf("[deviation] creating job for pod %s, rate: %v", maxPodName, rate)
							if err := r.createJob(ctx, *jobTpl, config.DEVIATION, maxPodName, ns, rate); err != nil {
								log.Errorf("[deviation] failed to create job: %s", err)
								continue
							}
						} else {
							log.Debugf("[deviation] job is already running")
						}
					}

				default:
					return errors.ErrInvalidAgentResourceType(r.agentResourceType.String())
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

func (r *rebalancer) getPodByAgentName(agentName string) (*pod.Pod, error) {
	agentPods, ok := r.pods.Load().([]pod.Pod)
	if !ok {
		return nil, errors.New("cannot read agent pod list")
	}
	// 2. get target agent pod by name
	for _, p := range agentPods {
		if p.Name == agentName {
			return &p, nil
		}
	}
	return nil, errors.New("pod not found")
}

func (r *rebalancer) createJob(ctx context.Context, jobTpl job.Job, reason config.RebalanceReason, agentName, agentNs string, rate float64) error {
	// check indexing or not
	if reason == config.BIAS {
		p, err := r.getPodByAgentName(agentName)
		if err != nil {
			return err
		}
		c, err := agent.New(agent.WithAddrs(p.IP))
		if err != nil {
			return err
		}
		res, err := c.IndexInfo(ctx, new(payload.Empty))
		if err != nil {
			return err
		}

		if res.GetIndexing() {
			return errors.New("pod is indexing, job will not be created")
		}

		// if saving flag = true
		//     return error
	}
	jobTpl.Name += "-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	if len(r.jobNamespace) != 0 {
		jobTpl.Namespace = r.jobNamespace
	}

	if jobTpl.Labels == nil {
		jobTpl.Labels = make(map[string]string)
	}
	jobTpl.Labels[qualifiedNamePrefix+"reason"] = reason.String()
	jobTpl.Labels[qualifiedNamePrefix+"target_agent_name"] = agentName
	jobTpl.Labels[qualifiedNamePrefix+"target_agent_namespace"] = agentNs

	if jobTpl.Annotations == nil {
		jobTpl.Annotations = make(map[string]string)
	}
	jobTpl.Annotations[qualifiedNamePrefix+"controller_name"] = r.podName
	jobTpl.Annotations[qualifiedNamePrefix+"controller_namespace"] = r.podNamespace
	if rate > 0 {
		jobTpl.Annotations[qualifiedNamePrefix+"rate"] = strconv.FormatFloat(rate, 'f', 4, 64)
	}

	log.Debugf("jobTpl.Labels: %#v\n", jobTpl.Labels)
	log.Debugf("jobTpl.Annotations: %#v\n", jobTpl.Annotations)

	jobTpl.Spec.Template.Labels = jobTpl.Labels
	jobTpl.Spec.Template.ObjectMeta.Annotations = jobTpl.Annotations

	log.Debugf("jobTpl.Spec.Template.Labels: %#v\n", jobTpl.Spec.Template.Labels)
	log.Debugf("jobTpl.Spec.Template.Annotations: %#v\n", jobTpl.Spec.Template.Annotations)

	jobTpl.Spec.Template.ObjectMeta.Labels[qualifiedNamePrefix+"target_agent_name"] = agentName
	if rate > 0 {
		jobTpl.Spec.Template.ObjectMeta.Annotations[qualifiedNamePrefix+"rate"] = strconv.FormatFloat(rate, 'f', 4, 64)
	}

	c := r.ctrl.GetManager().GetClient()
	if err := c.Create(ctx, &jobTpl); err != nil {
		return errors.ErrK8sFailedToCreateJob(err)
	}

	return nil
}

func (r *rebalancer) genPodModels() (podModels map[string][]*model.Pod, err error) {
	mpods, ok := r.podMetrics.Load().(map[string]mpod.Pod)
	if !ok {
		return nil, errors.ErrEmptyReconcileResult("pod metrics")
	}

	pods, ok := r.pods.Load().([]pod.Pod)
	if !ok {
		return nil, errors.ErrEmptyReconcileResult("pod")
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

func (r *rebalancer) namespaceByJobs() (jobmap map[string][]job.Job, err error) {
	jobs, ok := r.jobs.Load().([]job.Job)
	if !ok {
		return nil, errors.ErrEmptyReconcileResult("job")
	}

	jobmap = make(map[string][]job.Job)
	for _, j := range jobs {
		if _, ok := jobmap[j.Namespace]; !ok {
			jobmap[j.Namespace] = make([]job.Job, 0)
		}
		jobmap[j.Namespace] = append(jobmap[j.Namespace], j)
	}

	return
}

func (r *rebalancer) genJobTpl() (jobTpl *job.Job, err error) {
	tmpl, ok := r.jobTemplate.Load().(string)
	if !ok {
		return nil, errors.ErrJobTemplateNotFound()
	}
	jobTpl = &job.Job{}
	err = r.decoder.DecodeInto([]byte(tmpl), jobTpl)
	if err != nil {
		return nil, errors.ErrFailedToDecodeJobTemplate(err)
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

func (r *rebalancer) isJobRunning(jobsmap map[string][]job.Job, ns string) bool {
	for _, jobs := range jobsmap {
		for _, job := range jobs {
			if job.Labels[qualifiedNamePrefix+"reason"] != config.MANUAL.String() && job.Status.Active != 0 && job.Labels[qualifiedNamePrefix+"target_agent_namespace"] == ns {
				return true
			}
		}
	}
	return false
}

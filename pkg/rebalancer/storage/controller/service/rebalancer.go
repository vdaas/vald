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
	"unsafe"

	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/backoff"
	agent "github.com/vdaas/vald/internal/client/v1/client/agent/core"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/decoder"
	"github.com/vdaas/vald/internal/k8s/job"
	mpod "github.com/vdaas/vald/internal/k8s/metrics/pod"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/k8s/statefulset"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/rebalancer/storage/controller/config"
	"github.com/vdaas/vald/pkg/rebalancer/storage/controller/model"
)

const (
	qualifiedNamePrefix string = "rebalancer.vald.vdaas.org/"
)

var (
	defaultGRPCOps = []grpc.Option{
		grpc.WithConnectionPoolSize(1),
		grpc.WithInsecure(true),
		grpc.WithResolveDNS(false),
	}
)

// Rebalancer represents the rebalancer interface.
type Rebalancer interface {
	Start(ctx context.Context) (<-chan error, error)
	Close(ctx context.Context) error
}

type rebalancer struct {
	podName      string
	podNamespace string

	jobs                    atomic.Value
	jobName                 string
	jobNamespace            string
	jobTemplate             string   // row manifest template data of rebalance job.
	jobObject               *job.Job // object generated from template.
	currentDeviationJobName atomic.Value

	agentName         string
	agentPort         int
	agentNamespace    string
	agentResourceType config.AgentResourceType
	pods              atomic.Value
	podMetrics        atomic.Value

	statefulSets atomic.Value

	leaderElectionID string

	rcd           time.Duration // reconcile check duration
	eg            errgroup.Group
	ctrl          k8s.Controller
	tolerance     float64
	rateThreshold float64
	decoder       *decoder.Decoder

	prevDesiredPods int // prev desired replica number in statefulset

	agentClient grpc.Client
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

	r.jobObject = new(job.Job)
	err = r.decoder.DecodeInto(*(*[]byte)(unsafe.Pointer(&r.jobTemplate)), r.jobObject)
	if err != nil {
		return nil, errors.ErrFailedToDecodeJobTemplate(err)
	}

	if r.agentClient == nil {
		r.agentClient = grpc.New(
			grpc.WithBackoff(backoff.New()),
		)
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
		job.WithOnReconcileFunc(r.jobReconcile),
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
			statefulset.WithOnReconcileFunc(r.statefulsetReconcile),
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
		k8s.WithResourceController(rc),
		k8s.WithResourceController(pod.New(
			pod.WithControllerName("pod discover"),
			pod.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			pod.WithOnReconcileFunc(r.podReconcile),
		)),
		k8s.WithResourceController(mpod.New(
			mpod.WithControllerName("pod metrics discover"),
			mpod.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			mpod.WithOnReconcileFunc(r.podMetricsReconcile),
		)),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *rebalancer) podReconcile(podList map[string][]pod.Pod) {
	log.Debugf("[reconcile pod] length podList[%s]: %d", r.agentName, len(podList[r.agentName]))
	pods, ok := podList[r.agentName]
	if !ok {
		log.Infof("pod not found: %s", r.agentName)
		return
	}

	for _, pod := range pods {
		addr := net.JoinHostPort(pod.IP, uint16(r.agentPort))
		if !r.agentClient.IsConnected(context.Background(), addr) { // TODO
			_, err := r.agentClient.Connect(context.Background(), addr)
			if err != nil {
				// error handling
			}
		}
	}
	r.pods.Store(pods)
}

func (r *rebalancer) jobReconcile(jobList map[string][]job.Job) {
	log.Debugf("[reconcile Job] length Joblist: %d", len(jobList))
	jobs, ok := jobList[r.jobName]
	if ok {
		r.jobs.Store(jobs)
		return
	}

	r.jobs.Store(make([]job.Job, 0))
	log.Infof("job not found: %s", r.jobName)
}

func (r *rebalancer) statefulsetReconcile(ctx context.Context, statefulSetList map[string][]statefulset.StatefulSet) {
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

	var (
		desiredReplica = int(*sss[0].Spec.Replicas)
		currentReplica = int(sss[0].Status.Replicas)
	)

	log.Debugf("[reconcile StatefulSet] StatefulSet[%s]: desired replica: %d, current replica: %d", r.agentName, desiredReplica, currentReplica)

	// skip job creation when current pod status is not desired
	if currentReplica != desiredReplica {
		r.statefulSets.Store(sss[0])
		return
	}

	// if current desired replica number is less than previous desired replica number, create recovery job
	if currentReplica < r.prevDesiredPods {
		// If the number of replicas is reduced from 3 to 2, get the name of the pod below.
		// vald-agent-0
		// vald-agent-1
		// vald-agent-2 <- deteced
		for i := r.prevDesiredPods; i > currentReplica; i-- {
			name := r.agentName + "-" + strconv.Itoa(i-1)

			if err := r.deleteDeviationJob(ctx); err != nil {
				log.Error(err)
			}
			log.Debugf("[recovery] creating job for pod %s", name)
			if err := r.createJob(ctx, *r.jobObject, config.RECOVERY, name, r.agentNamespace, 1); err != nil {
				log.Errorf("[recovery] failed to create job: %s", err)
			}
		}
	}
	r.prevDesiredPods = currentReplica
	r.statefulSets.Store(sss[0])
}

func (r *rebalancer) podMetricsReconcile(podList map[string]mpod.Pod) {
	if len(podList) > 0 {
		r.podMetrics.Store(podList)
		return
	}
	log.Info("pod metrics not found")
}

func (r *rebalancer) Close(ctx context.Context) error {
	return r.agentClient.Close(ctx)
}

// Start starts the rebalancer controller loop for the Vald agent index rebalancer.
func (r *rebalancer) Start(ctx context.Context) (<-chan error, error) {
	cech, err := r.ctrl.Start(ctx) // close
	if err != nil {
		return nil, err
	}

	gech, err := r.agentClient.StartConnectionMonitor(ctx)
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
			case err := <-gech:
				ech <- err
			case <-dt.C:
				podsByNamespace, err := r.podsByNamespace()
				if err != nil {
					log.Infof("error generating pod models: %s", err.Error())
					continue
				}

				jobsByNamespace, err := r.getJobsByNamespace()
				if err != nil {
					log.Infof("error generating job models: %s", err.Error())
				}

				// TODO: cache specified reconciled result based on agentResourceType.
				switch r.agentResourceType {
				case config.STATEFULSET:
					ss, ok := r.statefulSets.Load().(statefulset.StatefulSet)
					if !ok {
						log.Info("statefulset is empty")
						continue
					}

					// Because there is a possibility that it will be supported for each namespace in the future,
					// declared it with a map type using the namespace as a key.
					ssModelByNamespace := make(map[string]*model.StatefulSet)
					ssModelByNamespace[ss.Namespace] = &model.StatefulSet{
						Name:            ss.Name,
						Namespace:       ss.Namespace,
						DesiredReplicas: ss.Spec.Replicas,
						Replicas:        ss.Status.Replicas,
					}

					for ns, sm := range ssModelByNamespace {

						// If the current number of replicas does not match the desired number of replicas,
						// the number of pods is changing and no processing is performed.
						if *sm.DesiredReplicas != sm.Replicas {
							log.Debugf("[not desired] desired replica: %d, current replica: %d", *sm.DesiredReplicas, sm.Replicas)
							continue
						}

						maxPodName, rate := r.getBiasOverDetail(podsByNamespace[ns])
						if len(maxPodName) == 0 || rate < r.rateThreshold {
							log.Debugf("[rate/podname checking] pod name, rate, rateThreshold: %s, %.3f, %f", maxPodName, rate, r.rateThreshold)
							continue
						}

						if !r.isJobRunning(jobsByNamespace, ns) {
							log.Debugf("[deviation] creating job for pod %s, rate: %v", maxPodName, rate)
							if err := r.createJob(ctx, *r.jobObject, config.DEVIATION, maxPodName, ns, rate); err != nil {
								log.Errorf("[deviation] failed to create job: %s", err)
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

func (r *rebalancer) deleteDeviationJob(ctx context.Context) error {
	jobName, ok := r.currentDeviationJobName.Load().(string)
	if !ok || jobName == "" {
		return nil
	}
	job := new(job.Job)

	c := r.ctrl.GetManager().GetClient()
	err := c.Get(ctx, k8s.ObjectKey{
		Namespace: r.jobNamespace,
		Name:      jobName,
	}, job)
	if err != nil {
		return err
	}

	// if the deviation job is not running, skip
	if job.Status.Active == 0 {
		return nil
	}

	if err := c.Delete(ctx, job); err != nil {
		return errors.ErrK8sFailedToDeleteJob(err)
	}
	r.currentDeviationJobName.Store("")
	return nil
}

func (r *rebalancer) createJob(ctx context.Context, jobTpl job.Job, reason config.RebalanceReason, agentName, agentNs string, rate float64) (err error) {
	if reason == config.DEVIATION {
		p, err := r.getPodByAgentName(agentName)
		if err != nil {
			return err
		}

		agentAddr := net.JoinHostPort(p.IP, uint16(r.agentPort))
		var res *payload.Info_Index_Count

		_, err = r.agentClient.Do(ctx, agentAddr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			res, err = agent.NewAgentClient(conn).IndexInfo(ctx, new(payload.Empty))
			return res, err
		})
		if err != nil {
			return err
		}

		if res.GetIndexing() || res.GetSaving() {
			return errors.New("pod is indexing, job will not be created")
		}
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

	if reason == config.DEVIATION {
		r.currentDeviationJobName.Store(jobTpl.Name)
	}

	return nil
}

func (r *rebalancer) podsByNamespace() (podModels map[string][]*model.Pod, err error) {
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

func (r *rebalancer) getJobsByNamespace() (jobmap map[string][]job.Job, err error) {
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

func getDecreasedPodNames(prev, cur []pod.Pod, ns string) (podNames []string) {
	podNames = make([]string, 0, len(prev))
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

func (r *rebalancer) isJobRunning(jobsByNamespace map[string][]job.Job, ns string) bool {
	for _, jobs := range jobsByNamespace {
		for _, job := range jobs {
			if job.Labels[qualifiedNamePrefix+"reason"] != config.MANUAL.String() && job.Status.Active != 0 && job.Labels[qualifiedNamePrefix+"target_agent_namespace"] == ns {
				return true
			}
		}
	}
	return false
}

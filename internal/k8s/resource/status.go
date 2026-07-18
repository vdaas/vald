//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package resource

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/vdaas/vald/internal/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceStatus int

const (
	StatusUnknown       ResourceStatus = iota // Unknown state
	StatusPending                             // Resource is initializing or waiting for update observation
	StatusUpdating                            // Resource is in the process of updating/rolling out changes
	StatusAvailable                           // Resource is fully operational
	StatusDegraded                            // Resource is operational but with some issues
	StatusFailed                              // Resource has failed
	StatusCompleted                           // For jobs: execution completed successfully
	StatusScheduled                           // For jobs: scheduled but not yet started
	StatusScaling                             // Resource is scaling up or down
	StatusPaused                              // Resource update/rollout is paused
	StatusTerminating                         // Resource (e.g., Pod) is in the process of termination
	StatusNotReady                            // Resource (e.g., Pod) is running but not yet ready
	StatusLoadBalancing                       // Service is still provisioning a load balancer
)

// WaitForStatus polls the named object (or, when name is empty, every object
// matching labelSelector) until its evaluated status is one of statuses, the
// context is canceled, or the default timeout expires. The polling skeleton is
// shared with ObjectClient.Wait through waitLoop; only the per-tick step
// differs, and the name/labelSelector split is resolved once instead of per
// tick.
func WaitForStatus[T Object, L ObjectList, C NamedObject, I ResourceInterface[T, L, C]](
	ctx context.Context, client I, name string, labelSelector string, statuses ...ResourceStatus,
) (matched bool, err error) {
	var obj T
	if !slices.ContainsFunc(possibleStatuses(obj), func(st ResourceStatus) bool {
		return slices.Contains(statuses, st)
	}) {
		return false, errors.ErrStatusPatternNeverMatched
	}

	var step func(context.Context) (bool, error)
	if name != "" {
		step = func(ctx context.Context) (bool, error) {
			obj, err := client.Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			status, info, err := checkResourceState(obj)
			if err != nil {
				return false, errors.Wrap(err, info)
			}
			return slices.Contains(statuses, status), nil
		}
	} else {
		opts := metav1.ListOptions{}
		if labelSelector != "" {
			opts.LabelSelector = labelSelector
		}
		step = func(ctx context.Context) (bool, error) {
			l, err := client.List(ctx, opts)
			if err != nil {
				return false, err
			}
			objs, err := apimeta.ExtractList(l)
			if err != nil {
				return false, errors.Wrap(err, "failed to extract items")
			}
			if len(objs) == 0 {
				return false, nil
			}
			for _, obj := range objs {
				status, info, err := checkResourceState(obj)
				if err != nil {
					return false, errors.Wrap(err, info)
				}
				if !slices.Contains(statuses, status) {
					return false, nil
				}
			}
			return true, nil
		}
	}
	return waitLoop(ctx, errors.ErrTimeoutWaitingForResourceStatus, step)
}

// kindStatusInfo bundles the statuses a Kind can report with the evaluator
// that computes them, so possibleStatuses and checkResourceState share one
// per-Kind entry instead of enumerating the same Kind set in two separate
// type switches that must be kept in sync by hand.
type kindStatusInfo struct {
	evaluate func(obj any) (ResourceStatus, string, error)
	possible []ResourceStatus
}

// evaluatorFor wraps a Kind-specific evaluator so it can be stored in
// kindStatusTable under an any-typed signature. registerKind derives the table
// key and the evaluator's argument type from the same type parameter R, so a
// mismatched key/evaluator pair cannot be registered; the ok check remains as
// a defensive fallback for direct table manipulation.
func evaluatorFor[R any](
	f func(R) (ResourceStatus, string, error),
) func(any) (ResourceStatus, string, error) {
	return func(obj any) (ResourceStatus, string, error) {
		res, ok := obj.(R)
		if !ok {
			return StatusUnknown, "Unsupported resource type", errors.ErrUnsupportedKubernetesResourceType(obj)
		}
		return f(res)
	}
}

// kindStatusTable is a read-only-after-init dispatch table shared by
// possibleStatuses/checkResourceState; it is built by newKindStatusTable
// (rather than a var + init()) so the registration order stays a plain,
// testable function instead of hidden package-init side effects.
//
//nolint:gochecknoglobals
var kindStatusTable = newKindStatusTable()

// registerKind ties the statuses a Kind can report to its evaluator under the
// single type parameter R (inferred from eval), so the table key and the
// evaluator's argument type cannot drift apart the way hand-paired map
// literal entries could.
func registerKind[R Object](
	t map[reflect.Type]kindStatusInfo,
	possible []ResourceStatus,
	eval func(R) (ResourceStatus, string, error),
) {
	t[reflect.TypeFor[R]()] = kindStatusInfo{
		possible: possible,
		evaluate: evaluatorFor(eval),
	}
}

func newKindStatusTable() map[reflect.Type]kindStatusInfo {
	t := make(map[reflect.Type]kindStatusInfo)
	registerKind(t, []ResourceStatus{StatusPending, StatusUpdating, StatusAvailable, StatusDegraded, StatusFailed, StatusPaused}, evaluateDeployment)
	registerKind(t, []ResourceStatus{StatusPending, StatusUpdating, StatusAvailable, StatusDegraded, StatusFailed}, evaluateStatefulSet)
	registerKind(t, []ResourceStatus{StatusPending, StatusUpdating, StatusAvailable, StatusDegraded, StatusFailed}, evaluateDaemonSet)
	registerKind(t, []ResourceStatus{StatusUpdating, StatusFailed, StatusCompleted, StatusScheduled}, evaluateJob)
	registerKind(t, []ResourceStatus{StatusPaused, StatusPending, StatusAvailable}, evaluateCronJob)
	registerKind(t, []ResourceStatus{StatusUnknown, StatusAvailable, StatusPending, StatusCompleted, StatusFailed, StatusTerminating, StatusNotReady}, evaluatePod)
	registerKind(t, []ResourceStatus{StatusAvailable, StatusLoadBalancing}, evaluateService)
	registerKind(t, []ResourceStatus{StatusPending, StatusAvailable}, evaluateSecret)
	registerKind(t, []ResourceStatus{StatusPending, StatusAvailable}, evaluateConfigMap)
	registerKind(t, []ResourceStatus{StatusPending, StatusAvailable}, evaluateMutatingWebhookConfiguration)
	registerKind(t, []ResourceStatus{StatusPending, StatusAvailable}, evaluateValidatingWebhookConfiguration)
	return t
}

func possibleStatuses[T Object](obj T) []ResourceStatus {
	if info, ok := kindStatusTable[reflect.TypeOf(obj)]; ok {
		return info.possible
	}
	return []ResourceStatus{StatusUnknown}
}

// checkResourceState dispatches obj to its Kind's evaluator via
// kindStatusTable. It takes any because dispatch happens on the dynamic type;
// callers hold either a concrete T or a runtime.Object from apimeta.ExtractList.
func checkResourceState(obj any) (ResourceStatus, string, error) {
	if info, ok := kindStatusTable[reflect.TypeOf(obj)]; ok {
		return info.evaluate(obj)
	}
	return StatusUnknown, "Unsupported resource type", errors.ErrUnsupportedKubernetesResourceType(obj)
}

func evaluateDeployment(deploy *appsv1.Deployment) (ResourceStatus, string, error) {
	desired := int32(1)
	if deploy.Spec.Replicas != nil {
		desired = *deploy.Spec.Replicas
	}

	// details builds the status string lazily so the fully-operational path
	// pays no formatting cost.
	details := func(msg string) string {
		return fmt.Sprintf("Name: %s, Generation: %d, ObservedGeneration: %d, Spec.Replicas: %d, Status.Replicas: %d, UpdatedReplicas: %d, AvailableReplicas: %d.%s",
			deploy.GetName(), deploy.GetGeneration(), deploy.Status.ObservedGeneration, desired, deploy.Status.Replicas, deploy.Status.UpdatedReplicas, deploy.Status.AvailableReplicas, msg)
	}

	if deploy.Spec.Paused {
		return StatusPaused, details("Deployment is paused."), nil
	}

	if deploy.Status.ObservedGeneration < deploy.Generation {
		return StatusPending, details("Update not yet observed by controller."), nil
	}

	var progressingCond, availableCond *appsv1.DeploymentCondition
	for i := range deploy.Status.Conditions {
		switch deploy.Status.Conditions[i].Type {
		case appsv1.DeploymentProgressing:
			progressingCond = &deploy.Status.Conditions[i]
		case appsv1.DeploymentAvailable:
			availableCond = &deploy.Status.Conditions[i]
		case appsv1.DeploymentReplicaFailure:
			// Not evaluated for status; listed explicitly so this switch
			// stays exhaustive over DeploymentConditionType.
		}
	}
	if progressingCond != nil && progressingCond.Status == corev1.ConditionFalse {
		return StatusFailed, details(fmt.Sprintf("Progressing condition: %s, Status: %s.", progressingCond.Reason, progressingCond.Status)), nil
	}
	if availableCond != nil && availableCond.Status == corev1.ConditionFalse {
		return StatusDegraded, details(fmt.Sprintf("Available condition: %s, Status: %s.", availableCond.Reason, availableCond.Status)), nil
	}

	if deploy.Status.UpdatedReplicas < desired {
		return StatusUpdating, details(fmt.Sprintf("Only %d out of %d replicas updated.", deploy.Status.UpdatedReplicas, desired)), nil
	}
	if deploy.Status.UpdatedReplicas < deploy.Status.Replicas {
		return StatusUpdating, details(fmt.Sprintf("There are %d total replicas but only %d replicas updated.", deploy.Status.Replicas, deploy.Status.UpdatedReplicas)), nil
	}
	if deploy.Status.AvailableReplicas < desired {
		return StatusDegraded, details(fmt.Sprintf("Only %d out of %d replicas available.", deploy.Status.AvailableReplicas, desired)), nil
	}

	return StatusAvailable, "Deployment is fully operational.", nil
}

func evaluateStatefulSet(sts *appsv1.StatefulSet) (ResourceStatus, string, error) {
	desired := int32(1)
	if sts.Spec.Replicas != nil {
		desired = *sts.Spec.Replicas
	}

	details := func(msg string) string {
		return fmt.Sprintf(
			"Name: %s, Generation: %d, ObservedGeneration: %d, Spec.Replicas: %d, CurrentReplicas: %d, UpdatedReplicas: %d, ReadyReplicas: %d, CurrentRevision: %s, UpdateRevision: %s.%s",
			sts.GetName(),
			sts.GetGeneration(),
			sts.Status.ObservedGeneration,
			desired,
			sts.Status.CurrentReplicas,
			sts.Status.UpdatedReplicas,
			sts.Status.ReadyReplicas,
			sts.Status.CurrentRevision,
			sts.Status.UpdateRevision,
			msg,
		)
	}

	if sts.Status.ObservedGeneration < sts.Generation {
		return StatusPending, details("Update not yet observed by controller."), nil
	}

	if sts.Status.UpdatedReplicas < desired {
		return StatusUpdating, details(fmt.Sprintf("Only %d out of %d replicas updated.", sts.Status.UpdatedReplicas, desired)), nil
	}

	if sts.Status.CurrentReplicas < desired {
		return StatusUpdating, details(fmt.Sprintf("Only %d out of %d replicas are currently running.", sts.Status.CurrentReplicas, desired)), nil
	}

	if sts.Status.ReadyReplicas < desired {
		return StatusDegraded, details(fmt.Sprintf("Only %d out of %d replicas are ready.", sts.Status.ReadyReplicas, desired)), nil
	}

	if sts.Status.UpdateRevision != sts.Status.CurrentRevision {
		return StatusUpdating, details(fmt.Sprintf("Revision mismatch: CurrentRevision=%s, UpdateRevision=%s.", sts.Status.CurrentRevision, sts.Status.UpdateRevision)), nil
	}

	return StatusAvailable, "StatefulSet is fully operational.", nil
}

func evaluateDaemonSet(ds *appsv1.DaemonSet) (ResourceStatus, string, error) {
	details := func(msg string) string {
		return fmt.Sprintf("Name: %s, Generation: %d, ObservedGeneration: %d, DesiredNumberScheduled: %d, UpdatedNumberScheduled: %d, NumberAvailable: %d, NumberReady: %d.%s",
			ds.GetName(), ds.GetGeneration(), ds.Status.ObservedGeneration, ds.Status.DesiredNumberScheduled, ds.Status.UpdatedNumberScheduled, ds.Status.NumberAvailable, ds.Status.NumberReady, msg)
	}

	if ds.Status.ObservedGeneration < ds.Generation {
		return StatusPending, details("Update not yet observed by controller."), nil
	}

	// Check DaemonSet conditions if present (not always available)
	for _, cond := range ds.Status.Conditions {
		if cond.Type == appsv1.DaemonSetConditionType("") && cond.Status == corev1.ConditionFalse {
			return StatusFailed, details(fmt.Sprintf("Condition %s is false: %s.", cond.Type, cond.Reason)), nil
		}
	}

	if ds.Status.UpdatedNumberScheduled < ds.Status.DesiredNumberScheduled {
		return StatusUpdating, details(fmt.Sprintf("Only %d out of %d pods updated.", ds.Status.UpdatedNumberScheduled, ds.Status.DesiredNumberScheduled)), nil
	}

	if ds.Status.NumberAvailable < ds.Status.DesiredNumberScheduled {
		return StatusDegraded, details(fmt.Sprintf("Only %d out of %d pods available.", ds.Status.NumberAvailable, ds.Status.DesiredNumberScheduled)), nil
	}

	return StatusAvailable, "DaemonSet is fully operational.", nil
}

func evaluateJob(job *batchv1.Job) (ResourceStatus, string, error) {
	details := func(msg string) string {
		return fmt.Sprintf("Name: %s, Active: %d, Succeeded: %d, Failed: %d.%s", job.GetName(), job.Status.Active, job.Status.Succeeded, job.Status.Failed, msg)
	}

	for _, cond := range job.Status.Conditions {
		if cond.Status != corev1.ConditionTrue {
			continue
		}
		switch cond.Type {
		case batchv1.JobFailed:
			return StatusFailed, details(fmt.Sprintf("Condition Type: %s, Status: %s, Reason: %s.", cond.Type, cond.Status, cond.Reason)), nil
		case batchv1.JobComplete:
			return StatusCompleted, "Job has completed.", nil
		case batchv1.JobSuspended, batchv1.JobFailureTarget, batchv1.JobSuccessCriteriaMet:
			// Not evaluated for status; listed explicitly so this switch stays
			// exhaustive over JobConditionType.
		}
	}

	if job.Status.Succeeded > 0 {
		return StatusCompleted, "Job has succeeded.", nil
	}
	if job.Status.Active > 0 {
		return StatusUpdating, details("Job is currently running."), nil
	}

	// If none of the above, the job may be scheduled but not yet started.
	return StatusScheduled, details("Job is scheduled but not yet started."), nil
}

func evaluateCronJob(cronjob *batchv1.CronJob) (ResourceStatus, string, error) {
	if cronjob.Spec.Suspend != nil && *cronjob.Spec.Suspend {
		return StatusPaused, fmt.Sprintf("CronJob Name: %s.CronJob is suspended.", cronjob.GetName()), nil
	}

	if cronjob.Status.LastScheduleTime == nil {
		return StatusPending, fmt.Sprintf("CronJob Name: %s.CronJob has not yet scheduled any jobs.", cronjob.GetName()), nil
	}

	return StatusAvailable, "CronJob is scheduled.", nil
}

func evaluatePod(pod *corev1.Pod) (ResourceStatus, string, error) {
	details := func(msg string) string {
		return fmt.Sprintf("Pod %s Phase: %s.%s", pod.GetName(), pod.Status.Phase, msg)
	}

	if pod.DeletionTimestamp != nil {
		return StatusTerminating, details(fmt.Sprintf("Pod is terminating (DeletionTimestamp: %s).", pod.DeletionTimestamp.Format(time.RFC3339))), nil
	}

	switch pod.Status.Phase {
	case corev1.PodPending:
		return StatusPending, details("Pod is pending scheduling or initialization."), nil
	case corev1.PodRunning:
		ready := false
		var readyCond *corev1.PodCondition
		for i := range pod.Status.Conditions {
			if pod.Status.Conditions[i].Type == corev1.PodReady {
				readyCond = &pod.Status.Conditions[i]
				ready = readyCond.Status == corev1.ConditionTrue
			}
		}
		if !ready {
			msg := "Pod is running but not ready."
			if readyCond != nil {
				msg = fmt.Sprintf("PodReady condition: %s (Reason: %s).%s", readyCond.Status, readyCond.Reason, msg)
			}
			return StatusNotReady, details(msg), nil
		}
		return StatusAvailable, "Pod is running and ready.", nil
	case corev1.PodSucceeded:
		return StatusCompleted, "Pod has completed successfully.", nil
	case corev1.PodFailed:
		return StatusFailed, details("Pod execution has failed."), nil
	case corev1.PodUnknown:
		return StatusUnknown, details("Pod status is unknown."), nil
	default:
		return StatusUnknown, details("Pod status unrecognized."), nil
	}
}

func evaluateService(svc *corev1.Service) (ResourceStatus, string, error) {
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer && len(svc.Status.LoadBalancer.Ingress) == 0 {
		return StatusLoadBalancing,
			fmt.Sprintf("Name %s, Service Type: %s. ClusterIP: %s.LoadBalancer ingress not yet assigned.",
				svc.GetName(), svc.Spec.Type, svc.Spec.ClusterIP), nil
	}
	return StatusAvailable, "Service is operational.", nil
}

func evaluateSecret(secret *corev1.Secret) (ResourceStatus, string, error) {
	if len(secret.Data) == 0 {
		return StatusPending, fmt.Sprintf("Secret Name: %s, Data entries: 0.Secret has no data yet.", secret.GetName()), nil
	}
	return StatusAvailable, "Secret data is populated.", nil
}

func evaluateConfigMap(cm *corev1.ConfigMap) (ResourceStatus, string, error) {
	if len(cm.Data)+len(cm.BinaryData) == 0 {
		return StatusPending, fmt.Sprintf("ConfigMap Name: %s, Data entries: 0, BinaryData entries: 0.ConfigMap has no data yet.", cm.GetName()), nil
	}
	return StatusAvailable, "ConfigMap data is populated.", nil
}

func evaluateWebhookConfig[W any](
	kind, name string, webhooks []W, info func(W) (string, []byte),
) (ResourceStatus, string, error) {
	if len(webhooks) == 0 {
		return StatusPending, fmt.Sprintf("%s Name: %s, Webhooks: 0.No webhooks registered yet.", kind, name), nil
	}
	for _, wh := range webhooks {
		whName, bundle := info(wh)
		if len(bundle) == 0 {
			return StatusPending, fmt.Sprintf("%s Name: %s, Webhooks: %d.CABundle is not yet injected for webhook %s.", kind, name, len(webhooks), whName), nil
		}
	}
	return StatusAvailable, "All webhooks have CABundle injected.", nil
}

func evaluateMutatingWebhookConfiguration(
	mwc *admissionregistrationv1.MutatingWebhookConfiguration,
) (ResourceStatus, string, error) {
	return evaluateWebhookConfig("MutatingWebhookConfiguration", mwc.GetName(), mwc.Webhooks,
		func(wh admissionregistrationv1.MutatingWebhook) (string, []byte) {
			return wh.Name, wh.ClientConfig.CABundle
		})
}

func evaluateValidatingWebhookConfiguration(
	vwc *admissionregistrationv1.ValidatingWebhookConfiguration,
) (ResourceStatus, string, error) {
	return evaluateWebhookConfig("ValidatingWebhookConfiguration", vwc.GetName(), vwc.Webhooks,
		func(wh admissionregistrationv1.ValidatingWebhook) (string, []byte) {
			return wh.Name, wh.ClientConfig.CABundle
		})
}

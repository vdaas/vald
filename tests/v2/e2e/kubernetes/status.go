//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package kubernetes

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/vdaas/vald/internal/errors"
	appsv1 "k8s.io/api/apps/v1"             // For Deployment, StatefulSet, DaemonSet
	batchv1 "k8s.io/api/batch/v1"           // For Job and CronJob
	corev1 "k8s.io/api/core/v1"             // For Pod, PersistentVolumeClaim, Service
	networkingv1 "k8s.io/api/networking/v1" // For Ingress
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// --------------------------------------------------------------------------------
// ResourceStatus enum with extended states for detailed status reporting.
// --------------------------------------------------------------------------------
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
	StatusBound                               // PVC is bound to a volume
	StatusLoadBalancing                       // Service is still provisioning a load balancer
)

// Human-readable mapping for ResourceStatus values.
var ResourceStatusMap = map[ResourceStatus]string{
	StatusUnknown:       "Unknown state",
	StatusPending:       "Initializing or waiting for update observation",
	StatusUpdating:      "Updating / Rolling out a new version",
	StatusAvailable:     "Fully operational",
	StatusDegraded:      "Degraded state",
	StatusFailed:        "Failed",
	StatusCompleted:     "Completed successfully",
	StatusScheduled:     "Scheduled but not started",
	StatusScaling:       "Scaling in progress",
	StatusPaused:        "Rollout paused",
	StatusTerminating:   "Terminating",
	StatusNotReady:      "Running but not ready",
	StatusBound:         "PVC is bound",
	StatusLoadBalancing: "Load balancer provisioning in progress",
}

func extractItems[T any](obj any) ([]T, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	items := v.FieldByName("Items")
	if !items.IsValid() || items.Kind() != reflect.Slice {
		return nil, fmt.Errorf("field 'Items' not found or not a slice")
	}

	out := make([]T, items.Len())
	for i := 0; i < items.Len(); i++ {
		v := items.Index(i)
		if v.CanAddr() {
			if ptr, ok := v.Addr().Interface().(T); ok {
				out[i] = ptr
				continue
			}
		}
		val, ok := v.Interface().(T) // fallback to value match
		if !ok {
			return nil, fmt.Errorf("item at index %d is not of type T", i)
		}
		out[i] = val
	}
	return out, nil
}

// --------------------------------------------------------------------------------
// WaitForStatus waits for specific Kubernetes resources to reach a specific status.
// The function checks the status of the resources at regular intervals and returns
// objects, a boolean indicating if the status matched, and an error (if any).
// The function supports Deployment, StatefulSet, DaemonSet, Job, CronJob, Pod,
// PersistentVolumeClaim, Service, and Ingress.
// --------------------------------------------------------------------------------
func WaitForStatus[T Object, L ObjectList, C NamedObject, I ResourceInterface[T, L, C]](
	ctx context.Context, client I, name string, labelSelector string, statuses ...ResourceStatus,
) (matched bool, err error) {
	var obj T
	if !slices.ContainsFunc(PossibleStatuses(obj), func(st ResourceStatus) bool {
		return slices.Contains(statuses, st)
	}) {
		return false, errors.ErrStatusPatternNeverMatched
	}

	ticker := time.NewTicker(5 * time.Second)
	timeout := time.NewTimer(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <- timeout.C:
			return false, errors.New("timeout waiting for resource status")
		case <-ticker.C:
			opts := metav1.ListOptions{}
			if name != "" {
				obj, err := client.Get(ctx, name, metav1.GetOptions{})
				if err != nil {
					return false, err
				}
				status, info, err := CheckResourceState(obj)
				if err != nil {
					return false, errors.Wrap(err, info)
				}
				return slices.Contains(statuses, status), nil
			}
			if labelSelector != "" {
				opts.LabelSelector = labelSelector
			}
			l, err := client.List(ctx, opts)
			if err != nil {
				return false, err
			}
			matched = true
			items, err := extractItems[T](l)
			// if no resources found yet, keep polling
			if len(items) == 0 {
				continue
			}
			if err != nil {
				return false, errors.Wrap(err, "failed to extract items")
			}
			for _, obj := range items {
				status, info, err := CheckResourceState(obj)
				if err != nil {
					return false, errors.Wrap(err, info)
				}
				if !slices.Contains(statuses, status) {
					matched = false
				}
			}
			if matched {
				return true, nil
			}
		}
	}
}

// --------------------------------------------------------------------------------
// PossibleStatuses returns a list of possible ResourceStatus values for a given object.
// The function supports Deployment, StatefulSet, DaemonSet, Job, CronJob, Pod,
// PersistentVolumeClaim, Service, and Ingress.
// --------------------------------------------------------------------------------
func PossibleStatuses[T Object](obj T) []ResourceStatus {
	switch any(obj).(type) {
	case *appsv1.Deployment:
		return []ResourceStatus{StatusPending, StatusUpdating, StatusAvailable, StatusDegraded, StatusFailed, StatusPaused}
	case *appsv1.StatefulSet:
		return []ResourceStatus{StatusPending, StatusUpdating, StatusAvailable, StatusDegraded, StatusFailed}
	case *appsv1.DaemonSet:
		return []ResourceStatus{StatusPending, StatusUpdating, StatusAvailable, StatusDegraded, StatusFailed}
	case *batchv1.Job:
		return []ResourceStatus{StatusUpdating, StatusFailed, StatusCompleted, StatusScheduled}
	case *batchv1.CronJob:
		return []ResourceStatus{StatusPaused, StatusPending, StatusAvailable}
	case *corev1.Pod:
		return []ResourceStatus{StatusUnknown, StatusAvailable, StatusPending, StatusCompleted, StatusFailed, StatusTerminating, StatusNotReady}
	case *corev1.PersistentVolumeClaim:
		return []ResourceStatus{StatusUnknown, StatusPending, StatusFailed, StatusBound}
	case *corev1.Service:
		return []ResourceStatus{StatusAvailable, StatusLoadBalancing}
	case *networkingv1.Ingress:
		return []ResourceStatus{StatusPending, StatusAvailable}
	default:
		return []ResourceStatus{StatusUnknown}
	}
}

// --------------------------------------------------------------------------------
// checkResourceState determines the detailed state of a Kubernetes resource.
// It returns a ResourceStatus enum, a detailed string message, and an error (if any).
// This function supports Deployment, StatefulSet, DaemonSet, Job, CronJob, Pod,
// PersistentVolumeClaim, Service, and Ingress.
// --------------------------------------------------------------------------------
func CheckResourceState[T Object](obj T) (ResourceStatus, string, error) {
	switch res := any(obj).(type) {
	case *appsv1.Deployment:
		return evaluateDeployment(res)
	case *appsv1.StatefulSet:
		return evaluateStatefulSet(res)
	case *appsv1.DaemonSet:
		return evaluateDaemonSet(res)
	case *batchv1.Job:
		return evaluateJob(res)
	case *batchv1.CronJob:
		return evaluateCronJob(res)
	case *corev1.Pod:
		return evaluatePod(res)
	case *corev1.PersistentVolumeClaim:
		return evaluatePVC(res)
	case *corev1.Service:
		return evaluateService(res)
	case *networkingv1.Ingress:
		return evaluateIngress(res)
	default:
		return StatusUnknown, "Unsupported resource type", errors.ErrUnsupportedKubernetesResourceType(obj)
	}
}

// --------------------------------------------------------------------------------
// evaluateDeployment evaluates the status of a Deployment resource.
// It checks:
// - Generation vs ObservedGeneration
// - Spec.Replicas (desired) vs Status.Replicas, UpdatedReplicas, and AvailableReplicas
// - Conditions: DeploymentProgressing and DeploymentAvailable
// - Whether the deployment is paused
// --------------------------------------------------------------------------------
func evaluateDeployment(deploy *appsv1.Deployment) (ResourceStatus, string, error) {
	desired := int32(1)
	if deploy.Spec.Replicas != nil {
		desired = *deploy.Spec.Replicas
	}

	// Build a detailed status string.
	statusDetails := fmt.Sprintf("Name: %s, Generation: %d, ObservedGeneration: %d, Spec.Replicas: %d, Status.Replicas: %d, UpdatedReplicas: %d, AvailableReplicas: %d.",
		deploy.GetName(), deploy.GetGeneration(), deploy.Status.ObservedGeneration, desired, deploy.Status.Replicas, deploy.Status.UpdatedReplicas, deploy.Status.AvailableReplicas)

	// Check if the Deployment is paused.
	if deploy.Spec.Paused {
		statusDetails += "Deployment is paused."
		return StatusPaused, statusDetails, nil
	}

	// Ensure the controller has observed the latest update.
	if deploy.Status.ObservedGeneration < deploy.Generation {
		statusDetails += "Update not yet observed by controller."
		return StatusPending, statusDetails, nil
	}

	// Inspect Deployment conditions.
	var progressingCond *appsv1.DeploymentCondition
	var availableCond *appsv1.DeploymentCondition
	for _, cond := range deploy.Status.Conditions {
		if cond.Type == appsv1.DeploymentProgressing {
			progressingCond = &cond
		} else if cond.Type == appsv1.DeploymentAvailable {
			availableCond = &cond
		}
	}
	if progressingCond != nil {
		statusDetails += fmt.Sprintf("Progressing condition: %s, Status: %s.", progressingCond.Reason, progressingCond.Status)
		if progressingCond.Status == corev1.ConditionFalse {
			return StatusFailed, statusDetails, nil
		}
	}
	if availableCond != nil {
		statusDetails += fmt.Sprintf("Available condition: %s, Status: %s.", availableCond.Reason, availableCond.Status)
		if availableCond.Status == corev1.ConditionFalse {
			return StatusDegraded, statusDetails, nil
		}
	}

	// Check if the number of updated and available replicas meets the desired count.
	if deploy.Status.UpdatedReplicas < desired {
		statusDetails += fmt.Sprintf("Only %d out of %d replicas updated.", deploy.Status.UpdatedReplicas, desired)
		return StatusUpdating, statusDetails, nil
	}
	if deploy.Status.UpdatedReplicas < deploy.Status.Replicas {
		statusDetails += fmt.Sprintf("There are %d total replicas but only %d replicas updated.", deploy.Status.UpdatedReplicas, deploy.Status.Replicas)
		return StatusUpdating, statusDetails, nil
	}
	if deploy.Status.AvailableReplicas < desired {
		statusDetails += fmt.Sprintf("Only %d out of %d replicas available.", deploy.Status.AvailableReplicas, desired)
		return StatusDegraded, statusDetails, nil
	}

	statusDetails += "Deployment is fully operational."
	return StatusAvailable, statusDetails, nil
}

// --------------------------------------------------------------------------------
// evaluateStatefulSet evaluates the status of a StatefulSet resource.
// It checks:
// - Generation vs ObservedGeneration
// - Spec.Replicas (desired) vs UpdatedReplicas, CurrentReplicas, and ReadyReplicas
// - Whether UpdateRevision equals CurrentRevision
// --------------------------------------------------------------------------------
func evaluateStatefulSet(sts *appsv1.StatefulSet) (ResourceStatus, string, error) {
	desired := int32(1)
	if sts.Spec.Replicas != nil {
		desired = *sts.Spec.Replicas
	}

	statusDetails := fmt.Sprintf(
		"Name: %s, Generation: %d, ObservedGeneration: %d, Spec.Replicas: %d, CurrentReplicas: %d, UpdatedReplicas: %d, ReadyReplicas: %d, CurrentRevision: %s, UpdateRevision: %s.",
		sts.GetName(),
		sts.GetGeneration(),
		sts.Status.ObservedGeneration,
		desired,
		sts.Status.CurrentReplicas,
		sts.Status.UpdatedReplicas,
		sts.Status.ReadyReplicas,
		sts.Status.CurrentRevision,
		sts.Status.UpdateRevision,
	)

	if sts.Status.ObservedGeneration < sts.Generation {
		statusDetails += "Update not yet observed by controller."
		return StatusPending, statusDetails, nil
	}

	if sts.Status.UpdatedReplicas < desired {
		statusDetails += fmt.Sprintf("Only %d out of %d replicas updated.", sts.Status.UpdatedReplicas, desired)
		return StatusUpdating, statusDetails, nil
	}

	if sts.Status.CurrentReplicas < desired {
		statusDetails += fmt.Sprintf("Only %d out of %d replicas are currently running.", sts.Status.CurrentReplicas, desired)
		return StatusUpdating, statusDetails, nil
	}

	if sts.Status.ReadyReplicas < desired {
		statusDetails += fmt.Sprintf("Only %d out of %d replicas are ready.", sts.Status.ReadyReplicas, desired)
		return StatusDegraded, statusDetails, nil
	}

	if sts.Status.UpdateRevision != sts.Status.CurrentRevision {
		statusDetails += fmt.Sprintf("Revision mismatch: CurrentRevision=%s, UpdateRevision=%s.", sts.Status.CurrentRevision, sts.Status.UpdateRevision)
		return StatusUpdating, statusDetails, nil
	}

	statusDetails += "StatefulSet is fully operational."
	return StatusAvailable, statusDetails, nil
}

// --------------------------------------------------------------------------------
// evaluateDaemonSet evaluates the status of a DaemonSet resource.
// It checks:
// - Generation vs ObservedGeneration
// - DesiredNumberScheduled vs UpdatedNumberScheduled and NumberAvailable
// - Conditions if available for additional insights
// --------------------------------------------------------------------------------
func evaluateDaemonSet(ds *appsv1.DaemonSet) (ResourceStatus, string, error) {
	statusDetails := fmt.Sprintf("Name: %s, Generation: %d, ObservedGeneration: %d, DesiredNumberScheduled: %d, UpdatedNumberScheduled: %d, NumberAvailable: %d, NumberReady: %d.",
		ds.GetName(), ds.GetGeneration(), ds.Status.ObservedGeneration, ds.Status.DesiredNumberScheduled, ds.Status.UpdatedNumberScheduled, ds.Status.NumberAvailable, ds.Status.NumberReady)

	if ds.Status.ObservedGeneration < ds.Generation {
		statusDetails += "Update not yet observed by controller."
		return StatusPending, statusDetails, nil
	}

	// Check DaemonSet conditions if present (not always available)
	for _, cond := range ds.Status.Conditions {
		// Using a generic condition check similar to DeploymentProgressing.
		if cond.Type == appsv1.DaemonSetConditionType("") && cond.Status == corev1.ConditionFalse {
			statusDetails += fmt.Sprintf("Condition %s is false: %s.", cond.Type, cond.Reason)
			return StatusFailed, statusDetails, nil
		}
	}

	if ds.Status.UpdatedNumberScheduled < ds.Status.DesiredNumberScheduled {
		statusDetails += fmt.Sprintf("Only %d out of %d pods updated.", ds.Status.UpdatedNumberScheduled, ds.Status.DesiredNumberScheduled)
		return StatusUpdating, statusDetails, nil
	}

	if ds.Status.NumberAvailable < ds.Status.DesiredNumberScheduled {
		statusDetails += fmt.Sprintf("Only %d out of %d pods available.", ds.Status.NumberAvailable, ds.Status.DesiredNumberScheduled)
		return StatusDegraded, statusDetails, nil
	}

	statusDetails += "DaemonSet is fully operational."
	return StatusAvailable, statusDetails, nil
}

// --------------------------------------------------------------------------------
// evaluateJob evaluates the status of a Job resource.
// It checks:
// - Active, Succeeded, and Failed counts
// - Conditions (e.g., JobFailed, JobComplete)
// --------------------------------------------------------------------------------
func evaluateJob(job *batchv1.Job) (ResourceStatus, string, error) {
	statusDetails := fmt.Sprintf("Name: %s, Active: %d, Succeeded: %d, Failed: %d.", job.GetName(), job.Status.Active, job.Status.Succeeded, job.Status.Failed)

	// Check job conditions for additional information.
	for _, cond := range job.Status.Conditions {
		statusDetails += fmt.Sprintf("Condition Type: %s, Status: %s, Reason: %s.", cond.Type, cond.Status, cond.Reason)
		if cond.Type == batchv1.JobFailed && cond.Status == corev1.ConditionTrue {
			return StatusFailed, statusDetails, nil
		}
		if cond.Type == batchv1.JobComplete && cond.Status == corev1.ConditionTrue {
			return StatusCompleted, statusDetails, nil
		}
	}

	if job.Status.Succeeded > 0 {
		statusDetails += "Job has succeeded."
		return StatusCompleted, statusDetails, nil
	}
	if job.Status.Active > 0 {
		statusDetails += "Job is currently running."
		return StatusUpdating, statusDetails, nil
	}

	// If none of the above, the job may be scheduled but not yet started.
	statusDetails += "Job is scheduled but not yet started."
	return StatusScheduled, statusDetails, nil
}

// --------------------------------------------------------------------------------
// evaluateCronJob evaluates the status of a CronJob resource.
// It checks:
// - Whether the CronJob is suspended
// - LastScheduleTime and its timing
// --------------------------------------------------------------------------------
func evaluateCronJob(cronjob *batchv1.CronJob) (ResourceStatus, string, error) {
	statusDetails := fmt.Sprintf("CronJob Name: %s.", cronjob.GetName())

	// Check if the CronJob is suspended.
	if cronjob.Spec.Suspend != nil && *cronjob.Spec.Suspend {
		statusDetails += "CronJob is suspended."
		return StatusPaused, statusDetails, nil
	}

	// Check the last schedule time.
	if cronjob.Status.LastScheduleTime == nil {
		statusDetails += "CronJob has not yet scheduled any jobs."
		return StatusPending, statusDetails, nil
	}

	lastSchedule := cronjob.Status.LastScheduleTime.Format(time.RFC3339)
	statusDetails += fmt.Sprintf("Last scheduled at %s.", lastSchedule)
	return StatusAvailable, statusDetails, nil
}

// --------------------------------------------------------------------------------
// evaluatePod evaluates the status of a Pod resource.
// It checks:
// - Pod Phase (Pending, Running, Succeeded, Failed, Unknown)
// - Pod conditions for readiness (especially PodReady)
// - DeletionTimestamp to detect termination
// --------------------------------------------------------------------------------
func evaluatePod(pod *corev1.Pod) (ResourceStatus, string, error) {
	statusDetails := fmt.Sprintf("Pod %s Phase: %s.", pod.GetName(), pod.Status.Phase)

	// Check if the pod is being terminated.
	if pod.DeletionTimestamp != nil {
		statusDetails += fmt.Sprintf("Pod is terminating (DeletionTimestamp: %s).", pod.DeletionTimestamp.Format(time.RFC3339))
		return StatusTerminating, statusDetails, nil
	}

	// Evaluate based on pod phase.
	switch pod.Status.Phase {
	case corev1.PodPending:
		statusDetails += "Pod is pending scheduling or initialization."
		return StatusPending, statusDetails, nil
	case corev1.PodRunning:
		// Check PodReady condition.
		ready := false
		for _, cond := range pod.Status.Conditions {
			if cond.Type == corev1.PodReady {
				ready = (cond.Status == corev1.ConditionTrue)
				statusDetails += fmt.Sprintf("PodReady condition: %s (Reason: %s).", cond.Status, cond.Reason)
			}
		}
		if !ready {
			statusDetails += "Pod is running but not ready."
			return StatusNotReady, statusDetails, nil
		}
		statusDetails += "Pod is running and ready."
		return StatusAvailable, statusDetails, nil
	case corev1.PodSucceeded:
		statusDetails += "Pod has completed successfully."
		return StatusCompleted, statusDetails, nil
	case corev1.PodFailed:
		statusDetails += "Pod execution has failed."
		return StatusFailed, statusDetails, nil
	case corev1.PodUnknown:
		statusDetails += "Pod status is unknown."
		return StatusUnknown, statusDetails, nil
	default:
		statusDetails += "Pod status unrecognized."
		return StatusUnknown, statusDetails, nil
	}
}

// --------------------------------------------------------------------------------
// evaluatePVC evaluates the status of a PersistentVolumeClaim (PVC).
// It checks:
// - PVC Phase (Bound, Pending, Lost)
// - Provides details on volume binding.
// --------------------------------------------------------------------------------
func evaluatePVC(pvc *corev1.PersistentVolumeClaim) (ResourceStatus, string, error) {
	statusDetails := fmt.Sprintf("PVC %s Phase: %s.", pvc.GetName(), pvc.Status.Phase)
	switch pvc.Status.Phase {
	case corev1.ClaimBound:
		statusDetails += fmt.Sprintf("PVC is bound to volume: %s.", pvc.Spec.VolumeName)
		return StatusBound, statusDetails, nil
	case corev1.ClaimPending:
		statusDetails += "PVC is pending binding to a volume."
		return StatusPending, statusDetails, nil
	case corev1.ClaimLost:
		statusDetails += "PVC has lost its bound volume."
		return StatusFailed, statusDetails, nil
	default:
		statusDetails += "PVC status unrecognized."
		return StatusUnknown, statusDetails, nil
	}
}

// --------------------------------------------------------------------------------
// evaluateService evaluates the status of a Service resource.
// It checks:
// - Service Type (ClusterIP, NodePort, LoadBalancer, ExternalName)
// - For LoadBalancer services, it verifies if ingress information is available.
// --------------------------------------------------------------------------------
func evaluateService(svc *corev1.Service) (ResourceStatus, string, error) {
	statusDetails := fmt.Sprintf("Name %s, Service Type: %s. ClusterIP: %s.", svc.GetName(), svc.Spec.Type, svc.Spec.ClusterIP)
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(svc.Status.LoadBalancer.Ingress) == 0 {
			statusDetails += "LoadBalancer ingress not yet assigned."
			return StatusLoadBalancing, statusDetails, nil
		}
		ingresses := ""
		for _, ingress := range svc.Status.LoadBalancer.Ingress {
			ingresses += ingress.IP + " " + ingress.Hostname + " "
		}
		statusDetails += fmt.Sprintf("LoadBalancer ingress assigned: %s.", ingresses)
	}
	return StatusAvailable, statusDetails + "Service is operational.", nil
}

// --------------------------------------------------------------------------------
// evaluateIngress evaluates the status of an Ingress resource.
// It checks:
// - Whether the Ingress has been assigned an external IP (via LoadBalancer)
// - Provides details on the number of ingress points.
// --------------------------------------------------------------------------------
func evaluateIngress(ing *networkingv1.Ingress) (ResourceStatus, string, error) {
	statusDetails := fmt.Sprintf("Ingress Name: %s.", ing.GetName())
	if len(ing.Status.LoadBalancer.Ingress) == 0 {
		statusDetails += "No external ingress IP assigned yet."
		return StatusPending, statusDetails, nil
	}
	ingresses := ""
	for _, lb := range ing.Status.LoadBalancer.Ingress {
		ingresses += lb.IP + " " + lb.Hostname + " "
	}
	statusDetails += fmt.Sprintf("External ingress IP(s) assigned: %s.", ingresses)
	return StatusAvailable, statusDetails, nil
}

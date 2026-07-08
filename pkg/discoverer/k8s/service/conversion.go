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

// Package service manages the main logic of server.
package service

import (
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	corev1 "k8s.io/api/core/v1"
	metrics "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

// Pod is the discoverer's domain view of a running Kubernetes Pod with its
// averaged per-container resource requests/limits.
type Pod struct {
	Labels      map[string]string
	Annotations map[string]string
	Name        string
	NodeName    string
	Namespace   string
	IP          string
	CPULimit    float64
	CPURequest  float64
	MemLimit    float64
	MemRequest  float64
}

// Node is the discoverer's domain view of a Kubernetes Node with its
// capacity and remaining allocatable resources.
type Node struct {
	Name         string
	InternalAddr string
	ExternalAddr string
	CPUCapacity  float64
	CPURemain    float64
	MemCapacity  float64
	MemRemain    float64
}

// podsByAppName converts a PodList into the discoverer's Pod domain model
// grouped by application name. Terminating pods, pods outside the given
// namespace (when non-empty), and non-running pods are skipped.
// skipcq: CRT-P0006
func podsByAppName(list *corev1.PodList, namespace string) map[string][]Pod {
	var (
		cpuLimit   float64
		cpuRequest float64
		memLimit   float64
		memRequest float64
		pods       = make(map[string][]Pod, len(list.Items))
	)

	// skipcq: CRT-P0006
	for _, pod := range list.Items {
		if pod.GetObjectMeta().GetDeletionTimestamp() != nil ||
			(namespace != "" && !strings.EqualFold(pod.GetNamespace(), namespace)) ||
			pod.Status.Phase != corev1.PodRunning {
			continue
		}
		cpuLimit = 0.0
		memLimit = 0.0
		cpuRequest = 0.0
		memRequest = 0.0
		for _, container := range pod.Spec.Containers {
			request := container.Resources.Requests
			limit := container.Resources.Limits
			cpuLimit += float64(limit.Cpu().Value())
			memLimit += float64(limit.Memory().Value())
			cpuRequest += float64(request.Cpu().Value())
			memRequest += float64(request.Memory().Value())
		}
		cpuLimit /= float64(len(pod.Spec.Containers))
		memLimit /= float64(len(pod.Spec.Containers))
		cpuRequest /= float64(len(pod.Spec.Containers))
		memRequest /= float64(len(pod.Spec.Containers))
		podName, ok := pod.GetObjectMeta().GetLabels()["app"]
		if !ok {
			pns := strings.Split(pod.GetName(), "-")
			podName = strings.Join(pns[:len(pns)-1], "-")
		}

		// No per-group capacity hint: reserving len(list.Items) for every
		// group overallocates by the number of groups (measured ~14x); let
		// append grow each group's slice naturally.
		pods[podName] = append(pods[podName], Pod{
			Name:        pod.GetName(),
			NodeName:    pod.Spec.NodeName,
			Namespace:   pod.GetNamespace(),
			IP:          pod.Status.PodIP,
			CPULimit:    cpuLimit,
			CPURequest:  cpuRequest,
			MemLimit:    memLimit,
			MemRequest:  memRequest,
			Labels:      pod.GetLabels(),
			Annotations: pod.GetAnnotations(),
		})
	}
	return pods
}

// toNodes converts a NodeList into the discoverer's Node domain model.
// Terminating nodes are skipped.
func toNodes(list *corev1.NodeList) []Node {
	nodes := make([]Node, 0, len(list.Items))
	for _, node := range list.Items {
		if node.GetDeletionTimestamp() != nil {
			log.Debugf("reconcile process will be skipped for node: %s, status: %s, deletion timestamp: %s",
				node.GetName(),
				node.Status.Phase,
				node.GetDeletionTimestamp())
			continue
		}
		remain := node.Status.Allocatable
		limit := node.Status.Capacity
		var eip, iip string
		for _, addr := range node.Status.Addresses {
			switch addr.Type {
			case corev1.NodeInternalIP:
				iip = addr.Address
			case corev1.NodeInternalDNS:
				if iip == "" {
					iip = addr.Address
				}
			case corev1.NodeExternalIP:
				eip = addr.Address
			case corev1.NodeExternalDNS:
				if eip == "" {
					eip = addr.Address
				}
			}
		}
		nodes = append(nodes, Node{
			Name:         node.GetName(),
			ExternalAddr: eip,
			InternalAddr: iip,
			CPUCapacity:  float64(limit.Cpu().Value()),
			CPURemain:    float64(remain.Cpu().Value()),
			MemCapacity:  float64(limit.Memory().Value()),
			MemRemain:    float64(remain.Memory().Value()),
		})
	}
	return nodes
}

// podStatusPhaseIndexer extracts the status.phase cache index value for Pods.
func podStatusPhaseIndexer(obj k8s.Object) []string {
	pod, ok := obj.(*corev1.Pod)
	if !ok || pod.GetDeletionTimestamp() != nil {
		return nil
	}
	return []string{string(pod.Status.Phase)}
}

// nodeStatusPhaseIndexer extracts the status.phase cache index value for Nodes.
func nodeStatusPhaseIndexer(obj k8s.Object) []string {
	node, ok := obj.(*corev1.Node)
	if !ok || node.GetDeletionTimestamp() != nil {
		return nil
	}
	return []string{string(node.Status.Phase)}
}

// PodMetrics is the discoverer's domain view of a Pod's averaged
// per-container resource usage from the metrics API.
type PodMetrics struct {
	Name      string
	Namespace string
	CPU       float64
	Mem       float64
}

// NodeMetrics is the discoverer's domain view of a Node's resource usage
// from the metrics API.
type NodeMetrics struct {
	Name    string
	CPU     float64
	Mem     float64
	Pods    int64
	Storage int64
}

// Service is the discoverer's domain view of a Kubernetes Service exposing
// its gRPC/REST API ports.
type Service struct {
	Labels      map[string]string
	Annotations map[string]string
	Name        string
	ClusterIP   string
	ClusterIPs  []string
	Ports       []servicePort
}

type servicePort struct {
	Name string
	Port int32
}

// toPodMetricsMap converts a PodMetricsList into the discoverer's PodMetrics
// domain model keyed by pod name, averaging usage over the containers.
func toPodMetricsMap(list *metrics.PodMetricsList) map[string]PodMetrics {
	var (
		cpuUsage float64
		memUsage float64
		pods     = make(map[string]PodMetrics, len(list.Items))
	)

	for _, pod := range list.Items {
		cpuUsage = 0.0
		memUsage = 0.0
		for _, container := range pod.Containers {
			cpuUsage += float64(container.Usage.Cpu().Value())
			memUsage += float64(container.Usage.Memory().Value())
		}

		cpuUsage /= float64(len(pod.Containers))
		memUsage /= float64(len(pod.Containers))

		pods[pod.GetObjectMeta().GetName()] = PodMetrics{
			Name:      pod.GetName(),
			Namespace: pod.GetNamespace(),
			CPU:       cpuUsage,
			Mem:       memUsage,
		}
	}
	return pods
}

// toNodeMetricsMap converts a NodeMetricsList into the discoverer's
// NodeMetrics domain model keyed by node name.
func toNodeMetricsMap(list *metrics.NodeMetricsList) map[string]NodeMetrics {
	nodes := make(map[string]NodeMetrics, len(list.Items))
	for _, node := range list.Items {
		nodeName := node.GetName()
		nodes[nodeName] = NodeMetrics{
			Name:    nodeName,
			CPU:     float64(node.Usage.Cpu().Value()),
			Mem:     float64(node.Usage.Memory().Value()),
			Storage: node.Usage.StorageEphemeral().Value(),
			Pods:    node.Usage.Pods().Value(),
		}
	}
	return nodes
}

// extractAPIPorts filters the service ports down to the vald API ports.
func extractAPIPorts(ports []corev1.ServicePort) []servicePort {
	var apiPorts []servicePort
	for _, port := range ports {
		if port.Name == "grpc" || port.Name == "rest" {
			apiPorts = append(apiPorts, servicePort{
				Name: port.Name,
				Port: port.Port,
			})
		}
	}
	return apiPorts
}

// toServices converts a ServiceList into the discoverer's Service domain
// model. Terminating services are skipped.
func toServices(list *corev1.ServiceList) []Service {
	svcs := make([]Service, 0, len(list.Items))
	for i := range list.Items {
		svc := &list.Items[i]
		if svc.GetDeletionTimestamp() != nil {
			log.Debugf("reconcile process will be skipped for svc: %s, status: %v, deletion timestamp: %s",
				svc.GetName(),
				svc.Status,
				svc.GetDeletionTimestamp())
			continue
		}
		svcs = append(svcs, Service{
			Name:        svc.GetName(),
			ClusterIP:   svc.Spec.ClusterIP,
			ClusterIPs:  svc.Spec.ClusterIPs,
			Ports:       extractAPIPorts(svc.Spec.Ports),
			Labels:      svc.GetLabels(),
			Annotations: svc.GetAnnotations(),
		})
	}
	return svcs
}

// podMetricsContainersNameIndexer extracts the containers.name cache index
// values for PodMetrics so that field selectors on containers.name work.
func podMetricsContainersNameIndexer(obj k8s.Object) []string {
	pod, ok := obj.(*metrics.PodMetrics)
	if !ok {
		return nil
	}
	res := make([]string, 0, len(pod.Containers))
	for _, pc := range pod.Containers {
		res = append(res, pc.Name)
	}
	return res
}

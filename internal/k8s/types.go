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

package k8s

import (
	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type (
	Object                    = client.Object
	ObjectKey                 = client.ObjectKey
	DeleteAllOfOptions        = client.DeleteAllOfOptions
	DeleteOptions             = client.DeleteOptions
	ListOptions               = client.ListOptions
	ListOption                = client.ListOption
	CreateOption              = client.CreateOption
	CreateOptions             = client.CreateOptions
	GetOption                 = client.GetOption
	GetOptions                = client.GetOptions
	UpdateOptions             = client.UpdateOptions
	MatchingLabels            = client.MatchingLabels
	InNamespace               = client.InNamespace
	VolumeSnapshot            = snapshotv1.VolumeSnapshot
	VolumeSnapshotList        = snapshotv1.VolumeSnapshotList
	Pod                       = corev1.Pod
	PodList                   = corev1.PodList
	Node                      = corev1.Node
	NodeList                  = corev1.NodeList
	Service                   = corev1.Service
	ServiceList               = corev1.ServiceList
	ServicePort               = corev1.ServicePort
	Deployment                = appsv1.Deployment
	DeploymentList            = appsv1.DeploymentList
	ObjectMeta                = metav1.ObjectMeta
	EnvVar                    = corev1.EnvVar
	Job                       = batchv1.Job
	JobList                   = batchv1.JobList
	JobStatus                 = batchv1.JobStatus
	CronJob                   = batchv1.CronJob
	Result                    = reconcile.Result
	OwnerReference            = metav1.OwnerReference
	PersistentVolumeClaim     = corev1.PersistentVolumeClaim
	PersistentVolumeClaimList = corev1.PersistentVolumeClaimList
	PersistentVolumeClaimSpec = corev1.PersistentVolumeClaimSpec
	TypedLocalObjectReference = corev1.TypedLocalObjectReference
	Manager                   = manager.Manager
)

const (
	DeletePropagationBackground = metav1.DeletePropagationBackground
	WatchDeletedEvent           = watch.Deleted
	SelectionOpEquals           = selection.Equals
	SelectionOpExists           = selection.Exists
	PodIndexLabel               = appsv1.PodIndexLabel
	PodRunning                  = corev1.PodRunning
	NodeInternalIP              = corev1.NodeInternalIP
	NodeInternalDNS             = corev1.NodeInternalDNS
	NodeExternalIP              = corev1.NodeExternalIP
	NodeExternalDNS             = corev1.NodeExternalDNS
)

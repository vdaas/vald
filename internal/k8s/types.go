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
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/util/jsonpath"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/yaml"
)

type (
	Client                    = client.Client
	Object                    = client.Object
	ObjectList                = client.ObjectList
	ObjectKey                 = client.ObjectKey
	DeleteAllOfOptions        = client.DeleteAllOfOptions
	DeleteOptions             = client.DeleteOptions
	ListOptions               = client.ListOptions
	ListOption                = client.ListOption
	CreateOption              = client.CreateOption
	GetOption                 = client.GetOption
	MatchingLabels            = client.MatchingLabels
	MatchingFields            = client.MatchingFields
	InNamespace               = client.InNamespace
	VolumeSnapshot            = snapshotv1.VolumeSnapshot
	VolumeSnapshotList        = snapshotv1.VolumeSnapshotList
	Pod                       = corev1.Pod
	PodList                   = corev1.PodList
	PodSpec                   = corev1.PodSpec
	PodStatus                 = corev1.PodStatus
	Container                 = corev1.Container
	Node                      = corev1.Node
	NodeList                  = corev1.NodeList
	Deployment                = appsv1.Deployment
	DeploymentList            = appsv1.DeploymentList
	ObjectMeta                = metav1.ObjectMeta
	Job                       = batchv1.Job
	JobList                   = batchv1.JobList
	JobStatus                 = batchv1.JobStatus
	Result                    = reconcile.Result
	OwnerReference            = metav1.OwnerReference
	PersistentVolumeClaim     = corev1.PersistentVolumeClaim
	PersistentVolumeClaimList = corev1.PersistentVolumeClaimList
	PersistentVolumeClaimSpec = corev1.PersistentVolumeClaimSpec
	TypedLocalObjectReference = corev1.TypedLocalObjectReference
	Manager                   = manager.Manager

	// reconciler interface types
	Request       = ctrl.Request
	Reconciler    = reconcile.Reconciler
	ForOption     = builder.ForOption
	OwnsOption    = builder.OwnsOption
	WatchesOption = builder.WatchesOption
	EventHandler  = handler.EventHandler

	// condition types
	Condition       = metav1.Condition
	ConditionStatus = metav1.ConditionStatus

	// k8s workload types not yet aliased
	Toleration           = corev1.Toleration
	ResourceList         = corev1.ResourceList
	ResourceRequirements = corev1.ResourceRequirements
	ServiceType          = corev1.ServiceType

	// networking types
	PathType = networkingv1.PathType

	// unstructured / schema / scheme types
	Unstructured         = unstructured.Unstructured
	GroupVersionKind     = schema.GroupVersionKind
	GroupVersionResource = schema.GroupVersionResource
	NamespacedName       = types.NamespacedName
	PatchOptions         = metav1.PatchOptions
	Scheme               = runtime.Scheme

	// metrics API types
	APIPodMetrics      = metricsv1beta1.PodMetrics
	APIPodMetricsList  = metricsv1beta1.PodMetricsList
	APINodeMetrics     = metricsv1beta1.NodeMetrics
	APINodeMetricsList = metricsv1beta1.NodeMetricsList
)

const (
	ApplyPatchType              = types.ApplyPatchType
	DeletePropagationBackground = metav1.DeletePropagationBackground
	WatchDeletedEvent           = watch.Deleted
	SelectionOpEquals           = selection.Equals
	SelectionOpExists           = selection.Exists
	PodIndexLabel               = appsv1.PodIndexLabel
	PodRunning                  = corev1.PodRunning
	PodPending                  = corev1.PodPending
	ResourceMemory              = corev1.ResourceMemory
	NodeInternalIP              = corev1.NodeInternalIP
	NodeInternalDNS             = corev1.NodeInternalDNS
	NodeExternalIP              = corev1.NodeExternalIP
	NodeExternalDNS             = corev1.NodeExternalDNS

	// condition status constants
	ConditionTrue    = metav1.ConditionTrue
	ConditionFalse   = metav1.ConditionFalse
	ConditionUnknown = metav1.ConditionUnknown

	// toleration constants
	TolerationOpEqual     = corev1.TolerationOpEqual
	TaintEffectNoSchedule = corev1.TaintEffectNoSchedule

	// service type constants
	ServiceTypeClusterIP    = corev1.ServiceTypeClusterIP
	ServiceTypeLoadBalancer = corev1.ServiceTypeLoadBalancer
	ServiceTypeNodePort     = corev1.ServiceTypeNodePort

	// path type constants
	PathTypeExact                  = networkingv1.PathTypeExact
	PathTypeImplementationSpecific = networkingv1.PathTypeImplementationSpecific
	PathTypePrefix                 = networkingv1.PathTypePrefix
)

// MetricsAddToScheme registers the metrics.k8s.io API types with a scheme.
var MetricsAddToScheme = metricsv1beta1.AddToScheme

var (
	IntOrStringFrom      = intstr.FromString
	IsNotFound           = apierrors.IsNotFound
	NestedFloat64        = unstructured.NestedFloat64
	NewJSONPath          = jsonpath.New
	NewScheme            = runtime.NewScheme
	NewYAMLOrJSONDecoder = utilyaml.NewYAMLOrJSONDecoder
	YAMLMarshal          = yaml.Marshal
	YAMLUnmarshal        = yaml.Unmarshal
)

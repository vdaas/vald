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
)

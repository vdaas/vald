// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package valdrelease

import (
	config "github.com/vdaas/vald/internal/config"
	corev1 "k8s.io/api/core/v1"
)

// Defines values for AgentAlgorithm.
const (
	Faiss AgentAlgorithm = "faiss"
	Ngt   AgentAlgorithm = "ngt"
)

// Defines values for AgentKind.
const (
	AgentKindDaemonSet   AgentKind = "DaemonSet"
	AgentKindDeployment  AgentKind = "Deployment"
	AgentKindStatefulSet AgentKind = "StatefulSet"
)

// Defines values for AgentPodManagementPolicy.
const (
	OrderedReady AgentPodManagementPolicy = "OrderedReady"
	Parallel     AgentPodManagementPolicy = "Parallel"
)

// Defines values for AgentServiceType.
const (
	AgentServiceTypeClusterIP    AgentServiceType = "ClusterIP"
	AgentServiceTypeLoadBalancer AgentServiceType = "LoadBalancer"
	AgentServiceTypeNodePort     AgentServiceType = "NodePort"
)

// Defines values for AgentUnhealthyPodEvictionPolicy.
const (
	AgentUnhealthyPodEvictionPolicyAlwaysAllow     AgentUnhealthyPodEvictionPolicy = "AlwaysAllow"
	AgentUnhealthyPodEvictionPolicyIfHealthyBudget AgentUnhealthyPodEvictionPolicy = "IfHealthyBudget"
)

// Defines values for AgentFaissMethodType.
const (
	Binaryindex AgentFaissMethodType = "binaryindex"
	Ivfpq       AgentFaissMethodType = "ivfpq"
)

// Defines values for AgentFaissMetricType.
const (
	AgentFaissMetricTypeInnerproduct AgentFaissMetricType = "innerproduct"
	AgentFaissMetricTypeL2           AgentFaissMetricType = "l2"
)

// Defines values for AgentNgtDistanceType.
const (
	AgentNgtDistanceTypeAng              AgentNgtDistanceType = "ang"
	AgentNgtDistanceTypeAngle            AgentNgtDistanceType = "angle"
	AgentNgtDistanceTypeCos              AgentNgtDistanceType = "cos"
	AgentNgtDistanceTypeCosine           AgentNgtDistanceType = "cosine"
	AgentNgtDistanceTypeDotproduct       AgentNgtDistanceType = "dotproduct"
	AgentNgtDistanceTypeDp               AgentNgtDistanceType = "dp"
	AgentNgtDistanceTypeHam              AgentNgtDistanceType = "ham"
	AgentNgtDistanceTypeHamming          AgentNgtDistanceType = "hamming"
	AgentNgtDistanceTypeInnerproduct     AgentNgtDistanceType = "innerproduct"
	AgentNgtDistanceTypeIp               AgentNgtDistanceType = "ip"
	AgentNgtDistanceTypeJac              AgentNgtDistanceType = "jac"
	AgentNgtDistanceTypeJaccard          AgentNgtDistanceType = "jaccard"
	AgentNgtDistanceTypeL1               AgentNgtDistanceType = "l1"
	AgentNgtDistanceTypeL2               AgentNgtDistanceType = "l2"
	AgentNgtDistanceTypeLoren            AgentNgtDistanceType = "loren"
	AgentNgtDistanceTypeLorentz          AgentNgtDistanceType = "lorentz"
	AgentNgtDistanceTypeNormalizedangle  AgentNgtDistanceType = "normalizedangle"
	AgentNgtDistanceTypeNormalizedcosine AgentNgtDistanceType = "normalizedcosine"
	AgentNgtDistanceTypeNormalizedl2     AgentNgtDistanceType = "normalizedl2"
	AgentNgtDistanceTypeNormang          AgentNgtDistanceType = "normang"
	AgentNgtDistanceTypeNormcos          AgentNgtDistanceType = "normcos"
	AgentNgtDistanceTypeNorml2           AgentNgtDistanceType = "norml2"
	AgentNgtDistanceTypePoinc            AgentNgtDistanceType = "poinc"
	AgentNgtDistanceTypePoincare         AgentNgtDistanceType = "poincare"
	AgentNgtDistanceTypeSparsejaccard    AgentNgtDistanceType = "sparsejaccard"
	AgentNgtDistanceTypeSpjac            AgentNgtDistanceType = "spjac"
)

// Defines values for AgentNgtObjectType.
const (
	Float   AgentNgtObjectType = "float"
	Float16 AgentNgtObjectType = "float16"
	Uint8   AgentNgtObjectType = "uint8"
)

// Defines values for AgentSidecarConfigBlobStorageStorageType.
const (
	CloudStorage AgentSidecarConfigBlobStorageStorageType = "cloud_storage"
	S3           AgentSidecarConfigBlobStorageStorageType = "s3"
)

// Defines values for AgentSidecarConfigCompressCompressAlgorithm.
const (
	Gob  AgentSidecarConfigCompressCompressAlgorithm = "gob"
	Gzip AgentSidecarConfigCompressCompressAlgorithm = "gzip"
	Lz4  AgentSidecarConfigCompressCompressAlgorithm = "lz4"
	Zstd AgentSidecarConfigCompressCompressAlgorithm = "zstd"
)

// Defines values for AgentSidecarServiceType.
const (
	AgentSidecarServiceTypeClusterIP    AgentSidecarServiceType = "ClusterIP"
	AgentSidecarServiceTypeLoadBalancer AgentSidecarServiceType = "LoadBalancer"
	AgentSidecarServiceTypeNodePort     AgentSidecarServiceType = "NodePort"
)

// Defines values for DefaultsGrpcClientDialOptionInterceptors.
const (
	MetricInterceptor DefaultsGrpcClientDialOptionInterceptors = "MetricInterceptor"
	TraceInterceptor  DefaultsGrpcClientDialOptionInterceptors = "TraceInterceptor"
)

// Defines values for DefaultsObservabilityMetricsVersionInfoLabels.
const (
	AlgorithmInfo     DefaultsObservabilityMetricsVersionInfoLabels = "algorithm_info"
	BuildCpuInfoFlags DefaultsObservabilityMetricsVersionInfoLabels = "build_cpu_info_flags"
	BuildTime         DefaultsObservabilityMetricsVersionInfoLabels = "build_time"
	CgoEnabled        DefaultsObservabilityMetricsVersionInfoLabels = "cgo_enabled"
	GitCommit         DefaultsObservabilityMetricsVersionInfoLabels = "git_commit"
	GoArch            DefaultsObservabilityMetricsVersionInfoLabels = "go_arch"
	GoOs              DefaultsObservabilityMetricsVersionInfoLabels = "go_os"
	GoVersion         DefaultsObservabilityMetricsVersionInfoLabels = "go_version"
	ServerName        DefaultsObservabilityMetricsVersionInfoLabels = "server_name"
	ValdVersion       DefaultsObservabilityMetricsVersionInfoLabels = "vald_version"
)

// Defines values for DiscovererKind.
const (
	DiscovererKindDaemonSet  DiscovererKind = "DaemonSet"
	DiscovererKindDeployment DiscovererKind = "Deployment"
)

// Defines values for DiscovererServiceType.
const (
	DiscovererServiceTypeClusterIP    DiscovererServiceType = "ClusterIP"
	DiscovererServiceTypeLoadBalancer DiscovererServiceType = "LoadBalancer"
	DiscovererServiceTypeNodePort     DiscovererServiceType = "NodePort"
)

// Defines values for DiscovererUnhealthyPodEvictionPolicy.
const (
	DiscovererUnhealthyPodEvictionPolicyAlwaysAllow     DiscovererUnhealthyPodEvictionPolicy = "AlwaysAllow"
	DiscovererUnhealthyPodEvictionPolicyIfHealthyBudget DiscovererUnhealthyPodEvictionPolicy = "IfHealthyBudget"
)

// Defines values for GatewayFilterKind.
const (
	GatewayFilterKindDaemonSet  GatewayFilterKind = "DaemonSet"
	GatewayFilterKindDeployment GatewayFilterKind = "Deployment"
)

// Defines values for GatewayFilterServiceType.
const (
	GatewayFilterServiceTypeClusterIP    GatewayFilterServiceType = "ClusterIP"
	GatewayFilterServiceTypeLoadBalancer GatewayFilterServiceType = "LoadBalancer"
	GatewayFilterServiceTypeNodePort     GatewayFilterServiceType = "NodePort"
)

// Defines values for GatewayFilterUnhealthyPodEvictionPolicy.
const (
	GatewayFilterUnhealthyPodEvictionPolicyAlwaysAllow     GatewayFilterUnhealthyPodEvictionPolicy = "AlwaysAllow"
	GatewayFilterUnhealthyPodEvictionPolicyIfHealthyBudget GatewayFilterUnhealthyPodEvictionPolicy = "IfHealthyBudget"
)

// Defines values for GatewayLbKind.
const (
	GatewayLbKindDaemonSet  GatewayLbKind = "DaemonSet"
	GatewayLbKindDeployment GatewayLbKind = "Deployment"
)

// Defines values for GatewayLbServiceType.
const (
	GatewayLbServiceTypeClusterIP    GatewayLbServiceType = "ClusterIP"
	GatewayLbServiceTypeLoadBalancer GatewayLbServiceType = "LoadBalancer"
	GatewayLbServiceTypeNodePort     GatewayLbServiceType = "NodePort"
)

// Defines values for GatewayLbUnhealthyPodEvictionPolicy.
const (
	GatewayLbUnhealthyPodEvictionPolicyAlwaysAllow     GatewayLbUnhealthyPodEvictionPolicy = "AlwaysAllow"
	GatewayLbUnhealthyPodEvictionPolicyIfHealthyBudget GatewayLbUnhealthyPodEvictionPolicy = "IfHealthyBudget"
)

// Defines values for GatewayMirrorKind.
const (
	GatewayMirrorKindDaemonSet  GatewayMirrorKind = "DaemonSet"
	GatewayMirrorKindDeployment GatewayMirrorKind = "Deployment"
)

// Defines values for GatewayMirrorServiceType.
const (
	GatewayMirrorServiceTypeClusterIP    GatewayMirrorServiceType = "ClusterIP"
	GatewayMirrorServiceTypeLoadBalancer GatewayMirrorServiceType = "LoadBalancer"
	GatewayMirrorServiceTypeNodePort     GatewayMirrorServiceType = "NodePort"
)

// Defines values for GatewayMirrorUnhealthyPodEvictionPolicy.
const (
	GatewayMirrorUnhealthyPodEvictionPolicyAlwaysAllow     GatewayMirrorUnhealthyPodEvictionPolicy = "AlwaysAllow"
	GatewayMirrorUnhealthyPodEvictionPolicyIfHealthyBudget GatewayMirrorUnhealthyPodEvictionPolicy = "IfHealthyBudget"
)

// Defines values for ImagePullPolicy.
const (
	Always       ImagePullPolicy = "Always"
	IfNotPresent ImagePullPolicy = "IfNotPresent"
	Never        ImagePullPolicy = "Never"
)

// Defines values for ManagerIndexKind.
const (
	ManagerIndexKindDaemonSet  ManagerIndexKind = "DaemonSet"
	ManagerIndexKindDeployment ManagerIndexKind = "Deployment"
)

// Defines values for ManagerIndexServiceType.
const (
	ClusterIP    ManagerIndexServiceType = "ClusterIP"
	LoadBalancer ManagerIndexServiceType = "LoadBalancer"
	NodePort     ManagerIndexServiceType = "NodePort"
)

// Defines values for ManagerIndexUnhealthyPodEvictionPolicy.
const (
	ManagerIndexUnhealthyPodEvictionPolicyAlwaysAllow     ManagerIndexUnhealthyPodEvictionPolicy = "AlwaysAllow"
	ManagerIndexUnhealthyPodEvictionPolicyIfHealthyBudget ManagerIndexUnhealthyPodEvictionPolicy = "IfHealthyBudget"
)

// Defines values for ManagerIndexOperatorKind.
const (
	DaemonSet  ManagerIndexOperatorKind = "DaemonSet"
	Deployment ManagerIndexOperatorKind = "Deployment"
)

// Values defines model for Values.
type Values struct {
	Agent      *Agent      `json:"agent,omitempty"`
	Defaults   *Defaults   `json:"defaults,omitempty"`
	Discoverer *Discoverer `json:"discoverer,omitempty"`
	Gateway    *Gateway    `json:"gateway,omitempty"`
	Manager    *Manager    `json:"manager,omitempty"`
}

// Affinity defines model for affinity.
type Affinity = corev1.Affinity

// Agent defines model for agent.
type Agent struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// Algorithm agent algorithm type. it should be `ngt` or `faiss`.
	Algorithm *AgentAlgorithm `json:"algorithm,omitempty"`

	// Annotations deployment annotations
	Annotations        *map[string]any          `json:"annotations,omitempty"`
	ClusterRole        *AgentClusterRole        `json:"clusterRole,omitempty"`
	ClusterRoleBinding *AgentClusterRoleBinding `json:"clusterRoleBinding,omitempty"`

	// Enabled agent enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env *Env `json:"env,omitempty"`

	// ExternalTrafficPolicy external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy *string     `json:"externalTrafficPolicy,omitempty"`
	Faiss                 *AgentFaiss `json:"faiss,omitempty"`
	Hpa                   *Hpa        `json:"hpa,omitempty"`
	Image                 *Image      `json:"image,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// Kind deployment kind: Deployment, DaemonSet or StatefulSet
	Kind    *AgentKind `json:"kind,omitempty"`
	Logging *Logging   `json:"logging,omitempty"`

	// MaxReplicas maximum number of replicas. if HPA is disabled, this value will be ignored.
	MaxReplicas *int `json:"maxReplicas,omitempty"`

	// MaxUnavailable maximum number of unavailable replicas
	MaxUnavailable *string `json:"maxUnavailable,omitempty"`

	// MinReplicas minimum number of replicas. if HPA is disabled, the replicas will be set to this value
	MinReplicas *int `json:"minReplicas,omitempty"`

	// Name name of agent deployment
	Name *string   `json:"name,omitempty"`
	Ngt  *AgentNgt `json:"ngt,omitempty"`

	// NodeName node name
	NodeName *string `json:"nodeName,omitempty"`

	// NodeSelector node selector
	NodeSelector     *NodeSelector          `json:"nodeSelector,omitempty"`
	Observability    *Observability         `json:"observability,omitempty"`
	PersistentVolume *AgentPersistentVolume `json:"persistentVolume,omitempty"`

	// PodAnnotations pod annotations
	PodAnnotations *map[string]any `json:"podAnnotations,omitempty"`

	// PodManagementPolicy pod management policy: OrderedReady or Parallel
	PodManagementPolicy *AgentPodManagementPolicy `json:"podManagementPolicy,omitempty"`
	PodPriority         *PodPriority              `json:"podPriority,omitempty"`

	// PodSecurityContext security context for pod
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ProgressDeadlineSeconds progress deadline seconds
	ProgressDeadlineSeconds *int `json:"progressDeadlineSeconds,omitempty"`

	// Readreplica readreplica deployment annotations
	Readreplica *AgentReadreplica `json:"readreplica,omitempty"`

	// Resources compute resources
	Resources *Resources `json:"resources,omitempty"`

	// RevisionHistoryLimit number of old history to retain to allow rollback
	RevisionHistoryLimit *int                `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *AgentRollingUpdate `json:"rollingUpdate,omitempty"`

	// SecurityContext security context for container
	SecurityContext    *corev1.SecurityContext `json:"securityContext,omitempty"`
	ServerConfig       *ServerConfig           `json:"server_config,omitempty"`
	Service            *Service                `json:"service,omitempty"`
	ServiceAccountName *string                 `json:"serviceAccountName,omitempty"`

	// ServiceType service type: ClusterIP, LoadBalancer or NodePort
	ServiceType *AgentServiceType `json:"serviceType,omitempty"`
	Sidecar     *AgentSidecar     `json:"sidecar,omitempty"`

	// TerminationGracePeriodSeconds duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds *int `json:"terminationGracePeriodSeconds,omitempty"`

	// TimeZone Time zone
	TimeZone *string `json:"time_zone,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TopologySpreadConstraints topology spread constraints of gateway pods
	TopologySpreadConstraints *TopologySpreadConstraints `json:"topologySpreadConstraints,omitempty"`

	// UnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either cautious or permissive eviction.
	UnhealthyPodEvictionPolicy *AgentUnhealthyPodEvictionPolicy `json:"unhealthyPodEvictionPolicy,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`

	// VolumeMounts volume mounts
	VolumeMounts *VolumeMounts `json:"volumeMounts,omitempty"`

	// Volumes volumes
	Volumes *Volumes `json:"volumes,omitempty"`
}

// AgentAlgorithm agent algorithm type. it should be `ngt` or `faiss`.
type AgentAlgorithm string

// AgentKind deployment kind: Deployment, DaemonSet or StatefulSet
type AgentKind string

// AgentPodManagementPolicy pod management policy: OrderedReady or Parallel
type AgentPodManagementPolicy string

// AgentServiceType service type: ClusterIP, LoadBalancer or NodePort
type AgentServiceType string

// AgentUnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either cautious or
// permissive eviction.
type AgentUnhealthyPodEvictionPolicy string

// AgentClusterRole defines model for agent_clusterRole.
type AgentClusterRole struct {
	// Enabled creates clusterRole resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRole
	Name *string `json:"name,omitempty"`
}

// AgentClusterRoleBinding defines model for agent_clusterRoleBinding.
type AgentClusterRoleBinding struct {
	// Enabled creates clusterRoleBinding resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRoleBinding
	Name *string `json:"name,omitempty"`
}

// AgentFaiss defines model for agent_faiss.
type AgentFaiss struct {
	// AutoIndexCheckDuration check duration of automatic indexing
	AutoIndexCheckDuration *string `json:"auto_index_check_duration,omitempty"`

	// AutoIndexDurationLimit limit duration of automatic indexing
	AutoIndexDurationLimit *string `json:"auto_index_duration_limit,omitempty"`

	// AutoIndexLength number of cache to trigger automatic indexing
	AutoIndexLength *int `json:"auto_index_length,omitempty"`

	// AutoSaveIndexDuration duration of automatic save index
	AutoSaveIndexDuration *string `json:"auto_save_index_duration,omitempty"`

	// Dimension vector dimension
	Dimension *int `json:"dimension,omitempty"`

	// EnableCopyOnWrite enable copy on write saving for more stable backup
	EnableCopyOnWrite *bool `json:"enable_copy_on_write,omitempty"`

	// EnableInMemoryMode in-memory mode enabled
	EnableInMemoryMode *bool `json:"enable_in_memory_mode,omitempty"`

	// EnableProactiveGc enable proactive GC call for reducing heap memory allocation
	EnableProactiveGc *bool `json:"enable_proactive_gc,omitempty"`

	// IndexPath path to index data
	IndexPath *string `json:"index_path,omitempty"`

	// InitialDelayMaxDuration maximum duration for initial delay
	InitialDelayMaxDuration *string          `json:"initial_delay_max_duration,omitempty"`
	Kvsdb                   *AgentFaissKvsdb `json:"kvsdb,omitempty"`

	// LoadIndexTimeoutFactor a factor of load index timeout. timeout duration will be calculated by (index count to be loaded) * (factor).
	LoadIndexTimeoutFactor *string `json:"load_index_timeout_factor,omitempty"`

	// M m
	M *int `json:"m,omitempty"`

	// MaxLoadIndexTimeout maximum duration of load index timeout
	MaxLoadIndexTimeout *string `json:"max_load_index_timeout,omitempty"`

	// MethodType method type it should be `ivfpq` or `binaryindex`
	MethodType *AgentFaissMethodType `json:"method_type,omitempty"`

	// MetricType metric type it should be `innerproduct` or `l2`
	MetricType *AgentFaissMetricType `json:"metric_type,omitempty"`

	// MinLoadIndexTimeout minimum duration of load index timeout
	MinLoadIndexTimeout *string `json:"min_load_index_timeout,omitempty"`

	// Namespace namespace of myself
	Namespace *string `json:"namespace,omitempty"`

	// NbitsPerIdx nbits_per_idx
	NbitsPerIdx *int `json:"nbits_per_idx,omitempty"`

	// Nlist nlist
	Nlist *int `json:"nlist,omitempty"`

	// PodName pod name of myself
	PodName *string           `json:"pod_name,omitempty"`
	Vqueue  *AgentFaissVqueue `json:"vqueue,omitempty"`
}

// AgentFaissMethodType method type it should be `ivfpq` or `binaryindex`
type AgentFaissMethodType string

// AgentFaissMetricType metric type it should be `innerproduct` or `l2`
type AgentFaissMetricType string

// AgentFaissKvsdb defines model for agent_faiss_kvsdb.
type AgentFaissKvsdb struct {
	// Concurrency kvsdb processing concurrency
	Concurrency *int `json:"concurrency,omitempty"`
}

// AgentFaissVqueue defines model for agent_faiss_vqueue.
type AgentFaissVqueue struct {
	// DeleteBufferPoolSize delete slice pool buffer size
	DeleteBufferPoolSize *int `json:"delete_buffer_pool_size,omitempty"`

	// InsertBufferPoolSize insert slice pool buffer size
	InsertBufferPoolSize *int `json:"insert_buffer_pool_size,omitempty"`
}

// AgentNgt defines model for agent_ngt.
type AgentNgt struct {
	// AutoCreateIndexPoolSize batch process pool size of automatic create index operation
	AutoCreateIndexPoolSize *int `json:"auto_create_index_pool_size,omitempty"`

	// AutoIndexCheckDuration check duration of automatic indexing
	AutoIndexCheckDuration *string `json:"auto_index_check_duration,omitempty"`

	// AutoIndexDurationLimit limit duration of automatic indexing
	AutoIndexDurationLimit *string `json:"auto_index_duration_limit,omitempty"`

	// AutoIndexLength number of cache to trigger automatic indexing
	AutoIndexLength *int `json:"auto_index_length,omitempty"`

	// AutoSaveIndexDuration duration of automatic save index
	AutoSaveIndexDuration *string `json:"auto_save_index_duration,omitempty"`

	// BrokenIndexHistoryLimit maximum number of broken index generations to backup
	BrokenIndexHistoryLimit *int `json:"broken_index_history_limit,omitempty"`

	// BulkInsertChunkSize bulk insert chunk size
	BulkInsertChunkSize *int `json:"bulk_insert_chunk_size,omitempty"`

	// CreationEdgeSize creation edge size
	CreationEdgeSize *int `json:"creation_edge_size,omitempty"`

	// DefaultEpsilon default epsilon used for search
	DefaultEpsilon *float32 `json:"default_epsilon,omitempty"`

	// DefaultPoolSize default create index batch pool size
	DefaultPoolSize *int `json:"default_pool_size,omitempty"`

	// DefaultRadius default radius used for search
	DefaultRadius *float32 `json:"default_radius,omitempty"`

	// Dimension vector dimension
	Dimension *int `json:"dimension,omitempty"`

	// DistanceType distance type. it should be `l1`, `l2`, `angle`, `hamming`, `cosine`,`poincare`, `lorentz`, `jaccard`, `sparsejaccard`, `normalizedangle` or `normalizedcosine` or `innerproduct`. for further details about NGT libraries supported distance is https://github.com/NGT-labs/NGT/wiki/Command-Quick-Reference and vald agent's supported NGT distance type is https://pkg.go.dev/github.com/vdaas/vald/internal/core/algorithm/ngt#pkg-constants
	DistanceType *AgentNgtDistanceType `json:"distance_type,omitempty"`

	// EnableCopyOnWrite enable copy on write saving for more stable backup
	EnableCopyOnWrite *bool `json:"enable_copy_on_write,omitempty"`

	// EnableExportIndexInfoToK8s enable export index info to k8s
	EnableExportIndexInfoToK8s *bool `json:"enable_export_index_info_to_k8s,omitempty"`

	// EnableInMemoryMode in-memory mode enabled
	EnableInMemoryMode *bool `json:"enable_in_memory_mode,omitempty"`

	// EnableProactiveGc enable proactive GC call for reducing heap memory allocation
	EnableProactiveGc *bool `json:"enable_proactive_gc,omitempty"`

	// EnableStatistics enable index statistics loading
	EnableStatistics *bool `json:"enable_statistics,omitempty"`

	// EpsilonForCreation the epsilon used for creation
	EpsilonForCreation *float32 `json:"epsilon_for_creation,omitempty"`

	// ErrorBufferLimit maximum number of core ngt error buffer pool size limit
	ErrorBufferLimit *int `json:"error_buffer_limit,omitempty"`

	// ExportIndexInfoDuration duration of exporting index info
	ExportIndexInfoDuration *string `json:"export_index_info_duration,omitempty"`

	// IndexPath path to index data
	IndexPath *string `json:"index_path,omitempty"`

	// InitialDelayMaxDuration maximum duration for initial delay
	InitialDelayMaxDuration *string        `json:"initial_delay_max_duration,omitempty"`
	Kvsdb                   *AgentNgtKvsdb `json:"kvsdb,omitempty"`

	// LoadIndexTimeoutFactor a factor of load index timeout. timeout duration will be calculated by (index count to be loaded) * (factor).
	LoadIndexTimeoutFactor *string `json:"load_index_timeout_factor,omitempty"`

	// MaxLoadIndexTimeout maximum duration of load index timeout
	MaxLoadIndexTimeout *string `json:"max_load_index_timeout,omitempty"`

	// MinLoadIndexTimeout minimum duration of load index timeout
	MinLoadIndexTimeout *string `json:"min_load_index_timeout,omitempty"`

	// Namespace namespace of myself
	Namespace *string `json:"namespace,omitempty"`

	// ObjectType object type. it should be `float` or `uint8` or `float16`. for further details: https://github.com/NGT-labs/NGT/wiki/Command-Quick-Reference
	ObjectType *AgentNgtObjectType `json:"object_type,omitempty"`

	// PodName pod name of myself
	PodName *string `json:"pod_name,omitempty"`

	// SearchEdgeSize search edge size
	SearchEdgeSize *int            `json:"search_edge_size,omitempty"`
	Vqueue         *AgentNgtVqueue `json:"vqueue,omitempty"`
}

// AgentNgtDistanceType distance type. it should be `l1`, `l2`, `angle`, `hamming`, `cosine`,`poincare`, `lorentz`, `jaccard`, `sparsejaccard`,
// `normalizedangle` or `normalizedcosine` or `innerproduct`. for further details about NGT libraries supported distance is
// https://github.com/NGT-labs/NGT/wiki/Command-Quick-Reference and vald agent's supported NGT distance type is
// https://pkg.go.dev/github.com/vdaas/vald/internal/core/algorithm/ngt#pkg-constants
type AgentNgtDistanceType string

// AgentNgtObjectType object type. it should be `float` or `uint8` or `float16`. for further details:
// https://github.com/NGT-labs/NGT/wiki/Command-Quick-Reference
type AgentNgtObjectType string

// AgentNgtKvsdb defines model for agent_ngt_kvsdb.
type AgentNgtKvsdb struct {
	// Concurrency kvsdb processing concurrency
	Concurrency *int `json:"concurrency,omitempty"`
}

// AgentNgtVqueue defines model for agent_ngt_vqueue.
type AgentNgtVqueue struct {
	// DeleteBufferPoolSize delete slice pool buffer size
	DeleteBufferPoolSize *int `json:"delete_buffer_pool_size,omitempty"`

	// InsertBufferPoolSize insert slice pool buffer size
	InsertBufferPoolSize *int `json:"insert_buffer_pool_size,omitempty"`
}

// AgentPersistentVolume defines model for agent_persistentVolume.
type AgentPersistentVolume struct {
	// AccessMode agent pod storage accessMode
	AccessMode *string `json:"accessMode,omitempty"`

	// Enabled enables PVC. It is required to enable if agent pod's file store functionality is enabled with non in-memory mode
	Enabled *bool `json:"enabled,omitempty"`

	// MountPropagation agent pod storage mountPropagation
	MountPropagation *string `json:"mountPropagation,omitempty"`

	// Size size of agent pod volume
	Size *string `json:"size,omitempty"`

	// StorageClass storageClass name for agent pod volume
	StorageClass *string `json:"storageClass,omitempty"`
}

// AgentReadreplica readreplica deployment annotations
type AgentReadreplica struct {
	// ComponentName app.kubernetes.io/component name of agent readreplica
	ComponentName *string `json:"component_name,omitempty"`

	// Enabled [This feature is WORK IN PROGRESS]enable agent readreplica
	Enabled *bool `json:"enabled,omitempty"`
	Hpa     *Hpa  `json:"hpa,omitempty"`

	// LabelKey label key to identify read replica resources
	LabelKey *string `json:"label_key,omitempty"`

	// MaxReplicas maximum number of replicas. if HPA is disabled, this value will be ignored.
	MaxReplicas *int `json:"maxReplicas,omitempty"`

	// MinReplicas minimum number of replicas. if HPA is disabled, the replicas will be set to this value
	MinReplicas *int `json:"minReplicas,omitempty"`

	// Name name of agent readreplica
	Name *string `json:"name,omitempty"`

	// Service service settings for read replica service resources
	Service *AgentReadreplicaService `json:"service,omitempty"`

	// SnapshotClassname snapshot class name for snapshotter used for read replica
	SnapshotClassname *string `json:"snapshot_classname,omitempty"`

	// VolumeName name of clone volume of agent pvc for read replica
	VolumeName *string `json:"volume_name,omitempty"`
}

// AgentReadreplicaService service settings for read replica service resources
type AgentReadreplicaService struct {
	// Annotations readreplica deployment annotations
	Annotations *map[string]any `json:"annotations,omitempty"`
}

// AgentRollingUpdate defines model for agent_rollingUpdate.
type AgentRollingUpdate struct {
	// MaxSurge max surge of rolling update
	MaxSurge *string `json:"maxSurge,omitempty"`

	// MaxUnavailable max unavailable of rolling update
	MaxUnavailable *string `json:"maxUnavailable,omitempty"`

	// Partition StatefulSet partition
	Partition *int `json:"partition,omitempty"`
}

// AgentSidecar defines model for agent_sidecar.
type AgentSidecar struct {
	Config *AgentSidecarConfig `json:"config,omitempty"`

	// Enabled sidecar enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env   *Env   `json:"env,omitempty"`
	Image *Image `json:"image,omitempty"`

	// InitContainerEnabled sidecar on initContainer mode enabled.
	InitContainerEnabled *bool    `json:"initContainerEnabled,omitempty"`
	Logging              *Logging `json:"logging,omitempty"`

	// Name name of agent sidecar
	Name          *string        `json:"name,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// Resources compute resources
	Resources    *Resources           `json:"resources,omitempty"`
	ServerConfig *ServerConfig        `json:"server_config,omitempty"`
	Service      *AgentSidecarService `json:"service,omitempty"`

	// TimeZone Time zone
	TimeZone *string `json:"time_zone,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`
}

// AgentSidecarConfig defines model for agent_sidecar_config.
type AgentSidecarConfig struct {
	// AutoBackupDuration auto backup duration
	AutoBackupDuration *string `json:"auto_backup_duration,omitempty"`

	// AutoBackupEnabled auto backup triggered by timer is enabled
	AutoBackupEnabled *bool                          `json:"auto_backup_enabled,omitempty"`
	BlobStorage       *AgentSidecarConfigBlobStorage `json:"blob_storage,omitempty"`
	Client            *AgentSidecarConfigClient      `json:"client,omitempty"`
	Compress          *AgentSidecarConfigCompress    `json:"compress,omitempty"`

	// Filename backup filename
	Filename *string `json:"filename,omitempty"`

	// FilenameSuffix suffix for backup filename
	FilenameSuffix *string `json:"filename_suffix,omitempty"`

	// PostStopTimeout timeout for observing file changes during post stop
	PostStopTimeout *string  `json:"post_stop_timeout,omitempty"`
	RestoreBackoff  *Backoff `json:"restore_backoff,omitempty"`

	// RestoreBackoffEnabled restore backoff enabled
	RestoreBackoffEnabled *bool `json:"restore_backoff_enabled,omitempty"`

	// WatchEnabled auto backup triggered by file changes is enabled
	WatchEnabled *bool `json:"watch_enabled,omitempty"`
}

// AgentSidecarConfigBlobStorage defines model for agent_sidecar_config_blob_storage.
type AgentSidecarConfigBlobStorage struct {
	// Bucket bucket name
	Bucket       *string                                    `json:"bucket,omitempty"`
	CloudStorage *AgentSidecarConfigBlobStorageCloudStorage `json:"cloud_storage,omitempty"`
	S3           *AgentSidecarConfigBlobStorageS3           `json:"s3,omitempty"`

	// StorageType storage type
	StorageType *AgentSidecarConfigBlobStorageStorageType `json:"storage_type,omitempty"`
}

// AgentSidecarConfigBlobStorageStorageType storage type
type AgentSidecarConfigBlobStorageStorageType string

// AgentSidecarConfigBlobStorageCloudStorage defines model for agent_sidecar_config_blob_storage_cloud_storage.
type AgentSidecarConfigBlobStorageCloudStorage struct {
	Client *AgentSidecarConfigBlobStorageCloudStorageClient `json:"client,omitempty"`

	// Url cloud storage url
	Url *string `json:"url,omitempty"`

	// WriteBufferSize bytes of the chunks for upload
	WriteBufferSize *int `json:"write_buffer_size,omitempty"`

	// WriteCacheControl Cache-Control of HTTP Header
	WriteCacheControl *string `json:"write_cache_control,omitempty"`

	// WriteContentDisposition Content-Disposition of HTTP Header
	WriteContentDisposition *string `json:"write_content_disposition,omitempty"`

	// WriteContentEncoding the encoding of the blob's content
	WriteContentEncoding *string `json:"write_content_encoding,omitempty"`

	// WriteContentLanguage the language of blob's content
	WriteContentLanguage *string `json:"write_content_language,omitempty"`

	// WriteContentType MIME type of the blob
	WriteContentType *string `json:"write_content_type,omitempty"`
}

// AgentSidecarConfigBlobStorageCloudStorageClient defines model for agent_sidecar_config_blob_storage_cloud_storage_client.
type AgentSidecarConfigBlobStorageCloudStorageClient struct {
	// CredentialsFilePath credentials file path
	CredentialsFilePath *string `json:"credentials_file_path,omitempty"`

	// CredentialsJson credentials json
	CredentialsJson *string `json:"credentials_json,omitempty"`
}

// AgentSidecarConfigBlobStorageS3 defines model for agent_sidecar_config_blob_storage_s3.
type AgentSidecarConfigBlobStorageS3 struct {
	// AccessKey s3 access key
	AccessKey *string `json:"access_key,omitempty"`

	// Enable100Continue enable AWS SDK adding the 'Expect: 100-Continue' header to PUT requests over 2MB of content.
	Enable100Continue *bool `json:"enable_100_continue,omitempty"`

	// EnableContentMd5Validation enable the S3 client to add MD5 checksum to upload API calls.
	EnableContentMd5Validation *bool `json:"enable_content_md5_validation,omitempty"`

	// EnableEndpointDiscovery enable endpoint discovery
	EnableEndpointDiscovery *bool `json:"enable_endpoint_discovery,omitempty"`

	// EnableEndpointHostPrefix enable prefixing request endpoint hosts with modeled information
	EnableEndpointHostPrefix *bool `json:"enable_endpoint_host_prefix,omitempty"`

	// EnableParamValidation enables semantic parameter validation
	EnableParamValidation *bool `json:"enable_param_validation,omitempty"`

	// EnableSsl enable ssl for s3 session
	EnableSsl *bool `json:"enable_ssl,omitempty"`

	// Endpoint s3 endpoint
	Endpoint *string `json:"endpoint,omitempty"`

	// ForcePathStyle use path-style addressing
	ForcePathStyle *bool `json:"force_path_style,omitempty"`

	// MaxChunkSize s3 download max chunk size
	MaxChunkSize *string `json:"max_chunk_size,omitempty"`

	// MaxPartSize s3 multipart upload max part size
	MaxPartSize *string `json:"max_part_size,omitempty"`

	// MaxRetries maximum number of retries of s3 client
	MaxRetries *int `json:"max_retries,omitempty"`

	// Region s3 region
	Region *string `json:"region,omitempty"`

	// SecretAccessKey s3 secret access key
	SecretAccessKey *string `json:"secret_access_key,omitempty"`

	// Token s3 token
	Token *string `json:"token,omitempty"`

	// UseAccelerate enable s3 accelerate feature
	UseAccelerate *bool `json:"use_accelerate,omitempty"`

	// UseArnRegion s3 service client to use the region specified in the ARN
	UseArnRegion *bool `json:"use_arn_region,omitempty"`

	// UseDualStack use dual stack
	UseDualStack *bool `json:"use_dual_stack,omitempty"`
}

// AgentSidecarConfigClient defines model for agent_sidecar_config_client.
type AgentSidecarConfigClient struct {
	Net       *Net                               `json:"net,omitempty"`
	Transport *AgentSidecarConfigClientTransport `json:"transport,omitempty"`
}

// AgentSidecarConfigClientTransport defines model for agent_sidecar_config_client_transport.
type AgentSidecarConfigClientTransport struct {
	Backoff      *Backoff                                       `json:"backoff,omitempty"`
	RoundTripper *AgentSidecarConfigClientTransportRoundTripper `json:"round_tripper,omitempty"`
}

// AgentSidecarConfigClientTransportRoundTripper defines model for agent_sidecar_config_client_transport_round_tripper.
type AgentSidecarConfigClientTransportRoundTripper struct {
	// ExpectContinueTimeout expect continue timeout
	ExpectContinueTimeout *string `json:"expect_continue_timeout,omitempty"`

	// ForceAttemptHttp2 force attempt HTTP2
	ForceAttemptHttp2 *bool `json:"force_attempt_http_2,omitempty"`

	// IdleConnTimeout timeout for idle connections
	IdleConnTimeout *string `json:"idle_conn_timeout,omitempty"`

	// MaxConnsPerHost maximum count of connections per host
	MaxConnsPerHost *int `json:"max_conns_per_host,omitempty"`

	// MaxIdleConns maximum count of idle connections
	MaxIdleConns *int `json:"max_idle_conns,omitempty"`

	// MaxIdleConnsPerHost maximum count of idle connections per host
	MaxIdleConnsPerHost *int `json:"max_idle_conns_per_host,omitempty"`

	// MaxResponseHeaderSize maximum response header size
	MaxResponseHeaderSize *int `json:"max_response_header_size,omitempty"`

	// ReadBufferSize read buffer size
	ReadBufferSize *int `json:"read_buffer_size,omitempty"`

	// ResponseHeaderTimeout timeout for response header
	ResponseHeaderTimeout *string `json:"response_header_timeout,omitempty"`

	// TlsHandshakeTimeout TLS handshake timeout
	TlsHandshakeTimeout *string `json:"tls_handshake_timeout,omitempty"`

	// WriteBufferSize write buffer size
	WriteBufferSize *int `json:"write_buffer_size,omitempty"`
}

// AgentSidecarConfigCompress defines model for agent_sidecar_config_compress.
type AgentSidecarConfigCompress struct {
	// CompressAlgorithm compression algorithm. must be `gob`, `gzip`, `lz4` or `zstd`
	CompressAlgorithm *AgentSidecarConfigCompressCompressAlgorithm `json:"compress_algorithm,omitempty"`

	// CompressionLevel compression level. value range relies on which algorithm is used. `gob`: level will be ignored. `gzip`: -1 (default compression), 0 (no compression), or 1 (best speed) to 9 (best compression). `lz4`: >= 0, higher is better compression. `zstd`: 1 (fastest) to 22 (best), however implementation relies on klauspost/compress.
	CompressionLevel *int `json:"compression_level,omitempty"`
}

// AgentSidecarConfigCompressCompressAlgorithm compression algorithm. must be `gob`, `gzip`, `lz4` or `zstd`
type AgentSidecarConfigCompressCompressAlgorithm string

// AgentSidecarService defines model for agent_sidecar_service.
type AgentSidecarService struct {
	// Annotations agent sidecar service annotations
	Annotations *map[string]any `json:"annotations,omitempty"`

	// Enabled agent sidecar service enabled
	Enabled *bool `json:"enabled,omitempty"`

	// ExternalTrafficPolicy external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy *string `json:"externalTrafficPolicy,omitempty"`

	// Labels agent sidecar service labels
	Labels *map[string]any `json:"labels,omitempty"`

	// Type service type: ClusterIP, LoadBalancer or NodePort
	Type *AgentSidecarServiceType `json:"type,omitempty"`
}

// AgentSidecarServiceType service type: ClusterIP, LoadBalancer or NodePort
type AgentSidecarServiceType string

// Backoff defines model for backoff.
type Backoff = config.Backoff

// Defaults defines model for defaults.
type Defaults struct {
	Grpc          *DefaultsGrpc  `json:"grpc,omitempty"`
	Image         *DefaultsImage `json:"image,omitempty"`
	Logging       *Logging       `json:"logging,omitempty"`
	NetworkPolicy *NetworkPolicy `json:"networkPolicy,omitempty"`
	Observability *Observability `json:"observability,omitempty"`
	ServerConfig  *ServerConfig  `json:"server_config,omitempty"`

	// TimeZone Time zone
	TimeZone *string `json:"time_zone,omitempty"`
}

// DefaultsGrpc defines model for defaults_grpc.
type DefaultsGrpc struct {
	Client *GrpcClient `json:"client,omitempty"`
}

// DefaultsGrpcClientCircuitBreaker defines model for defaults_grpc_client_circuit_breaker.
type DefaultsGrpcClientCircuitBreaker struct {
	// ClosedErrorRate gRPC client circuitbreaker closed error rate
	ClosedErrorRate *float32 `json:"closed_error_rate,omitempty"`

	// ClosedRefreshTimeout gRPC client circuitbreaker closed refresh timeout
	ClosedRefreshTimeout *string `json:"closed_refresh_timeout,omitempty"`

	// HalfOpenErrorRate gRPC client circuitbreaker half-open error rate
	HalfOpenErrorRate *float32 `json:"half_open_error_rate,omitempty"`

	// MinSamples gRPC client circuitbreaker minimum sampling count
	MinSamples *int `json:"min_samples,omitempty"`

	// OpenTimeout gRPC client circuitbreaker open timeout
	OpenTimeout *string `json:"open_timeout,omitempty"`
}

// DefaultsGrpcClientConnectionPool defines model for defaults_grpc_client_connection_pool.
type DefaultsGrpcClientConnectionPool struct {
	// EnableDnsResolver enables gRPC client connection pool dns resolver, when enabled vald uses ip handshake exclude dns discovery which improves network performance
	EnableDnsResolver *bool `json:"enable_dns_resolver,omitempty"`

	// EnableMetrics enables gRPC client connection pool metrics
	EnableMetrics *bool `json:"enable_metrics,omitempty"`

	// EnableRebalance enables gRPC client connection pool rebalance
	EnableRebalance *bool `json:"enable_rebalance,omitempty"`

	// OldConnCloseDuration makes delay before gRPC client connection closing during connection pool rebalance
	OldConnCloseDuration *string `json:"old_conn_close_duration,omitempty"`

	// RebalanceDuration gRPC client connection pool rebalance duration
	RebalanceDuration *string `json:"rebalance_duration,omitempty"`

	// Size gRPC client connection pool size
	Size *int `json:"size,omitempty"`
}

// DefaultsGrpcClientDialOption defines model for defaults_grpc_client_dial_option.
type DefaultsGrpcClientDialOption struct {
	// Authority gRPC client dial option authority
	Authority *string `json:"authority,omitempty"`

	// BackoffBaseDelay gRPC client dial option base backoff delay
	BackoffBaseDelay *string `json:"backoff_base_delay,omitempty"`

	// BackoffJitter gRPC client dial option base backoff delay
	BackoffJitter *float32 `json:"backoff_jitter,omitempty"`

	// BackoffMaxDelay gRPC client dial option max backoff delay
	BackoffMaxDelay *string `json:"backoff_max_delay,omitempty"`

	// BackoffMultiplier gRPC client dial option base backoff delay
	BackoffMultiplier *float32 `json:"backoff_multiplier,omitempty"`

	// DisableRetry gRPC client dial option disables retry
	DisableRetry *bool `json:"disable_retry,omitempty"`

	// EnableBackoff gRPC client dial option backoff enabled
	EnableBackoff *bool `json:"enable_backoff,omitempty"`

	// IdleTimeout gRPC client dial option idle_timeout
	IdleTimeout *string `json:"idle_timeout,omitempty"`

	// InitialConnectionWindowSize gRPC client dial option initial connection window size
	InitialConnectionWindowSize *int `json:"initial_connection_window_size,omitempty"`

	// InitialWindowSize gRPC client dial option initial window size
	InitialWindowSize *int `json:"initial_window_size,omitempty"`

	// Insecure gRPC client dial option insecure enabled
	Insecure *bool `json:"insecure,omitempty"`

	// Interceptors gRPC client interceptors
	Interceptors *[]DefaultsGrpcClientDialOptionInterceptors `json:"interceptors,omitempty"`
	Keepalive    *DefaultsGrpcClientDialOptionKeepalive      `json:"keepalive,omitempty"`

	// MaxCallAttempts gRPC client dial option number of max call attempts
	MaxCallAttempts *int `json:"max_call_attempts,omitempty"`

	// MaxHeaderListSize gRPC client dial option max header list size
	MaxHeaderListSize *int `json:"max_header_list_size,omitempty"`

	// MaxMsgSize gRPC client dial option max message size
	MaxMsgSize *int `json:"max_msg_size,omitempty"`

	// MinConnectionTimeout gRPC client dial option minimum connection timeout
	MinConnectionTimeout *string `json:"min_connection_timeout,omitempty"`
	Net                  *Net    `json:"net,omitempty"`

	// ReadBufferSize gRPC client dial option read buffer size
	ReadBufferSize *int `json:"read_buffer_size,omitempty"`

	// SharedWriteBuffer gRPC client dial option sharing write buffer
	SharedWriteBuffer *bool `json:"shared_write_buffer,omitempty"`

	// Timeout gRPC client dial option timeout
	Timeout *string `json:"timeout,omitempty"`

	// UserAgent gRPC client dial option user_agent
	UserAgent *string `json:"user_agent,omitempty"`

	// WriteBufferSize gRPC client dial option write buffer size
	WriteBufferSize *int `json:"write_buffer_size,omitempty"`
}

// DefaultsGrpcClientDialOptionInterceptors defines model for DefaultsGrpcClientDialOption.Interceptors.
type DefaultsGrpcClientDialOptionInterceptors string

// DefaultsGrpcClientDialOptionKeepalive defines model for defaults_grpc_client_dial_option_keepalive.
type DefaultsGrpcClientDialOptionKeepalive struct {
	// PermitWithoutStream gRPC client keep alive permit without stream
	PermitWithoutStream *bool `json:"permit_without_stream,omitempty"`

	// Time gRPC client keep alive time
	Time *string `json:"time,omitempty"`

	// Timeout gRPC client keep alive timeout
	Timeout *string `json:"timeout,omitempty"`
}

// DefaultsGrpcClientDialOptionNetDialer defines model for defaults_grpc_client_dial_option_net_dialer.
type DefaultsGrpcClientDialOptionNetDialer struct {
	// DualStackEnabled gRPC client TCP dialer dual stack enabled
	DualStackEnabled *bool `json:"dual_stack_enabled,omitempty"`

	// Keepalive gRPC client TCP dialer keep alive
	Keepalive *string `json:"keepalive,omitempty"`

	// Timeout gRPC client TCP dialer timeout
	Timeout *string `json:"timeout,omitempty"`
}

// DefaultsGrpcClientDialOptionNetDns defines model for defaults_grpc_client_dial_option_net_dns.
type DefaultsGrpcClientDialOptionNetDns struct {
	// CacheEnabled gRPC client DNS cache enabled
	CacheEnabled *bool `json:"cache_enabled,omitempty"`

	// CacheExpiration gRPC client DNS cache expiration
	CacheExpiration *string `json:"cache_expiration,omitempty"`

	// RefreshDuration gRPC client DNS cache refresh duration
	RefreshDuration *string `json:"refresh_duration,omitempty"`
}

// DefaultsImage defines model for defaults_image.
type DefaultsImage struct {
	// Registry default docker image registry (applied to all images; override per-component via image.registry)
	Registry *string `json:"registry,omitempty"`

	// Tag docker image tag
	Tag *string `json:"tag,omitempty"`
}

// DefaultsNetworkPolicyCustom custom network policies that a user can add
type DefaultsNetworkPolicyCustom struct {
	// Egress custom egress network policies that a user can add
	Egress *[]map[string]any `json:"egress,omitempty"`

	// Ingress custom ingress network policies that a user can add
	Ingress *[]map[string]any `json:"ingress,omitempty"`
}

// DefaultsObservabilityMetrics defines model for defaults_observability_metrics.
type DefaultsObservabilityMetrics struct {
	// EnableCgo CGO metrics enabled
	EnableCgo *bool `json:"enable_cgo,omitempty"`

	// EnableGoroutine goroutine metrics enabled
	EnableGoroutine *bool `json:"enable_goroutine,omitempty"`

	// EnableMemory memory metrics enabled
	EnableMemory *bool `json:"enable_memory,omitempty"`

	// EnableVersionInfo version info metrics enabled
	EnableVersionInfo *bool `json:"enable_version_info,omitempty"`

	// VersionInfoLabels enabled label names of version info
	VersionInfoLabels *[]DefaultsObservabilityMetricsVersionInfoLabels `json:"version_info_labels,omitempty"`
}

// DefaultsObservabilityMetricsVersionInfoLabels defines model for DefaultsObservabilityMetrics.VersionInfoLabels.
type DefaultsObservabilityMetricsVersionInfoLabels string

// DefaultsObservabilityOtlp defines model for defaults_observability_otlp.
type DefaultsObservabilityOtlp struct {
	// Attribute default resource attribute
	Attribute *DefaultsObservabilityOtlpAttribute `json:"attribute,omitempty"`

	// CollectorEndpoint OpenTelemetry Collector endpoint
	CollectorEndpoint *string `json:"collector_endpoint,omitempty"`

	// MetricsExportInterval metrics export interval
	MetricsExportInterval *string `json:"metrics_export_interval,omitempty"`

	// MetricsExportTimeout metrics export timeout
	MetricsExportTimeout *string `json:"metrics_export_timeout,omitempty"`

	// TraceBatchTimeout trace batch timeout
	TraceBatchTimeout *string `json:"trace_batch_timeout,omitempty"`

	// TraceExportTimeout trace export timeout
	TraceExportTimeout *string `json:"trace_export_timeout,omitempty"`

	// TraceMaxExportBatchSize trace maximum export batch size
	TraceMaxExportBatchSize *int `json:"trace_max_export_batch_size,omitempty"`

	// TraceMaxQueueSize trace maximum queue size
	TraceMaxQueueSize *int `json:"trace_max_queue_size,omitempty"`
}

// DefaultsObservabilityOtlpAttribute default resource attribute
type DefaultsObservabilityOtlpAttribute struct {
	// Namespace namespace
	Namespace *string `json:"namespace,omitempty"`

	// NodeName node name
	NodeName *string `json:"node_name,omitempty"`

	// PodName pod name
	PodName *string `json:"pod_name,omitempty"`

	// ServiceName service name
	ServiceName *string `json:"service_name,omitempty"`
}

// DefaultsObservabilityTrace defines model for defaults_observability_trace.
type DefaultsObservabilityTrace struct {
	// Enabled trace enabled
	Enabled *bool `json:"enabled,omitempty"`
}

// DefaultsServerConfigHealths defines model for defaults_server_config_healths.
type DefaultsServerConfigHealths struct {
	Liveness  *DefaultsServerConfigHealthsLiveness  `json:"liveness,omitempty"`
	Readiness *DefaultsServerConfigHealthsReadiness `json:"readiness,omitempty"`
	Startup   *DefaultsServerConfigHealthsStartup   `json:"startup,omitempty"`
}

// DefaultsServerConfigHealthsLiveness defines model for defaults_server_config_healths_liveness.
type DefaultsServerConfigHealthsLiveness struct {
	// Enabled liveness server enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Host liveness server host
	Host          *string                                           `json:"host,omitempty"`
	LivenessProbe *DefaultsServerConfigHealthsLivenessLivenessProbe `json:"livenessProbe,omitempty"`

	// Port liveness server port
	Port   *int        `json:"port,omitempty"`
	Server *RestServer `json:"server,omitempty"`

	// ServicePort liveness server service port
	ServicePort *int `json:"servicePort,omitempty"`
}

// DefaultsServerConfigHealthsLivenessLivenessProbe defines model for defaults_server_config_healths_liveness_livenessProbe.
type DefaultsServerConfigHealthsLivenessLivenessProbe struct {
	// FailureThreshold liveness probe failure threshold
	FailureThreshold *int                                                     `json:"failureThreshold,omitempty"`
	HttpGet          *DefaultsServerConfigHealthsLivenessLivenessProbeHttpGet `json:"httpGet,omitempty"`

	// InitialDelaySeconds liveness probe initial delay seconds
	InitialDelaySeconds *int `json:"initialDelaySeconds,omitempty"`

	// PeriodSeconds liveness probe period seconds
	PeriodSeconds *int `json:"periodSeconds,omitempty"`

	// SuccessThreshold liveness probe success threshold
	SuccessThreshold *int `json:"successThreshold,omitempty"`

	// TimeoutSeconds liveness probe timeout seconds
	TimeoutSeconds *int `json:"timeoutSeconds,omitempty"`
}

// DefaultsServerConfigHealthsLivenessLivenessProbeHttpGet defines model for defaults_server_config_healths_liveness_livenessProbe_httpGet.
type DefaultsServerConfigHealthsLivenessLivenessProbeHttpGet struct {
	// Path liveness probe path
	Path *string `json:"path,omitempty"`

	// Port liveness probe port
	Port *string `json:"port,omitempty"`

	// Scheme liveness probe scheme
	Scheme *string `json:"scheme,omitempty"`
}

// DefaultsServerConfigHealthsReadiness defines model for defaults_server_config_healths_readiness.
type DefaultsServerConfigHealthsReadiness struct {
	// Enabled readiness server enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Host readiness server host
	Host *string `json:"host,omitempty"`

	// Port readiness server port
	Port           *int                                                `json:"port,omitempty"`
	ReadinessProbe *DefaultsServerConfigHealthsReadinessReadinessProbe `json:"readinessProbe,omitempty"`
	Server         *RestServer                                         `json:"server,omitempty"`

	// ServicePort readiness server service port
	ServicePort *int `json:"servicePort,omitempty"`
}

// DefaultsServerConfigHealthsReadinessReadinessProbe defines model for defaults_server_config_healths_readiness_readinessProbe.
type DefaultsServerConfigHealthsReadinessReadinessProbe struct {
	// FailureThreshold readiness probe failure threshold
	FailureThreshold *int                                                       `json:"failureThreshold,omitempty"`
	HttpGet          *DefaultsServerConfigHealthsReadinessReadinessProbeHttpGet `json:"httpGet,omitempty"`

	// InitialDelaySeconds readiness probe initial delay seconds
	InitialDelaySeconds *int `json:"initialDelaySeconds,omitempty"`

	// PeriodSeconds readiness probe period seconds
	PeriodSeconds *int `json:"periodSeconds,omitempty"`

	// SuccessThreshold readiness probe success threshold
	SuccessThreshold *int `json:"successThreshold,omitempty"`

	// TimeoutSeconds readiness probe timeout seconds
	TimeoutSeconds *int `json:"timeoutSeconds,omitempty"`
}

// DefaultsServerConfigHealthsReadinessReadinessProbeHttpGet defines model for defaults_server_config_healths_readiness_readinessProbe_httpGet.
type DefaultsServerConfigHealthsReadinessReadinessProbeHttpGet struct {
	// Path readiness probe path
	Path *string `json:"path,omitempty"`

	// Port readiness probe port
	Port *string `json:"port,omitempty"`

	// Scheme readiness probe scheme
	Scheme *string `json:"scheme,omitempty"`
}

// DefaultsServerConfigHealthsStartup defines model for defaults_server_config_healths_startup.
type DefaultsServerConfigHealthsStartup struct {
	// Enabled startup server enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Port startup server port
	Port         *int                                            `json:"port,omitempty"`
	StartupProbe *DefaultsServerConfigHealthsStartupStartupProbe `json:"startupProbe,omitempty"`
}

// DefaultsServerConfigHealthsStartupStartupProbe defines model for defaults_server_config_healths_startup_startupProbe.
type DefaultsServerConfigHealthsStartupStartupProbe struct {
	// FailureThreshold startup probe failure threshold
	FailureThreshold *int                                                   `json:"failureThreshold,omitempty"`
	HttpGet          *DefaultsServerConfigHealthsStartupStartupProbeHttpGet `json:"httpGet,omitempty"`

	// InitialDelaySeconds startup probe initial delay seconds
	InitialDelaySeconds *int `json:"initialDelaySeconds,omitempty"`

	// PeriodSeconds startup probe period seconds
	PeriodSeconds *int `json:"periodSeconds,omitempty"`

	// SuccessThreshold startup probe success threshold
	SuccessThreshold *int `json:"successThreshold,omitempty"`

	// TimeoutSeconds startup probe timeout seconds
	TimeoutSeconds *int `json:"timeoutSeconds,omitempty"`
}

// DefaultsServerConfigHealthsStartupStartupProbeHttpGet defines model for defaults_server_config_healths_startup_startupProbe_httpGet.
type DefaultsServerConfigHealthsStartupStartupProbeHttpGet struct {
	// Path startup probe path
	Path *string `json:"path,omitempty"`

	// Port startup probe port
	Port *string `json:"port,omitempty"`

	// Scheme startup probe scheme
	Scheme *string `json:"scheme,omitempty"`
}

// DefaultsServerConfigMetrics defines model for defaults_server_config_metrics.
type DefaultsServerConfigMetrics struct {
	Pprof *DefaultsServerConfigMetricsPprof `json:"pprof,omitempty"`
}

// DefaultsServerConfigMetricsPprof defines model for defaults_server_config_metrics_pprof.
type DefaultsServerConfigMetricsPprof struct {
	// Enabled pprof server enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Host pprof server host
	Host *string `json:"host,omitempty"`

	// Port pprof server port
	Port   *int        `json:"port,omitempty"`
	Server *RestServer `json:"server,omitempty"`

	// ServicePort pprof server service port
	ServicePort *int `json:"servicePort,omitempty"`
}

// DefaultsServerConfigServers defines model for defaults_server_config_servers.
type DefaultsServerConfigServers struct {
	Grpc *DefaultsServerConfigServersGrpc `json:"grpc,omitempty"`
	Rest *DefaultsServerConfigServersRest `json:"rest,omitempty"`
}

// DefaultsServerConfigServersGrpc defines model for defaults_server_config_servers_grpc.
type DefaultsServerConfigServersGrpc struct {
	// Enabled gRPC server enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Host gRPC server host
	Host *string `json:"host,omitempty"`

	// Port gRPC server port
	Port   *int        `json:"port,omitempty"`
	Server *GrpcServer `json:"server,omitempty"`

	// ServicePort gRPC server service port
	ServicePort *int `json:"servicePort,omitempty"`
}

// DefaultsServerConfigServersGrpcServerGrpcKeepalive defines model for defaults_server_config_servers_grpc_server_grpc_keepalive.
type DefaultsServerConfigServersGrpcServerGrpcKeepalive struct {
	// MaxConnAge gRPC server keep alive max connection age
	MaxConnAge *string `json:"max_conn_age,omitempty"`

	// MaxConnAgeGrace gRPC server keep alive max connection age grace
	MaxConnAgeGrace *string `json:"max_conn_age_grace,omitempty"`

	// MaxConnIdle gRPC server keep alive max connection idle
	MaxConnIdle *string `json:"max_conn_idle,omitempty"`

	// MinTime gRPC server keep alive min_time
	MinTime *string `json:"min_time,omitempty"`

	// PermitWithoutStream gRPC server keep alive permit_without_stream
	PermitWithoutStream *bool `json:"permit_without_stream,omitempty"`

	// Time gRPC server keep alive time
	Time *string `json:"time,omitempty"`

	// Timeout gRPC server keep alive timeout
	Timeout *string `json:"timeout,omitempty"`
}

// DefaultsServerConfigServersRest defines model for defaults_server_config_servers_rest.
type DefaultsServerConfigServersRest struct {
	// Enabled REST server enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Host REST server host
	Host *string `json:"host,omitempty"`

	// Port REST server port
	Port   *int        `json:"port,omitempty"`
	Server *RestServer `json:"server,omitempty"`

	// ServicePort REST server service port
	ServicePort *int `json:"servicePort,omitempty"`
}

// Discoverer defines model for discoverer.
type Discoverer struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// Annotations deployment annotations
	Annotations        *map[string]any               `json:"annotations,omitempty"`
	ClusterRole        *DiscovererClusterRole        `json:"clusterRole,omitempty"`
	ClusterRoleBinding *DiscovererClusterRoleBinding `json:"clusterRoleBinding,omitempty"`
	Discoverer         *DiscovererDiscoverer         `json:"discoverer,omitempty"`

	// Enabled discoverer enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env *Env `json:"env,omitempty"`

	// ExternalTrafficPolicy external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy *string `json:"externalTrafficPolicy,omitempty"`
	Hpa                   *Hpa    `json:"hpa,omitempty"`
	Image                 *Image  `json:"image,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// InternalTrafficPolicy internal traffic policy : Cluster or Local
	InternalTrafficPolicy *string `json:"internalTrafficPolicy,omitempty"`

	// Kind deployment kind: Deployment or DaemonSet
	Kind    *DiscovererKind `json:"kind,omitempty"`
	Logging *Logging        `json:"logging,omitempty"`

	// MaxReplicas maximum number of replicas. if HPA is disabled, this value will be ignored.
	MaxReplicas *int `json:"maxReplicas,omitempty"`

	// MaxUnavailable maximum number of unavailable replicas
	MaxUnavailable *string `json:"maxUnavailable,omitempty"`

	// MinReplicas minimum number of replicas. if HPA is disabled, the replicas will be set to this value
	MinReplicas *int `json:"minReplicas,omitempty"`

	// Name name of discoverer deployment
	Name *string `json:"name,omitempty"`

	// NodeName node name
	NodeName *string `json:"nodeName,omitempty"`

	// NodeSelector node selector
	NodeSelector  *NodeSelector  `json:"nodeSelector,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// PodAnnotations pod annotations
	PodAnnotations *map[string]any `json:"podAnnotations,omitempty"`
	PodPriority    *PodPriority    `json:"podPriority,omitempty"`

	// PodSecurityContext security context for pod
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ProgressDeadlineSeconds progress deadline seconds
	ProgressDeadlineSeconds *int `json:"progressDeadlineSeconds,omitempty"`

	// Resources compute resources
	Resources *Resources `json:"resources,omitempty"`

	// RevisionHistoryLimit number of old history to retain to allow rollback
	RevisionHistoryLimit *int           `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate `json:"rollingUpdate,omitempty"`

	// SecurityContext security context for container
	SecurityContext    *corev1.SecurityContext `json:"securityContext,omitempty"`
	ServerConfig       *ServerConfig           `json:"server_config,omitempty"`
	Service            *Service                `json:"service,omitempty"`
	ServiceAccountName *string                 `json:"serviceAccountName,omitempty"`

	// ServiceType service type: ClusterIP, LoadBalancer or NodePort
	ServiceType *DiscovererServiceType `json:"serviceType,omitempty"`

	// TerminationGracePeriodSeconds duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds *int `json:"terminationGracePeriodSeconds,omitempty"`

	// TimeZone Time zone
	TimeZone *string `json:"time_zone,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TopologySpreadConstraints topology spread constraints of gateway pods
	TopologySpreadConstraints *TopologySpreadConstraints `json:"topologySpreadConstraints,omitempty"`

	// UnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either cautious or permissive eviction.
	UnhealthyPodEvictionPolicy *DiscovererUnhealthyPodEvictionPolicy `json:"unhealthyPodEvictionPolicy,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`

	// VolumeMounts volume mounts
	VolumeMounts *VolumeMounts `json:"volumeMounts,omitempty"`

	// Volumes volumes
	Volumes *Volumes `json:"volumes,omitempty"`
}

// DiscovererKind deployment kind: Deployment or DaemonSet
type DiscovererKind string

// DiscovererServiceType service type: ClusterIP, LoadBalancer or NodePort
type DiscovererServiceType string

// DiscovererUnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either cautious
// or permissive eviction.
type DiscovererUnhealthyPodEvictionPolicy string

// DiscovererClusterRole defines model for discoverer_clusterRole.
type DiscovererClusterRole struct {
	// Enabled creates clusterRole resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRole
	Name *string `json:"name,omitempty"`
}

// DiscovererClusterRoleBinding defines model for discoverer_clusterRoleBinding.
type DiscovererClusterRoleBinding struct {
	// Enabled creates clusterRoleBinding resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRoleBinding
	Name *string `json:"name,omitempty"`
}

// DiscovererDiscoverer defines model for discoverer_discoverer.
type DiscovererDiscoverer struct {
	// DiscoveryDuration duration to discovery
	DiscoveryDuration *string `json:"discovery_duration,omitempty"`

	// Name name to discovery
	Name *string `json:"name,omitempty"`

	// Namespace namespace to discovery
	Namespace *string `json:"namespace,omitempty"`
	Net       *Net    `json:"net,omitempty"`

	// Selectors k8s resource selectors
	Selectors *DiscovererDiscovererSelectors `json:"selectors,omitempty"`
}

// DiscovererDiscovererSelectors k8s resource selectors
type DiscovererDiscovererSelectors struct {
	// Node k8s resource selectors for node discovery
	Node *DiscovererDiscovererSelectorsNode `json:"node,omitempty"`

	// NodeMetrics k8s resource selectors for node_metrics discovery
	NodeMetrics *DiscovererDiscovererSelectorsNodeMetrics `json:"node_metrics,omitempty"`

	// Pod k8s resource selectors for pod discovery
	Pod *DiscovererDiscovererSelectorsPod `json:"pod,omitempty"`

	// PodMetrics k8s resource selectors for pod_metrics discovery
	PodMetrics *DiscovererDiscovererSelectorsPodMetrics `json:"pod_metrics,omitempty"`

	// Service k8s resource selectors for service discovery
	Service *DiscovererDiscovererSelectorsService `json:"service,omitempty"`
}

// DiscovererDiscovererSelectorsNode k8s resource selectors for node discovery
type DiscovererDiscovererSelectorsNode struct {
	// Fields k8s field selectors for node discovery
	Fields *map[string]any `json:"fields,omitempty"`

	// Labels k8s label selectors for node discovery
	Labels *map[string]any `json:"labels,omitempty"`
}

// DiscovererDiscovererSelectorsNodeMetrics k8s resource selectors for node_metrics discovery
type DiscovererDiscovererSelectorsNodeMetrics struct {
	// Fields k8s field selectors for node_metrics discovery
	Fields *map[string]any `json:"fields,omitempty"`

	// Labels k8s label selectors for node_metrics discovery
	Labels *map[string]any `json:"labels,omitempty"`
}

// DiscovererDiscovererSelectorsPod k8s resource selectors for pod discovery
type DiscovererDiscovererSelectorsPod struct {
	// Fields k8s field selectors for pod discovery
	Fields *map[string]any `json:"fields,omitempty"`

	// Labels k8s label selectors for pod discovery
	Labels *map[string]any `json:"labels,omitempty"`
}

// DiscovererDiscovererSelectorsPodMetrics k8s resource selectors for pod_metrics discovery
type DiscovererDiscovererSelectorsPodMetrics struct {
	// Fields k8s field selectors for pod_metrics discovery
	Fields *map[string]any `json:"fields,omitempty"`

	// Labels k8s label selectors for pod_metrics discovery
	Labels *map[string]any `json:"labels,omitempty"`
}

// DiscovererDiscovererSelectorsService k8s resource selectors for service discovery
type DiscovererDiscovererSelectorsService struct {
	// Fields k8s field selectors for service discovery
	Fields *map[string]any `json:"fields,omitempty"`

	// Labels k8s label selectors for service discovery
	Labels *map[string]any `json:"labels,omitempty"`
}

// Env environment variables
type Env = []corev1.EnvVar

// Gateway defines model for gateway.
type Gateway struct {
	Filter *GatewayFilter `json:"filter,omitempty"`
	Lb     *GatewayLb     `json:"lb,omitempty"`
	Mirror *GatewayMirror `json:"mirror,omitempty"`
}

// GatewayFilter defines model for gateway_filter.
type GatewayFilter struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// Annotations deployment annotations
	Annotations *map[string]any `json:"annotations,omitempty"`

	// Enabled gateway enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env *Env `json:"env,omitempty"`

	// ExternalTrafficPolicy external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy *string                     `json:"externalTrafficPolicy,omitempty"`
	GatewayConfig         *GatewayFilterGatewayConfig `json:"gateway_config,omitempty"`
	Hpa                   *Hpa                        `json:"hpa,omitempty"`
	Image                 *Image                      `json:"image,omitempty"`
	Ingress               *GatewayFilterIngress       `json:"ingress,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// InternalTrafficPolicy internal traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	InternalTrafficPolicy *string `json:"internalTrafficPolicy,omitempty"`

	// Kind deployment kind: Deployment or DaemonSet
	Kind    *GatewayFilterKind `json:"kind,omitempty"`
	Logging *Logging           `json:"logging,omitempty"`

	// MaxReplicas maximum number of replicas. if HPA is disabled, this value will be ignored.
	MaxReplicas *int `json:"maxReplicas,omitempty"`

	// MaxUnavailable maximum number of unavailable replicas
	MaxUnavailable *string `json:"maxUnavailable,omitempty"`

	// MinReplicas minimum number of replicas. if HPA is disabled, the replicas will be set to this value
	MinReplicas *int `json:"minReplicas,omitempty"`

	// Name name of filter gateway deployment
	Name *string `json:"name,omitempty"`

	// NodeName node name
	NodeName *string `json:"nodeName,omitempty"`

	// NodeSelector node selector
	NodeSelector  *NodeSelector  `json:"nodeSelector,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// PodAnnotations pod annotations
	PodAnnotations *map[string]any `json:"podAnnotations,omitempty"`
	PodPriority    *PodPriority    `json:"podPriority,omitempty"`

	// PodSecurityContext security context for pod
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ProgressDeadlineSeconds progress deadline seconds
	ProgressDeadlineSeconds *int `json:"progressDeadlineSeconds,omitempty"`

	// Resources compute resources
	Resources *Resources `json:"resources,omitempty"`

	// RevisionHistoryLimit number of old history to retain to allow rollback
	RevisionHistoryLimit *int           `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate `json:"rollingUpdate,omitempty"`

	// SecurityContext security context for container
	SecurityContext    *corev1.SecurityContext `json:"securityContext,omitempty"`
	ServerConfig       *ServerConfig           `json:"server_config,omitempty"`
	Service            *Service                `json:"service,omitempty"`
	ServiceAccountName *string                 `json:"serviceAccountName,omitempty"`

	// ServiceType service type: ClusterIP, LoadBalancer or NodePort
	ServiceType *GatewayFilterServiceType `json:"serviceType,omitempty"`

	// TerminationGracePeriodSeconds duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds *int `json:"terminationGracePeriodSeconds,omitempty"`

	// TimeZone Time zone
	TimeZone *string `json:"time_zone,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TopologySpreadConstraints topology spread constraints of gateway pods
	TopologySpreadConstraints *TopologySpreadConstraints `json:"topologySpreadConstraints,omitempty"`

	// UnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either cautious or permissive eviction.
	UnhealthyPodEvictionPolicy *GatewayFilterUnhealthyPodEvictionPolicy `json:"unhealthyPodEvictionPolicy,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`

	// VolumeMounts volume mounts
	VolumeMounts *VolumeMounts `json:"volumeMounts,omitempty"`

	// Volumes volumes
	Volumes *Volumes `json:"volumes,omitempty"`
}

// GatewayFilterKind deployment kind: Deployment or DaemonSet
type GatewayFilterKind string

// GatewayFilterServiceType service type: ClusterIP, LoadBalancer or NodePort
type GatewayFilterServiceType string

// GatewayFilterUnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either
// cautious or permissive eviction.
type GatewayFilterUnhealthyPodEvictionPolicy string

// GatewayFilterGatewayConfig defines model for gateway_filter_gateway_config.
type GatewayFilterGatewayConfig struct {
	// EgressFilter gRPC client config for egress filter
	EgressFilter  *GatewayFilterGatewayConfigEgressFilter `json:"egress_filter,omitempty"`
	GatewayClient *GrpcClient                             `json:"gateway_client,omitempty"`

	// IngressFilter gRPC client config for ingress filter
	IngressFilter *GatewayFilterGatewayConfigIngressFilter `json:"ingress_filter,omitempty"`
}

// GatewayFilterGatewayConfigEgressFilter gRPC client config for egress filter
type GatewayFilterGatewayConfigEgressFilter struct {
	Client *GrpcClient `json:"client,omitempty"`

	// DistanceFilters distance egress vector filter targets
	DistanceFilters *[]string `json:"distance_filters,omitempty"`

	// ObjectFilters object egress vector filter targets
	ObjectFilters *[]string `json:"object_filters,omitempty"`
}

// GatewayFilterGatewayConfigIngressFilter gRPC client config for ingress filter
type GatewayFilterGatewayConfigIngressFilter struct {
	Client *GrpcClient `json:"client,omitempty"`

	// InsertFilters insert ingress vector filter targets
	InsertFilters *[]string `json:"insert_filters,omitempty"`

	// SearchFilters search ingress vector filter targets
	SearchFilters *[]string `json:"search_filters,omitempty"`

	// UpdateFilters update ingress vector filter targets
	UpdateFilters *[]string `json:"update_filters,omitempty"`

	// UpsertFilters upsert ingress vector filter targets
	UpsertFilters *[]string `json:"upsert_filters,omitempty"`

	// Vectorizer object ingress vectorize filter targets
	Vectorizer *string `json:"vectorizer,omitempty"`
}

// GatewayFilterIngress defines model for gateway_filter_ingress.
type GatewayFilterIngress struct {
	// Annotations annotations for ingress
	Annotations *map[string]any `json:"annotations,omitempty"`

	// DefaultBackend defaultBackend config
	DefaultBackend *GatewayFilterIngressDefaultBackend `json:"defaultBackend,omitempty"`

	// Enabled gateway ingress enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Host ingress hostname
	Host *string `json:"host,omitempty"`

	// PathType gateway ingress pathType
	PathType *string `json:"pathType,omitempty"`

	// ServicePort service port to be exposed by ingress
	ServicePort *string `json:"servicePort,omitempty"`

	// Tls ingress tls config
	Tls *[]map[string]any `json:"tls,omitempty"`
}

// GatewayFilterIngressDefaultBackend defaultBackend config
type GatewayFilterIngressDefaultBackend struct {
	// Enabled gateway ingress defaultBackend enabled
	Enabled *bool `json:"enabled,omitempty"`
}

// GatewayLb defines model for gateway_lb.
type GatewayLb struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// Annotations deployment annotations
	Annotations *map[string]any `json:"annotations,omitempty"`

	// Enabled gateway enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env *Env `json:"env,omitempty"`

	// ExternalTrafficPolicy external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy *string                 `json:"externalTrafficPolicy,omitempty"`
	GatewayConfig         *GatewayLbGatewayConfig `json:"gateway_config,omitempty"`
	Hpa                   *Hpa                    `json:"hpa,omitempty"`
	Image                 *Image                  `json:"image,omitempty"`
	Ingress               *GatewayLbIngress       `json:"ingress,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// InternalTrafficPolicy internal traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	InternalTrafficPolicy *string `json:"internalTrafficPolicy,omitempty"`

	// Kind deployment kind: Deployment or DaemonSet
	Kind    *GatewayLbKind `json:"kind,omitempty"`
	Logging *Logging       `json:"logging,omitempty"`

	// MaxReplicas maximum number of replicas. if HPA is disabled, this value will be ignored.
	MaxReplicas *int `json:"maxReplicas,omitempty"`

	// MaxUnavailable maximum number of unavailable replicas
	MaxUnavailable *string `json:"maxUnavailable,omitempty"`

	// MinReplicas minimum number of replicas. if HPA is disabled, the replicas will be set to this value
	MinReplicas *int `json:"minReplicas,omitempty"`

	// Name name of gateway deployment
	Name *string `json:"name,omitempty"`

	// NodeName node name
	NodeName *string `json:"nodeName,omitempty"`

	// NodeSelector node selector
	NodeSelector  *NodeSelector  `json:"nodeSelector,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// PodAnnotations pod annotations
	PodAnnotations *map[string]any `json:"podAnnotations,omitempty"`
	PodPriority    *PodPriority    `json:"podPriority,omitempty"`

	// PodSecurityContext security context for pod
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ProgressDeadlineSeconds progress deadline seconds
	ProgressDeadlineSeconds *int `json:"progressDeadlineSeconds,omitempty"`

	// Resources compute resources
	Resources *Resources `json:"resources,omitempty"`

	// RevisionHistoryLimit number of old history to retain to allow rollback
	RevisionHistoryLimit *int           `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate `json:"rollingUpdate,omitempty"`

	// SecurityContext security context for container
	SecurityContext    *corev1.SecurityContext `json:"securityContext,omitempty"`
	ServerConfig       *ServerConfig           `json:"server_config,omitempty"`
	Service            *Service                `json:"service,omitempty"`
	ServiceAccountName *string                 `json:"serviceAccountName,omitempty"`

	// ServiceType service type: ClusterIP, LoadBalancer or NodePort
	ServiceType *GatewayLbServiceType `json:"serviceType,omitempty"`

	// TerminationGracePeriodSeconds duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds *int `json:"terminationGracePeriodSeconds,omitempty"`

	// TimeZone Time zone
	TimeZone *string `json:"time_zone,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TopologySpreadConstraints topology spread constraints of gateway pods
	TopologySpreadConstraints *TopologySpreadConstraints `json:"topologySpreadConstraints,omitempty"`

	// UnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either cautious or permissive eviction.
	UnhealthyPodEvictionPolicy *GatewayLbUnhealthyPodEvictionPolicy `json:"unhealthyPodEvictionPolicy,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`

	// VolumeMounts volume mounts
	VolumeMounts *VolumeMounts `json:"volumeMounts,omitempty"`

	// Volumes volumes
	Volumes *Volumes `json:"volumes,omitempty"`
}

// GatewayLbKind deployment kind: Deployment or DaemonSet
type GatewayLbKind string

// GatewayLbServiceType service type: ClusterIP, LoadBalancer or NodePort
type GatewayLbServiceType string

// GatewayLbUnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either cautious
// or permissive eviction.
type GatewayLbUnhealthyPodEvictionPolicy string

// GatewayLbGatewayConfig defines model for gateway_lb_gateway_config.
type GatewayLbGatewayConfig struct {
	// AgentNamespace agent namespace
	AgentNamespace *string                           `json:"agent_namespace,omitempty"`
	Discoverer     *GatewayLbGatewayConfigDiscoverer `json:"discoverer,omitempty"`

	// IndexReplica number of index replica
	IndexReplica *int `json:"index_replica,omitempty"`

	// MultiOperationConcurrency number of concurrency of multiXXX api's operation
	MultiOperationConcurrency *int `json:"multi_operation_concurrency,omitempty"`

	// NodeName node name
	NodeName *string `json:"node_name,omitempty"`
}

// GatewayLbGatewayConfigDiscoverer defines model for gateway_lb_gateway_config_discoverer.
type GatewayLbGatewayConfigDiscoverer struct {
	AgentClientOptions *GrpcClient `json:"agent_client_options,omitempty"`
	Client             *GrpcClient `json:"client,omitempty"`
	Duration           *string     `json:"duration,omitempty"`
	ReadClient         *GrpcClient `json:"read_client,omitempty"`
}

// GatewayLbIngress defines model for gateway_lb_ingress.
type GatewayLbIngress struct {
	// Annotations annotations for ingress
	Annotations *map[string]any `json:"annotations,omitempty"`

	// DefaultBackend defaultBackend config
	DefaultBackend *GatewayLbIngressDefaultBackend `json:"defaultBackend,omitempty"`

	// Enabled gateway ingress enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Host ingress hostname
	Host *string `json:"host,omitempty"`

	// PathType gateway ingress pathType
	PathType *string `json:"pathType,omitempty"`

	// ServicePort service port to be exposed by ingress
	ServicePort *string `json:"servicePort,omitempty"`

	// Tls ingress tls config
	Tls *[]map[string]any `json:"tls,omitempty"`
}

// GatewayLbIngressDefaultBackend defaultBackend config
type GatewayLbIngressDefaultBackend struct {
	// Enabled gateway ingress defaultBackend enabled
	Enabled *bool `json:"enabled,omitempty"`
}

// GatewayMirror defines model for gateway_mirror.
type GatewayMirror struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// Annotations deployment annotations
	Annotations        *map[string]any                  `json:"annotations,omitempty"`
	ClusterRole        *GatewayMirrorClusterRole        `json:"clusterRole,omitempty"`
	ClusterRoleBinding *GatewayMirrorClusterRoleBinding `json:"clusterRoleBinding,omitempty"`

	// Enabled gateway enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env *Env `json:"env,omitempty"`

	// ExternalTrafficPolicy external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy *string                     `json:"externalTrafficPolicy,omitempty"`
	GatewayConfig         *GatewayMirrorGatewayConfig `json:"gateway_config,omitempty"`
	Hpa                   *Hpa                        `json:"hpa,omitempty"`
	Image                 *Image                      `json:"image,omitempty"`
	Ingress               *GatewayMirrorIngress       `json:"ingress,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// InternalTrafficPolicy internal traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	InternalTrafficPolicy *string `json:"internalTrafficPolicy,omitempty"`

	// Kind deployment kind: Deployment or DaemonSet
	Kind    *GatewayMirrorKind `json:"kind,omitempty"`
	Logging *Logging           `json:"logging,omitempty"`

	// MaxReplicas maximum number of replicas. if HPA is disabled, this value will be ignored.
	MaxReplicas *int `json:"maxReplicas,omitempty"`

	// MaxUnavailable maximum number of unavailable replicas
	MaxUnavailable *string `json:"maxUnavailable,omitempty"`

	// MinReplicas minimum number of replicas. if HPA is disabled, the replicas will be set to this value
	MinReplicas *int `json:"minReplicas,omitempty"`

	// Name name of gateway deployment
	Name *string `json:"name,omitempty"`

	// NodeName node name
	NodeName *string `json:"nodeName,omitempty"`

	// NodeSelector node selector
	NodeSelector  *NodeSelector  `json:"nodeSelector,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// PodAnnotations pod annotations
	PodAnnotations *map[string]any `json:"podAnnotations,omitempty"`
	PodPriority    *PodPriority    `json:"podPriority,omitempty"`

	// PodSecurityContext security context for pod
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ProgressDeadlineSeconds progress deadline seconds
	ProgressDeadlineSeconds *int `json:"progressDeadlineSeconds,omitempty"`

	// Resources compute resources
	Resources *Resources `json:"resources,omitempty"`

	// RevisionHistoryLimit number of old history to retain to allow rollback
	RevisionHistoryLimit *int           `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate `json:"rollingUpdate,omitempty"`

	// SecurityContext security context for container
	SecurityContext    *corev1.SecurityContext `json:"securityContext,omitempty"`
	ServerConfig       *ServerConfig           `json:"server_config,omitempty"`
	Service            *Service                `json:"service,omitempty"`
	ServiceAccountName *string                 `json:"serviceAccountName,omitempty"`

	// ServiceType service type: ClusterIP, LoadBalancer or NodePort
	ServiceType *GatewayMirrorServiceType `json:"serviceType,omitempty"`

	// TerminationGracePeriodSeconds duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds *int `json:"terminationGracePeriodSeconds,omitempty"`

	// TimeZone Time zone
	TimeZone *string `json:"time_zone,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TopologySpreadConstraints topology spread constraints of gateway pods
	TopologySpreadConstraints *TopologySpreadConstraints `json:"topologySpreadConstraints,omitempty"`

	// UnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either cautious or permissive eviction.
	UnhealthyPodEvictionPolicy *GatewayMirrorUnhealthyPodEvictionPolicy `json:"unhealthyPodEvictionPolicy,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`

	// VolumeMounts volume mounts
	VolumeMounts *VolumeMounts `json:"volumeMounts,omitempty"`

	// Volumes volumes
	Volumes *Volumes `json:"volumes,omitempty"`
}

// GatewayMirrorKind deployment kind: Deployment or DaemonSet
type GatewayMirrorKind string

// GatewayMirrorServiceType service type: ClusterIP, LoadBalancer or NodePort
type GatewayMirrorServiceType string

// GatewayMirrorUnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either
// cautious or permissive eviction.
type GatewayMirrorUnhealthyPodEvictionPolicy string

// GatewayMirrorAffinityNodeAffinity defines model for gateway_mirror_affinity_nodeAffinity.
type GatewayMirrorAffinityNodeAffinity struct {
	// PreferredDuringSchedulingIgnoredDuringExecution node affinity preferred scheduling terms
	PreferredDuringSchedulingIgnoredDuringExecution *[]map[string]any                                                                `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
	RequiredDuringSchedulingIgnoredDuringExecution  *GatewayMirrorAffinityNodeAffinityRequiredDuringSchedulingIgnoredDuringExecution `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// GatewayMirrorAffinityNodeAffinityRequiredDuringSchedulingIgnoredDuringExecution defines model for
// gateway_mirror_affinity_nodeAffinity_requiredDuringSchedulingIgnoredDuringExecution.
type GatewayMirrorAffinityNodeAffinityRequiredDuringSchedulingIgnoredDuringExecution struct {
	// NodeSelectorTerms node affinity required node selectors
	NodeSelectorTerms *[]map[string]any `json:"nodeSelectorTerms,omitempty"`
}

// GatewayMirrorAffinityPodAffinity defines model for gateway_mirror_affinity_podAffinity.
type GatewayMirrorAffinityPodAffinity struct {
	// PreferredDuringSchedulingIgnoredDuringExecution pod affinity preferred scheduling terms
	PreferredDuringSchedulingIgnoredDuringExecution *[]map[string]any `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`

	// RequiredDuringSchedulingIgnoredDuringExecution pod affinity required scheduling terms
	RequiredDuringSchedulingIgnoredDuringExecution *[]map[string]any `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// GatewayMirrorAffinityPodAntiAffinity defines model for gateway_mirror_affinity_podAntiAffinity.
type GatewayMirrorAffinityPodAntiAffinity struct {
	// PreferredDuringSchedulingIgnoredDuringExecution pod anti-affinity preferred scheduling terms
	PreferredDuringSchedulingIgnoredDuringExecution *[]map[string]any `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`

	// RequiredDuringSchedulingIgnoredDuringExecution pod anti-affinity required scheduling terms
	RequiredDuringSchedulingIgnoredDuringExecution *[]map[string]any `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// GatewayMirrorClusterRole defines model for gateway_mirror_clusterRole.
type GatewayMirrorClusterRole struct {
	// Enabled creates clusterRole resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRole
	Name *string `json:"name,omitempty"`
}

// GatewayMirrorClusterRoleBinding defines model for gateway_mirror_clusterRoleBinding.
type GatewayMirrorClusterRoleBinding struct {
	// Enabled creates clusterRoleBinding resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRoleBinding
	Name *string `json:"name,omitempty"`
}

// GatewayMirrorGatewayConfig defines model for gateway_mirror_gateway_config.
type GatewayMirrorGatewayConfig struct {
	Client *GrpcClient `json:"client,omitempty"`

	// Colocation colocation name
	Colocation *string `json:"colocation,omitempty"`

	// DiscoveryDuration duration to discovery
	DiscoveryDuration *string `json:"discovery_duration,omitempty"`

	// GatewayAddr address for lb-gateway
	GatewayAddr *string `json:"gateway_addr,omitempty"`

	// Group mirror group name
	Group *string `json:"group,omitempty"`

	// Namespace namespace to discovery
	Namespace *string `json:"namespace,omitempty"`
	Net       *Net    `json:"net,omitempty"`

	// PodName self mirror gateway pod name
	PodName *string `json:"pod_name,omitempty"`

	// RegisterDuration duration to register mirror-gateway.
	RegisterDuration *string `json:"register_duration,omitempty"`

	// SelfMirrorAddr address for self mirror-gateway
	SelfMirrorAddr *string `json:"self_mirror_addr,omitempty"`
}

// GatewayMirrorIngress defines model for gateway_mirror_ingress.
type GatewayMirrorIngress struct {
	// Annotations annotations for ingress
	Annotations *map[string]any `json:"annotations,omitempty"`

	// DefaultBackend defaultBackend config
	DefaultBackend *GatewayMirrorIngressDefaultBackend `json:"defaultBackend,omitempty"`

	// Enabled gateway ingress enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Host ingress hostname
	Host *string `json:"host,omitempty"`

	// PathType gateway ingress pathType
	PathType *string `json:"pathType,omitempty"`

	// ServicePort service port to be exposed by ingress
	ServicePort *string           `json:"servicePort,omitempty"`
	Tls         *[]map[string]any `json:"tls,omitempty"`
}

// GatewayMirrorIngressDefaultBackend defaultBackend config
type GatewayMirrorIngressDefaultBackend struct {
	// Enabled gateway ingress defaultBackend enabled
	Enabled *bool `json:"enabled,omitempty"`
}

// GrpcClient defines model for grpc_client.
type GrpcClient struct {
	// Addrs gRPC client addresses
	Addrs          *[]string                         `json:"addrs,omitempty"`
	Backoff        *Backoff                          `json:"backoff,omitempty"`
	CallOption     *map[string]any                   `json:"call_option,omitempty"`
	CircuitBreaker *DefaultsGrpcClientCircuitBreaker `json:"circuit_breaker,omitempty"`
	ConnectionPool *DefaultsGrpcClientConnectionPool `json:"connection_pool,omitempty"`
	ContentSubtype *string                           `json:"content_subtype,omitempty"`
	DialOption     *DefaultsGrpcClientDialOption     `json:"dial_option,omitempty"`

	// HealthCheckDuration gRPC client health check duration
	HealthCheckDuration   *string `json:"health_check_duration,omitempty"`
	MaxRecvMsgSize        *int    `json:"max_recv_msg_size,omitempty"`
	MaxRetryRpcBufferSize *int    `json:"max_retry_rpc_buffer_size,omitempty"`
	MaxSendMsgSize        *int    `json:"max_send_msg_size,omitempty"`
	Tls                   *Tls    `json:"tls,omitempty"`
	WaitForReady          *bool   `json:"wait_for_ready,omitempty"`
}

// GrpcServer defines model for grpc_server.
type GrpcServer = config.Server

// GrpcServerConfig defines model for grpc_server_config.
type GrpcServerConfig = config.GRPC

// Hpa defines model for hpa.
type Hpa struct {
	// Enabled HPA enabled
	Enabled *bool `json:"enabled,omitempty"`

	// TargetCPUUtilizationPercentage HPA CPU utilization percentage
	TargetCPUUtilizationPercentage *int `json:"targetCPUUtilizationPercentage,omitempty"`
}

// Http2ServerConfig defines model for http2_server_config.
type Http2ServerConfig = config.HTTP2

// HttpServerConfig defines model for http_server_config.
type HttpServerConfig = config.HTTP

// Image defines model for image.
type Image struct {
	// PullPolicy image pull policy
	PullPolicy *ImagePullPolicy `json:"pullPolicy,omitempty"`

	// Repository image repository
	Repository *string `json:"repository,omitempty"`

	// Tag image tag (overrides defaults.image.tag)
	Tag *string `json:"tag,omitempty"`
}

// ImagePullPolicy image pull policy
type ImagePullPolicy string

// InitContainers init containers
type InitContainers = []corev1.Container

// Logging defines model for logging.
type Logging = config.Logging

// Manager defines model for manager.
type Manager struct {
	Index *ManagerIndex `json:"index,omitempty"`
}

// ManagerIndex defines model for manager_index.
type ManagerIndex struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// Annotations deployment annotations
	Annotations *map[string]any        `json:"annotations,omitempty"`
	Corrector   *ManagerIndexCorrector `json:"corrector,omitempty"`
	Creator     *ManagerIndexCreator   `json:"creator,omitempty"`
	Deleter     *ManagerIndexDeleter   `json:"deleter,omitempty"`

	// Enabled index manager enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env *Env `json:"env,omitempty"`

	// ExternalTrafficPolicy external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy *string              `json:"externalTrafficPolicy,omitempty"`
	Image                 *Image               `json:"image,omitempty"`
	Indexer               *ManagerIndexIndexer `json:"indexer,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// Kind deployment kind: Deployment or DaemonSet
	Kind    *ManagerIndexKind `json:"kind,omitempty"`
	Logging *Logging          `json:"logging,omitempty"`

	// MaxUnavailable maximum number of unavailable replicas
	MaxUnavailable *string `json:"maxUnavailable,omitempty"`

	// Name name of index manager deployment
	Name *string `json:"name,omitempty"`

	// NodeName node name
	NodeName *string `json:"nodeName,omitempty"`

	// NodeSelector node selector
	NodeSelector  *NodeSelector  `json:"nodeSelector,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// Operator [THIS FEATURE IS WIP] operator that manages vald index
	Operator *ManagerIndexOperator `json:"operator,omitempty"`

	// PodAnnotations pod annotations
	PodAnnotations *map[string]any `json:"podAnnotations,omitempty"`
	PodPriority    *PodPriority    `json:"podPriority,omitempty"`

	// PodSecurityContext security context for pod
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ProgressDeadlineSeconds progress deadline seconds
	ProgressDeadlineSeconds *int                     `json:"progressDeadlineSeconds,omitempty"`
	Readreplica             *ManagerIndexReadreplica `json:"readreplica,omitempty"`

	// Replicas number of replicas
	Replicas *int `json:"replicas,omitempty"`

	// Resources compute resources
	Resources *Resources `json:"resources,omitempty"`

	// RevisionHistoryLimit number of old history to retain to allow rollback
	RevisionHistoryLimit *int               `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate     `json:"rollingUpdate,omitempty"`
	Saver                *ManagerIndexSaver `json:"saver,omitempty"`

	// SecurityContext security context for container
	SecurityContext    *corev1.SecurityContext `json:"securityContext,omitempty"`
	ServerConfig       *ServerConfig           `json:"server_config,omitempty"`
	Service            *Service                `json:"service,omitempty"`
	ServiceAccountName *string                 `json:"serviceAccountName,omitempty"`

	// ServiceType service type: ClusterIP, LoadBalancer or NodePort
	ServiceType *ManagerIndexServiceType `json:"serviceType,omitempty"`

	// TerminationGracePeriodSeconds duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds *int `json:"terminationGracePeriodSeconds,omitempty"`

	// TimeZone Time zone
	TimeZone *string `json:"time_zone,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TopologySpreadConstraints topology spread constraints of gateway pods
	TopologySpreadConstraints *TopologySpreadConstraints `json:"topologySpreadConstraints,omitempty"`

	// UnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either cautious or permissive eviction.
	UnhealthyPodEvictionPolicy *ManagerIndexUnhealthyPodEvictionPolicy `json:"unhealthyPodEvictionPolicy,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`

	// VolumeMounts volume mounts
	VolumeMounts *VolumeMounts `json:"volumeMounts,omitempty"`

	// Volumes volumes
	Volumes *Volumes `json:"volumes,omitempty"`
}

// ManagerIndexKind deployment kind: Deployment or DaemonSet
type ManagerIndexKind string

// ManagerIndexServiceType service type: ClusterIP, LoadBalancer or NodePort
type ManagerIndexServiceType string

// ManagerIndexUnhealthyPodEvictionPolicy controls whether unhealthy pods can be evicted based on the application's healthy pod count, supporting either
// cautious or permissive eviction.
type ManagerIndexUnhealthyPodEvictionPolicy string

// ManagerIndexCorrector defines model for manager_index_corrector.
type ManagerIndexCorrector struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// AgentNamespace namespace of agent pods to manage
	AgentNamespace *string                          `json:"agent_namespace,omitempty"`
	Discoverer     *ManagerIndexCorrectorDiscoverer `json:"discoverer,omitempty"`

	// Enabled enable index correction CronJob
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env     *Env        `json:"env,omitempty"`
	Gateway *GrpcClient `json:"gateway,omitempty"`
	Image   *Image      `json:"image,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// KvsBackgroundCompactionInterval interval of checked id list kvs compaction
	KvsBackgroundCompactionInterval *string `json:"kvs_background_compaction_interval,omitempty"`

	// KvsBackgroundSyncInterval interval of checked id list kvs sync
	KvsBackgroundSyncInterval *string `json:"kvs_background_sync_interval,omitempty"`

	// Name name of index correction job
	Name *string `json:"name,omitempty"`

	// NodeSelector node selector
	NodeSelector *NodeSelector `json:"nodeSelector,omitempty"`

	// NodeName node name
	NodeName      *string        `json:"node_name,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// Schedule CronJob schedule setting for index correction
	Schedule           *string       `json:"schedule,omitempty"`
	ServerConfig       *ServerConfig `json:"server_config,omitempty"`
	ServiceAccountName *string       `json:"serviceAccountName,omitempty"`

	// StartingDeadlineSeconds startingDeadlineSeconds setting for K8s completed jobs
	StartingDeadlineSeconds *int `json:"startingDeadlineSeconds,omitempty"`

	// StreamListConcurrency concurrency for stream list object rpc
	StreamListConcurrency *int `json:"stream_list_concurrency,omitempty"`

	// Suspend CronJob suspend setting for index correction
	Suspend *bool `json:"suspend,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TtlSecondsAfterFinished ttl setting for K8s completed jobs
	TtlSecondsAfterFinished *int `json:"ttlSecondsAfterFinished,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`
}

// ManagerIndexCorrectorDiscoverer defines model for manager_index_corrector_discoverer.
type ManagerIndexCorrectorDiscoverer struct {
	AgentClientOptions *GrpcClient `json:"agent_client_options,omitempty"`
	Client             *GrpcClient `json:"client,omitempty"`

	// Duration refresh duration to discover
	Duration *string `json:"duration,omitempty"`
}

// ManagerIndexCreator defines model for manager_index_creator.
type ManagerIndexCreator struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// AgentNamespace namespace of agent pods to manage
	AgentNamespace *string `json:"agent_namespace,omitempty"`

	// Concurrency concurrency for indexing
	Concurrency *int `json:"concurrency,omitempty"`

	// CreationPoolSize number of pool size of create index processing
	CreationPoolSize *int                           `json:"creation_pool_size,omitempty"`
	Discoverer       *ManagerIndexCreatorDiscoverer `json:"discoverer,omitempty"`

	// Enabled enable index creation CronJob
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env   *Env   `json:"env,omitempty"`
	Image *Image `json:"image,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// Name name of index creation job
	Name *string `json:"name,omitempty"`

	// NodeSelector node selector
	NodeSelector *NodeSelector `json:"nodeSelector,omitempty"`

	// NodeName node name
	NodeName      *string        `json:"node_name,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// Schedule CronJob schedule setting for index creation
	Schedule           *string       `json:"schedule,omitempty"`
	ServerConfig       *ServerConfig `json:"server_config,omitempty"`
	ServiceAccountName *string       `json:"serviceAccountName,omitempty"`

	// StartingDeadlineSeconds startingDeadlineSeconds setting for K8s completed jobs
	StartingDeadlineSeconds *int `json:"startingDeadlineSeconds,omitempty"`

	// Suspend CronJob suspend setting for index creation
	Suspend *bool `json:"suspend,omitempty"`

	// TargetAddrs indexing target addresses
	TargetAddrs *[]string `json:"target_addrs,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TtlSecondsAfterFinished ttl setting for K8s completed jobs
	TtlSecondsAfterFinished *int `json:"ttlSecondsAfterFinished,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`
}

// ManagerIndexCreatorDiscoverer defines model for manager_index_creator_discoverer.
type ManagerIndexCreatorDiscoverer struct {
	AgentClientOptions *GrpcClient `json:"agent_client_options,omitempty"`
	Client             *GrpcClient `json:"client,omitempty"`

	// Duration refresh duration to discover
	Duration *string `json:"duration,omitempty"`
}

// ManagerIndexDeleter defines model for manager_index_deleter.
type ManagerIndexDeleter struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// AgentNamespace namespace of agent pods to manage
	AgentNamespace *string `json:"agent_namespace,omitempty"`

	// Concurrency concurrency for indexing
	Concurrency *int                           `json:"concurrency,omitempty"`
	Discoverer  *ManagerIndexDeleterDiscoverer `json:"discoverer,omitempty"`

	// Enabled enable index deletion CronJob
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env   *Env   `json:"env,omitempty"`
	Image *Image `json:"image,omitempty"`

	// IndexId index id for deletion
	IndexId *string `json:"index_id,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// Name name of index deletion job
	Name *string `json:"name,omitempty"`

	// NodeSelector node selector
	NodeSelector *NodeSelector `json:"nodeSelector,omitempty"`

	// NodeName node name
	NodeName      *string        `json:"node_name,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// Schedule CronJob schedule setting for index deletion
	Schedule           *string       `json:"schedule,omitempty"`
	ServerConfig       *ServerConfig `json:"server_config,omitempty"`
	ServiceAccountName *string       `json:"serviceAccountName,omitempty"`

	// StartingDeadlineSeconds startingDeadlineSeconds setting for K8s completed jobs
	StartingDeadlineSeconds *int `json:"startingDeadlineSeconds,omitempty"`

	// Suspend CronJob suspend setting for index deletion
	Suspend *bool `json:"suspend,omitempty"`

	// TargetAddrs indexing target addresses
	TargetAddrs *[]string `json:"target_addrs,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TtlSecondsAfterFinished ttl setting for K8s completed jobs
	TtlSecondsAfterFinished *int `json:"ttlSecondsAfterFinished,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`
}

// ManagerIndexDeleterDiscoverer defines model for manager_index_deleter_discoverer.
type ManagerIndexDeleterDiscoverer struct {
	AgentClientOptions *GrpcClient `json:"agent_client_options,omitempty"`
	Client             *GrpcClient `json:"client,omitempty"`

	// Duration refresh duration to discover
	Duration *string `json:"duration,omitempty"`
}

// ManagerIndexIndexer defines model for manager_index_indexer.
type ManagerIndexIndexer struct {
	// AgentNamespace namespace of agent pods to manage
	AgentNamespace *string `json:"agent_namespace,omitempty"`

	// AutoIndexCheckDuration check duration of automatic indexing
	AutoIndexCheckDuration *string `json:"auto_index_check_duration,omitempty"`

	// AutoIndexDurationLimit limit duration of automatic indexing
	AutoIndexDurationLimit *string `json:"auto_index_duration_limit,omitempty"`

	// AutoIndexLength number of cache to trigger automatic indexing
	AutoIndexLength *int `json:"auto_index_length,omitempty"`

	// AutoSaveIndexDurationLimit limit duration of automatic index saving
	AutoSaveIndexDurationLimit *string `json:"auto_save_index_duration_limit,omitempty"`

	// AutoSaveIndexWaitDuration duration of automatic index saving wait duration for next saving
	AutoSaveIndexWaitDuration *string `json:"auto_save_index_wait_duration,omitempty"`

	// Concurrency concurrency
	Concurrency *int `json:"concurrency,omitempty"`

	// CreationPoolSize number of pool size of create index processing
	CreationPoolSize *int                           `json:"creation_pool_size,omitempty"`
	Discoverer       *ManagerIndexIndexerDiscoverer `json:"discoverer,omitempty"`

	// NodeName node name
	NodeName *string `json:"node_name,omitempty"`
}

// ManagerIndexIndexerDiscoverer defines model for manager_index_indexer_discoverer.
type ManagerIndexIndexerDiscoverer struct {
	AgentClientOptions *GrpcClient `json:"agent_client_options,omitempty"`
	Client             *GrpcClient `json:"client,omitempty"`

	// Duration refresh duration to discover
	Duration *string `json:"duration,omitempty"`
}

// ManagerIndexOperator [THIS FEATURE IS WIP] operator that manages vald index
type ManagerIndexOperator struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// Annotations deployment annotations
	Annotations        *map[string]any                         `json:"annotations,omitempty"`
	ClusterRole        *ManagerIndexOperatorClusterRole        `json:"clusterRole,omitempty"`
	ClusterRoleBinding *ManagerIndexOperatorClusterRoleBinding `json:"clusterRoleBinding,omitempty"`

	// Enabled index operator enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env   *Env   `json:"env,omitempty"`
	Image *Image `json:"image,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// Kind deployment kind: Deployment or DaemonSet
	Kind    *ManagerIndexOperatorKind `json:"kind,omitempty"`
	Logging *Logging                  `json:"logging,omitempty"`

	// Name name of manager.index.operator deployment
	Name *string `json:"name,omitempty"`

	// Namespace namespace to discovery
	Namespace *string `json:"namespace,omitempty"`

	// NodeName node name
	NodeName *string `json:"nodeName,omitempty"`

	// NodeSelector node selector
	NodeSelector  *NodeSelector  `json:"nodeSelector,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// PodAnnotations pod annotations
	PodAnnotations *map[string]any `json:"podAnnotations,omitempty"`
	PodPriority    *PodPriority    `json:"podPriority,omitempty"`

	// PodSecurityContext security context for pod
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// ProgressDeadlineSeconds progress deadline seconds
	ProgressDeadlineSeconds *int `json:"progressDeadlineSeconds,omitempty"`

	// Replicas number of replicas.
	Replicas *int `json:"replicas,omitempty"`

	// Resources compute resources
	Resources *Resources `json:"resources,omitempty"`

	// RevisionHistoryLimit number of old history to retain to allow rollback
	RevisionHistoryLimit *int           `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate `json:"rollingUpdate,omitempty"`

	// RotationJobConcurrency maximum concurrent rotator job run.
	RotationJobConcurrency *int `json:"rotation_job_concurrency,omitempty"`

	// SecurityContext security context for container
	SecurityContext    *corev1.SecurityContext `json:"securityContext,omitempty"`
	ServerConfig       *ServerConfig           `json:"server_config,omitempty"`
	ServiceAccountName *string                 `json:"serviceAccountName,omitempty"`

	// TerminationGracePeriodSeconds duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds *int `json:"terminationGracePeriodSeconds,omitempty"`

	// TimeZone Time zone
	TimeZone *string `json:"time_zone,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TopologySpreadConstraints topology spread constraints of gateway pods
	TopologySpreadConstraints *TopologySpreadConstraints `json:"topologySpreadConstraints,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`

	// VolumeMounts volume mounts
	VolumeMounts *VolumeMounts `json:"volumeMounts,omitempty"`

	// Volumes volumes
	Volumes *Volumes `json:"volumes,omitempty"`
}

// ManagerIndexOperatorKind deployment kind: Deployment or DaemonSet
type ManagerIndexOperatorKind string

// ManagerIndexOperatorClusterRole defines model for manager_index_operator_clusterRole.
type ManagerIndexOperatorClusterRole struct {
	// Enabled creates clusterRole resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRole
	Name *string `json:"name,omitempty"`
}

// ManagerIndexOperatorClusterRoleBinding defines model for manager_index_operator_clusterRoleBinding.
type ManagerIndexOperatorClusterRoleBinding struct {
	// Enabled creates clusterRoleBinding resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRoleBinding
	Name *string `json:"name,omitempty"`
}

// ManagerIndexReadreplica defines model for manager_index_readreplica.
type ManagerIndexReadreplica struct {
	// Rotator [This feature is work in progress] readreplica agents rotation job
	Rotator *ManagerIndexReadreplicaRotator `json:"rotator,omitempty"`
}

// ManagerIndexReadreplicaRotator [This feature is work in progress] readreplica agents rotation job
type ManagerIndexReadreplicaRotator struct {
	// AgentNamespace namespace of agent pods to manage
	AgentNamespace     *string                                           `json:"agent_namespace,omitempty"`
	ClusterRole        *ManagerIndexReadreplicaRotatorClusterRole        `json:"clusterRole,omitempty"`
	ClusterRoleBinding *ManagerIndexReadreplicaRotatorClusterRoleBinding `json:"clusterRoleBinding,omitempty"`

	// Env environment variables
	Env   *Env   `json:"env,omitempty"`
	Image *Image `json:"image,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// Name name of readreplica rotator job
	Name          *string        `json:"name,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// PodSecurityContext security context for pod
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// SecurityContext security context for container
	SecurityContext    *corev1.SecurityContext `json:"securityContext,omitempty"`
	ServerConfig       *ServerConfig           `json:"server_config,omitempty"`
	ServiceAccountName *string                 `json:"serviceAccountName,omitempty"`

	// TargetReadReplicaIdAnnotationsKey name of annotations key for target read replica id
	TargetReadReplicaIdAnnotationsKey *string `json:"target_read_replica_id_annotations_key,omitempty"`

	// TtlSecondsAfterFinished ttl setting for K8s completed jobs
	TtlSecondsAfterFinished *int `json:"ttlSecondsAfterFinished,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`
}

// ManagerIndexReadreplicaRotatorClusterRole defines model for manager_index_readreplica_rotator_clusterRole.
type ManagerIndexReadreplicaRotatorClusterRole struct {
	// Enabled creates clusterRole resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRole
	Name *string `json:"name,omitempty"`
}

// ManagerIndexReadreplicaRotatorClusterRoleBinding defines model for manager_index_readreplica_rotator_clusterRoleBinding.
type ManagerIndexReadreplicaRotatorClusterRoleBinding struct {
	// Enabled creates clusterRoleBinding resource
	Enabled *bool `json:"enabled,omitempty"`

	// Name name of clusterRoleBinding
	Name *string `json:"name,omitempty"`
}

// ManagerIndexSaver defines model for manager_index_saver.
type ManagerIndexSaver struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// AgentNamespace namespace of agent pods to manage
	AgentNamespace *string `json:"agent_namespace,omitempty"`

	// Concurrency concurrency for index saving
	Concurrency *int                         `json:"concurrency,omitempty"`
	Discoverer  *ManagerIndexSaverDiscoverer `json:"discoverer,omitempty"`

	// Enabled enable index save CronJob
	Enabled *bool `json:"enabled,omitempty"`

	// Env environment variables
	Env   *Env   `json:"env,omitempty"`
	Image *Image `json:"image,omitempty"`

	// InitContainers init containers
	InitContainers *InitContainers `json:"initContainers,omitempty"`

	// Name name of index save job
	Name *string `json:"name,omitempty"`

	// NodeSelector node selector
	NodeSelector *NodeSelector `json:"nodeSelector,omitempty"`

	// NodeName node name
	NodeName      *string        `json:"node_name,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// Schedule CronJob schedule setting for index save
	Schedule           *string       `json:"schedule,omitempty"`
	ServerConfig       *ServerConfig `json:"server_config,omitempty"`
	ServiceAccountName *string       `json:"serviceAccountName,omitempty"`

	// StartingDeadlineSeconds startingDeadlineSeconds setting for K8s completed jobs
	StartingDeadlineSeconds *int `json:"startingDeadlineSeconds,omitempty"`

	// Suspend CronJob suspend setting for index creation
	Suspend *bool `json:"suspend,omitempty"`

	// TargetAddrs index saving target addresses
	TargetAddrs *[]string `json:"target_addrs,omitempty"`

	// Tolerations tolerations
	Tolerations *Tolerations `json:"tolerations,omitempty"`

	// TtlSecondsAfterFinished ttl setting for K8s completed jobs
	TtlSecondsAfterFinished *int `json:"ttlSecondsAfterFinished,omitempty"`

	// Version version of gateway config
	Version *Version `json:"version,omitempty"`
}

// ManagerIndexSaverDiscoverer defines model for manager_index_saver_discoverer.
type ManagerIndexSaverDiscoverer struct {
	AgentClientOptions *GrpcClient `json:"agent_client_options,omitempty"`
	Client             *GrpcClient `json:"client,omitempty"`

	// Duration refresh duration to discover
	Duration *string `json:"duration,omitempty"`
}

// Net defines model for net.
type Net = config.Net

// NetworkPolicy defines model for networkPolicy.
type NetworkPolicy struct {
	// Custom custom network policies that a user can add
	Custom *DefaultsNetworkPolicyCustom `json:"custom,omitempty"`

	// Enabled if network policy enabled
	Enabled *bool `json:"enabled,omitempty"`
}

// NodeSelector node selector
type NodeSelector = map[string]string

// Observability defines model for observability.
type Observability = config.Observability

// PodPriority defines model for podPriority.
type PodPriority struct {
	// Enabled gateway pod PriorityClass enabled
	Enabled *bool `json:"enabled,omitempty"`

	// Value gateway pod PriorityClass value
	Value *int `json:"value,omitempty"`
}

// Resources compute resources
type Resources = corev1.ResourceRequirements

// RestServer defines model for rest_server.
type RestServer = config.Server

// RollingUpdate defines model for rollingUpdate.
type RollingUpdate struct {
	// MaxSurge max surge of rolling update
	MaxSurge *string `json:"maxSurge,omitempty"`

	// MaxUnavailable max unavailable of rolling update
	MaxUnavailable *string `json:"maxUnavailable,omitempty"`
}

// ServerConfig defines model for server_config.
type ServerConfig struct {
	// FullShutdownDuration server full shutdown duration
	FullShutdownDuration *string                      `json:"full_shutdown_duration,omitempty"`
	Healths              *DefaultsServerConfigHealths `json:"healths,omitempty"`
	Metrics              *DefaultsServerConfigMetrics `json:"metrics,omitempty"`
	Servers              *DefaultsServerConfigServers `json:"servers,omitempty"`
	Tls                  *Tls                         `json:"tls,omitempty"`
}

// Service defines model for service.
type Service struct {
	// Annotations service annotations
	Annotations *map[string]any `json:"annotations,omitempty"`

	// Labels service labels
	Labels *map[string]any `json:"labels,omitempty"`
}

// SocketOption defines model for socket_option.
type SocketOption = config.SocketOption

// Tls defines model for tls.
type Tls = config.TLS

// Tolerations tolerations
type Tolerations = []corev1.Toleration

// TopologySpreadConstraints topology spread constraints of gateway pods
type TopologySpreadConstraints = []corev1.TopologySpreadConstraint

// Version version of gateway config
type Version = string

// VolumeMounts volume mounts
type VolumeMounts = []corev1.VolumeMount

// Volumes volumes
type Volumes = []corev1.Volume

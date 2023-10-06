// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package vald

import (
	"github.com/vdaas/vald/internal/config"
)

// Affinity
type Affinity struct {
	NodeAffinity    *NodeAffinity    `json:"nodeAffinity,omitempty"`
	PodAffinity     *PodAffinity     `json:"podAffinity,omitempty"`
	PodAntiAffinity *PodAntiAffinity `json:"podAntiAffinity,omitempty"`
}

// Agent
type Agent struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// deployment annotations
	Annotations *Annotations `json:"annotations,omitempty"`

	// agent enabled
	Enabled bool `json:"enabled,omitempty"`

	// environment variables
	Env []*EnvItems `json:"env,omitempty"`

	// external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy string `json:"externalTrafficPolicy,omitempty"`
	Hpa                   *Hpa   `json:"hpa,omitempty"`
	Image                 *Image `json:"image,omitempty"`

	// init containers
	InitContainers []*InitContainersItems `json:"initContainers,omitempty"`

	// deployment kind: Deployment, DaemonSet or StatefulSet
	Kind    string   `json:"kind,omitempty"`
	Logging *Logging `json:"logging,omitempty"`

	// maximum number of replicas. if HPA is disabled, this value will be ignored.
	MaxReplicas int `json:"maxReplicas,omitempty"`

	// maximum number of unavailable replicas
	MaxUnavailable string `json:"maxUnavailable,omitempty"`

	// minimum number of replicas. if HPA is disabled, the replicas will be set to this value
	MinReplicas int `json:"minReplicas,omitempty"`

	// name of agent deployment
	Name string      `json:"name,omitempty"`
	Ngt  *config.NGT `json:"ngt,omitempty"`

	// node name
	NodeName string `json:"nodeName,omitempty"`

	// node selector
	NodeSelector     *NodeSelector     `json:"nodeSelector,omitempty"`
	Observability    *Observability    `json:"observability,omitempty"`
	PersistentVolume *PersistentVolume `json:"persistentVolume,omitempty"`

	// pod annotations
	PodAnnotations *PodAnnotations `json:"podAnnotations,omitempty"`

	// pod management policy: OrderedReady or Parallel
	PodManagementPolicy string       `json:"podManagementPolicy,omitempty"`
	PodPriority         *PodPriority `json:"podPriority,omitempty"`

	// security context for pod
	PodSecurityContext *PodSecurityContext `json:"podSecurityContext,omitempty"`

	// progress deadline seconds
	ProgressDeadlineSeconds int `json:"progressDeadlineSeconds,omitempty"`

	// compute resources
	Resources *Resources `json:"resources,omitempty"`

	// number of old history to retain to allow rollback
	RevisionHistoryLimit int            `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate `json:"rollingUpdate,omitempty"`

	// security context for container
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	ServerConfig    *ServerConfig    `json:"server_config,omitempty"`
	Service         *Service         `json:"service,omitempty"`

	// service type: ClusterIP, LoadBalancer or NodePort
	ServiceType string   `json:"serviceType,omitempty"`
	Sidecar     *Sidecar `json:"sidecar,omitempty"`

	// duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds int `json:"terminationGracePeriodSeconds,omitempty"`

	// Time zone
	TimeZone string `json:"time_zone,omitempty"`

	// tolerations
	Tolerations []*TolerationsItems `json:"tolerations,omitempty"`

	// topology spread constraints of gateway pods
	TopologySpreadConstraints []*TopologySpreadConstraintsItems `json:"topologySpreadConstraints,omitempty"`

	// version of gateway config
	Version string `json:"version,omitempty"`

	// volume mounts
	VolumeMounts []*VolumeMountsItems `json:"volumeMounts,omitempty"`

	// volumes
	Volumes []*VolumesItems `json:"volumes,omitempty"`
}

// Annotations deployment annotations
type Annotations map[string]string

// BlobStorage
type BlobStorage struct {
	// bucket name
	Bucket       string        `json:"bucket,omitempty"`
	CloudStorage *CloudStorage `json:"cloud_storage,omitempty"`
	S3           *S3           `json:"s3,omitempty"`

	// storage type
	StorageType string `json:"storage_type,omitempty"`
}

// CloudStorage
type CloudStorage struct {
	Client *config.GRPCClient `json:"client,omitempty"`

	// cloud storage url
	Url string `json:"url,omitempty"`

	// bytes of the chunks for upload
	WriteBufferSize int `json:"write_buffer_size,omitempty"`

	// Cache-Control of HTTP Header
	WriteCacheControl string `json:"write_cache_control,omitempty"`

	// Content-Disposition of HTTP Header
	WriteContentDisposition string `json:"write_content_disposition,omitempty"`

	// the encoding of the blob's content
	WriteContentEncoding string `json:"write_content_encoding,omitempty"`

	// the language of blob's content
	WriteContentLanguage string `json:"write_content_language,omitempty"`

	// MIME type of the blob
	WriteContentType string `json:"write_content_type,omitempty"`
}

// ClusterRole
type ClusterRole struct {
	// creates clusterRole resource
	Enabled bool `json:"enabled,omitempty"`

	// name of clusterRole
	Name string `json:"name,omitempty"`
}

// ClusterRoleBinding
type ClusterRoleBinding struct {
	// creates clusterRoleBinding resource
	Enabled bool `json:"enabled,omitempty"`

	// name of clusterRoleBinding
	Name string `json:"name,omitempty"`
}

// Collector
type Collector struct {
	// metrics collect duration. if it is set as 5s, enabled metrics are collected every 5 seconds.
	Duration string   `json:"duration,omitempty"`
	Metrics  *Metrics `json:"metrics,omitempty"`
}

// Compress
type Compress struct {
	// compression algorithm. must be `gob`, `gzip`, `lz4` or `zstd`
	CompressAlgorithm string `json:"compress_algorithm,omitempty"`

	// compression level. value range relies on which algorithm is used. `gob`: level will be ignored. `gzip`: -1 (default compression), 0 (no compression), or 1 (best speed) to 9 (best compression). `lz4`: >= 0, higher is better compression. `zstd`: 1 (fastest) to 22 (best), however implementation relies on klauspost/compress.
	CompressionLevel int `json:"compression_level,omitempty"`
}

// Config
type Config struct {
	// auto backup duration
	AutoBackupDuration string `json:"auto_backup_duration,omitempty"`

	// auto backup triggered by timer is enabled
	AutoBackupEnabled bool               `json:"auto_backup_enabled,omitempty"`
	BlobStorage       *BlobStorage       `json:"blob_storage,omitempty"`
	Client            *config.GRPCClient `json:"client,omitempty"`
	Compress          *Compress          `json:"compress,omitempty"`

	// backup filename
	Filename string `json:"filename,omitempty"`

	// suffix for backup filename
	FilenameSuffix string `json:"filename_suffix,omitempty"`

	// timeout for observing file changes during post stop
	PostStopTimeout string          `json:"post_stop_timeout,omitempty"`
	RestoreBackoff  *config.Backoff `json:"restore_backoff,omitempty"`

	// restore backoff enabled
	RestoreBackoffEnabled bool `json:"restore_backoff_enabled,omitempty"`

	// auto backup triggered by file changes is enabled
	WatchEnabled bool `json:"watch_enabled,omitempty"`
}

// Defaults
type Defaults struct {
	Grpc          *Grpc          `json:"grpc,omitempty"`
	Image         *Image         `json:"image,omitempty"`
	Ingress       *Ingress       `json:"ingress,omitempty"`
	Logging       *Logging       `json:"logging,omitempty"`
	Observability *Observability `json:"observability,omitempty"`
	ServerConfig  *ServerConfig  `json:"server_config,omitempty"`

	// Time zone
	TimeZone string `json:"time_zone,omitempty"`
}

// Dialer
type Dialer struct {
	// gRPC client TCP dialer dual stack enabled
	DualStackEnabled bool `json:"dual_stack_enabled,omitempty"`

	// gRPC client TCP dialer keep alive
	Keepalive string `json:"keepalive,omitempty"`

	// gRPC client TCP dialer timeout
	Timeout string `json:"timeout,omitempty"`
}

// Discoverer
type Discoverer struct {
	AgentClientOptions *config.GRPCClient `json:"agent_client_options,omitempty"`
	Client             *config.GRPCClient `json:"client,omitempty"`

	// refresh duration to discover
	Duration string `json:"duration,omitempty"`
}

// Dns
type Dns struct {
	// gRPC client TCP DNS cache enabled
	CacheEnabled bool `json:"cache_enabled,omitempty"`

	// gRPC client TCP DNS cache expiration
	CacheExpiration string `json:"cache_expiration,omitempty"`

	// gRPC client TCP DNS cache refresh duration
	RefreshDuration string `json:"refresh_duration,omitempty"`
}

// EgressFilter gRPC client config for egress filter
type EgressFilter struct {
	Client *config.GRPCClient `json:"client,omitempty"`

	// distance egress vector filter targets
	DistanceFilters []string `json:"distance_filters,omitempty"`

	// object egress vector filter targets
	ObjectFilters []string `json:"object_filters,omitempty"`
}

// EnvItems
type EnvItems struct{}

// Fields k8s field selectors for pod discovery
type Fields struct{}

// Filter
type Filter struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// deployment annotations
	Annotations *Annotations `json:"annotations,omitempty"`

	// gateway enabled
	Enabled bool `json:"enabled,omitempty"`

	// environment variables
	Env []*EnvItems `json:"env,omitempty"`

	// external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy string         `json:"externalTrafficPolicy,omitempty"`
	GatewayConfig         *GatewayConfig `json:"gateway_config,omitempty"`
	Hpa                   *Hpa           `json:"hpa,omitempty"`
	Image                 *Image         `json:"image,omitempty"`
	Ingress               *Ingress       `json:"ingress,omitempty"`

	// init containers
	InitContainers []*InitContainersItems `json:"initContainers,omitempty"`

	// deployment kind: Deployment or DaemonSet
	Kind    string   `json:"kind,omitempty"`
	Logging *Logging `json:"logging,omitempty"`

	// maximum number of replicas. if HPA is disabled, this value will be ignored.
	MaxReplicas int `json:"maxReplicas,omitempty"`

	// maximum number of unavailable replicas
	MaxUnavailable string `json:"maxUnavailable,omitempty"`

	// minimum number of replicas. if HPA is disabled, the replicas will be set to this value
	MinReplicas int `json:"minReplicas,omitempty"`

	// name of filter gateway deployment
	Name string `json:"name,omitempty"`

	// node name
	NodeName string `json:"nodeName,omitempty"`

	// node selector
	NodeSelector  *NodeSelector  `json:"nodeSelector,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// pod annotations
	PodAnnotations *PodAnnotations `json:"podAnnotations,omitempty"`
	PodPriority    *PodPriority    `json:"podPriority,omitempty"`

	// security context for pod
	PodSecurityContext *PodSecurityContext `json:"podSecurityContext,omitempty"`

	// progress deadline seconds
	ProgressDeadlineSeconds int `json:"progressDeadlineSeconds,omitempty"`

	// compute resources
	Resources *Resources `json:"resources,omitempty"`

	// number of old history to retain to allow rollback
	RevisionHistoryLimit int            `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate `json:"rollingUpdate,omitempty"`

	// security context for container
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	ServerConfig    *ServerConfig    `json:"server_config,omitempty"`
	Service         *Service         `json:"service,omitempty"`

	// service type: ClusterIP, LoadBalancer or NodePort
	ServiceType string `json:"serviceType,omitempty"`

	// duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds int `json:"terminationGracePeriodSeconds,omitempty"`

	// Time zone
	TimeZone string `json:"time_zone,omitempty"`

	// tolerations
	Tolerations []*TolerationsItems `json:"tolerations,omitempty"`

	// topology spread constraints of gateway pods
	TopologySpreadConstraints []*TopologySpreadConstraintsItems `json:"topologySpreadConstraints,omitempty"`

	// version of gateway config
	Version string `json:"version,omitempty"`

	// volume mounts
	VolumeMounts []*VolumeMountsItems `json:"volumeMounts,omitempty"`

	// volumes
	Volumes []*VolumesItems `json:"volumes,omitempty"`
}

// Gateway
type Gateway struct {
	Filter *Filter `json:"filter,omitempty"`
	Lb     *Lb     `json:"lb,omitempty"`
}

// GatewayConfig
type GatewayConfig struct {
	// gRPC client config for egress filter
	EgressFilter  *EgressFilter      `json:"egress_filter,omitempty"`
	GatewayClient *config.GRPCClient `json:"gateway_client,omitempty"`

	// gRPC client config for ingress filter
	IngressFilter *IngressFilter `json:"ingress_filter,omitempty"`
}

// Grpc
type Grpc struct {
	// gRPC server enabled
	Enabled bool `json:"enabled,omitempty"`

	// gRPC server host
	Host string `json:"host,omitempty"`

	// gRPC server port
	Port   int     `json:"port,omitempty"`
	Server *Server `json:"server,omitempty"`

	// gRPC server service port
	ServicePort int `json:"servicePort,omitempty"`
}

// Healths
type Healths struct {
	Liveness  *Liveness  `json:"liveness,omitempty"`
	Readiness *Readiness `json:"readiness,omitempty"`
	Startup   *Startup   `json:"startup,omitempty"`
}

// Hpa
type Hpa struct {
	// HPA enabled
	Enabled bool `json:"enabled,omitempty"`

	// HPA CPU utilization percentage
	TargetCPUUtilizationPercentage int `json:"targetCPUUtilizationPercentage,omitempty"`
}

// Http
type Http struct {
	// REST server handler timeout
	HandlerTimeout string `json:"handler_timeout,omitempty"`

	// REST server idle timeout
	IdleTimeout string `json:"idle_timeout,omitempty"`

	// REST server read header timeout
	ReadHeaderTimeout string `json:"read_header_timeout,omitempty"`

	// REST server read timeout
	ReadTimeout string `json:"read_timeout,omitempty"`

	// REST server shutdown duration
	ShutdownDuration string `json:"shutdown_duration,omitempty"`

	// REST server write timeout
	WriteTimeout string `json:"write_timeout,omitempty"`
}

// HttpGet
type HttpGet struct {
	// startup probe path
	Path string `json:"path,omitempty"`

	// startup probe port
	Port string `json:"port,omitempty"`

	// startup probe scheme
	Scheme string `json:"scheme,omitempty"`
}

// Image
type Image struct {
	// image pull policy
	PullPolicy string `json:"pullPolicy,omitempty"`

	// image repository
	Repository string `json:"repository,omitempty"`

	// image tag (overrides defaults.image.tag)
	Tag string `json:"tag,omitempty"`
}

// Index
type Index struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// deployment annotations
	Annotations *Annotations `json:"annotations,omitempty"`

	// index manager enabled
	Enabled bool `json:"enabled,omitempty"`

	// environment variables
	Env []*EnvItems `json:"env,omitempty"`

	// external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy string   `json:"externalTrafficPolicy,omitempty"`
	Image                 *Image   `json:"image,omitempty"`
	Indexer               *Indexer `json:"indexer,omitempty"`

	// init containers
	InitContainers []*InitContainersItems `json:"initContainers,omitempty"`

	// deployment kind: Deployment or DaemonSet
	Kind    string   `json:"kind,omitempty"`
	Logging *Logging `json:"logging,omitempty"`

	// maximum number of unavailable replicas
	MaxUnavailable string `json:"maxUnavailable,omitempty"`

	// name of index manager deployment
	Name string `json:"name,omitempty"`

	// node name
	NodeName string `json:"nodeName,omitempty"`

	// node selector
	NodeSelector  *NodeSelector  `json:"nodeSelector,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// pod annotations
	PodAnnotations *PodAnnotations `json:"podAnnotations,omitempty"`
	PodPriority    *PodPriority    `json:"podPriority,omitempty"`

	// security context for pod
	PodSecurityContext *PodSecurityContext `json:"podSecurityContext,omitempty"`

	// progress deadline seconds
	ProgressDeadlineSeconds int `json:"progressDeadlineSeconds,omitempty"`

	// number of replicas
	Replicas int `json:"replicas,omitempty"`

	// compute resources
	Resources *Resources `json:"resources,omitempty"`

	// number of old history to retain to allow rollback
	RevisionHistoryLimit int            `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate `json:"rollingUpdate,omitempty"`

	// security context for container
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	ServerConfig    *ServerConfig    `json:"server_config,omitempty"`
	Service         *Service         `json:"service,omitempty"`

	// service type: ClusterIP, LoadBalancer or NodePort
	ServiceType string `json:"serviceType,omitempty"`

	// duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds int `json:"terminationGracePeriodSeconds,omitempty"`

	// Time zone
	TimeZone string `json:"time_zone,omitempty"`

	// tolerations
	Tolerations []*TolerationsItems `json:"tolerations,omitempty"`

	// topology spread constraints of gateway pods
	TopologySpreadConstraints []*TopologySpreadConstraintsItems `json:"topologySpreadConstraints,omitempty"`

	// version of gateway config
	Version string `json:"version,omitempty"`

	// volume mounts
	VolumeMounts []*VolumeMountsItems `json:"volumeMounts,omitempty"`

	// volumes
	Volumes []*VolumesItems `json:"volumes,omitempty"`
}

// Indexer
type Indexer struct {
	// namespace of agent pods to manage
	AgentNamespace string `json:"agent_namespace,omitempty"`

	// check duration of automatic indexing
	AutoIndexCheckDuration string `json:"auto_index_check_duration,omitempty"`

	// limit duration of automatic indexing
	AutoIndexDurationLimit string `json:"auto_index_duration_limit,omitempty"`

	// number of cache to trigger automatic indexing
	AutoIndexLength int `json:"auto_index_length,omitempty"`

	// limit duration of automatic index saving
	AutoSaveIndexDurationLimit string `json:"auto_save_index_duration_limit,omitempty"`

	// duration of automatic index saving wait duration for next saving
	AutoSaveIndexWaitDuration string `json:"auto_save_index_wait_duration,omitempty"`

	// concurrency
	Concurrency int `json:"concurrency,omitempty"`

	// number of pool size of create index processing
	CreationPoolSize int         `json:"creation_pool_size,omitempty"`
	Discoverer       *Discoverer `json:"discoverer,omitempty"`

	// node name
	NodeName string `json:"node_name,omitempty"`
}

// Ingress
type Ingress struct {
	// annotations for ingress
	Annotations *Annotations `json:"annotations,omitempty"`

	// gateway ingress enabled
	Enabled bool `json:"enabled,omitempty"`

	// ingress hostname
	Host string `json:"host,omitempty"`

	// gateway ingress pathType
	PathType string `json:"pathType,omitempty"`

	// service port to be exposed by ingress
	ServicePort string `json:"servicePort,omitempty"`
}

// IngressFilter gRPC client config for ingress filter
type IngressFilter struct {
	Client *config.GRPCClient `json:"client,omitempty"`

	// insert ingress vector filter targets
	InsertFilters []string `json:"insert_filters,omitempty"`

	// search ingress vector filter targets
	SearchFilters []string `json:"search_filters,omitempty"`

	// update ingress vector filter targets
	UpdateFilters []string `json:"update_filters,omitempty"`

	// upsert ingress vector filter targets
	UpsertFilters []string `json:"upsert_filters,omitempty"`

	// object ingress vectorize filter targets
	Vectorizer string `json:"vectorizer,omitempty"`
}

// InitContainersItems
type InitContainersItems struct {
	Type          string
	Name          string
	Target        string
	Image         string
	SleepDuration int
}

// Initializer
type Initializer struct{}

// Jaeger
type Jaeger struct {
	// Jaeger agent endpoint
	AgentEndpoint string `json:"agent_endpoint,omitempty"`

	// Jaeger agent reconnect interval
	AgentReconnectInterval string `json:"agent_reconnect_interval,omitempty"`

	// Jaeger Agent max packet size
	AgentMaxPacketSize int `json:"agent_max_packet_size,omitempty"`

	// Jaeger collector endpoint
	CollectorEndpoint string `json:"collector_endpoint,omitempty"`

	// Jaeger exporter enabled
	Enabled bool `json:"enabled,omitempty"`

	// Jaeger password
	Password string `json:"password,omitempty"`

	// Jaeger service name
	ServiceName string `json:"service_name,omitempty"`

	// Jaeger username
	Username string `json:"username,omitempty"`

	// Jaeger export batch timeout
	BatchTimeout string `json:"batch_timeout,omitempty"`

	// Jaeger export timeout
	ExportTimeout string `json:"export_timeout,omitempty"`

	// Jaeger max export batch size
	MaxExportBatchSize int `json:"max_export_batch_size,omitempty"`

	// Jaeger max queue size
	MaxQueueSize int `json:"max_queue_size,omitempty"`
}

// Keepalive
type Keepalive struct {
	// gRPC client keep alive permit without stream
	PermitWithoutStream bool `json:"permit_without_stream,omitempty"`

	// gRPC client keep alive time
	Time string `json:"time,omitempty"`

	// gRPC client keep alive timeout
	Timeout string `json:"timeout,omitempty"`
}

// Kvsdb
type Kvsdb struct {
	// kvsdb processing concurrency
	Concurrency int `json:"concurrency,omitempty"`
}

// Labels service labels
type Labels map[string]string

// Lb
type Lb struct {
	Affinity *Affinity `json:"affinity,omitempty"`

	// deployment annotations
	Annotations *Annotations `json:"annotations,omitempty"`

	// gateway enabled
	Enabled bool `json:"enabled,omitempty"`

	// environment variables
	Env []*EnvItems `json:"env,omitempty"`

	// external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local
	ExternalTrafficPolicy string         `json:"externalTrafficPolicy,omitempty"`
	GatewayConfig         *GatewayConfig `json:"gateway_config,omitempty"`
	Hpa                   *Hpa           `json:"hpa,omitempty"`
	Image                 *Image         `json:"image,omitempty"`
	Ingress               *Ingress       `json:"ingress,omitempty"`

	// init containers
	InitContainers []*InitContainersItems `json:"initContainers,omitempty"`

	// deployment kind: Deployment or DaemonSet
	Kind    string   `json:"kind,omitempty"`
	Logging *Logging `json:"logging,omitempty"`

	// maximum number of replicas. if HPA is disabled, this value will be ignored.
	MaxReplicas int `json:"maxReplicas,omitempty"`

	// maximum number of unavailable replicas
	MaxUnavailable string `json:"maxUnavailable,omitempty"`

	// minimum number of replicas. if HPA is disabled, the replicas will be set to this value
	MinReplicas int `json:"minReplicas,omitempty"`

	// name of gateway deployment
	Name string `json:"name,omitempty"`

	// node name
	NodeName string `json:"nodeName,omitempty"`

	// node selector
	NodeSelector  *NodeSelector  `json:"nodeSelector,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// pod annotations
	PodAnnotations *PodAnnotations `json:"podAnnotations,omitempty"`
	PodPriority    *PodPriority    `json:"podPriority,omitempty"`

	// security context for pod
	PodSecurityContext *PodSecurityContext `json:"podSecurityContext,omitempty"`

	// progress deadline seconds
	ProgressDeadlineSeconds int `json:"progressDeadlineSeconds,omitempty"`

	// compute resources
	Resources *Resources `json:"resources,omitempty"`

	// number of old history to retain to allow rollback
	RevisionHistoryLimit int            `json:"revisionHistoryLimit,omitempty"`
	RollingUpdate        *RollingUpdate `json:"rollingUpdate,omitempty"`

	// security context for container
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	ServerConfig    *ServerConfig    `json:"server_config,omitempty"`
	Service         *Service         `json:"service,omitempty"`

	// service type: ClusterIP, LoadBalancer or NodePort
	ServiceType string `json:"serviceType,omitempty"`

	// duration in seconds pod needs to terminate gracefully
	TerminationGracePeriodSeconds int `json:"terminationGracePeriodSeconds,omitempty"`

	// Time zone
	TimeZone string `json:"time_zone,omitempty"`

	// tolerations
	Tolerations []*TolerationsItems `json:"tolerations,omitempty"`

	// topology spread constraints of gateway pods
	TopologySpreadConstraints []*TopologySpreadConstraintsItems `json:"topologySpreadConstraints,omitempty"`

	// version of gateway config
	Version string `json:"version,omitempty"`

	// volume mounts
	VolumeMounts []*VolumeMountsItems `json:"volumeMounts,omitempty"`

	// volumes
	Volumes []*VolumesItems `json:"volumes,omitempty"`
}

// Limits
type Limits struct{}

// Liveness
type Liveness struct {
	// liveness server enabled
	Enabled bool `json:"enabled,omitempty"`

	// liveness server host
	Host          string         `json:"host,omitempty"`
	LivenessProbe *LivenessProbe `json:"livenessProbe,omitempty"`

	// liveness server port
	Port   int     `json:"port,omitempty"`
	Server *Server `json:"server,omitempty"`

	// liveness server service port
	ServicePort int `json:"servicePort,omitempty"`
}

// LivenessProbe
type LivenessProbe struct {
	// liveness probe failure threshold
	FailureThreshold int      `json:"failureThreshold,omitempty"`
	HttpGet          *HttpGet `json:"httpGet,omitempty"`

	// liveness probe initial delay seconds
	InitialDelaySeconds int `json:"initialDelaySeconds,omitempty"`

	// liveness probe period seconds
	PeriodSeconds int `json:"periodSeconds,omitempty"`

	// liveness probe success threshold
	SuccessThreshold int `json:"successThreshold,omitempty"`

	// liveness probe timeout seconds
	TimeoutSeconds int `json:"timeoutSeconds,omitempty"`
}

// Logging
type Logging struct {
	// logging format. logging format must be `raw` or `json`
	Format string `json:"format,omitempty"`

	// logging level. logging level must be `debug`, `info`, `warn`, `error` or `fatal`.
	Level string `json:"level,omitempty"`

	// logger name. currently logger must be `glg` or `zap`.
	Logger string `json:"logger,omitempty"`
}

// Manager
type Manager struct {
	Index *Index `json:"index,omitempty"`
}

// Metrics
type Metrics struct {
	// CGO metrics enabled
	EnableCgo bool `json:"enable_cgo,omitempty"`

	// goroutine metrics enabled
	EnableGoroutine bool `json:"enable_goroutine,omitempty"`

	// memory metrics enabled
	EnableMemory bool `json:"enable_memory,omitempty"`

	// version info metrics enabled
	EnableVersionInfo bool `json:"enable_version_info,omitempty"`

	// enabled label names of version info
	VersionInfoLabels []string `json:"version_info_labels,omitempty"`
}

// Node k8s resource selectors for node discovery
type Node struct {
	// k8s field selectors for node discovery
	Fields *Fields `json:"fields,omitempty"`
	// k8s label selectors for node discovery
	Labels *Labels `json:"labels,omitempty"`
}

// NodeAffinity
type NodeAffinity struct {
	// node affinity preferred scheduling terms
	PreferredDuringSchedulingIgnoredDuringExecution []*PreferredDuringSchedulingIgnoredDuringExecutionItems `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
	RequiredDuringSchedulingIgnoredDuringExecution  *RequiredDuringSchedulingIgnoredDuringExecution         `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// NodeMetrics k8s resource selectors for node_metrics discovery
type NodeMetrics struct {
	// k8s field selectors for node_metrics discovery
	Fields *Fields `json:"fields,omitempty"`

	// k8s label selectors for node_metrics discovery
	Labels *Labels `json:"labels,omitempty"`
}

// NodeSelector node selector
type NodeSelector struct{}

// NodeSelectorTermsItems
type NodeSelectorTermsItems struct{}

// Observability
type Observability struct {
	Metrics *Metrics `json:"metrics,omitempty"`

	// observability features enabled
	Enabled    bool        `json:"enabled,omitempty"`
	Jaeger     *Jaeger     `json:"jaeger,omitempty"`
	Prometheus *Prometheus `json:"prometheus,omitempty"`
	Trace      *Trace      `json:"trace,omitempty"`
}

// PersistentVolume
type PersistentVolume struct {
	// agent pod storage accessMode
	AccessMode string `json:"accessMode,omitempty"`

	// enables PVC. It is required to enable if agent pod's file store functionality is enabled with non in-memory mode
	Enabled bool `json:"enabled,omitempty"`

	// size of agent pod volume
	Size string `json:"size,omitempty"`

	// storageClass name for agent pod volume
	StorageClass string `json:"storageClass,omitempty"`
}

// Pod k8s resource selectors for pod discovery
type Pod struct {
	// k8s field selectors for pod discovery
	Fields *Fields `json:"fields,omitempty"`

	// k8s label selectors for pod discovery
	Labels *Labels `json:"labels,omitempty"`
}

// PodAffinity
type PodAffinity struct {
	// pod affinity preferred scheduling terms
	PreferredDuringSchedulingIgnoredDuringExecution []*PreferredDuringSchedulingIgnoredDuringExecutionItems `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`

	// pod affinity required scheduling terms
	RequiredDuringSchedulingIgnoredDuringExecution []*RequiredDuringSchedulingIgnoredDuringExecutionItems `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// PodAnnotations pod annotations
type PodAnnotations struct{}

// PodAntiAffinity
type PodAntiAffinity struct {
	// pod anti-affinity preferred scheduling terms
	PreferredDuringSchedulingIgnoredDuringExecution []*PreferredDuringSchedulingIgnoredDuringExecutionItems `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`

	// pod anti-affinity required scheduling terms
	RequiredDuringSchedulingIgnoredDuringExecution []*RequiredDuringSchedulingIgnoredDuringExecutionItems `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// PodMetrics k8s resource selectors for pod_metrics discovery
type PodMetrics struct {
	// k8s field selectors for pod_metrics discovery
	Fields *Fields `json:"fields,omitempty"`

	// k8s label selectors for pod_metrics discovery
	Labels *Labels `json:"labels,omitempty"`
}

// PodPriority
type PodPriority struct {
	// gateway pod PriorityClass enabled
	Enabled bool `json:"enabled,omitempty"`

	// gateway pod PriorityClass value
	Value int `json:"value,omitempty"`
}

// PodSecurityContext security context for pod
type PodSecurityContext struct{}

// Pprof
type Pprof struct {
	// pprof server enabled
	Enabled bool `json:"enabled,omitempty"`

	// pprof server host
	Host string `json:"host,omitempty"`

	// pprof server port
	Port   int     `json:"port,omitempty"`
	Server *Server `json:"server,omitempty"`

	// pprof server service port
	ServicePort int `json:"servicePort,omitempty"`
}

// PreferredDuringSchedulingIgnoredDuringExecutionItems
type PreferredDuringSchedulingIgnoredDuringExecutionItems struct{}

// Prometheus
type Prometheus struct {
	// Prometheus exporter enabled
	Enabled bool `json:"enabled,omitempty"`

	// Prometheus exporter endpoint
	Endpoint string `json:"endpoint,omitempty"`

	// service namespace for metrics
	Namespace string `json:"namespace,omitempty"`

	// Prometheus collect interval
	CollectInterval string `json:"collect_interval,omitempty"`

	// Prometheus collect timeout
	CollectTimeout string `json:"collect_timeout,omitempty"`

	// Prometheus collect with in memory
	EnableInMemoryMode bool `json:"enable_in_memory_mode,omitempty"`
}

// Readiness
type Readiness struct {
	// readiness server enabled
	Enabled bool `json:"enabled,omitempty"`

	// readiness server host
	Host string `json:"host,omitempty"`

	// readiness server port
	Port           int             `json:"port,omitempty"`
	ReadinessProbe *ReadinessProbe `json:"readinessProbe,omitempty"`
	Server         *Server         `json:"server,omitempty"`

	// readiness server service port
	ServicePort int `json:"servicePort,omitempty"`
}

// ReadinessProbe
type ReadinessProbe struct {
	// readiness probe failure threshold
	FailureThreshold int      `json:"failureThreshold,omitempty"`
	HttpGet          *HttpGet `json:"httpGet,omitempty"`

	// readiness probe initial delay seconds
	InitialDelaySeconds int `json:"initialDelaySeconds,omitempty"`

	// readiness probe period seconds
	PeriodSeconds int `json:"periodSeconds,omitempty"`

	// readiness probe success threshold
	SuccessThreshold int `json:"successThreshold,omitempty"`

	// readiness probe timeout seconds
	TimeoutSeconds int `json:"timeoutSeconds,omitempty"`
}

// Requests
type Requests struct{}

// RequiredDuringSchedulingIgnoredDuringExecution
type RequiredDuringSchedulingIgnoredDuringExecution struct {
	// node affinity required node selectors
	NodeSelectorTerms []*NodeSelectorTermsItems `json:"nodeSelectorTerms,omitempty"`
}

// RequiredDuringSchedulingIgnoredDuringExecutionItems
type RequiredDuringSchedulingIgnoredDuringExecutionItems struct{}

// Resources compute resources
type Resources struct {
	Limits   *Limits   `json:"limits,omitempty"`
	Requests *Requests `json:"requests,omitempty"`
}

// Rest
type Rest struct {
	// REST server enabled
	Enabled bool `json:"enabled,omitempty"`

	// REST server host
	Host string `json:"host,omitempty"`

	// REST server port
	Port   int     `json:"port,omitempty"`
	Server *Server `json:"server,omitempty"`

	// REST server service port
	ServicePort int `json:"servicePort,omitempty"`
}

// RollingUpdate
type RollingUpdate struct {
	// max surge of rolling update
	MaxSurge string `json:"maxSurge,omitempty"`

	// max unavailable of rolling update
	MaxUnavailable string `json:"maxUnavailable,omitempty"`
}

// S3
type S3 struct {
	// s3 access key
	AccessKey string `json:"access_key,omitempty"`

	// enable AWS SDK adding the 'Expect: 100-Continue' header to PUT requests over 2MB of content.
	Enable100Continue bool `json:"enable_100_continue,omitempty"`

	// enable the S3 client to add MD5 checksum to upload API calls.
	EnableContentMd5Validation bool `json:"enable_content_md5_validation,omitempty"`

	// enable endpoint discovery
	EnableEndpointDiscovery bool `json:"enable_endpoint_discovery,omitempty"`

	// enable prefixing request endpoint hosts with modeled information
	EnableEndpointHostPrefix bool `json:"enable_endpoint_host_prefix,omitempty"`

	// enables semantic parameter validation
	EnableParamValidation bool `json:"enable_param_validation,omitempty"`

	// enable ssl for s3 session
	EnableSsl bool `json:"enable_ssl,omitempty"`

	// s3 endpoint
	Endpoint string `json:"endpoint,omitempty"`

	// use path-style addressing
	ForcePathStyle bool `json:"force_path_style,omitempty"`

	// s3 download max chunk size
	MaxChunkSize string `json:"max_chunk_size,omitempty"`

	// s3 multipart upload max part size
	MaxPartSize string `json:"max_part_size,omitempty"`

	// maximum number of retries of s3 client
	MaxRetries int `json:"max_retries,omitempty"`

	// s3 region
	Region string `json:"region,omitempty"`

	// s3 secret access key
	SecretAccessKey string `json:"secret_access_key,omitempty"`

	// s3 token
	Token string `json:"token,omitempty"`

	// enable s3 accelerate feature
	UseAccelerate bool `json:"use_accelerate,omitempty"`

	// s3 service client to use the region specified in the ARN
	UseArnRegion bool `json:"use_arn_region,omitempty"`

	// use dual stack
	UseDualStack bool `json:"use_dual_stack,omitempty"`
}

// SecurityContext security context for container
type SecurityContext struct{}

// Selectors k8s resource selectors
type Selectors struct {
	// k8s resource selectors for node discovery
	Node *Node `json:"node,omitempty"`

	// k8s resource selectors for node_metrics discovery
	NodeMetrics *NodeMetrics `json:"node_metrics,omitempty"`

	// k8s resource selectors for pod discovery
	Pod *Pod `json:"pod,omitempty"`

	// k8s resource selectors for pod_metrics discovery
	PodMetrics *PodMetrics `json:"pod_metrics,omitempty"`
}

// Server
type Server struct {
	Http *Http `json:"http,omitempty"`

	// REST server server mode
	Mode string `json:"mode,omitempty"`

	// mysql network
	Network string `json:"network,omitempty"`

	// REST server probe wait time
	ProbeWaitTime string               `json:"probe_wait_time,omitempty"`
	SocketOption  *config.SocketOption `json:"socket_option,omitempty"`

	// mysql socket_path
	SocketPath string `json:"socket_path,omitempty"`
}

// ServerConfig
type ServerConfig struct {
	// server full shutdown duration
	FullShutdownDuration string   `json:"full_shutdown_duration,omitempty"`
	Healths              *Healths `json:"healths,omitempty"`
	Metrics              *Metrics `json:"metrics,omitempty"`
	Servers              *Servers `json:"servers,omitempty"`
	Tls                  *Tls     `json:"tls,omitempty"`
}

// Servers
type Servers struct {
	Grpc *Grpc `json:"grpc,omitempty"`
	Rest *Rest `json:"rest,omitempty"`
}

// Service
type Service struct {
	// service annotations
	Annotations *Annotations `json:"annotations,omitempty"`

	// service labels
	Labels *Labels `json:"labels,omitempty"`
}

// ServiceAccount
type ServiceAccount struct {
	// creates service account
	Enabled bool `json:"enabled,omitempty"`

	// name of service account
	Name string `json:"name,omitempty"`
}

// Sidecar
type Sidecar struct {
	Config *Config `json:"config,omitempty"`

	// sidecar enabled
	Enabled bool `json:"enabled,omitempty"`

	// environment variables
	Env   []*EnvItems `json:"env,omitempty"`
	Image *Image      `json:"image,omitempty"`

	// sidecar on initContainer mode enabled.
	InitContainerEnabled bool     `json:"initContainerEnabled,omitempty"`
	Logging              *Logging `json:"logging,omitempty"`

	// name of agent sidecar
	Name          string         `json:"name,omitempty"`
	Observability *Observability `json:"observability,omitempty"`

	// compute resources
	Resources    *Resources    `json:"resources,omitempty"`
	ServerConfig *ServerConfig `json:"server_config,omitempty"`
	Service      *Service      `json:"service,omitempty"`

	// Time zone
	TimeZone string `json:"time_zone,omitempty"`

	// version of gateway config
	Version string `json:"version,omitempty"`
}

// Startup
type Startup struct {
	// startup server enabled
	Enabled bool `json:"enabled,omitempty"`

	// startup server port
	Port         int           `json:"port,omitempty"`
	StartupProbe *StartupProbe `json:"startupProbe,omitempty"`
}

// StartupProbe
type StartupProbe struct {
	// startup probe failure threshold
	FailureThreshold int      `json:"failureThreshold,omitempty"`
	HttpGet          *HttpGet `json:"httpGet,omitempty"`

	// startup probe initial delay seconds
	InitialDelaySeconds int `json:"initialDelaySeconds,omitempty"`

	// startup probe period seconds
	PeriodSeconds int `json:"periodSeconds,omitempty"`

	// startup probe success threshold
	SuccessThreshold int `json:"successThreshold,omitempty"`

	// startup probe timeout seconds
	TimeoutSeconds int `json:"timeoutSeconds,omitempty"`
}

// Tls
type Tls struct {
	// TLS ca path
	Ca string `json:"ca,omitempty"`

	// TLS cert path
	Cert string `json:"cert,omitempty"`

	// TLS enabled
	Enabled bool `json:"enabled,omitempty"`

	// enable/disable skip SSL certificate verification
	InsecureSkipVerify bool `json:"insecure_skip_verify,omitempty"`

	// TLS key path
	Key string `json:"key,omitempty"`
}

// TolerationsItems
type TolerationsItems struct{}

// TopologySpreadConstraintsItems
type TopologySpreadConstraintsItems struct{}

// Trace
type Trace struct {
	// trace enabled
	Enabled bool `json:"enabled,omitempty"`
}

// Values
type Values struct {
	Agent       *Agent       `json:"agent,omitempty"`
	Defaults    *Defaults    `json:"defaults,omitempty"`
	Discoverer  *Discoverer  `json:"discoverer,omitempty"`
	Gateway     *Gateway     `json:"gateway,omitempty"`
	Initializer *Initializer `json:"initializer,omitempty"`
	Manager     *Manager     `json:"manager,omitempty"`
}

// VolumeMountsItems
type VolumeMountsItems struct{}

// VolumesItems
type VolumesItems struct{}

// Vqueue
type VQueue struct {
	// delete slice pool buffer size
	DeleteBufferPoolSize int `json:"delete_buffer_pool_size,omitempty"`
	// insert slice pool buffer size
	InsertBufferPoolSize int `json:"insert_buffer_pool_size,omitempty"`
}

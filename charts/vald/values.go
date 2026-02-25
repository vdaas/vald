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
package vald

import "github.com/vdaas/vald/internal/config"

// Affinity.
type Affinity struct {
	NodeAffinity    *NodeAffinity    `json:"nodeAffinity,omitempty"`
	PodAffinity     *PodAffinity     `json:"podAffinity,omitempty"`
	PodAntiAffinity *PodAntiAffinity `json:"podAntiAffinity,omitempty"`
}

// Agent.
type Agent struct {
	PersistentVolume              *PersistentVolume                 `json:"persistentVolume,omitempty"`
	ServerConfig                  *ServerConfig                     `json:"server_config,omitempty"`
	Sidecar                       *Sidecar                          `json:"sidecar,omitempty"`
	Service                       *Service                          `json:"service,omitempty"`
	Affinity                      *Affinity                         `json:"affinity,omitempty"`
	Hpa                           *Hpa                              `json:"hpa,omitempty"`
	Image                         *Image                            `json:"image,omitempty"`
	SecurityContext               *SecurityContext                  `json:"securityContext,omitempty"`
	RollingUpdate                 *RollingUpdate                    `json:"rollingUpdate,omitempty"`
	Logging                       *Logging                          `json:"logging,omitempty"`
	Resources                     *Resources                        `json:"resources,omitempty"`
	PodSecurityContext            *PodSecurityContext               `json:"podSecurityContext,omitempty"`
	PodPriority                   *PodPriority                      `json:"podPriority,omitempty"`
	NodeSelector                  *NodeSelector                     `json:"nodeSelector,omitempty"`
	Ngt                           *config.NGT                       `json:"ngt,omitempty"`
	Annotations                   *Annotations                      `json:"annotations,omitempty"`
	PodAnnotations                *PodAnnotations                   `json:"podAnnotations,omitempty"`
	Observability                 *Observability                    `json:"observability,omitempty"`
	ExternalTrafficPolicy         string                            `json:"externalTrafficPolicy,omitempty"`
	NodeName                      string                            `json:"nodeName,omitempty"`
	Version                       string                            `json:"version,omitempty"`
	ServiceType                   string                            `json:"serviceType,omitempty"`
	MaxUnavailable                string                            `json:"maxUnavailable,omitempty"`
	PodManagementPolicy           string                            `json:"podManagementPolicy,omitempty"`
	TimeZone                      string                            `json:"time_zone,omitempty"`
	Kind                          string                            `json:"kind,omitempty"`
	Name                          string                            `json:"name,omitempty"`
	Tolerations                   []*TolerationsItems               `json:"tolerations,omitempty"`
	Env                           []*EnvItems                       `json:"env,omitempty"`
	InitContainers                []*InitContainersItems            `json:"initContainers,omitempty"`
	TopologySpreadConstraints     []*TopologySpreadConstraintsItems `json:"topologySpreadConstraints,omitempty"`
	VolumeMounts                  []*VolumeMountsItems              `json:"volumeMounts,omitempty"`
	Volumes                       []*VolumesItems                   `json:"volumes,omitempty"`
	RevisionHistoryLimit          int                               `json:"revisionHistoryLimit,omitempty"`
	TerminationGracePeriodSeconds int                               `json:"terminationGracePeriodSeconds,omitempty"`
	MaxReplicas                   int                               `json:"maxReplicas,omitempty"`
	MinReplicas                   int                               `json:"minReplicas,omitempty"`
	ProgressDeadlineSeconds       int                               `json:"progressDeadlineSeconds,omitempty"`
	Enabled                       bool                              `json:"enabled,omitempty"`
}

// Annotations deployment annotations.
type Annotations map[string]string

// BlobStorage.
type BlobStorage struct {
	// bucket name
	Bucket       string        `json:"bucket,omitempty"`
	CloudStorage *CloudStorage `json:"cloud_storage,omitempty"`
	S3           *S3           `json:"s3,omitempty"`

	// storage type
	StorageType string `json:"storage_type,omitempty"`
}

// CloudStorage.
type CloudStorage struct {
	Client                  *config.GRPCClient `json:"client,omitempty"`
	Url                     string             `json:"url,omitempty"`
	WriteCacheControl       string             `json:"write_cache_control,omitempty"`
	WriteContentDisposition string             `json:"write_content_disposition,omitempty"`
	WriteContentEncoding    string             `json:"write_content_encoding,omitempty"`
	WriteContentLanguage    string             `json:"write_content_language,omitempty"`
	WriteContentType        string             `json:"write_content_type,omitempty"`
	WriteBufferSize         int                `json:"write_buffer_size,omitempty"`
}

// ClusterRole.
type ClusterRole struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// ClusterRoleBinding.
type ClusterRoleBinding struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// Collector.
type Collector struct {
	Metrics  *Metrics `json:"metrics,omitempty"`
	Duration string   `json:"duration,omitempty"`
}

// Compress.
type Compress struct {
	// compression algorithm. must be `gob`, `gzip`, `lz4` or `zstd`
	CompressAlgorithm string `json:"compress_algorithm,omitempty"`

	// compression level. value range relies on which algorithm is used. `gob`: level will be ignored. `gzip`: -1 (default compression), 0 (no compression), or 1 (best speed) to 9 (best compression). `lz4`: >= 0, higher is better compression. `zstd`: 1 (fastest) to 22 (best), however implementation relies on klauspost/compress.
	CompressionLevel int `json:"compression_level,omitempty"`
}

// Config.
type Config struct {
	BlobStorage           *BlobStorage       `json:"blob_storage,omitempty"`
	Client                *config.GRPCClient `json:"client,omitempty"`
	Compress              *Compress          `json:"compress,omitempty"`
	RestoreBackoff        *config.Backoff    `json:"restore_backoff,omitempty"`
	AutoBackupDuration    string             `json:"auto_backup_duration,omitempty"`
	Filename              string             `json:"filename,omitempty"`
	FilenameSuffix        string             `json:"filename_suffix,omitempty"`
	PostStopTimeout       string             `json:"post_stop_timeout,omitempty"`
	AutoBackupEnabled     bool               `json:"auto_backup_enabled,omitempty"`
	RestoreBackoffEnabled bool               `json:"restore_backoff_enabled,omitempty"`
	WatchEnabled          bool               `json:"watch_enabled,omitempty"`
}

// Defaults.
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

// Dialer.
type Dialer struct {
	Keepalive        string `json:"keepalive,omitempty"`
	Timeout          string `json:"timeout,omitempty"`
	DualStackEnabled bool   `json:"dual_stack_enabled,omitempty"`
}

// Discoverer.
type Discoverer struct {
	AgentClientOptions *config.GRPCClient `json:"agent_client_options,omitempty"`
	Client             *config.GRPCClient `json:"client,omitempty"`

	// refresh duration to discover
	Duration string `json:"duration,omitempty"`
}

// Dns.
type Dns struct {
	CacheExpiration string `json:"cache_expiration,omitempty"`
	RefreshDuration string `json:"refresh_duration,omitempty"`
	CacheEnabled    bool   `json:"cache_enabled,omitempty"`
}

// EgressFilter gRPC client config for egress filter.
type EgressFilter struct {
	Client *config.GRPCClient `json:"client,omitempty"`

	// distance egress vector filter targets
	DistanceFilters []string `json:"distance_filters,omitempty"`

	// object egress vector filter targets
	ObjectFilters []string `json:"object_filters,omitempty"`
}

// EnvItems.
type EnvItems struct{}

// Fields k8s field selectors for pod discovery.
type Fields struct{}

// Filter.
type Filter struct {
	ServerConfig                  *ServerConfig                     `json:"server_config,omitempty"`
	Service                       *Service                          `json:"service,omitempty"`
	Observability                 *Observability                    `json:"observability,omitempty"`
	Affinity                      *Affinity                         `json:"affinity,omitempty"`
	SecurityContext               *SecurityContext                  `json:"securityContext,omitempty"`
	PodPriority                   *PodPriority                      `json:"podPriority,omitempty"`
	Hpa                           *Hpa                              `json:"hpa,omitempty"`
	Image                         *Image                            `json:"image,omitempty"`
	Ingress                       *Ingress                          `json:"ingress,omitempty"`
	RollingUpdate                 *RollingUpdate                    `json:"rollingUpdate,omitempty"`
	Resources                     *Resources                        `json:"resources,omitempty"`
	Logging                       *Logging                          `json:"logging,omitempty"`
	PodSecurityContext            *PodSecurityContext               `json:"podSecurityContext,omitempty"`
	GatewayConfig                 *GatewayConfig                    `json:"gateway_config,omitempty"`
	PodAnnotations                *PodAnnotations                   `json:"podAnnotations,omitempty"`
	Annotations                   *Annotations                      `json:"annotations,omitempty"`
	NodeSelector                  *NodeSelector                     `json:"nodeSelector,omitempty"`
	NodeName                      string                            `json:"nodeName,omitempty"`
	Name                          string                            `json:"name,omitempty"`
	Version                       string                            `json:"version,omitempty"`
	MaxUnavailable                string                            `json:"maxUnavailable,omitempty"`
	TimeZone                      string                            `json:"time_zone,omitempty"`
	ServiceType                   string                            `json:"serviceType,omitempty"`
	Kind                          string                            `json:"kind,omitempty"`
	ExternalTrafficPolicy         string                            `json:"externalTrafficPolicy,omitempty"`
	Tolerations                   []*TolerationsItems               `json:"tolerations,omitempty"`
	InitContainers                []*InitContainersItems            `json:"initContainers,omitempty"`
	Env                           []*EnvItems                       `json:"env,omitempty"`
	TopologySpreadConstraints     []*TopologySpreadConstraintsItems `json:"topologySpreadConstraints,omitempty"`
	VolumeMounts                  []*VolumeMountsItems              `json:"volumeMounts,omitempty"`
	Volumes                       []*VolumesItems                   `json:"volumes,omitempty"`
	RevisionHistoryLimit          int                               `json:"revisionHistoryLimit,omitempty"`
	ProgressDeadlineSeconds       int                               `json:"progressDeadlineSeconds,omitempty"`
	TerminationGracePeriodSeconds int                               `json:"terminationGracePeriodSeconds,omitempty"`
	MaxReplicas                   int                               `json:"maxReplicas,omitempty"`
	MinReplicas                   int                               `json:"minReplicas,omitempty"`
	Enabled                       bool                              `json:"enabled,omitempty"`
}

// Gateway.
type Gateway struct {
	Filter *Filter `json:"filter,omitempty"`
	Lb     *Lb     `json:"lb,omitempty"`
}

// GatewayConfig.
type GatewayConfig struct {
	// gRPC client config for egress filter
	EgressFilter  *EgressFilter      `json:"egress_filter,omitempty"`
	GatewayClient *config.GRPCClient `json:"gateway_client,omitempty"`

	// gRPC client config for ingress filter
	IngressFilter *IngressFilter `json:"ingress_filter,omitempty"`
}

// Grpc.
type Grpc struct {
	Server      *Server `json:"server,omitempty"`
	Host        string  `json:"host,omitempty"`
	Port        int     `json:"port,omitempty"`
	ServicePort int     `json:"servicePort,omitempty"`
	Enabled     bool    `json:"enabled,omitempty"`
}

// Healths.
type Healths struct {
	Liveness  *Liveness  `json:"liveness,omitempty"`
	Readiness *Readiness `json:"readiness,omitempty"`
	Startup   *Startup   `json:"startup,omitempty"`
}

// Hpa.
type Hpa struct {
	// HPA enabled
	Enabled bool `json:"enabled,omitempty"`

	// HPA CPU utilization percentage
	TargetCPUUtilizationPercentage int `json:"targetCPUUtilizationPercentage,omitempty"`
}

// Http.
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

// HttpGet.
type HttpGet struct {
	// startup probe path
	Path string `json:"path,omitempty"`

	// startup probe port
	Port string `json:"port,omitempty"`

	// startup probe scheme
	Scheme string `json:"scheme,omitempty"`
}

// Image.
type Image struct {
	// image pull policy
	PullPolicy string `json:"pullPolicy,omitempty"`

	// image repository
	Repository string `json:"repository,omitempty"`

	// image tag (overrides defaults.image.tag)
	Tag string `json:"tag,omitempty"`
}

// Index.
type Index struct {
	ServerConfig                  *ServerConfig                     `json:"server_config,omitempty"`
	PodSecurityContext            *PodSecurityContext               `json:"podSecurityContext,omitempty"`
	Service                       *Service                          `json:"service,omitempty"`
	Affinity                      *Affinity                         `json:"affinity,omitempty"`
	SecurityContext               *SecurityContext                  `json:"securityContext,omitempty"`
	Image                         *Image                            `json:"image,omitempty"`
	Indexer                       *Indexer                          `json:"indexer,omitempty"`
	RollingUpdate                 *RollingUpdate                    `json:"rollingUpdate,omitempty"`
	Resources                     *Resources                        `json:"resources,omitempty"`
	Logging                       *Logging                          `json:"logging,omitempty"`
	PodPriority                   *PodPriority                      `json:"podPriority,omitempty"`
	Annotations                   *Annotations                      `json:"annotations,omitempty"`
	PodAnnotations                *PodAnnotations                   `json:"podAnnotations,omitempty"`
	NodeSelector                  *NodeSelector                     `json:"nodeSelector,omitempty"`
	Observability                 *Observability                    `json:"observability,omitempty"`
	ServiceType                   string                            `json:"serviceType,omitempty"`
	Name                          string                            `json:"name,omitempty"`
	TimeZone                      string                            `json:"time_zone,omitempty"`
	MaxUnavailable                string                            `json:"maxUnavailable,omitempty"`
	ExternalTrafficPolicy         string                            `json:"externalTrafficPolicy,omitempty"`
	Kind                          string                            `json:"kind,omitempty"`
	NodeName                      string                            `json:"nodeName,omitempty"`
	Version                       string                            `json:"version,omitempty"`
	InitContainers                []*InitContainersItems            `json:"initContainers,omitempty"`
	Tolerations                   []*TolerationsItems               `json:"tolerations,omitempty"`
	TopologySpreadConstraints     []*TopologySpreadConstraintsItems `json:"topologySpreadConstraints,omitempty"`
	Env                           []*EnvItems                       `json:"env,omitempty"`
	VolumeMounts                  []*VolumeMountsItems              `json:"volumeMounts,omitempty"`
	Volumes                       []*VolumesItems                   `json:"volumes,omitempty"`
	Replicas                      int                               `json:"replicas,omitempty"`
	TerminationGracePeriodSeconds int                               `json:"terminationGracePeriodSeconds,omitempty"`
	RevisionHistoryLimit          int                               `json:"revisionHistoryLimit,omitempty"`
	ProgressDeadlineSeconds       int                               `json:"progressDeadlineSeconds,omitempty"`
	Enabled                       bool                              `json:"enabled,omitempty"`
}

// Indexer.
type Indexer struct {
	Discoverer                 *Discoverer `json:"discoverer,omitempty"`
	AgentNamespace             string      `json:"agent_namespace,omitempty"`
	AutoIndexCheckDuration     string      `json:"auto_index_check_duration,omitempty"`
	AutoIndexDurationLimit     string      `json:"auto_index_duration_limit,omitempty"`
	AutoSaveIndexDurationLimit string      `json:"auto_save_index_duration_limit,omitempty"`
	AutoSaveIndexWaitDuration  string      `json:"auto_save_index_wait_duration,omitempty"`
	NodeName                   string      `json:"node_name,omitempty"`
	AutoIndexLength            int         `json:"auto_index_length,omitempty"`
	Concurrency                int         `json:"concurrency,omitempty"`
	CreationPoolSize           int         `json:"creation_pool_size,omitempty"`
}

// Ingress.
type Ingress struct {
	Annotations *Annotations `json:"annotations,omitempty"`
	Host        string       `json:"host,omitempty"`
	PathType    string       `json:"pathType,omitempty"`
	ServicePort string       `json:"servicePort,omitempty"`
	Enabled     bool         `json:"enabled,omitempty"`
}

// IngressFilter gRPC client config for ingress filter.
type IngressFilter struct {
	Client        *config.GRPCClient `json:"client,omitempty"`
	Vectorizer    string             `json:"vectorizer,omitempty"`
	InsertFilters []string           `json:"insert_filters,omitempty"`
	SearchFilters []string           `json:"search_filters,omitempty"`
	UpdateFilters []string           `json:"update_filters,omitempty"`
	UpsertFilters []string           `json:"upsert_filters,omitempty"`
}

// InitContainersItems.
type InitContainersItems struct {
	Type          string
	Name          string
	Target        string
	Image         string
	SleepDuration int
}

// Initializer.
type Initializer struct{}

// Jaeger.
type Jaeger struct {
	AgentEndpoint          string `json:"agent_endpoint,omitempty"`
	AgentReconnectInterval string `json:"agent_reconnect_interval,omitempty"`
	CollectorEndpoint      string `json:"collector_endpoint,omitempty"`
	Password               string `json:"password,omitempty"`
	ServiceName            string `json:"service_name,omitempty"`
	Username               string `json:"username,omitempty"`
	BatchTimeout           string `json:"batch_timeout,omitempty"`
	ExportTimeout          string `json:"export_timeout,omitempty"`
	AgentMaxPacketSize     int    `json:"agent_max_packet_size,omitempty"`
	MaxExportBatchSize     int    `json:"max_export_batch_size,omitempty"`
	MaxQueueSize           int    `json:"max_queue_size,omitempty"`
	Enabled                bool   `json:"enabled,omitempty"`
}

// Keepalive.
type Keepalive struct {
	Time                string `json:"time,omitempty"`
	Timeout             string `json:"timeout,omitempty"`
	PermitWithoutStream bool   `json:"permit_without_stream,omitempty"`
}

// Kvsdb.
type Kvsdb struct {
	// kvsdb processing concurrency
	Concurrency int `json:"concurrency,omitempty"`
}

// Labels service labels.
type Labels map[string]string

// Lb.
type Lb struct {
	ServerConfig                  *ServerConfig                     `json:"server_config,omitempty"`
	Service                       *Service                          `json:"service,omitempty"`
	Observability                 *Observability                    `json:"observability,omitempty"`
	Affinity                      *Affinity                         `json:"affinity,omitempty"`
	SecurityContext               *SecurityContext                  `json:"securityContext,omitempty"`
	PodPriority                   *PodPriority                      `json:"podPriority,omitempty"`
	Hpa                           *Hpa                              `json:"hpa,omitempty"`
	Image                         *Image                            `json:"image,omitempty"`
	Ingress                       *Ingress                          `json:"ingress,omitempty"`
	RollingUpdate                 *RollingUpdate                    `json:"rollingUpdate,omitempty"`
	Resources                     *Resources                        `json:"resources,omitempty"`
	Logging                       *Logging                          `json:"logging,omitempty"`
	PodSecurityContext            *PodSecurityContext               `json:"podSecurityContext,omitempty"`
	GatewayConfig                 *GatewayConfig                    `json:"gateway_config,omitempty"`
	PodAnnotations                *PodAnnotations                   `json:"podAnnotations,omitempty"`
	Annotations                   *Annotations                      `json:"annotations,omitempty"`
	NodeSelector                  *NodeSelector                     `json:"nodeSelector,omitempty"`
	NodeName                      string                            `json:"nodeName,omitempty"`
	Name                          string                            `json:"name,omitempty"`
	Version                       string                            `json:"version,omitempty"`
	MaxUnavailable                string                            `json:"maxUnavailable,omitempty"`
	TimeZone                      string                            `json:"time_zone,omitempty"`
	ServiceType                   string                            `json:"serviceType,omitempty"`
	Kind                          string                            `json:"kind,omitempty"`
	ExternalTrafficPolicy         string                            `json:"externalTrafficPolicy,omitempty"`
	Tolerations                   []*TolerationsItems               `json:"tolerations,omitempty"`
	InitContainers                []*InitContainersItems            `json:"initContainers,omitempty"`
	Env                           []*EnvItems                       `json:"env,omitempty"`
	TopologySpreadConstraints     []*TopologySpreadConstraintsItems `json:"topologySpreadConstraints,omitempty"`
	VolumeMounts                  []*VolumeMountsItems              `json:"volumeMounts,omitempty"`
	Volumes                       []*VolumesItems                   `json:"volumes,omitempty"`
	RevisionHistoryLimit          int                               `json:"revisionHistoryLimit,omitempty"`
	ProgressDeadlineSeconds       int                               `json:"progressDeadlineSeconds,omitempty"`
	TerminationGracePeriodSeconds int                               `json:"terminationGracePeriodSeconds,omitempty"`
	MaxReplicas                   int                               `json:"maxReplicas,omitempty"`
	MinReplicas                   int                               `json:"minReplicas,omitempty"`
	Enabled                       bool                              `json:"enabled,omitempty"`
}

// Limits.
type Limits struct{}

// Liveness.
type Liveness struct {
	LivenessProbe *LivenessProbe `json:"livenessProbe,omitempty"`
	Server        *Server        `json:"server,omitempty"`
	Host          string         `json:"host,omitempty"`
	Port          int            `json:"port,omitempty"`
	ServicePort   int            `json:"servicePort,omitempty"`
	Enabled       bool           `json:"enabled,omitempty"`
}

// LivenessProbe.
type LivenessProbe struct {
	HttpGet             *HttpGet `json:"httpGet,omitempty"`
	FailureThreshold    int      `json:"failureThreshold,omitempty"`
	InitialDelaySeconds int      `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int      `json:"periodSeconds,omitempty"`
	SuccessThreshold    int      `json:"successThreshold,omitempty"`
	TimeoutSeconds      int      `json:"timeoutSeconds,omitempty"`
}

// Logging.
type Logging struct {
	// logging format. logging format must be `raw` or `json`
	Format string `json:"format,omitempty"`

	// logging level. logging level must be `debug`, `info`, `warn`, `error` or `fatal`.
	Level string `json:"level,omitempty"`

	// logger name. currently logger must be `glg` or `zap`.
	Logger string `json:"logger,omitempty"`
}

// Manager.
type Manager struct {
	Index *Index `json:"index,omitempty"`
}

// Metrics.
type Metrics struct {
	VersionInfoLabels []string `json:"version_info_labels,omitempty"`
	EnableCgo         bool     `json:"enable_cgo,omitempty"`
	EnableGoroutine   bool     `json:"enable_goroutine,omitempty"`
	EnableMemory      bool     `json:"enable_memory,omitempty"`
	EnableVersionInfo bool     `json:"enable_version_info,omitempty"`
}

// Node k8s resource selectors for node discovery.
type Node struct {
	// k8s field selectors for node discovery
	Fields *Fields `json:"fields,omitempty"`
	// k8s label selectors for node discovery
	Labels *Labels `json:"labels,omitempty"`
}

// NodeAffinity.
type NodeAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  *RequiredDuringSchedulingIgnoredDuringExecution         `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
	PreferredDuringSchedulingIgnoredDuringExecution []*PreferredDuringSchedulingIgnoredDuringExecutionItems `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// NodeMetrics k8s resource selectors for node_metrics discovery.
type NodeMetrics struct {
	// k8s field selectors for node_metrics discovery
	Fields *Fields `json:"fields,omitempty"`

	// k8s label selectors for node_metrics discovery
	Labels *Labels `json:"labels,omitempty"`
}

// NodeSelector node selector.
type NodeSelector struct{}

// NodeSelectorTermsItems.
type NodeSelectorTermsItems struct{}

// Observability.
type Observability struct {
	Metrics    *Metrics    `json:"metrics,omitempty"`
	Jaeger     *Jaeger     `json:"jaeger,omitempty"`
	Prometheus *Prometheus `json:"prometheus,omitempty"`
	Trace      *Trace      `json:"trace,omitempty"`
	Enabled    bool        `json:"enabled,omitempty"`
}

// PersistentVolume.
type PersistentVolume struct {
	AccessMode   string `json:"accessMode,omitempty"`
	Size         string `json:"size,omitempty"`
	StorageClass string `json:"storageClass,omitempty"`
	Enabled      bool   `json:"enabled,omitempty"`
}

// Pod k8s resource selectors for pod discovery.
type Pod struct {
	// k8s field selectors for pod discovery
	Fields *Fields `json:"fields,omitempty"`

	// k8s label selectors for pod discovery
	Labels *Labels `json:"labels,omitempty"`
}

// PodAffinity.
type PodAffinity struct {
	// pod affinity preferred scheduling terms
	PreferredDuringSchedulingIgnoredDuringExecution []*PreferredDuringSchedulingIgnoredDuringExecutionItems `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`

	// pod affinity required scheduling terms
	RequiredDuringSchedulingIgnoredDuringExecution []*RequiredDuringSchedulingIgnoredDuringExecutionItems `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// PodAnnotations pod annotations.
type PodAnnotations struct{}

// PodAntiAffinity.
type PodAntiAffinity struct {
	// pod anti-affinity preferred scheduling terms
	PreferredDuringSchedulingIgnoredDuringExecution []*PreferredDuringSchedulingIgnoredDuringExecutionItems `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`

	// pod anti-affinity required scheduling terms
	RequiredDuringSchedulingIgnoredDuringExecution []*RequiredDuringSchedulingIgnoredDuringExecutionItems `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

// PodMetrics k8s resource selectors for pod_metrics discovery.
type PodMetrics struct {
	// k8s field selectors for pod_metrics discovery
	Fields *Fields `json:"fields,omitempty"`

	// k8s label selectors for pod_metrics discovery
	Labels *Labels `json:"labels,omitempty"`
}

// PodPriority.
type PodPriority struct {
	// gateway pod PriorityClass enabled
	Enabled bool `json:"enabled,omitempty"`

	// gateway pod PriorityClass value
	Value int `json:"value,omitempty"`
}

// PodSecurityContext security context for pod.
type PodSecurityContext struct{}

// Pprof.
type Pprof struct {
	Server      *Server `json:"server,omitempty"`
	Host        string  `json:"host,omitempty"`
	Port        int     `json:"port,omitempty"`
	ServicePort int     `json:"servicePort,omitempty"`
	Enabled     bool    `json:"enabled,omitempty"`
}

// PreferredDuringSchedulingIgnoredDuringExecutionItems.
type PreferredDuringSchedulingIgnoredDuringExecutionItems struct{}

// Prometheus.
type Prometheus struct {
	Endpoint           string `json:"endpoint,omitempty"`
	Namespace          string `json:"namespace,omitempty"`
	CollectInterval    string `json:"collect_interval,omitempty"`
	CollectTimeout     string `json:"collect_timeout,omitempty"`
	Enabled            bool   `json:"enabled,omitempty"`
	EnableInMemoryMode bool   `json:"enable_in_memory_mode,omitempty"`
}

// Readiness.
type Readiness struct {
	ReadinessProbe *ReadinessProbe `json:"readinessProbe,omitempty"`
	Server         *Server         `json:"server,omitempty"`
	Host           string          `json:"host,omitempty"`
	Port           int             `json:"port,omitempty"`
	ServicePort    int             `json:"servicePort,omitempty"`
	Enabled        bool            `json:"enabled,omitempty"`
}

// ReadinessProbe.
type ReadinessProbe struct {
	HttpGet             *HttpGet `json:"httpGet,omitempty"`
	FailureThreshold    int      `json:"failureThreshold,omitempty"`
	InitialDelaySeconds int      `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int      `json:"periodSeconds,omitempty"`
	SuccessThreshold    int      `json:"successThreshold,omitempty"`
	TimeoutSeconds      int      `json:"timeoutSeconds,omitempty"`
}

// Requests.
type Requests struct{}

// RequiredDuringSchedulingIgnoredDuringExecution.
type RequiredDuringSchedulingIgnoredDuringExecution struct {
	// node affinity required node selectors
	NodeSelectorTerms []*NodeSelectorTermsItems `json:"nodeSelectorTerms,omitempty"`
}

// RequiredDuringSchedulingIgnoredDuringExecutionItems.
type RequiredDuringSchedulingIgnoredDuringExecutionItems struct{}

// Resources compute resources.
type Resources struct {
	Limits   *Limits   `json:"limits,omitempty"`
	Requests *Requests `json:"requests,omitempty"`
}

// Rest.
type Rest struct {
	Server      *Server `json:"server,omitempty"`
	Host        string  `json:"host,omitempty"`
	Port        int     `json:"port,omitempty"`
	ServicePort int     `json:"servicePort,omitempty"`
	Enabled     bool    `json:"enabled,omitempty"`
}

// RollingUpdate.
type RollingUpdate struct {
	// max surge of rolling update
	MaxSurge string `json:"maxSurge,omitempty"`

	// max unavailable of rolling update
	MaxUnavailable string `json:"maxUnavailable,omitempty"`
}

// S3.
type S3 struct {
	MaxChunkSize               string `json:"max_chunk_size,omitempty"`
	Token                      string `json:"token,omitempty"`
	SecretAccessKey            string `json:"secret_access_key,omitempty"`
	Region                     string `json:"region,omitempty"`
	AccessKey                  string `json:"access_key,omitempty"`
	MaxPartSize                string `json:"max_part_size,omitempty"`
	Endpoint                   string `json:"endpoint,omitempty"`
	MaxRetries                 int    `json:"max_retries,omitempty"`
	EnableEndpointHostPrefix   bool   `json:"enable_endpoint_host_prefix,omitempty"`
	ForcePathStyle             bool   `json:"force_path_style,omitempty"`
	EnableSsl                  bool   `json:"enable_ssl,omitempty"`
	EnableParamValidation      bool   `json:"enable_param_validation,omitempty"`
	EnableEndpointDiscovery    bool   `json:"enable_endpoint_discovery,omitempty"`
	EnableContentMd5Validation bool   `json:"enable_content_md5_validation,omitempty"`
	Enable100Continue          bool   `json:"enable_100_continue,omitempty"`
	UseAccelerate              bool   `json:"use_accelerate,omitempty"`
	UseArnRegion               bool   `json:"use_arn_region,omitempty"`
	UseDualStack               bool   `json:"use_dual_stack,omitempty"`
}

// SecurityContext security context for container.
type SecurityContext struct{}

// Selectors k8s resource selectors.
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

// Server.
type Server struct {
	Http *Http `json:"http,omitempty"`

	// REST server mode
	Mode string `json:"mode,omitempty"`

	// mysql network
	Network string `json:"network,omitempty"`

	// REST server probe wait time
	ProbeWaitTime string               `json:"probe_wait_time,omitempty"`
	SocketOption  *config.SocketOption `json:"socket_option,omitempty"`

	// mysql socket_path
	SocketPath string `json:"socket_path,omitempty"`
}

// ServerConfig.
type ServerConfig struct {
	Healths              *Healths `json:"healths,omitempty"`
	Metrics              *Metrics `json:"metrics,omitempty"`
	Servers              *Servers `json:"servers,omitempty"`
	Tls                  *Tls     `json:"tls,omitempty"`
	FullShutdownDuration string   `json:"full_shutdown_duration,omitempty"`
}

// Servers.
type Servers struct {
	Grpc *Grpc `json:"grpc,omitempty"`
	Rest *Rest `json:"rest,omitempty"`
}

// Service.
type Service struct {
	// service annotations
	Annotations *Annotations `json:"annotations,omitempty"`

	// service labels
	Labels *Labels `json:"labels,omitempty"`
}

// ServiceAccount.
type ServiceAccount struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// Sidecar.
type Sidecar struct {
	ServerConfig         *ServerConfig  `json:"server_config,omitempty"`
	Image                *Image         `json:"image,omitempty"`
	Logging              *Logging       `json:"logging,omitempty"`
	Observability        *Observability `json:"observability,omitempty"`
	Resources            *Resources     `json:"resources,omitempty"`
	Config               *Config        `json:"config,omitempty"`
	Service              *Service       `json:"service,omitempty"`
	Name                 string         `json:"name,omitempty"`
	TimeZone             string         `json:"time_zone,omitempty"`
	Version              string         `json:"version,omitempty"`
	Env                  []*EnvItems    `json:"env,omitempty"`
	Enabled              bool           `json:"enabled,omitempty"`
	InitContainerEnabled bool           `json:"initContainerEnabled,omitempty"`
}

// Startup.
type Startup struct {
	StartupProbe *StartupProbe `json:"startupProbe,omitempty"`
	Port         int           `json:"port,omitempty"`
	Enabled      bool          `json:"enabled,omitempty"`
}

// StartupProbe.
type StartupProbe struct {
	HttpGet             *HttpGet `json:"httpGet,omitempty"`
	FailureThreshold    int      `json:"failureThreshold,omitempty"`
	InitialDelaySeconds int      `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int      `json:"periodSeconds,omitempty"`
	SuccessThreshold    int      `json:"successThreshold,omitempty"`
	TimeoutSeconds      int      `json:"timeoutSeconds,omitempty"`
}

// Tls.
type Tls struct {
	Ca                 string `json:"ca,omitempty"`
	Cert               string `json:"cert,omitempty"`
	Key                string `json:"key,omitempty"`
	Enabled            bool   `json:"enabled,omitempty"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify,omitempty"`
}

// TolerationsItems.
type TolerationsItems struct{}

// TopologySpreadConstraintsItems.
type TopologySpreadConstraintsItems struct{}

// Trace.
type Trace struct {
	// trace enabled
	Enabled bool `json:"enabled,omitempty"`
}

// Values.
type Values struct {
	Agent       *Agent       `json:"agent,omitempty"`
	Defaults    *Defaults    `json:"defaults,omitempty"`
	Discoverer  *Discoverer  `json:"discoverer,omitempty"`
	Gateway     *Gateway     `json:"gateway,omitempty"`
	Initializer *Initializer `json:"initializer,omitempty"`
	Manager     *Manager     `json:"manager,omitempty"`
}

// VolumeMountsItems.
type VolumeMountsItems struct{}

// VolumesItems.
type VolumesItems struct{}

// Vqueue.
type VQueue struct {
	// delete slice pool buffer size
	DeleteBufferPoolSize int `json:"delete_buffer_pool_size,omitempty"`
	// insert slice pool buffer size
	InsertBufferPoolSize int `json:"insert_buffer_pool_size,omitempty"`
}

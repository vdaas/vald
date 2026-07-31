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
package config

import (
	"time"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease"
	"github.com/vdaas/vald/internal/timeutil"
)

type GlobalConfig = config.GlobalConfig

type Data struct {
	Server              *config.Servers       `json:"server_config" yaml:"server_config"` //nolint:tagliatelle
	Observability       *config.Observability `json:"observability" yaml:"observability"`
	Operator            *Operator             `json:"operator"      yaml:"operator"`
	config.GlobalConfig `json:",inline" yaml:",inline"`
}

type LeaderElection struct {
	// ID represents the name of the leader election lock resource.
	ID string `json:"id" yaml:"id"`
	// Namespace represents the namespace in which the leader election lock resource is created.
	Namespace string `json:"namespace" yaml:"namespace"`
	// LeaseDuration represents how long non-leader candidates wait before acquiring an expired lease.
	LeaseDuration string `json:"lease_duration" yaml:"lease_duration"`
	// RenewDeadline represents how long the acting leader retries refreshing its lease before giving up.
	RenewDeadline string `json:"renew_deadline" yaml:"renew_deadline"`
	// RetryPeriod represents the wait between leader election actions.
	RetryPeriod string `json:"retry_period" yaml:"retry_period"`
	// Enabled enables the leader election.
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// Requeue represents the reconcile requeue intervals of the ValdOperatorRelease
// controller. Empty values keep the default behavior: no periodic requeue on
// success, exponential backoff on error and no requeue on NotFound.
type Requeue struct {
	// Success represents the periodic requeue interval after a successful reconcile.
	Success string `json:"success" yaml:"success"`
	// OnError represents the fixed retry interval when a reconcile fails.
	OnError string `json:"on_error" yaml:"on_error"`
	// NotFound represents the retry interval when the reconciled object is not found.
	NotFound string `json:"not_found" yaml:"not_found"`
}

type Controller struct {
	// LeaderElection represents the leader election configurations.
	LeaderElection *LeaderElection `json:"leader_election" yaml:"leader_election"`
	// Requeue represents the reconcile requeue intervals.
	Requeue *Requeue `json:"requeue" yaml:"requeue"`
	// MetricsAddress represents the bind address of the manager metrics server ("0" disables it).
	MetricsAddress string `json:"metrics_address" yaml:"metrics_address"`
	// SyncPeriod represents the informer cache resync interval (empty keeps the manager default).
	SyncPeriod string `json:"sync_period" yaml:"sync_period"`
	// CacheNamespaces restricts the manager cache to the given namespaces (empty watches the cluster).
	CacheNamespaces []string `json:"cache_namespaces" yaml:"cache_namespaces"`
	// MaxConcurrentReconciles represents the reconcile worker count of the ValdOperatorRelease controller.
	MaxConcurrentReconciles int `json:"max_concurrent_reconciles" yaml:"max_concurrent_reconciles"`
}

// VrsAgent represents the agent-related values the operator injects into the
// generated ValdRelease. These override the default VRS template, so they are
// exposed here rather than in the template.
type VrsAgent struct {
	EnableInMemoryMode *bool  `json:"enable_in_memory_mode" yaml:"enable_in_memory_mode"`
	MaxSurge           string `json:"max_surge"             yaml:"max_surge"`
	MaxUnavailable     string `json:"max_unavailable"       yaml:"max_unavailable"`
}

type VrsGateway struct {
	// IngressServicePort represents the service port name referenced by the gateway ingress.
	IngressServicePort string `json:"ingress_service_port" yaml:"ingress_service_port"`
	// IngressPathType represents the ingress path type (Prefix, Exact or ImplementationSpecific).
	IngressPathType string `json:"ingress_path_type" yaml:"ingress_path_type"`
	// HpaTargetCPUUtilization represents the gateway lb HPA target CPU utilization percentage.
	HpaTargetCPUUtilization int `json:"hpa_target_cpu_utilization" yaml:"hpa_target_cpu_utilization"`
}

type VrsIndexer struct {
	// AutoIndexDurationLimit represents the indexer auto index duration limit.
	AutoIndexDurationLimit string `json:"auto_index_duration_limit" yaml:"auto_index_duration_limit"`
	// AutoSaveIndexDurationLimit represents the indexer auto save index duration limit.
	AutoSaveIndexDurationLimit string `json:"auto_save_index_duration_limit" yaml:"auto_save_index_duration_limit"`
}

// Vrs represents the configurations that control how the operator renders
// ValdRelease resources. Component-level VRS contents (image, resources, ...)
// belong to the default VRS template; only values computed or overridden by
// the operator logic live here.
type Vrs struct {
	// Agent represents the agent-related builder values.
	Agent *VrsAgent `json:"agent" yaml:"agent"`
	// Gateway represents the gateway-related builder values.
	Gateway *VrsGateway `json:"gateway" yaml:"gateway"`
	// Indexer represents the index-manager builder values.
	Indexer *VrsIndexer `json:"indexer" yaml:"indexer"`
	// DefaultVrsPath represents the path of the template ValdRelease yaml.
	DefaultVrsPath string `json:"default_vrs_path" yaml:"default_vrs_path"`
	// LogLevel represents the log level passed through to the underlying vald deployment.
	LogLevel string `json:"log_level" yaml:"log_level"`
	// LogFormat represents the log format passed through to the underlying vald deployment.
	LogFormat string `json:"log_format" yaml:"log_format"`
	// Logger represents the logger passed through to the underlying vald deployment.
	Logger string `json:"logger" yaml:"logger"`
	// ManagedGenerationLabel represents the label key that records the owner generation on managed resources.
	ManagedGenerationLabel string `json:"managed_generation_label" yaml:"managed_generation_label"`
}

type NodePool struct {
	// LabelPrefix represents the optional prefix for nodepool-related label keys (e.g., "vald.vdaas.org").
	LabelPrefix string `json:"label_prefix" yaml:"label_prefix"`
	// AgentPodsPerNode represents the number of agent pods per node.
	AgentPodsPerNode int `json:"agent_pods_per_node" yaml:"agent_pods_per_node"`
	// RequireMatch generates VRS only when matching nodepools exist in the cluster.
	RequireMatch bool `json:"require_match" yaml:"require_match"`
}

type PersistentVolume struct {
	// DefaultStorageClass represents the fallback StorageClass when the CR omits it.
	DefaultStorageClass string `json:"default_storage_class" yaml:"default_storage_class"`
	// DefaultAccessMode represents the fallback AccessMode when the CR omits it.
	DefaultAccessMode string `json:"default_access_mode" yaml:"default_access_mode"`
	// BufferRatio represents the PV sizing buffer ratio: max(memoryRequest * BufferRatio, MinSizeBytes).
	BufferRatio float64 `json:"buffer_ratio" yaml:"buffer_ratio"`
	// MinSizeBytes represents the minimum PV size in bytes.
	MinSizeBytes int64 `json:"min_size_bytes" yaml:"min_size_bytes"`
}

type Networking struct {
	// GatewayIngressAnnotations represents the annotations applied to the gateway ingress.
	GatewayIngressAnnotations map[string]string `json:"gateway_ingress_annotations" yaml:"gateway_ingress_annotations"`
	// GatewayServiceType represents the service type of the gateway.
	GatewayServiceType string `json:"gateway_service_type" yaml:"gateway_service_type"`
	// DiscovererDaemonSetMaxSurge represents the discoverer DaemonSet rolling-update max surge (intstr percent string).
	DiscovererDaemonSetMaxSurge string `json:"discoverer_daemonset_max_surge" yaml:"discoverer_daemonset_max_surge"` //nolint:tagliatelle
	// DiscovererDaemonSetMaxUnavailable represents the discoverer DaemonSet rolling-update max unavailable (intstr percent string).
	DiscovererDaemonSetMaxUnavailable string `json:"discoverer_daemonset_max_unavailable" yaml:"discoverer_daemonset_max_unavailable"` //nolint:tagliatelle
	// EnableIngress enables the ingress for the gateway.
	EnableIngress bool `json:"enable_ingress" yaml:"enable_ingress"`
}

type Operator struct {
	// Controller represents the controller manager and reconcile loop configurations.
	Controller *Controller `json:"controller" yaml:"controller"`

	// Vrs represents the ValdRelease rendering configurations.
	Vrs *Vrs `json:"vrs" yaml:"vrs"`

	// NodePool represents the node-pool matching configurations.
	NodePool *NodePool `json:"node_pool" yaml:"node_pool"`

	// PersistentVolume represents the Agent.PersistentVolume fallback configurations.
	PersistentVolume *PersistentVolume `json:"persistent_volume" yaml:"persistent_volume"`

	// Networking represents the gateway/discoverer networking configurations.
	Networking *Networking `json:"networking" yaml:"networking"`

	// Name represents the controller name.
	Name string `json:"name" yaml:"name"`

	// Namespace represents the namespace the operator is deployed in.
	Namespace string `json:"namespace" yaml:"namespace"`
}

const (
	defaultVrsLogFormat            = "raw"
	defaultVrsLogger               = "glg"
	defaultManagedGenerationLabel  = "managed-generation"
	defaultGatewayIngressPort      = "grpc"
	defaultGatewayIngressPathType  = "Prefix"
	defaultGatewayHpaTargetCPU     = 80
	defaultAutoIndexDurationLimit  = "1h"
	defaultAutoSaveIndexDurLimit   = "-1h"
	defaultAgentEnableInMemoryMode = true
)

// Bind binds the actual data from the Operator receiver fields and fills the
// operator-owned defaults for every omitted value so that downstream
// consumers can rely on the sections being non-nil and populated.
func (o *Operator) Bind() *Operator {
	o.Name = config.GetActualValue(o.Name)
	o.Namespace = config.GetActualValue(o.Namespace)

	o.bindController()
	o.bindVrs()
	o.bindNodePool()
	o.bindPersistentVolume()
	o.bindNetworking()

	return o
}

func (o *Operator) bindController() {
	if o.Controller == nil {
		o.Controller = new(Controller)
	}
	c := o.Controller
	c.MetricsAddress = config.GetActualValue(c.MetricsAddress)
	c.SyncPeriod = config.GetActualValue(c.SyncPeriod)
	c.CacheNamespaces = config.GetActualValues(c.CacheNamespaces)

	if c.LeaderElection == nil {
		c.LeaderElection = new(LeaderElection)
	}
	le := c.LeaderElection
	le.ID = config.GetActualValue(le.ID)
	le.Namespace = config.GetActualValue(le.Namespace)
	le.LeaseDuration = config.GetActualValue(le.LeaseDuration)
	le.RenewDeadline = config.GetActualValue(le.RenewDeadline)
	le.RetryPeriod = config.GetActualValue(le.RetryPeriod)

	if c.Requeue == nil {
		c.Requeue = new(Requeue)
	}
	rq := c.Requeue
	rq.Success = config.GetActualValue(rq.Success)
	rq.OnError = config.GetActualValue(rq.OnError)
	rq.NotFound = config.GetActualValue(rq.NotFound)
}

func (o *Operator) bindVrs() {
	if o.Vrs == nil {
		o.Vrs = new(Vrs)
	}
	v := o.Vrs
	v.DefaultVrsPath = config.GetActualValue(v.DefaultVrsPath)
	v.LogLevel = config.GetActualValue(v.LogLevel)
	v.LogFormat = config.GetActualValue(v.LogFormat)
	v.Logger = config.GetActualValue(v.Logger)
	v.ManagedGenerationLabel = config.GetActualValue(v.ManagedGenerationLabel)
	if v.LogFormat == "" {
		v.LogFormat = defaultVrsLogFormat
	}
	if v.Logger == "" {
		v.Logger = defaultVrsLogger
	}
	if v.ManagedGenerationLabel == "" {
		v.ManagedGenerationLabel = defaultManagedGenerationLabel
	}

	if v.Agent == nil {
		v.Agent = new(VrsAgent)
	}
	va := v.Agent
	va.MaxSurge = config.GetActualValue(va.MaxSurge)
	va.MaxUnavailable = config.GetActualValue(va.MaxUnavailable)
	if va.MaxSurge == "" {
		va.MaxSurge = valdrelease.DefaultAgentMaxSurge
	}
	if va.MaxUnavailable == "" {
		va.MaxUnavailable = valdrelease.DefaultAgentMaxUnavailable
	}
	if va.EnableInMemoryMode == nil {
		enabled := defaultAgentEnableInMemoryMode
		va.EnableInMemoryMode = &enabled
	}

	if v.Gateway == nil {
		v.Gateway = new(VrsGateway)
	}
	vg := v.Gateway
	vg.IngressServicePort = config.GetActualValue(vg.IngressServicePort)
	vg.IngressPathType = config.GetActualValue(vg.IngressPathType)
	if vg.IngressServicePort == "" {
		vg.IngressServicePort = defaultGatewayIngressPort
	}
	if vg.IngressPathType == "" {
		vg.IngressPathType = defaultGatewayIngressPathType
	}
	if vg.HpaTargetCPUUtilization <= 0 {
		vg.HpaTargetCPUUtilization = defaultGatewayHpaTargetCPU
	}

	if v.Indexer == nil {
		v.Indexer = new(VrsIndexer)
	}
	vi := v.Indexer
	vi.AutoIndexDurationLimit = config.GetActualValue(vi.AutoIndexDurationLimit)
	vi.AutoSaveIndexDurationLimit = config.GetActualValue(vi.AutoSaveIndexDurationLimit)
	if vi.AutoIndexDurationLimit == "" {
		vi.AutoIndexDurationLimit = defaultAutoIndexDurationLimit
	}
	if vi.AutoSaveIndexDurationLimit == "" {
		vi.AutoSaveIndexDurationLimit = defaultAutoSaveIndexDurLimit
	}
}

func (o *Operator) bindNodePool() {
	if o.NodePool == nil {
		o.NodePool = new(NodePool)
	}
	o.NodePool.LabelPrefix = config.GetActualValue(o.NodePool.LabelPrefix)
}

func (o *Operator) bindPersistentVolume() {
	if o.PersistentVolume == nil {
		o.PersistentVolume = new(PersistentVolume)
	}
	pv := o.PersistentVolume
	pv.DefaultStorageClass = config.GetActualValue(pv.DefaultStorageClass)
	pv.DefaultAccessMode = config.GetActualValue(pv.DefaultAccessMode)
}

func (o *Operator) bindNetworking() {
	if o.Networking == nil {
		o.Networking = new(Networking)
	}
	nw := o.Networking
	nw.GatewayServiceType = config.GetActualValue(nw.GatewayServiceType)
	nw.DiscovererDaemonSetMaxSurge = config.GetActualValue(nw.DiscovererDaemonSetMaxSurge)
	nw.DiscovererDaemonSetMaxUnavailable = config.GetActualValue(nw.DiscovererDaemonSetMaxUnavailable)
	if nw.GatewayIngressAnnotations != nil {
		annotations := make(map[string]string, len(nw.GatewayIngressAnnotations))
		for k, v := range nw.GatewayIngressAnnotations {
			annotations[config.GetActualValue(k)] = config.GetActualValue(v)
		}
		nw.GatewayIngressAnnotations = annotations
	}
}

// Config is the runtime configuration consumed by the ValdOperatorRelease reconciler
// and the lifecycle builders. It bundles the parsed default VRS template with
// the operator settings loaded from the configuration file.
type Config struct {
	GatewayIngressAnnotations         map[string]string
	IndexerAutoIndexDurationLimit     string
	GatewayServiceType                string
	DiscovererDaemonSetMaxSurge       string
	NodePoolLabelPrefix               string
	DefaultStorageClass               string
	DefaultAccessMode                 string
	IndexerAutoSaveIndexDurationLimit string
	AgentMaxSurge                     string
	DiscovererDaemonSetMaxUnavailable string
	AgentMaxUnavailable               string
	ManagedGenerationLabel            string
	GatewayIngressServicePort         string
	GatewayIngressPathType            string
	VrsLogLevel                       string
	VrsLogFormat                      string
	VrsLogger                         string
	DefaultVrs                        DefaultVrsBundle
	GatewayHpaTargetCPUUtilization    int
	PvMinSizeBytes                    int64
	PvBufferRatio                     float64
	AgentPodsPerNode                  int
	MaxConcurrentReconciles           int
	RequeueAfterSuccess               time.Duration
	RequeueAfterError                 time.Duration
	RequeueAfterNotFound              time.Duration
	EnableIngress                     bool
	AgentEnableInMemoryMode           bool
	RequireNodePoolMatch              bool
}

// DefaultVrsBundle bundles the template ValdRelease yaml parsed into the
// unstructured form used by Build's overlay merge.
type DefaultVrsBundle struct {
	Us *k8s.Unstructured
}

// Load builds the runtime Config from the operator settings by reading and
// parsing the default VRS template yaml located at Vrs.DefaultVrsPath. Bind
// is applied first so that every omitted section falls back to its default.
func (o *Operator) Load() (*Config, error) {
	o = o.Bind()

	success, err := timeutil.Parse(o.Controller.Requeue.Success)
	if err != nil {
		return nil, errors.Wrap(err, "invalid controller.requeue.success")
	}
	onError, err := timeutil.Parse(o.Controller.Requeue.OnError)
	if err != nil {
		return nil, errors.Wrap(err, "invalid controller.requeue.on_error")
	}
	notFound, err := timeutil.Parse(o.Controller.Requeue.NotFound)
	if err != nil {
		return nil, errors.Wrap(err, "invalid controller.requeue.not_found")
	}

	c := &Config{
		AgentPodsPerNode:     o.NodePool.AgentPodsPerNode,
		RequireNodePoolMatch: o.NodePool.RequireMatch,
		NodePoolLabelPrefix:  o.NodePool.LabelPrefix,

		DefaultStorageClass: o.PersistentVolume.DefaultStorageClass,
		DefaultAccessMode:   o.PersistentVolume.DefaultAccessMode,
		PvBufferRatio:       o.PersistentVolume.BufferRatio,
		PvMinSizeBytes:      o.PersistentVolume.MinSizeBytes,

		EnableIngress:             o.Networking.EnableIngress,
		GatewayIngressAnnotations: o.Networking.GatewayIngressAnnotations,
		GatewayServiceType:        o.Networking.GatewayServiceType,

		GatewayIngressServicePort:      o.Vrs.Gateway.IngressServicePort,
		GatewayIngressPathType:         o.Vrs.Gateway.IngressPathType,
		GatewayHpaTargetCPUUtilization: o.Vrs.Gateway.HpaTargetCPUUtilization,

		VrsLogLevel:            o.Vrs.LogLevel,
		VrsLogFormat:           o.Vrs.LogFormat,
		VrsLogger:              o.Vrs.Logger,
		ManagedGenerationLabel: o.Vrs.ManagedGenerationLabel,

		AgentMaxSurge:           o.Vrs.Agent.MaxSurge,
		AgentMaxUnavailable:     o.Vrs.Agent.MaxUnavailable,
		AgentEnableInMemoryMode: *o.Vrs.Agent.EnableInMemoryMode,

		IndexerAutoIndexDurationLimit:     o.Vrs.Indexer.AutoIndexDurationLimit,
		IndexerAutoSaveIndexDurationLimit: o.Vrs.Indexer.AutoSaveIndexDurationLimit,

		DiscovererDaemonSetMaxSurge:       o.Networking.DiscovererDaemonSetMaxSurge,
		DiscovererDaemonSetMaxUnavailable: o.Networking.DiscovererDaemonSetMaxUnavailable,

		MaxConcurrentReconciles: o.Controller.MaxConcurrentReconciles,
		RequeueAfterSuccess:     success,
		RequeueAfterError:       onError,
		RequeueAfterNotFound:    notFound,
	}
	if c.GatewayIngressAnnotations == nil {
		c.GatewayIngressAnnotations = map[string]string{}
	}

	raw, err := file.ReadFile(o.Vrs.DefaultVrsPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read default vrs template")
	}

	var obj map[string]any
	if err := k8s.YAMLUnmarshal(raw, &obj); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal default vrs template")
	}
	c.DefaultVrs.Us = &k8s.Unstructured{Object: obj}

	return c, nil
}

func New(path string) (cfg *Data, err error) {
	cfg = new(Data)

	if err = config.Read(path, &cfg); err != nil {
		return nil, err
	}

	if cfg != nil {
		_ = cfg.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Server != nil {
		_ = cfg.Server.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	if cfg.Observability != nil {
		_ = cfg.Observability.Bind()
	} else {
		cfg.Observability = new(config.Observability).Bind()
	}

	if cfg.Operator != nil {
		_ = cfg.Operator.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}
	return cfg, nil
}

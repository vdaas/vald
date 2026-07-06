package config

import "github.com/vdaas/vald/hack/vald-operator/internal/infrastructure/env"

var (
	InternalHostDomain                = env.GetEnv("INTERNAL_HOST_DOMAIN", "")
	DefaultVrsPath                    = env.GetEnv("DEFAULT_VRS_PATH", "/opt/mvaldrelease/config/vrs.yaml")
	AgentPodsPerNode                  = env.GetEnvInt("AGENT_PODS_PER_NODE", "2")
	RequireNodePoolMatch              = env.GetEnvBool("REQUIRE_NODEPOOL_MATCH", "false")
	NodePoolLabelPrefix               = env.GetEnv("NODEPOOL_LABEL_PREFIX", "")
	DefaultStorageClass               = env.GetEnv("DEFAULT_STORAGE_CLASS", "standard")
	DefaultAccessMode                 = env.GetEnv("DEFAULT_ACCESS_MODE", "ReadWriteOnce")
	EnableIngress                     = env.GetEnvBool("ENABLE_INGRESS", "true")
	GatewayIngressAnnotations         = env.GetEnv("GATEWAY_INGRESS_ANNOTATIONS", "")
	GatewayServiceType                = env.GetEnv("GATEWAY_SERVICE_TYPE", "NodePort")
	VrsLogLevel                       = env.GetEnv("VRS_LOG_LEVEL", "warn")
	DiscovererDaemonSetMaxSurge       = env.GetEnv("DISCOVERER_DS_MAX_SURGE", "30%")
	DiscovererDaemonSetMaxUnavailable = env.GetEnv("DISCOVERER_DS_MAX_UNAVAILABLE", "0%")
	PvBufferRatio                     = env.GetEnvFloat64("PV_BUFFER_RATIO", "1.5")
	PvMinSizeBytes                    = env.GetEnvInt64("PV_MIN_SIZE_BYTES", "1073741824")
)

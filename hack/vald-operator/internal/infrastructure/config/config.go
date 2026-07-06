package config

import (
	"fmt"
	"io"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"
)

type Config struct {
	InternalHostDomain string

	// DefaultVrs bundles the template ValdRelease yaml: source path, raw bytes,
	// and the parsed unstructured form used by Build's overlay merge.
	DefaultVrs DefaultVrsBundle

	// NGT
	AgentPodsPerNode int

	// When true, only generate VRS if matching nodepools exist in the cluster.
	RequireNodePoolMatch bool

	// Optional prefix for nodepool-related label keys (e.g., "vald.vdaas.org").
	NodePoolLabelPrefix string

	// Fallback values for Agent.PersistentVolume when the CR omits StorageClass/AccessMode.
	DefaultStorageClass string
	DefaultAccessMode   string

	// PV sizing parameters: max(memoryRequest * PvBufferRatio, PvMinSizeBytes).
	PvBufferRatio  float64
	PvMinSizeBytes int64

	// Ingress settings for gateway.
	EnableIngress             bool
	GatewayIngressAnnotations map[string]string
	GatewayServiceType        string

	// Log level passed through to the underlying vald deployment.
	VrsLogLevel string

	// Discoverer DaemonSet rolling-update strategy values (intstr percent strings).
	DiscovererDaemonSetMaxSurge       string
	DiscovererDaemonSetMaxUnavailable string
}

type DefaultVrsBundle struct {
	Path string
	Raw  []byte
	Us   *unstructured.Unstructured
}

// New constructs a Config by reading environment-derived defaults from the
// osenv package and loading the default VRS YAML file. Callers should hold the
// returned instance as the sole source of configuration values and pass it
// explicitly to consumers rather than relying on package-level globals.
func New() (*Config, error) {
	c := &Config{}
	if err := c.load(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) load() error {
	c.InternalHostDomain = InternalHostDomain
	c.DefaultVrs.Path = DefaultVrsPath

	var err error
	c.DefaultVrs.Raw, err = c.readFile(c.DefaultVrs.Path)
	if err != nil {
		return err
	}
	var obj map[string]interface{}
	if err := yaml.Unmarshal(c.DefaultVrs.Raw, &obj); err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	c.DefaultVrs.Us = &unstructured.Unstructured{Object: obj}

	c.AgentPodsPerNode = AgentPodsPerNode
	c.RequireNodePoolMatch = RequireNodePoolMatch
	c.NodePoolLabelPrefix = NodePoolLabelPrefix
	c.DefaultStorageClass = DefaultStorageClass
	c.DefaultAccessMode = DefaultAccessMode
	c.EnableIngress = EnableIngress
	if GatewayIngressAnnotations != "" {
		var m map[string]string
		if err := yaml.Unmarshal([]byte(GatewayIngressAnnotations), &m); err != nil {
			return fmt.Errorf("failed to parse GATEWAY_INGRESS_ANNOTATIONS: %w", err)
		}
		c.GatewayIngressAnnotations = m
	} else {
		c.GatewayIngressAnnotations = map[string]string{}
	}
	c.GatewayServiceType = GatewayServiceType
	c.VrsLogLevel = VrsLogLevel
	c.DiscovererDaemonSetMaxSurge = DiscovererDaemonSetMaxSurge
	c.DiscovererDaemonSetMaxUnavailable = DiscovererDaemonSetMaxUnavailable
	c.PvBufferRatio = PvBufferRatio
	c.PvMinSizeBytes = PvMinSizeBytes
	return nil
}

func (c *Config) readFile(configFile string) ([]byte, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			logf.Log.Error(err, "error closing file", "file", configFile)
		}
	}()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

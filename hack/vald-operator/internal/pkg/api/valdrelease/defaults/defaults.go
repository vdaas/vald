// +k8s:deepcopy-gen=package
// +groupName=vald.vdaas.org
package defaults

type Defaults struct {
	TimeZone      string         `yaml:"time_zone,omitempty" json:"time_zone,omitempty"`
	Logging       *Logging       `yaml:"logging,omitempty" json:"logging,omitempty"`
	ServerConfig  *ServerConfig  `yaml:"server_config,omitempty" json:"server_config,omitempty"`
	Observability *Observability `yaml:"observability,omitempty" json:"observability,omitempty"`
}

type Logging struct {
	Format string `yaml:"format,omitempty" json:"format,omitempty"`
	Level  string `yaml:"level,omitempty" json:"level,omitempty"`
	Logger string `yaml:"logger,omitempty" json:"logger,omitempty"`
}

type ServerConfig struct {
	Servers *Servers `yaml:"servers,omitempty" json:"servers,omitempty"`
}

type Servers struct {
	GRPC *GRPC `yaml:"grpc,omitempty" json:"grpc,omitempty"`
}

type GRPC struct {
	Server *GRPCServer `yaml:"server,omitempty" json:"server,omitempty"`
}

type GRPCServer struct {
	GRPC *InterceptorConfig `yaml:"grpc,omitempty" json:"grpc,omitempty"`
}

type InterceptorConfig struct {
	Interceptors []string `yaml:"interceptors,omitempty" json:"interceptors,omitempty"`
}

type Observability struct {
	Enabled bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Trace   *Trace `yaml:"trace,omitempty" json:"trace,omitempty"`
	OTLP    *OTLP  `yaml:"otlp,omitempty" json:"otlp,omitempty"`
}

type Trace struct {
	Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
}

type OTLP struct {
	CollectorEndpoint string `yaml:"collector_endpoint,omitempty" json:"collector_endpoint,omitempty"`
}

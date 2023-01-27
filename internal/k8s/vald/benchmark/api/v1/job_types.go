//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type BenchmarkJobSpec struct {
	Target       *BenchmarkTarget    `json:"target,omitempty"`
	Dataset      *BenchmarkDataset   `json:"dataset,omitempty"`
	Dimension    int                 `json:"dimension,omitempty"`
	Replica      int                 `json:"replica,omitempty"`
	Repetition   int                 `json:"repetition,omitempty"`
	JobType      string              `json:"job_type,omitempty"`
	InsertConfig *InsertConfig       `json:"insert_config,omitempty"`
	UpdateConfig *UpdateConfig       `json:"update_config,omitempty"`
	UpsertConfig *UpsertConfig       `json:"upsert_config,omitempty"`
	SearchConfig *SearchConfig       `json:"search_config,omitempty"`
	RemoveConfig *RemoveConfig       `json:"remove_config,omitempty"`
	ClientConfig *ClientConfig       `json:"client_config,omitempty"`
	Rules        []*BenchmarkJobRule `json:"rules,omitempty"`
}

type BenchmarkJobStatus string

const (
	BenchmarkJobNotReady  = BenchmarkJobStatus("NotReady")
	BenchmarkJobAvailable = BenchmarkJobStatus("Available")
	BenchmarkJobHealthy   = BenchmarkJobStatus("Healthy")
)

// BenchmarkTarget defines the desired state of BenchmarkTarget
type BenchmarkTarget struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

// BenchmarkDataset defines the desired state of BenchmarkDateset
type BenchmarkDataset struct {
	Name    string                 `json:"name,omitempty"`
	Group   string                 `json:"group,omitempty"`
	Indexes int                    `json:"indexes,omitempty"`
	Range   *BenchmarkDatasetRange `json:"range,omitempty"`
}

// BenchmarkDatasetRange defines the desired state of BenchmarkDatesetRange
type BenchmarkDatasetRange struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

// BenchmarkJobRule defines the desired state of BenchmarkJobRule
type BenchmarkJobRule struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// InsertConfig defines the desired state of insert config
type InsertConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

// UpdateConfig defines the desired state of update config
type UpdateConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

// UpsertConfig defines the desired state of upsert config
type UpsertConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

// SearchConfig defines the desired state of search config
type SearchConfig struct {
	Epsilon float32 `json:"epsilon,omitempty"`
	Radius  float32 `json:"radius,omitempty"`
	Num     int32   `json:"num,omitempty"`
	MinNum  int32   `json:"min_num,omitempty"`
	Timeout string  `json:"timeout,omitempty"`
}

// RemoveConfig defines the desired state of remove config
type RemoveConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

// ClientConfig represents the configurations for gRPC client.
type ClientConfig struct {
	Addrs               []string        `json:"addrs"                 yaml:"addrs"`
	HealthCheckDuration string          `json:"health_check_duration" yaml:"health_check_duration"`
	ConnectionPool      *ConnectionPool `json:"connection_pool"       yaml:"connection_pool"`
	Backoff             *Backoff        `json:"backoff"               yaml:"backoff"`
	CircuitBreaker      *CircuitBreaker `json:"circuit_breaker"       yaml:"circuit_breaker"`
	CallOption          *CallOption     `json:"call_option"           yaml:"call_option"`
	DialOption          *DialOption     `json:"dial_option"           yaml:"dial_option"`
	TLS                 *TLS            `json:"tls"                   yaml:"tls"`
}

// CircuitBreaker represents the configuration for the internal circuitbreaker package.
type CircuitBreaker struct {
	ClosedErrorRate      float32 `yaml:"closed_error_rate"      json:"closed_error_rate,omitempty"`
	HalfOpenErrorRate    float32 `yaml:"half_open_error_rate"   json:"half_open_error_rate,omitempty"`
	MinSamples           int64   `yaml:"min_samples"            json:"min_samples,omitempty"`
	OpenTimeout          string  `yaml:"open_timeout"           json:"open_timeout,omitempty"`
	ClosedRefreshTimeout string  `yaml:"closed_refresh_timeout" json:"closed_refresh_timeout,omitempty"`
}

// TLS represent the TLS configuration for server.
type TLS struct {
	// Enable represent the server enable TLS or not.
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Cert represent the certificate environment variable key used to start server.
	Cert string `yaml:"cert" json:"cert"`

	// Key represent the private key environment variable key used to start server.
	Key string `yaml:"key" json:"key"`

	// CA represent the CA certificate environment variable key used to start server.
	CA string `yaml:"ca" json:"ca"`

	// InsecureSkipVerify represent enable/disable skip SSL certificate verification
	InsecureSkipVerify bool `yaml:"insecure_skip_verify" json:"insecure_skip_verify"`
}

// Backoff represents the configuration for the internal backoff package.
type Backoff struct {
	InitialDuration  string  `json:"initial_duration"   yaml:"initial_duration"`
	BackoffTimeLimit string  `json:"backoff_time_limit" yaml:"backoff_time_limit"`
	MaximumDuration  string  `json:"maximum_duration"   yaml:"maximum_duration"`
	JitterLimit      string  `json:"jitter_limit"       yaml:"jitter_limit"`
	BackoffFactor    float64 `json:"backoff_factor"     yaml:"backoff_factor"`
	RetryCount       int     `json:"retry_count"        yaml:"retry_count"`
	EnableErrorLog   bool    `json:"enable_error_log"   yaml:"enable_error_log"`
}

// CallOption represents the configurations for call option.
type CallOption struct {
	WaitForReady          bool `json:"wait_for_ready"            yaml:"wait_for_ready"`
	MaxRetryRPCBufferSize int  `json:"max_retry_rpc_buffer_size" yaml:"max_retry_rpc_buffer_size"`
	MaxRecvMsgSize        int  `json:"max_recv_msg_size"         yaml:"max_recv_msg_size"`
	MaxSendMsgSize        int  `json:"max_send_msg_size"         yaml:"max_send_msg_size"`
}

// DialOption represents the configurations for dial option.
type DialOption struct {
	WriteBufferSize             int                  `json:"write_buffer_size"              yaml:"write_buffer_size"`
	ReadBufferSize              int                  `json:"read_buffer_size"               yaml:"read_buffer_size"`
	InitialWindowSize           int                  `json:"initial_window_size"            yaml:"initial_window_size"`
	InitialConnectionWindowSize int                  `json:"initial_connection_window_size" yaml:"initial_connection_window_size"`
	MaxMsgSize                  int                  `json:"max_msg_size"                   yaml:"max_msg_size"`
	BackoffMaxDelay             string               `json:"backoff_max_delay"              yaml:"backoff_max_delay"`
	BackoffBaseDelay            string               `json:"backoff_base_delay"             yaml:"backoff_base_delay"`
	BackoffJitter               float64              `json:"backoff_jitter"                 yaml:"backoff_jitter"`
	BackoffMultiplier           float64              `json:"backoff_multiplier"             yaml:"backoff_multiplier"`
	MinimumConnectionTimeout    string               `json:"min_connection_timeout"         yaml:"min_connection_timeout"`
	EnableBackoff               bool                 `json:"enable_backoff"                 yaml:"enable_backoff"`
	Insecure                    bool                 `json:"insecure"                       yaml:"insecure"`
	Timeout                     string               `json:"timeout"                        yaml:"timeout"`
	Interceptors                []string             `json:"interceptors,omitempty"         yaml:"interceptors"`
	Net                         *Net                 `json:"net"                            yaml:"net"`
	Keepalive                   *GRPCClientKeepalive `json:"keepalive"                      yaml:"keepalive"`
}

// ConnectionPool represents the configurations for connection pool.
type ConnectionPool struct {
	ResolveDNS           bool   `json:"enable_dns_resolver"     yaml:"enable_dns_resolver"`
	EnableRebalance      bool   `json:"enable_rebalance"        yaml:"enable_rebalance"`
	RebalanceDuration    string `json:"rebalance_duration"      yaml:"rebalance_duration"`
	Size                 int    `json:"size"                    yaml:"size"`
	OldConnCloseDuration string `json:"old_conn_close_duration" yaml:"old_conn_close_duration"`
}

// GRPCClientKeepalive represents the configurations for gRPC keep-alive.
type GRPCClientKeepalive struct {
	Time                string `json:"time"                  yaml:"time"`
	Timeout             string `json:"timeout"               yaml:"timeout"`
	PermitWithoutStream bool   `json:"permit_without_stream" yaml:"permit_without_stream"`
}

// Net represents the network configuration tcp, udp, unix domain socket.
type Net struct {
	DNS          *DNS          `yaml:"dns"           json:"dns,omitempty"`
	Dialer       *Dialer       `yaml:"dialer"        json:"dialer,omitempty"`
	SocketOption *SocketOption `yaml:"socket_option" json:"socket_option,omitempty"`
	TLS          *TLS          `yaml:"tls"           json:"tls,omitempty"`
}

// Dialer represents the configuration for dial.
type Dialer struct {
	Timeout          string `yaml:"timeout"            json:"timeout,omitempty"`
	Keepalive        string `yaml:"keepalive"          json:"keepalive,omitempty"`
	FallbackDelay    string `yaml:"fallback_delay"     json:"fallback_delay,omitempty"`
	DualStackEnabled bool   `yaml:"dual_stack_enabled" json:"dual_stack_enabled,omitempty"`
}

// DNS represents the configuration for resolving DNS.
type DNS struct {
	CacheEnabled    bool   `yaml:"cache_enabled"    json:"cache_enabled,omitempty"`
	RefreshDuration string `yaml:"refresh_duration" json:"refresh_duration,omitempty"`
	CacheExpiration string `yaml:"cache_expiration" json:"cache_expiration,omitempty"`
}

// SocketOption represents the socket configurations.
type SocketOption struct {
	ReusePort                bool `json:"reuse_port,omitempty"                  yaml:"reuse_port"`
	ReuseAddr                bool `json:"reuse_addr,omitempty"                  yaml:"reuse_addr"`
	TCPFastOpen              bool `json:"tcp_fast_open,omitempty"               yaml:"tcp_fast_open"`
	TCPNoDelay               bool `json:"tcp_no_delay,omitempty"                yaml:"tcp_no_delay"`
	TCPCork                  bool `json:"tcp_cork,omitempty"                    yaml:"tcp_cork"`
	TCPQuickAck              bool `json:"tcp_quick_ack,omitempty"               yaml:"tcp_quick_ack"`
	TCPDeferAccept           bool `json:"tcp_defer_accept,omitempty"            yaml:"tcp_defer_accept"`
	IPTransparent            bool `json:"ip_transparent,omitempty"              yaml:"ip_transparent"`
	IPRecoverDestinationAddr bool `json:"ip_recover_destination_addr,omitempty" yaml:"ip_recover_destination_addr"`
}

type ValdBenchmarkJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   BenchmarkJobSpec   `json:"spec,omitempty"`
	Status BenchmarkJobStatus `json:"status,omitempty"`
}

type ValdBenchmarkJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []ValdBenchmarkJob `json:"items,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkDataset) DeepCopyInto(out *BenchmarkDataset) {
	*out = *in
	if in.Range != nil {
		in, out := &in.Range, &out.Range
		*out = new(BenchmarkDatasetRange)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkDataset.
func (in *BenchmarkDataset) DeepCopy() *BenchmarkDataset {
	if in == nil {
		return nil
	}
	out := new(BenchmarkDataset)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkDatasetRange) DeepCopyInto(out *BenchmarkDatasetRange) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkDatasetRange.
func (in *BenchmarkDatasetRange) DeepCopy() *BenchmarkDatasetRange {
	if in == nil {
		return nil
	}
	out := new(BenchmarkDatasetRange)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkJobRule) DeepCopyInto(out *BenchmarkJobRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkJobRule.
func (in *BenchmarkJobRule) DeepCopy() *BenchmarkJobRule {
	if in == nil {
		return nil
	}
	out := new(BenchmarkJobRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkJobSpec) DeepCopyInto(out *BenchmarkJobSpec) {
	*out = *in
	if in.Target != nil {
		in, out := &in.Target, &out.Target
		*out = new(BenchmarkTarget)
		**out = **in
	}
	if in.Dataset != nil {
		in, out := &in.Dataset, &out.Dataset
		*out = new(BenchmarkDataset)
		(*in).DeepCopyInto(*out)
	}
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]*BenchmarkJobRule, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(BenchmarkJobRule)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkJobSpec.
func (in *BenchmarkJobSpec) DeepCopy() *BenchmarkJobSpec {
	if in == nil {
		return nil
	}
	out := new(BenchmarkJobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkTarget) DeepCopyInto(out *BenchmarkTarget) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkTarget.
func (in *BenchmarkTarget) DeepCopy() *BenchmarkTarget {
	if in == nil {
		return nil
	}
	out := new(BenchmarkTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValdBenchmarkJob) DeepCopyInto(out *ValdBenchmarkJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkOperator.
func (in *ValdBenchmarkJob) DeepCopy() *ValdBenchmarkJob {
	if in == nil {
		return nil
	}
	out := new(ValdBenchmarkJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValdBenchmarkJob) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValdBenchmarkJobList) DeepCopyInto(out *ValdBenchmarkJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ValdBenchmarkJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkOperatorList.
func (in *ValdBenchmarkJobList) DeepCopy() *ValdBenchmarkJobList {
	if in == nil {
		return nil
	}
	out := new(ValdBenchmarkJobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValdBenchmarkJobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

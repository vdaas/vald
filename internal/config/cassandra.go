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

package config

import (
	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/tls"
)

// Cassandra represents the configuration for the internal cassandra package.
type Cassandra struct {
	// PoolConfig represents the pool configuration.
	PoolConfig *PoolConfig `json:"pool_config" yaml:"pool_config"`
	// Net represents the network configuration.
	Net *Net `json:"net" yaml:"net"`
	// TLS represents the TLS configuration.
	TLS *TLS `json:"tls" yaml:"tls"`
	// HostFilter represents the host filter configuration.
	HostFilter *HostFilter `json:"host_filter" yaml:"host_filter"`
	// ReconnectionPolicy represents the reconnection policy.
	ReconnectionPolicy *ReconnectionPolicy `json:"reconnection_policy" yaml:"reconnection_policy"`
	// RetryPolicy represents the retry policy.
	RetryPolicy *RetryPolicy `json:"retry_policy" yaml:"retry_policy"`
	// VKTable represents the Vector-Key table name.
	VKTable string `json:"vk_table" yaml:"vk_table"`
	// ReconnectInterval represents the reconnect interval.
	ReconnectInterval string `json:"reconnect_interval" yaml:"reconnect_interval"`
	// Consistency represents the consistency level.
	Consistency string `json:"consistency" yaml:"consistency"`
	// SerialConsistency represents the serial consistency level.
	SerialConsistency string `json:"serial_consistency" yaml:"serial_consistency"`
	// Username represents the username.
	Username string `json:"username" yaml:"username"`
	// Password represents the password.
	Password string `json:"password" yaml:"password"`
	// Keyspace represents the keyspace name.
	Keyspace string `json:"keyspace" yaml:"keyspace"`
	// KVTable represents the Key-Value table name.
	KVTable string `json:"kv_table" yaml:"kv_table"`
	// ConnectTimeout represents the connect timeout.
	ConnectTimeout string `json:"connect_timeout" yaml:"connect_timeout"`
	// Timeout represents the request timeout.
	Timeout string `json:"timeout" yaml:"timeout"`
	// SocketKeepalive represents the socket keepalive duration.
	SocketKeepalive string `json:"socket_keepalive" yaml:"socket_keepalive"`
	// WriteCoalesceWaitTime represents the write coalesce wait time.
	WriteCoalesceWaitTime string `json:"write_coalesce_wait_time" yaml:"write_coalesce_wait_time"`
	// CQLVersion represents the CQL version.
	CQLVersion string `json:"cql_version" yaml:"cql_version"`
	// MaxWaitSchemaAgreement represents the max wait schema agreement.
	MaxWaitSchemaAgreement string `json:"max_wait_schema_agreement" yaml:"max_wait_schema_agreement"`
	// VectorBackupTable represents the vector backup table name.
	VectorBackupTable string `json:"vector_backup_table" yaml:"vector_backup_table"`
	// Hosts represents the list of hosts.
	Hosts []string `json:"hosts" yaml:"hosts"`
	// MaxRoutingKeyInfo represents the max routing key info.
	MaxRoutingKeyInfo int `json:"max_routing_key_info" yaml:"max_routing_key_info"`
	// MaxPreparedStmts represents the max prepared statements.
	MaxPreparedStmts int `json:"max_prepared_stmts" yaml:"max_prepared_stmts"`
	// ProtoVersion represents the protocol version.
	ProtoVersion int `json:"proto_version" yaml:"proto_version"`
	// PageSize represents the page size.
	PageSize int `json:"page_size" yaml:"page_size"`
	// Port represents the port number.
	Port int `json:"port" yaml:"port"`
	// NumConns represents the number of connections.
	NumConns int `json:"num_conns" yaml:"num_conns"`
	// DisableTopologyEvents enables disabling topology events.
	DisableTopologyEvents bool `json:"disable_topology_events" yaml:"disable_topology_events"`
	// DefaultTimestamp enables default timestamp.
	DefaultTimestamp bool `json:"default_timestamp" yaml:"default_timestamp"`
	// DisableSchemaEvents enables disabling schema events.
	DisableSchemaEvents bool `json:"disable_schema_events" yaml:"disable_schema_events"`
	// DisableSkipMetadata enables disabling skip metadata.
	DisableSkipMetadata bool `json:"disable_skip_metadata" yaml:"disable_skip_metadata"`
	// DefaultIdempotence enables default idempotence.
	DefaultIdempotence bool `json:"default_idempotence" yaml:"default_idempotence"`
	// DisableNodeStatusEvents enables disabling node status events.
	DisableNodeStatusEvents bool `json:"disable_node_status_events" yaml:"disable_node_status_events"`
	// DisableInitialHostLookup enables disabling initial host lookup.
	DisableInitialHostLookup bool `json:"disable_initial_host_lookup" yaml:"disable_initial_host_lookup"`
	// IgnorePeerAddr enables ignoring peer address.
	IgnorePeerAddr bool `json:"ignore_peer_addr" yaml:"ignore_peer_addr"`
	// EnableHostVerification enables host verification.
	EnableHostVerification bool `json:"enable_host_verification" yaml:"enable_host_verification"`
}

// PoolConfig represents the configuration for the pool config.
type PoolConfig struct {
	// DataCenter represents the data center name.
	DataCenter string `json:"data_center" yaml:"data_center"`
	// DCAwareRouting enables DC aware routing.
	DCAwareRouting bool `json:"dc_aware_routing" yaml:"dc_aware_routing"`
	// NonLocalReplicasFallback enables non local replicas fallback.
	NonLocalReplicasFallback bool `json:"non_local_replicas_fallback" yaml:"non_local_replicas_fallback"`
	// ShuffleReplicas enables shuffle replicas.
	ShuffleReplicas bool `json:"shuffle_replicas" yaml:"shuffle_replicas"`
	// TokenAwareHostPolicy enables token aware host policy.
	TokenAwareHostPolicy bool `json:"token_aware_host_policy" yaml:"token_aware_host_policy"`
}

// Bind binds the actual data from the PoolConfig receiver fields.
func (pc *PoolConfig) Bind() *PoolConfig {
	pc.DataCenter = GetActualValue(pc.DataCenter)
	return pc
}

// RetryPolicy represents the configuration for the retry policy.
type RetryPolicy struct {
	// MinDuration represents the minimum duration.
	MinDuration string `json:"min_duration" yaml:"min_duration"`
	// MaxDuration represents the maximum duration.
	MaxDuration string `json:"max_duration" yaml:"max_duration"`
	// NumRetries represents the number of retries.
	NumRetries int `json:"num_retries" yaml:"num_retries"`
}

// Bind binds the actual data from the RetryPolicy receiver fields.
func (rp *RetryPolicy) Bind() *RetryPolicy {
	rp.MinDuration = GetActualValue(rp.MinDuration)
	rp.MaxDuration = GetActualValue(rp.MaxDuration)
	return rp
}

// ReconnectionPolicy represents the configuration for the reconnection policy.
type ReconnectionPolicy struct {
	// InitialInterval represents the initial interval.
	InitialInterval string `json:"initial_interval" yaml:"initial_interval"`
	// MaxRetries represents the maximum retry count.
	MaxRetries int `json:"max_retries" yaml:"max_retries"`
}

// Bind binds the actual data from the ReconnectionPolicy receiver fields.
func (rp *ReconnectionPolicy) Bind() *ReconnectionPolicy {
	rp.InitialInterval = GetActualValue(rp.InitialInterval)
	return rp
}

// HostFilter represents the configuration for the host filter.
type HostFilter struct {
	// DataCenter represents the data center name.
	DataCenter string `json:"data_center" yaml:"data_center"`
	// WhiteList represents the white list.
	WhiteList []string `json:"white_list" yaml:"white_list"`
	// Enabled enables host filter.
	Enabled bool `json:"enabled"`
}

// Bind binds the actual data from the HostFilter receiver fields.
func (hf *HostFilter) Bind() *HostFilter {
	hf.DataCenter = GetActualValue(hf.DataCenter)
	hf.WhiteList = GetActualValues(hf.WhiteList)
	return hf
}

// Bind binds the actual data from the Cassandra receiver fields.
func (c *Cassandra) Bind() *Cassandra {
	c.Hosts = GetActualValues(c.Hosts)
	c.CQLVersion = GetActualValue(c.CQLVersion)
	c.Timeout = GetActualValue(c.Timeout)
	c.ConnectTimeout = GetActualValue(c.ConnectTimeout)
	c.Keyspace = GetActualValue(c.Keyspace)
	c.Consistency = GetActualValue(c.Consistency)
	c.SerialConsistency = GetActualValue(c.SerialConsistency)
	c.Username = GetActualValue(c.Username)
	c.Password = GetActualValue(c.Password)

	if c.PoolConfig == nil {
		c.PoolConfig = new(PoolConfig)
	}
	c.PoolConfig.Bind()

	if c.RetryPolicy == nil {
		c.RetryPolicy = new(RetryPolicy)
	}
	c.RetryPolicy.Bind()

	if c.ReconnectionPolicy == nil {
		c.ReconnectionPolicy = new(ReconnectionPolicy)
	}
	c.ReconnectionPolicy.Bind()

	if c.HostFilter == nil {
		c.HostFilter = new(HostFilter)
	}
	c.HostFilter.Bind()

	c.SocketKeepalive = GetActualValue(c.SocketKeepalive)

	if c.TLS == nil {
		c.TLS = new(TLS)
	}
	c.TLS.Bind()

	if c.Net == nil {
		c.Net = new(Net)
	}
	c.Net.Bind()

	c.ReconnectInterval = GetActualValue(c.ReconnectInterval)
	c.MaxWaitSchemaAgreement = GetActualValue(c.MaxWaitSchemaAgreement)
	c.WriteCoalesceWaitTime = GetActualValue(c.WriteCoalesceWaitTime)

	c.KVTable = GetActualValue(c.KVTable)
	c.VKTable = GetActualValue(c.VKTable)

	c.VectorBackupTable = GetActualValue(c.VectorBackupTable)
	return c
}

// Opts creates and returns the slice with the functional options for the internal cassandra package.
// In addition, Opts sometimes returns the error when the any errors are occurred.
func (cfg *Cassandra) Opts() (opts []cassandra.Option, err error) {
	opts = []cassandra.Option{
		cassandra.WithHosts(cfg.Hosts...),
		cassandra.WithCQLVersion(cfg.CQLVersion),
		cassandra.WithProtoVersion(cfg.ProtoVersion),
		cassandra.WithTimeout(cfg.Timeout),
		cassandra.WithConnectTimeout(cfg.ConnectTimeout),
		cassandra.WithPort(cfg.Port),
		cassandra.WithKeyspace(cfg.Keyspace),
		cassandra.WithNumConns(cfg.NumConns),
		cassandra.WithConsistency(cfg.Consistency),
		cassandra.WithSerialConsistency(cfg.SerialConsistency),
		cassandra.WithUsername(cfg.Username),
		cassandra.WithPassword(cfg.Password),
		cassandra.WithSocketKeepalive(cfg.SocketKeepalive),
		cassandra.WithMaxPreparedStmts(cfg.MaxPreparedStmts),
		cassandra.WithMaxRoutingKeyInfo(cfg.MaxRoutingKeyInfo),
		cassandra.WithPageSize(cfg.PageSize),
		cassandra.WithEnableHostVerification(cfg.EnableHostVerification),
		cassandra.WithDefaultTimestamp(cfg.DefaultTimestamp),
		cassandra.WithReconnectInterval(cfg.ReconnectInterval),
		cassandra.WithMaxWaitSchemaAgreement(cfg.MaxWaitSchemaAgreement),
		cassandra.WithIgnorePeerAddr(cfg.IgnorePeerAddr),
		cassandra.WithDisableInitialHostLookup(cfg.DisableInitialHostLookup),
		cassandra.WithDisableNodeStatusEvents(cfg.DisableNodeStatusEvents),
		cassandra.WithDisableTopologyEvents(cfg.DisableTopologyEvents),
		cassandra.WithDisableSkipMetadata(cfg.DisableSkipMetadata),
		cassandra.WithDefaultIdempotence(cfg.DefaultIdempotence),
		cassandra.WithWriteCoalesceWaitTime(cfg.WriteCoalesceWaitTime),
	}
	if cfg.RetryPolicy != nil {
		opts = append(
			opts,
			cassandra.WithRetryPolicyNumRetries(cfg.RetryPolicy.NumRetries),
			cassandra.WithRetryPolicyMinDuration(cfg.RetryPolicy.MinDuration),
			cassandra.WithRetryPolicyMaxDuration(cfg.RetryPolicy.MaxDuration),
		)
	}
	if cfg.ReconnectionPolicy != nil {
		opts = append(
			opts,
			cassandra.WithReconnectionPolicyMaxRetries(cfg.ReconnectionPolicy.MaxRetries),
			cassandra.WithReconnectionPolicyInitialInterval(cfg.ReconnectionPolicy.InitialInterval),
		)
	}
	if cfg.PoolConfig != nil {
		opts = append(
			opts,
			cassandra.WithDC(cfg.PoolConfig.DataCenter),
			cassandra.WithDCAwareRouting(cfg.PoolConfig.DCAwareRouting),
			cassandra.WithNonLocalReplicasFallback(cfg.PoolConfig.NonLocalReplicasFallback),
			cassandra.WithShuffleReplicas(cfg.PoolConfig.ShuffleReplicas),
			cassandra.WithTokenAwareHostPolicy(cfg.PoolConfig.TokenAwareHostPolicy),
		)
	}
	if cfg.HostFilter != nil {
		opts = append(
			opts,
			cassandra.WithHostFilter(cfg.HostFilter.Enabled),
			cassandra.WithDCHostFilter(cfg.HostFilter.DataCenter),
			cassandra.WithWhiteListHostFilter(cfg.HostFilter.WhiteList),
		)
	}

	if cfg.Net != nil {
		netOpts, err := cfg.Net.Opts()
		if err != nil {
			return nil, err
		}
		der, err := net.NewDialer(netOpts...)
		if err != nil {
			return nil, err
		}

		opts = append(opts,
			cassandra.WithDialer(
				der,
			),
		)
	}

	if cfg.TLS != nil && cfg.TLS.Enabled {
		tcfg, err := tls.NewClientConfig(cfg.TLS.Opts()...)
		if err != nil {
			return nil, err
		}

		opts = append(
			opts,
			cassandra.WithTLS(tcfg),
			cassandra.WithTLSCertPath(cfg.TLS.Cert),
			cassandra.WithTLSKeyPath(cfg.TLS.Key),
			cassandra.WithTLSCAPath(cfg.TLS.CA),
		)
	}
	return opts, nil
}

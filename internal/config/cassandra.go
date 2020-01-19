//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package config providers configuration type and load configuration logic
package config

import (
	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/tls"
)

type Cassandra struct {
	Hosts          []string `json:"hosts" yaml:"hosts"`
	CQLVersion     string   `json:"cql_version" yaml:"cql_version"`
	ProtoVersion   int      `json:"proto_version" yaml:"proto_version"`
	Timeout        string   `json:"timeout" yaml:"timeout"`
	ConnectTimeout string   `json:"connect_timeout" yaml:"connect_timeout"`
	Port           int      `json:"port" yaml:"port"`
	Keyspace       string   `json:"keyspace" yaml:"keyspace"`
	NumConns       int      `json:"num_conns" yaml:"num_conns"`
	Consistency    string   `json:"consistency" yaml:"consistency"`

	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	PoolConfig         *PoolConfig         `json:"pool_config" yaml:"pool_config"`
	RetryPolicy        *RetryPolicy        `json:"retry_policy" yaml:"retry_policy"`
	ReconnectionPolicy *ReconnectionPolicy `json:"reconnection_policy" yaml:"reconnection_policy"`

	SocketKeepalive          string `json:"socket_keepalive" yaml:"socket_keepalive"`
	MaxPreparedStmts         int    `json:"max_prepared_stmts" yaml:"max_prepared_stmts"`
	MaxRoutingKeyInfo        int    `json:"max_routing_key_info" yaml:"max_routing_key_info"`
	PageSize                 int    `json:"page_size" yaml:"page_size"`
	TLS                      *TLS   `json:"tls" yaml:"tls"`
	EnableHostVerification   bool   `json:"enable_host_verification" yaml:"enable_host_verification"`
	DefaultTimestamp         bool   `json:"default_timestamp" yaml:"default_timestamp"`
	ReconnectInterval        string `json:"reconnect_interval" yaml:"reconnect_interval"`
	MaxWaitSchemaAgreement   string `json:"max_wait_schema_agreement" yaml:"max_wait_schema_agreement"`
	IgnorePeerAddr           bool   `json:"ignore_peer_addr" yaml:"ignore_peer_addr"`
	DisableInitialHostLookup bool   `json:"disable_initial_host_lookup" yaml:"disable_initial_host_lookup"`
	DisableNodeStatusEvents  bool   `json:"disable_node_status_events" yaml:"disable_node_status_events"`
	DisableTopologyEvents    bool   `json:"disable_topology_events" yaml:"disable_topology_events"`
	DisableSchemaEvents      bool   `json:"disable_schema_events" yaml:"disable_schema_events"`
	DisableSkipMetadata      bool   `json:"disable_skip_metadata" yaml:"disable_skip_metadata"`
	DefaultIdempotence       bool   `json:"default_idempotence" yaml:"default_idempotence"`
	WriteCoalesceWaitTime    string `json:"write_coalesce_wait_time" yaml:"write_coalesce_wait_time"`

	// meta
	KVTable string `json:"kv_table" yaml:"kv_table"`
	VKTable string `json:"vk_table" yaml:"vk_table"`

	// backup manager
	MetaTable string `json:"meta_table" yaml:"meta_table"`
}

type PoolConfig struct {
	DataCenter               string `json:"data_center" yaml:"data_center"`
	DCAwareRouting           bool   `json:"dc_aware_routing" yaml:"dc_aware_routing"`
	NonLocalReplicasFallback bool   `json:"non_local_replicas_fallback" yaml:"non_local_replicas_fallback"`
	ShuffleReplicas          bool   `json:"shuffle_replicas" yaml:"shuffle_replicas"`
}

type RetryPolicy struct {
	NumRetries  int    `json:"num_retries" yaml:"num_retries"`
	MinDuration string `json:"min_duration" yaml:"min_duration"`
	MaxDuration string `json:"max_duration" yaml:"max_duration"`
}

type ReconnectionPolicy struct {
	MaxRetries      int    `json:"max_retries" yaml:"max_retries"`
	InitialInterval string `json:"initial_interval" yaml:"initial_interval"`
}

func (c *Cassandra) Bind() *Cassandra {
	c.Hosts = GetActualValues(c.Hosts)
	c.CQLVersion = GetActualValue(c.CQLVersion)
	c.Timeout = GetActualValue(c.Timeout)
	c.ConnectTimeout = GetActualValue(c.ConnectTimeout)
	c.Keyspace = GetActualValue(c.Keyspace)
	c.Consistency = GetActualValue(c.Consistency)

	c.Username = GetActualValue(c.Username)
	c.Password = GetActualValue(c.Password)

	if c.RetryPolicy != nil {
		c.RetryPolicy.MinDuration = GetActualValue(c.RetryPolicy.MinDuration)
		c.RetryPolicy.MaxDuration = GetActualValue(c.RetryPolicy.MaxDuration)
	}
	if c.ReconnectionPolicy != nil {
		c.ReconnectionPolicy.InitialInterval = GetActualValue(c.ReconnectionPolicy.InitialInterval)
	}
	if c.PoolConfig != nil {
		c.PoolConfig.DataCenter = GetActualValue(c.PoolConfig.DataCenter)
	}
	c.SocketKeepalive = GetActualValue(c.SocketKeepalive)
	if c.TLS != nil {
		c.TLS.Bind()
	} else {
		c.TLS = new(TLS)
	}
	c.ReconnectInterval = GetActualValue(c.ReconnectInterval)
	c.MaxWaitSchemaAgreement = GetActualValue(c.MaxWaitSchemaAgreement)
	c.WriteCoalesceWaitTime = GetActualValue(c.WriteCoalesceWaitTime)

	c.KVTable = GetActualValue(c.KVTable)
	c.VKTable = GetActualValue(c.VKTable)

	c.MetaTable = GetActualValue(c.MetaTable)

	return c
}

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
		cassandra.WithUsername(cfg.Username),
		cassandra.WithPassword(cfg.Password),
		cassandra.WithRetryPolicyNumRetries(cfg.RetryPolicy.NumRetries),
		cassandra.WithRetryPolicyMinDuration(cfg.RetryPolicy.MinDuration),
		cassandra.WithRetryPolicyMaxDuration(cfg.RetryPolicy.MaxDuration),
		cassandra.WithReconnectionPolicyMaxRetries(cfg.ReconnectionPolicy.MaxRetries),
		cassandra.WithReconnectionPolicyInitialInterval(cfg.ReconnectionPolicy.InitialInterval),
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
		cassandra.WithKVTable(cfg.KVTable),
		cassandra.WithVKTable(cfg.VKTable),
	}

	if cfg.PoolConfig != nil {
		opts = append(opts, cassandra.WithDC(cfg.PoolConfig.DataCenter))
		if cfg.PoolConfig.DCAwareRouting {
			opts = append(opts, cassandra.WithEnableDCAwareRouting())
		}
		if cfg.PoolConfig.NonLocalReplicasFallback {
			opts = append(opts, cassandra.WithEnableNonLocalReplicasFallback())
		}
		if cfg.PoolConfig.ShuffleReplicas {
			opts = append(opts, cassandra.WithEnableShuffleReplicas())
		}
	}

	if cfg.TLS != nil && cfg.TLS.Enabled {
		tcfg, err := tls.New(
			tls.WithCert(cfg.TLS.Cert),
			tls.WithKey(cfg.TLS.Key),
			tls.WithCa(cfg.TLS.CA),
		)
		if err != nil {
			return opts, err
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

//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package config providers configuration type and load configuration logic
package config

type Cassandra struct {
	Hosts                  []string `json:"hosts" yaml:"hosts"`
	CQLVersion             string   `json:"cql_version" yaml:"cql_version"`
	Timeout                string   `json:"timeout" yaml:"timeout"`
	ConnectTimeout         string   `json:"connect_timeout" yaml:"connect_timeout"`
	Port                   int      `json:"port" yaml:"port"`
	NumConns               int      `json:"num_conns" yaml:"num_conns"`
	Consistency            string   `json:"consistency" yaml:"consistency"`
	MaxPreparedStmts       int      `json:"max_prepared_stmts" yaml:"max_prepared_stmts"`
	MaxRoutingKeyInfo      int      `json:"max_routing_key_info" yaml:"max_routing_key_info"`
	PageSize               int      `json:"page_size" yaml:"page_size"`
	DefaultTimestamp       bool     `json:"default_timestamp" yaml:"default_timestamp"`
	MaxWaitSchemaAgreement string   `json:"max_wait_schema_agreement" yaml:"max_wait_schema_agreement"`
	ReconnectInterval      string   `json:"reconnect_interval" yaml:"reconnect_interval"`
	// ConvictionPolicy:       &SimpleConvictionPolicy{},
	ReconnectionPolicy    *ReconnectionPolicy `json:"reconnection_policy" yaml:"reconnection_policy"`
	WriteCoalesceWaitTime string              `json:"write_coalesce_wait_time" yaml:"write_coalesce_wait_time"`

	Keyspace string `json:"keyspace" yaml:"keyspace"`
	KVTable  string `json:"kv_table" yaml:"kv_table"`
	VKTable  string `json:"vk_table" yaml:"vk_table"`

	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type ReconnectionPolicy struct {
	MaxRetries      int    `json:"max_retries" yaml:"max_retries"`
	InitialInterval string `json:"initial_interval" yaml:"initial_interval"`
}

func (c *Cassandra) Bind() *Cassandra {
	for i, addr := range c.Hosts {
		c.Hosts[i] = GetActualValue(addr)
	}
	c.CQLVersion = GetActualValue(c.CQLVersion)
	c.Timeout = GetActualValue(c.Timeout)
	c.ConnectTimeout = GetActualValue(c.ConnectTimeout)
	c.Consistency = GetActualValue(c.Consistency)
	c.MaxWaitSchemaAgreement = GetActualValue(c.MaxWaitSchemaAgreement)
	c.ReconnectInterval = GetActualValue(c.ReconnectInterval)
	c.WriteCoalesceWaitTime = GetActualValue(c.WriteCoalesceWaitTime)
	if c.ReconnectionPolicy != nil {
		c.ReconnectionPolicy.InitialInterval = GetActualValue(c.ReconnectionPolicy.InitialInterval)
	}

	c.Keyspace = GetActualValue(c.Keyspace)
	c.KVTable = GetActualValue(c.KVTable)
	c.VKTable = GetActualValue(c.VKTable)

	c.Username = GetActualValue(c.Username)
	c.Password = GetActualValue(c.Password)

	return c
}

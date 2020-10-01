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

// Package redis provides implementation of Go API for redis interface
package cassandra

import (
	"crypto/tls"
	"math"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/timeutil"
)

// Option represents the functional option for client.
// It wraps the gocql.ClusterConfig to the function option implementation.
// Please refer to the following link for more information.
// https://pkg.go.dev/github.com/gocql/gocql?tab=doc#ClusterConfig
type Option func(*client) error

var (
	defaultOpts = []Option{
		WithCQLVersion("3.0.0"),
		WithConnectTimeout("600ms"),
		WithConsistency(cQuorumKey),
		WithDCAwareRouting(false),
		WithDefaultIdempotence(false),
		WithDefaultTimestamp(true),
		WithDisableInitialHostLookup(false),
		WithDisableNodeStatusEvents(false),
		WithDisableSkipMetadata(false),
		WithDisableTopologyEvents(false),
		WithEnableHostVerification(false),
		WithIgnorePeerAddr(false),
		WithMaxPreparedStmts(1000),
		WithMaxRoutingKeyInfo(1000),
		WithMaxWaitSchemaAgreement("1m"),
		WithNonLocalReplicasFallback(false),
		WithNumConns(2),
		WithPageSize(5000),
		WithPort(9042),
		WithProtoVersion(0),
		WithReconnectInterval("1m"),
		WithSerialConsistency(scLocalSerialKey),
		WithShuffleReplicas(false),
		WithTimeout("600ms"),
		WithTokenAwareHostPolicy(true),
		WithWriteCoalesceWaitTime("200µs"),
	}
)

// WithHosts returns the option to set the hosts
func WithHosts(hosts ...string) Option {
	return func(c *client) error {
		if len(hosts) == 0 {
			return errors.NewErrInvalidOption("hosts", hosts)
		}
		if c.hosts == nil {
			c.hosts = hosts
		} else {
			c.hosts = append(c.hosts, hosts...)
		}
		return nil
	}
}

// WithDialer returns the option to set the dialer
func WithDialer(der gocql.Dialer) Option {
	return func(c *client) error {
		if der == nil {
			return errors.NewErrInvalidOption("dialer", der)
		}
		c.dialer = der
		return nil
	}
}

// WithCQLVersion returns the option to set the CQL version
func WithCQLVersion(version string) Option {
	return func(c *client) error {
		if len(version) == 0 {
			return errors.NewErrInvalidOption("cqlVersion", version)
		}
		c.cqlVersion = version
		return nil
	}
}

// WithProtoVersion returns the option to set the proto version
func WithProtoVersion(version int) Option {
	return func(c *client) error {
		if version < 0 {
			return errors.NewErrInvalidOption("protoVersion", version)
		}
		c.protoVersion = version
		return nil
	}
}

// WithTimeout returns the option to set the cassandra connect timeout time
func WithTimeout(dur string) Option {
	return func(c *client) error {
		if len(dur) == 0 {
			return errors.NewErrInvalidOption("timeout", dur)
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute // FIXME
		}
		c.timeout = d
		return nil
	}
}

// WithConnectTimeout returns the option to set the cassandra initial connection timeout
func WithConnectTimeout(dur string) Option {
	return func(c *client) error {
		if len(dur) == 0 {
			return errors.NewErrInvalidOption("connectTimeout", dur)
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return errors.NewErrCriticalOption("connectTimeout", dur, err)
		}

		c.connectTimeout = d
		return nil
	}
}

// WithPort returns the option to set the port number
func WithPort(port int) Option {
	return func(c *client) error {
		if port <= 0 || port > math.MaxUint16 {
			return errors.NewErrInvalidOption("port", port)
		}
		c.port = port
		return nil
	}
}

// WithKeyspace returns the option to set the keyspace
func WithKeyspace(keyspace string) Option {
	return func(c *client) error {
		if len(keyspace) == 0 {
			return errors.NewErrInvalidOption("keyspace", keyspace)
		}
		c.keyspace = keyspace
		return nil
	}
}

// WithNumConns returns the option to set the number of connection per host
func WithNumConns(numConns int) Option {
	return func(c *client) error {
		if numConns < 0 {
			return errors.NewErrInvalidOption("numConns", numConns)
		}
		c.numConns = numConns
		return nil
	}
}

var (
	cAnyKey         = "any"
	cOneKey         = "one"
	cTwoKey         = "two"
	cThreeKey       = "three"
	cAllKey         = "all"
	cQuorumKey      = "quorum"
	cLocalQuorumKey = "localquorum"
	cEachQuorumKey  = "eachquorum"
	cLocalOneKey    = "localone"

	consistenciesMap = map[string]gocql.Consistency{
		cAnyKey:         gocql.Any,
		cOneKey:         gocql.One,
		cTwoKey:         gocql.Two,
		cThreeKey:       gocql.Three,
		cQuorumKey:      gocql.Quorum,
		cAllKey:         gocql.All,
		cLocalQuorumKey: gocql.LocalQuorum,
		cEachQuorumKey:  gocql.EachQuorum,
		cLocalOneKey:    gocql.LocalOne,
	}
)

// WithConsistency returns the option to set the cassandra consistency level
func WithConsistency(consistency string) Option {
	return func(c *client) error {
		if len(consistency) == 0 {
			return errors.NewErrInvalidOption("consistency", consistency)
		}
		actual, ok := consistenciesMap[strings.TrimSpace(strings.Trim(strings.Trim(strings.ToLower(consistency), "_"), "-"))]
		if !ok {
			return errors.NewErrCriticalOption("consistency", consistency)
		}
		c.consistency = actual
		return nil
	}
}

var (
	scLocalSerialKey       = "localserial"
	scSerialKey            = "serial"
	serialConsistenciesMap = map[string]gocql.SerialConsistency{
		scLocalSerialKey: gocql.LocalSerial,
		scSerialKey:      gocql.Serial,
	}
)

// WithSerialConsistency returns the option to set the cassandra serial consistency level
func WithSerialConsistency(consistency string) Option {
	return func(c *client) error {
		if len(consistency) == 0 {
			return errors.NewErrInvalidOption("serialConsistency", consistency)
		}
		actual, ok := serialConsistenciesMap[strings.TrimSpace(strings.Trim(strings.Trim(strings.ToLower(consistency), "_"), "-"))]
		if !ok {
			return errors.NewErrCriticalOption("serialConsistency", consistency)
		}
		c.serialConsistency = actual
		return nil
	}
}

// WithCompressor returns the option to set the compressor
func WithCompressor(compressor gocql.Compressor) Option {
	return func(c *client) error {
		if compressor == nil {
			return errors.NewErrInvalidOption("compressor", compressor)
		}
		c.compressor = compressor
		return nil
	}
}

// WithUsername returns the option to set the username
func WithUsername(username string) Option {
	return func(c *client) error {
		if len(username) == 0 {
			return errors.NewErrInvalidOption("username", username)
		}
		c.username = username
		return nil
	}
}

// WithPassword returns the option to set the password
func WithPassword(password string) Option {
	return func(c *client) error {
		if len(password) == 0 {
			return errors.NewErrInvalidOption("password", password)
		}
		c.password = password
		return nil
	}
}

// WithAuthProvider returns the option to set the auth provider
func WithAuthProvider(authProvider func(h *gocql.HostInfo) (gocql.Authenticator, error)) Option {
	return func(c *client) error {
		if authProvider == nil {
			return errors.NewErrInvalidOption("authProvider", authProvider)
		}
		c.authProvider = authProvider
		return nil
	}
}

// WithRetryPolicyNumRetries returns the option to set the number of retries
func WithRetryPolicyNumRetries(n int) Option {
	return func(c *client) error {
		if n < 0 {
			return errors.NewErrInvalidOption("retryPolicyNumRetries", n)
		}
		c.retryPolicy.numRetries = n
		return nil
	}
}

// WithRetryPolicyMinDuration returns the option to set the retry min duration
func WithRetryPolicyMinDuration(minDuration string) Option {
	return func(c *client) error {
		if len(minDuration) == 0 {
			return errors.NewErrInvalidOption("retryPolicyMinDuration", minDuration)
		}
		d, err := timeutil.Parse(minDuration)
		if err != nil {
			return errors.NewErrCriticalOption("retryPolicyMinDuration", minDuration, err)
		}
		c.retryPolicy.minDuration = d
		return nil
	}
}

// WithRetryPolicyMaxDuration returns the option to set the retry max duration
func WithRetryPolicyMaxDuration(maxDuration string) Option {
	return func(c *client) error {
		if len(maxDuration) == 0 {
			return errors.NewErrInvalidOption("retryPolicyMaxDuration", maxDuration)
		}
		d, err := timeutil.Parse(maxDuration)
		if err != nil {
			return errors.NewErrCriticalOption("retryPolicyMaxDuration", maxDuration, err)
		}
		c.retryPolicy.maxDuration = d
		return nil
	}
}

// WithReconnectionPolicyInitialInterval returns the option to set the reconnect initial interval
func WithReconnectionPolicyInitialInterval(initialInterval string) Option {
	return func(c *client) error {
		if len(initialInterval) == 0 {
			return errors.NewErrInvalidOption("reconnectionPolicyInitialInterval", initialInterval)
		}
		d, err := timeutil.Parse(initialInterval)
		if err != nil {
			return errors.NewErrCriticalOption("reconnectionPolicyInitialInterval", initialInterval, err)
		}
		c.reconnectionPolicy.initialInterval = d
		return nil
	}
}

// WithReconnectionPolicyMaxRetries returns the option to set the reconnect max retries
func WithReconnectionPolicyMaxRetries(maxRetries int) Option {
	return func(c *client) error {
		if maxRetries < 0 {
			return errors.NewErrInvalidOption("maxRetries", maxRetries)
		}
		c.reconnectionPolicy.maxRetries = maxRetries
		return nil
	}
}

// WithSocketKeepalive returns the option to set the socket keepalive time
func WithSocketKeepalive(socketKeepalive string) Option {
	return func(c *client) error {
		if len(socketKeepalive) == 0 {
			return errors.NewErrInvalidOption("socketKeepalive", socketKeepalive)
		}
		d, err := timeutil.Parse(socketKeepalive)
		if err != nil {
			return errors.NewErrCriticalOption("socketKeepalive", socketKeepalive, err)
		}
		c.socketKeepalive = d
		return nil
	}
}

// WithMaxPreparedStmts returns the option to set the max prepared statement
func WithMaxPreparedStmts(maxPreparedStmts int) Option {
	return func(c *client) error {
		if maxPreparedStmts < 0 {
			return errors.NewErrInvalidOption("maxPreparedStmts", maxPreparedStmts)
		}
		c.maxPreparedStmts = maxPreparedStmts
		return nil
	}
}

// WithMaxRoutingKeyInfo returns the option to set the max routing key info
func WithMaxRoutingKeyInfo(maxRoutingKeyInfo int) Option {
	return func(c *client) error {
		if maxRoutingKeyInfo < 0 {
			return errors.NewErrInvalidOption("maxRoutingKeyInfo", maxRoutingKeyInfo)
		}
		c.maxRoutingKeyInfo = maxRoutingKeyInfo
		return nil
	}
}

// WithPageSize returns the option to set the page size
func WithPageSize(pageSize int) Option {
	return func(c *client) error {
		if pageSize < 0 {
			return errors.NewErrInvalidOption("pageSize", pageSize)
		}
		c.pageSize = pageSize
		return nil
	}
}

func WithTLS(tls *tls.Config) Option {
	return func(c *client) error {
		c.tls = tls
		return nil
	}
}

func WithTLSCertPath(certPath string) Option {
	return func(c *client) error {
		c.tlsCertPath = certPath
		return nil
	}
}

func WithTLSKeyPath(keyPath string) Option {
	return func(c *client) error {
		c.tlsKeyPath = keyPath
		return nil
	}
}

func WithTLSCAPath(caPath string) Option {
	return func(c *client) error {
		c.tlsCAPath = caPath
		return nil
	}
}

func WithEnableHostVerification(enableHostVerification bool) Option {
	return func(c *client) error {
		c.enableHostVerification = enableHostVerification
		return nil
	}
}

func WithDefaultTimestamp(defaultTimestamp bool) Option {
	return func(c *client) error {
		c.defaultTimestamp = defaultTimestamp
		return nil
	}
}

func WithDC(name string) Option {
	return func(c *client) error {
		c.poolConfig.dataCenterName = name
		return nil
	}
}

func WithDCAwareRouting(dcAwareRouting bool) Option {
	return func(c *client) error {
		c.poolConfig.enableDCAwareRouting = dcAwareRouting
		return nil
	}
}

func WithNonLocalReplicasFallback(nonLocalReplicasFallBack bool) Option {
	return func(c *client) error {
		c.poolConfig.enableNonLocalReplicasFallback = nonLocalReplicasFallBack
		return nil
	}
}

func WithShuffleReplicas(shuffleReplicas bool) Option {
	return func(c *client) error {
		c.poolConfig.enableShuffleReplicas = shuffleReplicas
		return nil
	}
}

func WithTokenAwareHostPolicy(tokenAwareHostPolicy bool) Option {
	return func(c *client) error {
		c.poolConfig.enableTokenAwareHostPolicy = tokenAwareHostPolicy
		return nil
	}
}

func WithMaxWaitSchemaAgreement(maxWaitSchemaAgreement string) Option {
	return func(c *client) error {
		d, err := timeutil.Parse(maxWaitSchemaAgreement)
		if err != nil {
			return err
		}
		c.maxWaitSchemaAgreement = d
		return nil
	}
}

func WithReconnectInterval(reconnectInterval string) Option {
	return func(c *client) error {
		d, err := timeutil.Parse(reconnectInterval)
		if err != nil {
			return err
		}
		c.reconnectInterval = d
		return nil
	}
}

func WithIgnorePeerAddr(ignorePeerAddr bool) Option {
	return func(c *client) error {
		c.ignorePeerAddr = ignorePeerAddr
		return nil
	}
}

func WithDisableInitialHostLookup(disableInitialHostLookup bool) Option {
	return func(c *client) error {
		c.disableInitialHostLookup = disableInitialHostLookup
		return nil
	}
}

func WithDisableNodeStatusEvents(disableNodeStatusEvents bool) Option {
	return func(c *client) error {
		c.disableNodeStatusEvents = disableNodeStatusEvents
		return nil
	}
}

func WithDisableTopologyEvents(disableTopologyEvents bool) Option {
	return func(c *client) error {
		c.disableTopologyEvents = disableTopologyEvents
		return nil
	}
}

func WithDisableSchemaEvents(disableSchemaEvents bool) Option {
	return func(c *client) error {
		c.disableSchemaEvents = disableSchemaEvents
		return nil
	}
}

func WithDisableSkipMetadata(disableSkipMetadata bool) Option {
	return func(c *client) error {
		c.disableSkipMetadata = disableSkipMetadata
		return nil
	}
}

func WithQueryObserver(obs QueryObserver) Option {
	return func(c *client) error {
		if obs != nil {
			c.queryObserver = obs
		}

		return nil
	}
}

func WithBatchObserver(obs BatchObserver) Option {
	return func(c *client) error {
		if obs != nil {
			c.batchObserver = obs
		}

		return nil
	}
}

func WithConnectObserver(obs ConnectObserver) Option {
	return func(c *client) error {
		if obs != nil {
			c.connectObserver = obs
		}

		return nil
	}
}

func WithFrameHeaderObserver(obs FrameHeaderObserver) Option {
	return func(c *client) error {
		if obs != nil {
			c.frameHeaderObserver = obs
		}

		return nil
	}
}

func WithDefaultIdempotence(defaultIdempotence bool) Option {
	return func(c *client) error {
		c.defaultIdempotence = defaultIdempotence
		return nil
	}
}

func WithWriteCoalesceWaitTime(writeCoalesceWaitTime string) Option {
	return func(c *client) error {
		d, err := timeutil.Parse(writeCoalesceWaitTime)
		if err != nil {
			return err
		}
		c.writeCoalesceWaitTime = d
		return nil
	}
}

func WithHostFilter(flg bool) Option {
	return func(c *client) error {
		c.hostFilter.enable = flg
		return nil
	}
}

func WithDCHostFilter(dc string) Option {
	return func(c *client) error {
		if len(dc) == 0 {
			return nil
		}
		c.hostFilter.dcHost = dc
		if !c.hostFilter.enable {
			WithHostFilter(true)(c)
		}
		return nil
	}
}

func WithWhiteListHostFilter(list []string) Option {
	return func(c *client) error {
		if len(list) <= 0 {
			return nil
		}
		c.hostFilter.whiteList = list
		if !c.hostFilter.enable {
			WithHostFilter(true)(c)
		}
		return nil
	}
}

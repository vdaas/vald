//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

type convictionPolicy struct{}

// NewConvictionPolicy returns the implementation of gocql.ConvictionPolicy.
func NewConvictionPolicy() gocql.ConvictionPolicy {
	return new(convictionPolicy)
}

// AddFailure implements gocql.ConvictionPolicy interface to handle failure and convicts all hosts.
func (c *convictionPolicy) AddFailure(err error, host *gocql.HostInfo) bool {
	log.Warn(errors.ErrCassandraHostDownDetected(err, host.String()))
	return true
}

// Reset clears the conviction state.
func (c *convictionPolicy) Reset(host *gocql.HostInfo) {
	log.Info("cassandra host %s reset detected\t%s", host.HostnameAndPort(), host.String())
}

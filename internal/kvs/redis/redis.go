//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
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

package redis

import (
	"context"
	"crypto/tls"
	"net"
	"reflect"
	"time"

	"github.com/go-redis/redis"
	"github.com/vdaas/vald/internal/errors"
)

type Redis interface {
}

type redisClient struct {
	addrs              []string
	clusterSlots       func() ([]redis.ClusterSlot, error)
	db                 int
	dialTimeout        time.Duration
	dialer             func(ctx context.Context, network, addr string) (net.Conn, error)
	idleCheckFrequency time.Duration
	idleTimeout        time.Duration
	keyPref            string
	maxConnAge         time.Duration
	maxRedirects       int
	maxRetries         int
	maxRetryBackoff    time.Duration
	minIdleConns       int
	minRetryBackoff    time.Duration
	onConnect          func(*redis.Conn) error
	onNewNode          func(*redis.Client)
	password           string
	poolSize           int
	poolTimeout        time.Duration
	readOnly           bool
	readTimeout        time.Duration
	routeByLatency     bool
	routeRandomly      bool
	tlsConfig          *tls.Config
	writeTimeout       time.Duration
}

func New(opts ...Option) (Redis, error) {
	r := new(redisClient)
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	switch len(r.addrs) {
	case 0:
		return nil, errors.ErrAddrsNotFound
	case 1:
	default:
	}
	return r, nil
}

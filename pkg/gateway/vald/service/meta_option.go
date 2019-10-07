//
// Copyright (C) 2019 kpango (Yusuke Kato)
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

// Package service
package service

import (
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
	"google.golang.org/grpc"
)

type MetaOption func(m *meta) error

var (
	defaultMetaOpts = []MetaOption{}
)

func WithMetaHost(host string) MetaOption {
	return func(m *meta) error {
		m.host = host
		return nil
	}
}

func WithMetaPort(port int) MetaOption {
	return func(m *meta) error {
		m.port = port
		return nil
	}
}

func WithHealthCheckDuration(dur string) MetaOption {
	return func(m *meta) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		m.hcDur = d
		return nil
	}
}

func WithMetaGRPCDialOption(opt grpc.DialOption) MetaOption {
	return func(m *meta) error {
		m.gopts = append(m.gopts, opt)
		return nil
	}
}

func WithMetaGRPCDialOptions(opts []grpc.DialOption) MetaOption {
	return func(m *meta) error {
		if m.gopts != nil && len(m.gopts) > 0 {
			m.gopts = append(m.gopts, opts...)
		} else {
			m.gopts = opts
		}
		return nil
	}
}

func WithMetaGRPCCallOption(opt grpc.CallOption) MetaOption {
	return func(m *meta) error {
		m.copts = append(m.copts, opt)
		return nil
	}
}

func WithMetaGRPCCallOptions(opts []grpc.CallOption) MetaOption {
	return func(m *meta) error {
		if m.copts != nil && len(m.copts) > 0 {
			m.copts = append(m.copts, opts...)
		} else {
			m.copts = opts
		}
		return nil
	}
}

func withMetaBackoff(bo backoff.Backoff) MetaOption {
	return func(m *meta) error {
		m.bo = bo
		return nil
	}
}

func withMetaErrGroup(eg errgroup.Group) MetaOption {
	return func(m *meta) error {
		m.eg = eg
		return nil
	}
}

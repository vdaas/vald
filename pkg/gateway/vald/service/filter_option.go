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

type FilterOption func(f *filter) error

var (
	defaultFilterOpts = []FilterOption{}
)

func WithFilterTargets(addrs ...string) FilterOption {
	return func(f *filter) error {
		f.addrs = addrs
		return nil
	}
}

func WithFilterHealthCheckDuration(dur string) FilterOption {
	return func(f *filter) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		f.hcDur = d
		return nil
	}
}

func WithFilterGRPCDialOption(opt grpc.DialOption) FilterOption {
	return func(f *filter) error {
		f.gopts = append(f.gopts, opt)
		return nil
	}
}

func WithFilterGRPCDialOptions(opts ...grpc.DialOption) FilterOption {
	return func(f *filter) error {
		if f.gopts != nil && len(f.gopts) > 0 {
			f.gopts = append(f.gopts, opts...)
		} else {
			f.gopts = opts
		}
		return nil
	}
}

func WithFilterGRPCCallOption(opt grpc.CallOption) FilterOption {
	return func(f *filter) error {
		f.copts = append(f.copts, opt)
		return nil
	}
}

func WithFilterGRPCCallOptions(opts ...grpc.CallOption) FilterOption {
	return func(f *filter) error {
		if f.copts != nil && len(f.copts) > 0 {
			f.copts = append(f.copts, opts...)
		} else {
			f.copts = opts
		}
		return nil
	}
}

func withFilterBackoff(bo backoff.Backoff) FilterOption {
	return func(f *filter) error {
		f.bo = bo
		return nil
	}
}

func withFilterErrGroup(eg errgroup.Group) FilterOption {
	return func(f *filter) error {
		f.eg = eg
		return nil
	}
}

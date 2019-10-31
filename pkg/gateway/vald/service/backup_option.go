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

type BackupOption func(b *backup) error

var (
	defaultBackupOpts = []BackupOption{}
)

func WithBackupHost(host string) BackupOption {
	return func(b *backup) error {
		b.host = host
		return nil
	}
}

func WithBackupPort(port int) BackupOption {
	return func(b *backup) error {
		b.port = port
		return nil
	}
}

func WithHealthCheckDuration(dur string) BackupOption {
	return func(b *backup) error {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second
		}
		b.hcDur = d
		return nil
	}
}

func WithBackupGRPCDialOption(opt grpc.DialOption) BackupOption {
	return func(b *backup) error {
		b.gopts = append(b.gopts, opt)
		return nil
	}
}

func WithBackupGRPCDialOptions(opts []grpc.DialOption) BackupOption {
	return func(b *backup) error {
		if b.gopts != nil && len(b.gopts) > 0 {
			b.gopts = append(b.gopts, opts...)
		} else {
			b.gopts = opts
		}
		return nil
	}
}

func WithBackupGRPCCallOption(opt grpc.CallOption) BackupOption {
	return func(b *backup) error {
		b.copts = append(b.copts, opt)
		return nil
	}
}

func WithBackupGRPCCallOptions(opts []grpc.CallOption) BackupOption {
	return func(b *backup) error {
		if b.copts != nil && len(b.copts) > 0 {
			b.copts = append(b.copts, opts...)
		} else {
			b.copts = opts
		}
		return nil
	}
}

func withBackupBackoff(bo backoff.Backoff) BackupOption {
	return func(b *backup) error {
		b.bo = bo
		return nil
	}
}

func withBackupErrGroup(eg errgroup.Group) BackupOption {
	return func(b *backup) error {
		b.eg = eg
		return nil
	}
}

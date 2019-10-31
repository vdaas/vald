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

package service

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	gback "github.com/vdaas/vald/apis/grpc/manager/backup"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Backup interface {
	GetObject(ctx context.Context, uuid string) (*payload.Object_MetaVector, error)
	GetLocation(ctx context.Context, uuid string) ([]string, error)
	Register(ctx context.Context, vec *payload.Object_MetaVector) error
	RegisterMultiple(ctx context.Context, vecs *payload.Object_MetaVectors) error
	Remove(ctx context.Context, uuid string) error
	RemoveMultiple(ctx context.Context, uuids ...string) error
}

type backup struct {
	hcDur time.Duration
	host  string
	port  int
	bc    atomic.Value
	eg    errgroup.Group
	bo    backoff.Backoff
	gopts []grpc.DialOption
	copts []grpc.CallOption
}

func NewBackup(opts ...BackupOption) (bu Backup, err error) {
	b := new(backup)
	for _, opt := range append(defaultBackupOpts, opts...) {
		err = opt(b)

		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (b *backup) Start(ctx context.Context) <-chan error {
	ech := make(chan error)
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", b.host, b.port), b.gopts...)
	if err != nil {
		ech <- err
	} else {
		b.bc.Store(gback.NewBackupClient(conn))
	}
	b.eg.Go(safety.RecoverFunc(func() (err error) {
		tick := time.NewTicker(b.hcDur)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				close(ech)
				return ctx.Err()
			case <-tick.C:
				if conn == nil ||
					conn.GetState() == connectivity.Shutdown ||
					conn.GetState() == connectivity.TransientFailure {
					if conn != nil {
						err = conn.Close()
						if err != nil {
							ech <- err
						}
					}
					conn, err = grpc.DialContext(ctx, fmt.Sprintf("%s:%d", b.host, b.port), b.gopts...)
					if err != nil {
						ech <- err
					} else {
						b.bc.Store(gback.NewBackupClient(conn))
						runtime.Gosched()
					}
				}
			}
		}
		return nil
	}))
	return ech
}

func (b *backup) GetObject(ctx context.Context, uuid string) (vec *payload.Object_MetaVector, err error) {
	vec, err = b.bc.Load().(gback.BackupClient).GetVector(ctx, &payload.Object_ID{
		Id: uuid,
	}, b.copts...)
	if err != nil {
		return nil, err
	}
	return vec, nil
}

func (b *backup) GetLocation(ctx context.Context, uuid string) (ipList []string, err error) {
	ips, err := b.bc.Load().(gback.BackupClient).Locations(ctx, &payload.Object_ID{
		Id: uuid,
	}, b.copts...)
	if err != nil {
		return nil, err
	}
	return ips.GetIp(), nil
}

func (b *backup) Register(ctx context.Context, vec *payload.Object_MetaVector) error {
	_, err := b.bc.Load().(gback.BackupClient).Register(ctx, vec, b.copts...)
	if err != nil {
		return err
	}
	return nil
}

func (b *backup) RegisterMultiple(ctx context.Context, vecs *payload.Object_MetaVectors) error {
	_, err := b.bc.Load().(gback.BackupClient).RegisterMulti(ctx, vecs, b.copts...)
	if err != nil {
		return err
	}
	return nil
}
func (b *backup) Remove(ctx context.Context, uuid string) error {
	_, err := b.bc.Load().(gback.BackupClient).Remove(ctx, &payload.Object_ID{
		Id: uuid,
	}, b.copts...)
	if err != nil {
		return err
	}
	return nil
}
func (b *backup) RemoveMultiple(ctx context.Context, uuids ...string) error {
	ids := new(payload.Object_IDs)
	ids.Ids = make([]string, 0, len(uuids))
	for _, uuid := range uuids {
		ids.Ids = append(ids.Ids, uuid)
	}
	_, err := b.bc.Load().(gback.BackupClient).RemoveMulti(ctx, ids, b.copts...)
	if err != nil {
		return err
	}
	return nil

}

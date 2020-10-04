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

// Package service
package service

import "github.com/vdaas/vald/internal/net/grpc"

type BackupOption func(b *backup) error

var (
	defaultBackupOpts = []BackupOption{}
)

func WithBackupAddr(addr string) BackupOption {
	return func(b *backup) error {
		b.addr = addr
		return nil
	}
}

func WithBackupClient(client grpc.Client) BackupOption {
	return func(b *backup) error {
		if client != nil {
			b.client = client
		}
		return nil
	}
}

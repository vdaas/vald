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

// Package observer provides storage observer
package observer

import (
	"context"
	"time"

	"github.com/vdaas/vald/pkg/agent/sidecar/service/storage"
)

type BackupInfo struct {
	StartTime time.Time
	EndTime   time.Time
	Bytes     int64

	*storage.StorageInfo
}

type Hook interface {
	BeforeProcess(ctx context.Context, info *BackupInfo) (context.Context, error)
	AfterProcess(ctx context.Context, info *BackupInfo) error
}

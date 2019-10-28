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

	"github.com/vdaas/vald/apis/grpc/payload"
)

type Backup interface {
	GetObject(ctx context.Context, uuid string) (*payload.Object_Vector, error)
	GetLocation(ctx context.Context, uuid string) ([]string, error)
	Register(ctx context.Context, vec *payload.Object_Vector, srvs ...string) error
	Remove(ctx context.Context, uuid string) error
}

type backup struct {
}

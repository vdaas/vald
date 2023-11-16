// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"reflect"

	snapshotclient "github.com/kubernetes-csi/external-snapshotter/client/v6/clientset/versioned"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"k8s.io/client-go/kubernetes"
)

const (
	apiName = "vald/index/job/readreplica/rotate"
)

// Rotator represents an interface for indexing.
type Rotator interface {
	Start(ctx context.Context) error
}

type rotator struct {
	replicaid        int
	namespace        string
	deploymentPrefix string
	snapshotPrefix   string
	pvcPrefix        string
	// TODO: この辺はconbenchがマージされたあと、GetClientとかで引っ張ってくる
	clientset  *kubernetes.Clientset
	sClientset *snapshotclient.Clientset
}

// New returns Indexer object if no error occurs.
func New(opts ...Option) (Rotator, error) {
	r := new(rotator)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(r); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(err)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return r, nil
}

// Start starts rotation process.
func (idx *rotator) Start(ctx context.Context) error {
	_, span := trace.StartSpan(ctx, apiName+"/service/rotator.Start")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	return nil
}

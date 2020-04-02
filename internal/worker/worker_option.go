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

// Package worker provides worker processes
package worker

import (
	"github.com/vdaas/vald/internal/errgroup"
)

type WorkerOption func(w *worker) error

var (
	defaultWorkerOpts = []WorkerOption{
		WithName("worker"),
		WithLimitation(0),
		WithBuffer(0),
		WithErrGroup(errgroup.Get()),
	}
)

func WithName(name string) WorkerOption {
	return func(w *worker) error {
		if name != "" {
			w.name = name
		}
		return nil
	}
}

func WithLimitation(limit int) WorkerOption {
	return func(w *worker) error {
		w.limitation = limit
		return nil
	}
}

func WithBuffer(b int) WorkerOption {
	return func(w *worker) error {
		w.buffer = b
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) WorkerOption {
	return func(w *worker) error {
		if eg != nil {
			w.eg = eg
		}
		return nil
	}
}

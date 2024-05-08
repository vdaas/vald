//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type WorkerOption func(w *worker) error

var defaultWorkerOpts = []WorkerOption{
	WithName("worker"),
	WithLimitation(10),
	WithErrGroup(errgroup.Get()),
}

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
		if limit > 0 {
			w.limitation = limit
		}
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

func WithQueueOption(opts ...QueueOption) WorkerOption {
	return func(w *worker) error {
		if opts == nil {
			return nil
		}
		if w.qopts != nil {
			w.qopts = append(w.qopts, opts...)
		} else {
			w.qopts = opts
		}
		return nil
	}
}

//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
)

// QueueOption represents the functional option for queue.
type QueueOption func(q *queue) error

var defaultQueueOpts = []QueueOption{
	WithQueueBuffer(10),
	WithQueueErrGroup(errgroup.Get()),
	WithQueueCheckDuration("200ms"),
}

// WithQueueBuffer returns the option to set the buffer for queue.
func WithQueueBuffer(buffer int) QueueOption {
	return func(q *queue) error {
		if buffer > 0 {
			q.buffer = buffer
		}
		return nil
	}
}

// WithQueueErrGroup returns the options to set the eg for queue.
func WithQueueErrGroup(eg errgroup.Group) QueueOption {
	return func(q *queue) error {
		if eg != nil {
			q.eg = eg
		}
		return nil
	}
}

// WithQueueCheckDuration returns the option to set the qcdur for queue.
// If dur is invalid string, it returns errror.
func WithQueueCheckDuration(dur string) QueueOption {
	return func(q *queue) error {
		if len(dur) == 0 {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		q.qcdur = d
		return nil
	}
}

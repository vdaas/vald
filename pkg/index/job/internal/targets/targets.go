// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// Package targets tracks the per-target outcome of an index job that fans one
// RPC out to every target agent.
package targets

import (
	"maps"
	"slices"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// Outcomes records the latest result of each target's callback. It is safe
// for concurrent use by the callbacks of OrderedRangeConcurrent.
type Outcomes struct {
	// errs keeps only the latest outcome per target: the client backoff can re-run the
	// callback for the same target, and a retry that succeeds must clear the earlier failure.
	errs map[string]error
	// pending keeps the targets whose callback has not run: the client skips the
	// callback for a target with no connection pool, a connection lookup failure, an
	// open circuit breaker or a done ctx and swallows that error, so such a target would
	// otherwise count as a success.
	pending map[string]struct{}
	mu      sync.Mutex
}

// New returns an Outcomes with every target in addrs pending.
func New(addrs []string) *Outcomes {
	o := &Outcomes{
		errs:    make(map[string]error, len(addrs)),
		pending: make(map[string]struct{}, len(addrs)),
	}
	for _, target := range addrs {
		o.pending[target] = struct{}{}
	}
	return o
}

// Record stores err as the latest outcome of target and marks it processed.
func (o *Outcomes) Record(target string, err error) {
	o.mu.Lock()
	delete(o.pending, target)
	if err != nil {
		o.errs[target] = err
	} else {
		delete(o.errs, target)
	}
	o.mu.Unlock()
}

// Err joins clientErr, the error returned by the client's range call, with the
// recorded target errors in target order and, only when clientErr is nil, an
// ErrGRPCClientConnNotFound for every target whose callback never ran.
func (o *Outcomes) Err(clientErr error) error {
	var errs error
	o.mu.Lock()
	for _, target := range slices.Sorted(maps.Keys(o.errs)) {
		errs = errors.Join(errs, o.errs[target])
	}
	// A non-nil clientErr already fails the job; report unprocessed targets only when the
	// client claimed success.
	if clientErr == nil {
		for _, target := range slices.Sorted(maps.Keys(o.pending)) {
			errs = errors.Join(errs, errors.ErrGRPCClientConnNotFound(target))
		}
	}
	o.mu.Unlock()
	return errors.Join(clientErr, errs)
}

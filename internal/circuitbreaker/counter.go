//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package circuitbreaker

import "sync/atomic"

type Counter interface {
	Successes() int64
	Fails() int64
}

type count struct {
	successes int64
	failures  int64
}

func (c *count) Successes() (n int64) {
	return atomic.LoadInt64(&c.successes)
}

func (c *count) Fails() (n int64) {
	return atomic.LoadInt64(&c.failures)
}

func (c *count) onSuccess() {
	atomic.AddInt64(&c.successes, 1)
}

func (c *count) onFail() {
	atomic.AddInt64(&c.failures, 1)
}

func (c *count) reset() {
	atomic.StoreInt64(&c.failures, 0)
	atomic.StoreInt64(&c.successes, 0)
}

var _ Counter = (*count)(nil)

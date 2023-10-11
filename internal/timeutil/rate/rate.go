//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

package rate

import (
	"context"
	"runtime"

	"go.uber.org/ratelimit"
	"golang.org/x/time/rate"
)

type Limiter interface {
	Wait(ctx context.Context) error
}

type limiter struct {
	isStd bool
	uber  ratelimit.Limiter
	std   *rate.Limiter
}

func NewLimiter(cnt int) Limiter {
	if runtime.GOMAXPROCS(0) >= 32 {
		return &limiter{
			isStd: true,
			std:   rate.NewLimiter(rate.Limit(cnt), 1),
		}
	}
	return &limiter{
		isStd: false,
		uber:  ratelimit.New(cnt, ratelimit.WithoutSlack),
	}
}

func (l *limiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if l.isStd {
			return l.std.Wait(ctx)
		}
		l.uber.Take()
	}
	return nil
}

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

package service

import (
	"math"
	"math/big"
	"os"
	"time"

	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/rand"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/timeutil"
)

// Option represent the functional option for ngt
type Option func(n *ngt) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithAutoIndexCheckDuration("30m"),
	WithAutoIndexDurationLimit("24h"),
	WithAutoSaveIndexDuration("35m"),
	WithAutoIndexLength(100),
	WithInitialDelayMaxDuration("3m"),
	WithMinLoadIndexTimeout("3m"),
	WithMaxLoadIndexTimeout("10m"),
	WithLoadIndexTimeoutFactor("1ms"),
	WithDefaultPoolSize(core.DefaultPoolSize),
	WithDefaultRadius(core.DefaultRadius),
	WithDefaultEpsilon(core.DefaultEpsilon),
	WithProactiveGC(true),
}

// WithErrGroup returns the functional option to set the error group.
func WithErrGroup(eg errgroup.Group) Option {
	return func(n *ngt) error {
		if eg != nil {
			n.eg = eg
		}

		return nil
	}
}

// WithEnableInMemoryMode returns the functional option to set the in memory mode flag.
func WithEnableInMemoryMode(enabled bool) Option {
	return func(n *ngt) error {
		n.inMem = enabled

		return nil
	}
}

// WithIndexPath returns the functional option to set the index path of the NGT.
func WithIndexPath(path string) Option {
	return func(n *ngt) error {
		if path == "" {
			return nil
		}
		n.path = file.Join(strings.TrimSuffix(path, string(os.PathSeparator)))
		return nil
	}
}

// WithAutoIndexCheckDuration returns the functional option to set the index check duration.
func WithAutoIndexCheckDuration(dur string) Option {
	return func(n *ngt) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		n.dur = d

		return nil
	}
}

// WithAutoIndexDurationLimit returns the functional option to set the auto index duration limit.
func WithAutoIndexDurationLimit(dur string) Option {
	return func(n *ngt) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		n.lim = d

		return nil
	}
}

// WithAutoSaveIndexDuration returns the functional option to set the auto save index duration.
func WithAutoSaveIndexDuration(dur string) Option {
	return func(n *ngt) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		n.sdur = d

		return nil
	}
}

// WithAutoIndexLength returns the functional option to set the auto index length.
func WithAutoIndexLength(l int) Option {
	return func(n *ngt) error {
		n.alen = l

		return nil
	}
}

const (
	defaultDurationLimit float64 = 1.1
	defaultRandDuration  int64   = 1
)

var (
	bigMaxFloat64 = big.NewFloat(math.MaxFloat64)
	bigMinFloat64 = big.NewFloat(math.SmallestNonzeroFloat64)
	bigMaxInt64   = big.NewInt(math.MaxInt64)
	bigMinInt64   = big.NewInt(math.MinInt64)
)

// WithInitialDelayMaxDuration returns the functional option to set the initial delay duration.
func WithInitialDelayMaxDuration(dur string) Option {
	return func(n *ngt) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		var dt time.Duration
		switch {
		case d <= time.Nanosecond:
			return nil
		case d <= time.Microsecond:
			dt = time.Nanosecond
		case d <= time.Millisecond:
			dt = time.Microsecond
		case d <= time.Second:
			dt = time.Millisecond
		default:
			dt = time.Second
		}

		dbs := math.Round(float64(d) / float64(dt))
		bdbs := big.NewFloat(dbs)
		if dbs <= 0 || bigMaxFloat64.Cmp(bdbs) <= 0 || bigMinFloat64.Cmp(bdbs) >= 0 {
			dbs = defaultDurationLimit
		}

		rnd := int64(rand.LimitedUint32(uint64(dbs)))
		brnd := big.NewInt(rnd)
		if rnd <= 0 || bigMaxInt64.Cmp(brnd) <= 0 || bigMinInt64.Cmp(brnd) >= 0 {
			rnd = defaultRandDuration
		}

		delay := time.Duration(rnd) * dt
		if delay <= 0 || delay >= math.MaxInt64 || delay <= math.MinInt64 {
			return WithInitialDelayMaxDuration(dur)(n)
		}

		n.idelay = delay

		return nil
	}
}

// WithMinLoadIndexTimeout returns the functional option to set the minimal load index timeout.
func WithMinLoadIndexTimeout(dur string) Option {
	return func(n *ngt) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		n.minLit = d

		return nil
	}
}

// WithMaxLoadIndexTimeout returns the functional option to set the maximum load index timeout.
func WithMaxLoadIndexTimeout(dur string) Option {
	return func(n *ngt) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		n.maxLit = d

		return nil
	}
}

// WithLoadIndexTimeoutFactor returns the functional option to set the factor of load index timeout.
func WithLoadIndexTimeoutFactor(dur string) Option {
	return func(n *ngt) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		n.litFactor = d

		return nil
	}
}

// WithDefaultPoolSize returns the functional option to set the default pool size for NGT.
func WithDefaultPoolSize(ps uint32) Option {
	return func(n *ngt) error {
		n.poolSize = ps

		return nil
	}
}

// WithDefaultRadius returns the functional option to set the default radius for NGT.
func WithDefaultRadius(rad float32) Option {
	return func(n *ngt) error {
		n.radius = rad

		return nil
	}
}

// WithDefaultEpsilon returns the functional option to set the default epsilon for NGT.
func WithDefaultEpsilon(epsilon float32) Option {
	return func(n *ngt) error {
		n.epsilon = epsilon

		return nil
	}
}

// WithProactiveGC returns the functional option to set the proactive GC enable flag.
func WithProactiveGC(enabled bool) Option {
	return func(n *ngt) error {
		n.enableProactiveGC = enabled
		return nil
	}
}

// WithCopyOnWrite returns the functional option to set the CoW enable flag.
func WithCopyOnWrite(enabled bool) Option {
	return func(n *ngt) error {
		n.enableCopyOnWrite = enabled
		return nil
	}
}

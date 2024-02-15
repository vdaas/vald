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

package service

import (
	"math"
	"math/big"
	"os"
	"time"

<<<<<<< HEAD
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/rand"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
=======
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/rand"
	"github.com/vdaas/vald/internal/strings"
>>>>>>> feature/agent/qbg
	"github.com/vdaas/vald/internal/timeutil"
)

// Option represent the functional option for faiss
type Option func(f *faiss) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithAutoIndexCheckDuration("30m"),
	WithAutoSaveIndexDuration("35m"),
	WithAutoIndexDurationLimit("24h"),
	WithAutoIndexLength(100),
	WithInitialDelayMaxDuration("3m"),
	WithMinLoadIndexTimeout("3m"),
	WithMaxLoadIndexTimeout("10m"),
	WithLoadIndexTimeoutFactor("1ms"),
	WithProactiveGC(true),
}

// WithErrGroup returns the functional option to set the error group.
func WithErrGroup(eg errgroup.Group) Option {
	return func(f *faiss) error {
		if eg != nil {
			f.eg = eg
		}

		return nil
	}
}

// WithEnableInMemoryMode returns the functional option to set the in memory mode flag.
func WithEnableInMemoryMode(enabled bool) Option {
	return func(f *faiss) error {
		f.inMem = enabled

		return nil
	}
}

// WithIndexPath returns the functional option to set the index path of the Faiss.
func WithIndexPath(path string) Option {
	return func(f *faiss) error {
		if path == "" {
			return nil
		}
		f.path = file.Join(strings.TrimSuffix(path, string(os.PathSeparator)))
		return nil
	}
}

// WithAutoIndexCheckDuration returns the functional option to set the index check duration.
func WithAutoIndexCheckDuration(dur string) Option {
	return func(f *faiss) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		f.dur = d

		return nil
	}
}

// WithAutoSaveIndexDuration returns the functional option to set the auto save index duration.
func WithAutoSaveIndexDuration(dur string) Option {
	return func(f *faiss) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		f.sdur = d

		return nil
	}
}

// WithAutoIndexDurationLimit returns the functional option to set the auto index duration limit.
func WithAutoIndexDurationLimit(dur string) Option {
	return func(f *faiss) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		f.lim = d

		return nil
	}
}

// WithAutoIndexLength returns the functional option to set the auto index length.
func WithAutoIndexLength(l int) Option {
	return func(f *faiss) error {
		f.alen = l

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
	return func(f *faiss) error {
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
			return WithInitialDelayMaxDuration(dur)(f)
		}

		f.idelay = delay

		return nil
	}
}

// WithMinLoadIndexTimeout returns the functional option to set the minimal load index timeout.
func WithMinLoadIndexTimeout(dur string) Option {
	return func(f *faiss) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		f.minLit = d

		return nil
	}
}

// WithMaxLoadIndexTimeout returns the functional option to set the maximum load index timeout.
func WithMaxLoadIndexTimeout(dur string) Option {
	return func(f *faiss) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		f.maxLit = d

		return nil
	}
}

// WithLoadIndexTimeoutFactor returns the functional option to set the factor of load index timeout.
func WithLoadIndexTimeoutFactor(dur string) Option {
	return func(f *faiss) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		f.litFactor = d

		return nil
	}
}

// WithProactiveGC returns the functional option to set the proactive GC enable flag.
func WithProactiveGC(enabled bool) Option {
	return func(f *faiss) error {
		f.enableProactiveGC = enabled
		return nil
	}
}

// WithCopyOnWrite returns the functional option to set the CoW enable flag.
func WithCopyOnWrite(enabled bool) Option {
	return func(f *faiss) error {
		f.enableCopyOnWrite = enabled
		return nil
	}
}

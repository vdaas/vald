//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"path/filepath"
	"strings"
	"time"

	core "github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/rand"
	"github.com/vdaas/vald/internal/timeutil"
)

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

func WithErrGroup(eg errgroup.Group) Option {
	return func(n *ngt) error {
		if eg != nil {
			n.eg = eg
		}

		return nil
	}
}

func WithEnableInMemoryMode(enabled bool) Option {
	return func(n *ngt) error {
		n.inMem = enabled

		return nil
	}
}

func WithIndexPath(path string) Option {
	return func(n *ngt) error {
		if path == "" {
			return nil
		}

		n.path = filepath.Clean(strings.TrimSuffix(path, "/"))

		return nil
	}
}

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

func WithAutoIndexLength(l int) Option {
	return func(n *ngt) error {
		n.alen = l

		return nil
	}
}

func WithInitialDelayMaxDuration(dur string) Option {
	return func(n *ngt) error {
		if dur == "" {
			return nil
		}

		d, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}

		n.idelay = time.Duration(int64(rand.LimitedUint32(uint64(d/time.Second)))) * time.Second

		return nil
	}
}

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

func WithDefaultPoolSize(ps uint32) Option {
	return func(n *ngt) error {
		n.poolSize = ps

		return nil
	}
}

func WithDefaultRadius(rad float32) Option {
	return func(n *ngt) error {
		n.radius = rad

		return nil
	}
}

func WithDefaultEpsilon(epsilon float32) Option {
	return func(n *ngt) error {
		n.epsilon = epsilon

		return nil
	}
}

func WithProactiveGC(enabled bool) Option {
	return func(n *ngt) error {
		n.enableProactiveGC = enabled
		return nil
	}
}

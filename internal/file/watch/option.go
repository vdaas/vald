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

package watch

import (
	"context"

	"github.com/vdaas/vald/internal/errgroup"
)

type Option func(w *watch) error

var defaultOptions = []Option{}

func WithErrGroup(eg errgroup.Group) Option {
	return func(w *watch) error {
		if eg != nil {
			w.eg = eg
		}
		return nil
	}
}

func WithDirs(dirs ...string) Option {
	return func(w *watch) error {
		if len(dirs) == 0 {
			return nil
		}

		if w.dirs == nil {
			w.dirs = make(map[string]struct{})
		}

		for _, dir := range dirs {
			w.dirs[dir] = struct{}{}
		}
		return nil
	}
}

func WithOnChange(f func(ctx context.Context, name string) error) Option {
	return func(w *watch) error {
		if f != nil {
			w.onChange = f
		}
		return nil
	}
}

func WithOnCreate(f func(ctx context.Context, name string) error) Option {
	return func(w *watch) error {
		if f != nil {
			w.onCreate = f
		}
		return nil
	}
}

func WithOnChmod(f func(ctx context.Context, name string) error) Option {
	return func(w *watch) error {
		if f != nil {
			w.onChmod = f
		}
		return nil
	}
}

func WithOnRename(f func(ctx context.Context, name string) error) Option {
	return func(w *watch) error {
		if f != nil {
			w.onRename = f
		}
		return nil
	}
}

func WithOnDelete(f func(ctx context.Context, name string) error) Option {
	return func(w *watch) error {
		if f != nil {
			w.onDelete = f
		}
		return nil
	}
}

func WithOnWrite(f func(ctx context.Context, name string) error) Option {
	return func(w *watch) error {
		if f != nil {
			w.onWrite = f
		}
		return nil
	}
}

func WithOnError(f func(ctx context.Context, err error) error) Option {
	return func(w *watch) error {
		if f != nil {
			w.onError = f
		}
		return nil
	}
}

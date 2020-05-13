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
package insert

import (
	"github.com/vdaas/vald/internal/client"
)

type InsertOption func(*insert) error

var (
	defaultInsertOpts = []InsertOption{
		WithConcurrency(100),
	}
)

func WithWriter(w client.Writer) InsertOption {
	return func(i *insert) error {
		i.w = w
		return nil
	}
}

func WithConcurrency(c int) InsertOption {
	return func(i *insert) error {
		i.c = c
		return nil
	}
}

func WithDataset(n string) InsertOption {
	return func(i *insert) error {
		i.n = n
		return nil
	}
}

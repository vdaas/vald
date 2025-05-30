// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package pogreb

import (
	"time"

	"github.com/akrylysov/pogreb"
)

var deafultOpts = []Option{
	WithPath("pogreb.db"),
	WithBackgroundSyncInterval(time.Second * 5),
}

// Option represents the functional option for database.
type Option func(*db) error

// WithPath returns the option to set path.
func WithPath(path string) Option {
	return func(d *db) error {
		if path != "" {
			d.path = path
		}
		return nil
	}
}

// WithBackgroundSyncInterval returns the option to sets the amount of time between background Sync() calls.
// Setting the value to 0 disables the automatic background synchronization.
// Setting the value to -1 or less makes the DB call Sync() after every write operation.
func WithBackgroundSyncInterval(dur time.Duration) Option {
	return func(d *db) error {
		if d.opts == nil {
			d.opts = new(pogreb.Options)
		}
		d.opts.BackgroundSyncInterval = dur
		return nil
	}
}

// WithBackgroundCompactionInterval returns the option to sets the amount of time between background Compact() calls.
// Setting the value to 0 or less disables the automatic background compaction.
func WithBackgroundCompactionInterval(dur time.Duration) Option {
	return func(d *db) error {
		if d.opts == nil {
			d.opts = new(pogreb.Options)
		}
		d.opts.BackgroundCompactionInterval = dur
		return nil
	}
}

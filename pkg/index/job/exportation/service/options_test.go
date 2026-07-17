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
package service

import (
	"testing"
	"time"

	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func TestWithStreamListConcurrency(t *testing.T) {
	t.Parallel()
	type test struct {
		name    string
		num     int
		wantErr bool
	}
	tests := []test{
		{name: "sets streamListConcurrency when num is positive", num: 5},
		{name: "returns invalid option error when num is zero", num: 0, wantErr: true},
		{name: "returns invalid option error when num is negative", num: -1, wantErr: true},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			e := new(export)
			err := WithStreamListConcurrency(test.num)(e)
			if test.wantErr {
				var invalidOpt *errors.ErrInvalidOption
				if !errors.As(err, &invalidOpt) {
					tt.Errorf("got_error: \"%v\", want a wrapped *errors.ErrInvalidOption", err)
				}
				return
			}
			if err != nil {
				tt.Errorf("err: %v", err)
			}
			if e.streamListConcurrency != test.num {
				tt.Errorf("got: %d, want: %d", e.streamListConcurrency, test.num)
			}
		})
	}
}

func TestWithKVSSyncInterval(t *testing.T) {
	t.Parallel()
	type test struct {
		name    string
		dur     string
		want    time.Duration
		wantErr bool
	}
	tests := []test{
		{name: "sets backgroundSyncInterval when dur is valid", dur: "5s", want: 5 * time.Second},
		{name: "no-op when dur is empty"},
		{name: "returns error when dur cannot be parsed", dur: "invalid", wantErr: true},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			e := new(export)
			err := WithKVSSyncInterval(test.dur)(e)
			if test.wantErr {
				if err == nil {
					tt.Error("expected an error but got nil")
				}
				return
			}
			if err != nil {
				tt.Errorf("err: %v", err)
			}
			if e.backgroundSyncInterval != test.want {
				tt.Errorf("got: %v, want: %v", e.backgroundSyncInterval, test.want)
			}
		})
	}
}

func TestWithKVSCompactionInterval(t *testing.T) {
	t.Parallel()
	type test struct {
		name    string
		dur     string
		want    time.Duration
		wantErr bool
	}
	tests := []test{
		{name: "sets backgroundCompactionInterval when dur is valid", dur: "5s", want: 5 * time.Second},
		{name: "no-op when dur is empty"},
		{name: "returns error when dur cannot be parsed", dur: "invalid", wantErr: true},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			e := new(export)
			err := WithKVSCompactionInterval(test.dur)(e)
			if test.wantErr {
				if err == nil {
					tt.Error("expected an error but got nil")
				}
				return
			}
			if err != nil {
				tt.Errorf("err: %v", err)
			}
			if e.backgroundCompactionInterval != test.want {
				tt.Errorf("got: %v, want: %v", e.backgroundCompactionInterval, test.want)
			}
		})
	}
}

func TestWithIndexPath(t *testing.T) {
	t.Parallel()
	type test struct {
		name    string
		path    string
		wantErr bool
	}
	tests := []test{
		{name: "sets indexPath when path is not empty", path: "/var/export/index"},
		{name: "returns invalid option error when path is empty", path: "", wantErr: true},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			e := new(export)
			err := WithIndexPath(test.path)(e)
			if test.wantErr {
				var invalidOpt *errors.ErrInvalidOption
				if !errors.As(err, &invalidOpt) {
					tt.Errorf("got_error: \"%v\", want a wrapped *errors.ErrInvalidOption", err)
				}
				return
			}
			if err != nil {
				tt.Errorf("err: %v", err)
			}
			if e.indexPath != test.path {
				tt.Errorf("got: %s, want: %s", e.indexPath, test.path)
			}
		})
	}
}

func TestWithGateway(t *testing.T) {
	t.Parallel()
	t.Run("sets gateway when client is not nil", func(tt *testing.T) {
		tt.Parallel()
		e := new(export)
		client := vald.NewValdClient(nil)
		if err := WithGateway(client)(e); err != nil {
			tt.Errorf("err: %v", err)
		}
		if e.gateway != client {
			tt.Error("gateway was not set to the given client")
		}
	})
	t.Run("returns critical option error when client is nil", func(tt *testing.T) {
		tt.Parallel()
		e := new(export)
		err := WithGateway(nil)(e)
		var critOpt *errors.ErrCriticalOption
		if !errors.As(err, &critOpt) {
			tt.Errorf("got_error: \"%v\", want a wrapped *errors.ErrCriticalOption", err)
		}
		if e.gateway != nil {
			tt.Error("gateway must remain nil")
		}
	})
}

func TestWithErrGroup(t *testing.T) {
	t.Parallel()
	t.Run("sets eg when eg is not nil", func(tt *testing.T) {
		tt.Parallel()
		e := new(export)
		eg := errgroup.Get()
		if err := WithErrGroup(eg)(e); err != nil {
			tt.Errorf("err: %v", err)
		}
		if e.eg != eg {
			tt.Error("eg was not set to the given errgroup")
		}
	})
	t.Run("no-op when eg is nil", func(tt *testing.T) {
		tt.Parallel()
		e := new(export)
		if err := WithErrGroup(nil)(e); err != nil {
			tt.Errorf("err: %v", err)
		}
		if e.eg != nil {
			tt.Error("eg must remain nil")
		}
	})
}

// NOT IMPLEMENTED BELOW

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
package usecase

import (
	"testing"

	iconf "github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/pkg/manager/index/config"
)

// auto_save_index_wait_duration must reach the indexer service: an unparsable value is
// rejected by service.WithSaveIndexDuration, which only happens when New passes it through.
func TestNew_wiresAutoSaveIndexWaitDuration(t *testing.T) {
	t.Parallel()
	cfg := &config.Data{
		Server:        new(iconf.Servers),
		Observability: new(iconf.Observability),
		Indexer: &iconf.Indexer{
			Discoverer:                new(iconf.DiscovererClient),
			AutoSaveIndexWaitDuration: "invalid",
		},
	}
	_, err := New(cfg)
	if err == nil {
		t.Fatal("New() error = nil, want the WithSaveIndexDuration option failure")
	}
	if want := "service.WithSaveIndexDuration.func"; !strings.Contains(err.Error(), want) {
		t.Errorf("New() error = %v, want it to mention %q", err, want)
	}
}

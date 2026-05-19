//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	err := os.WriteFile(path, []byte(`{
		"version": "v1.0.0",
		"server_config": {
			"full_shutdown_duration": "10ms"
		},
		"client": {
			"addrs": ["127.0.0.1:8081"],
			"dial_option": {
				"insecure": true
			}
		},
		"llm": {
			"openai": {
				"token": "test-token",
				"model": "small"
			}
		}
	}`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	cfg, err := NewConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.LLM.Provider != "openai" {
		t.Errorf("provider = %s, want openai", cfg.LLM.Provider)
	}
	if cfg.LLM.OpenAI.Token != "test-token" {
		t.Errorf("token = %s, want test-token", cfg.LLM.OpenAI.Token)
	}
}

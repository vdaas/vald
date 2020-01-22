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

// package config providers configuration type and load configuration logic
package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/vdaas/vald/internal/net/http/json"
	yaml "gopkg.in/yaml.v2"
)

// Default represent a application setting data content (config.yaml).
type Default struct {
	// Version represent configuration file version.
	Version string `json:"version" yaml:"version"`

	// TZ represent system time location .
	TZ string `json:"time_zone" yaml:"time_zone"`
}

func (c *Default) Bind() *Default {
	c.Version = GetActualValue(c.Version)
	c.TZ = GetActualValue(c.TZ)
	return c
}

func (c *Default) UnmarshalJSON(data []byte) (err error) {
	ic := new(struct {
		Ver string `json:"version"`
		TZ  string `json:"time_zone"`
	})
	err = json.Unmarshal(data, &ic)
	if err != nil {
		return err
	}
	c.Version = ic.Ver
	c.TZ = ic.TZ
	return nil
}

// New returns config struct or error when decode the configuration file to actually *Config struct.
func Read(path string, cfg interface{}) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	switch filepath.Ext(path) {
	case ".yaml":
		err = yaml.NewDecoder(f).Decode(cfg)
	case ".json":
		err = json.Decode(f, cfg)
	}
	return err
}

// GetActualValue returns the environment variable value if the val has prefix and suffix "_", otherwise the val will directly return.
func GetActualValue(val string) string {
	if checkPrefixAndSuffix(val, "_", "_") {
		return os.ExpandEnv(os.Getenv(strings.TrimPrefix(strings.TrimSuffix(val, "_"), "_")))
	}
	return os.ExpandEnv(val)
}

func GetActualValues(vals []string) []string {
	for i, val := range vals {
		vals[i] = GetActualValue(val)
	}
	return vals
}

// checkPrefixAndSuffix checks if the str has prefix and suffix
func checkPrefixAndSuffix(str, pref, suf string) bool {
	return strings.HasPrefix(str, pref) && strings.HasSuffix(str, suf)
}

func ToRawYaml(data interface{}) string {
	buf := bytes.NewBuffer(nil)
	yaml.NewEncoder(buf).Encode(data)
	return buf.String()
}

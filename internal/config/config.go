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

// package config providers configuration type and load configuration logic
package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io/ioutil"
	"github.com/vdaas/vald/internal/log"
	yaml "gopkg.in/yaml.v2"
)

// GlobalConfig represent a application setting data content (config.yaml).
type GlobalConfig struct {
	// Version represent configuration file version.
	Version string `json:"version" yaml:"version"`

	// TZ represent system time location .
	TZ string `json:"time_zone" yaml:"time_zone"`

	// Log represent log configuration.
	Logging *Logging `json:"logging,omitempty" yaml:"logging,omitempty"`
}

const (
	fileValuePrefix = "file://"
	envSymbol       = "_"
)

func (c *GlobalConfig) Bind() *GlobalConfig {
	c.Version = GetActualValue(c.Version)
	c.TZ = GetActualValue(c.TZ)

	if c.Logging != nil {
		c.Logging = c.Logging.Bind()
	}
	return c
}

func (c *GlobalConfig) UnmarshalJSON(data []byte) (err error) {
	ic := new(struct {
		Ver     string   `json:"version"`
		TZ      string   `json:"time_zone"`
		Logging *Logging `json:"logging"`
	})
	err = json.Unmarshal(data, &ic)
	if err != nil {
		return err
	}
	c.Version = ic.Ver
	c.TZ = ic.TZ
	c.Logging = ic.Logging
	return nil
}

// New returns config struct or error when decode the configuration file to actually *Config struct.
func Read(path string, cfg interface{}) (err error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Wrap(f.Close(), err.Error())
			return
		}
		err = f.Close()
	}()
	switch filepath.Ext(path) {
	case ".yaml":
		err = yaml.NewDecoder(f).Decode(cfg)
	case ".json":
		err = json.Decode(f, cfg)
	}
	return err
}

// GetActualValue returns the environment variable value if the val has prefix and suffix "_",
// if actual value start with file://{path} the return value will read from file
// otherwise the val will directly return.
func GetActualValue(val string) (res string) {
	if checkPrefixAndSuffix(val, envSymbol, envSymbol) {
		val = strings.TrimPrefix(strings.TrimSuffix(val, envSymbol), envSymbol)
		if !strings.HasPrefix(val, "$") {
			val = "$" + val
		}
	}
	res = os.ExpandEnv(val)
	if strings.HasPrefix(res, fileValuePrefix) {
		body, err := ioutil.ReadFile(strings.TrimPrefix(res, fileValuePrefix))
		if err != nil {
			return
		}
		res = *(*string)(unsafe.Pointer(&body))
	}
	return
}

func GetActualValues(vals []string) []string {
	for i, val := range vals {
		vals[i] = GetActualValue(val)
	}
	return vals
}

// checkPrefixAndSuffix checks if the str has prefix and suffix.
func checkPrefixAndSuffix(str, pref, suf string) bool {
	return strings.HasPrefix(str, pref) && strings.HasSuffix(str, suf)
}

func ToRawYaml(data interface{}) string {
	buf := bytes.NewBuffer(nil)
	err := yaml.NewEncoder(buf).Encode(data)
	if err != nil {
		log.Error(err)
	}
	return buf.String()
}

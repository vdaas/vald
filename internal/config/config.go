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

// package config providers configuration type and load configuration logic
package config

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
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

// Bind binds the actual data from the receiver field.
func (c *GlobalConfig) Bind() *GlobalConfig {
	c.Version = GetActualValue(c.Version)
	c.TZ = GetActualValue(c.TZ)

	if c.Logging != nil {
		c.Logging = c.Logging.Bind()
	}
	return c
}

// Read returns config struct or error when decoding the configuration file to actually *Config struct.
func Read(path string, cfg interface{}) (err error) {
	f, err := file.Open(path, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		if f != nil {
			if err != nil {
				err = errors.Join(f.Close(), err)
				return
			}
			err = f.Close()
		}
	}()
	switch ext := filepath.Ext(path); ext {
	case ".yaml", ".yml":
		err = yaml.NewDecoder(f).Decode(cfg)
	case ".json":
		err = json.Decode(f, cfg)
	default:
		err = errors.ErrUnsupportedConfigFileType(ext)
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
		body, err := file.ReadFile(strings.TrimPrefix(res, fileValuePrefix))
		if err != nil || body == nil {
			return
		}
		res = conv.Btoa(body)
	}
	return
}

// GetActualValues returns the environment variable values if the vals has string slice that has prefix and suffix "_",
// if actual value start with file://{path} the return value will read from file
// otherwise the val will directly return.
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

// ToRawYaml writes the YAML encoding of v to the stream and returns the string written to stream.
func ToRawYaml(data interface{}) string {
	buf := bytes.NewBuffer(nil)
	err := yaml.NewEncoder(buf).Encode(data)
	if err != nil {
		log.Error(err)
	}
	return buf.String()
}

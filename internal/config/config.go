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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"

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

// Merge merges multiple objects to one object.
// the value of each field is prioritized the value of last index of `objs`.
// if the length of `objs` is zero, it returns initial value of type T.
func Merge[T any](objs ...T) (dst T, err error) {
	switch len(objs) {
	case 0:
		return dst, nil
	case 1:
		dst = objs[0]
		return dst, nil
	default:
		dst = objs[0]
		visited := make(map[uintptr]bool)
		rdst := reflect.ValueOf(&dst)
		for _, src := range objs[1:] {
			err = deepMerge(rdst, reflect.ValueOf(&src), visited, "")
			if err != nil {
				return dst, err
			}
		}
	}
	return dst, err
}

func deepMerge(dst, src reflect.Value, visited map[uintptr]bool, fieldPath string) (err error) {
	if !src.IsValid() || src.IsZero() {
		return nil
	} else if !dst.IsValid() {
		dst = src
		log.Info(dst.Type(), dst, src)
	}
	dType := dst.Type()
	sType := src.Type()
	if dType != sType {
		return errors.ErrNotMatchFieldType(fieldPath, dType, sType)
	}
	sKind := src.Kind()
	if sKind == reflect.Ptr {
		src = src.Elem()
	}
	if sKind == reflect.Struct && src.CanAddr() {
		addr := src.Addr().Pointer()
		if visited[addr] {
			return nil
		}
		if src.NumField() > 1 {
			visited[addr] = true
		}
	}
	switch dst.Kind() {
	case reflect.Ptr:
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		return deepMerge(dst.Elem(), src, visited, fieldPath)
	case reflect.Struct:
		dnum := dst.NumField()
		snum := src.NumField()
		if dnum != snum {
			return errors.ErrNotMatchFieldNum(fieldPath, dnum, snum)
		}
		for i := 0; i < dnum; i++ {
			dstField := dst.Field(i)
			if dstField.CanSet() {
				nf := fmt.Sprintf("%s.%s(%d)", fieldPath, dType.Field(i).Name, i)
				if err = deepMerge(dstField, src.Field(i), visited, nf); err != nil {
					return errors.ErrDeepMergeKind(dst.Kind().String(), nf, err)
				}
			}
		}
	case reflect.Slice:
		srcLen := src.Len()
		if srcLen > 0 {
			if dst.IsNil() {
				dst.Set(reflect.MakeSlice(dType, srcLen, srcLen))
			} else {
				diffLen := srcLen - dst.Len()
				if diffLen > 0 {
					dst.Set(reflect.AppendSlice(dst, reflect.MakeSlice(dType, diffLen, diffLen)))
				}
			}
			for i := 0; i < srcLen; i++ {
				nf := fmt.Sprintf("%s[%d]", fieldPath, i)
				if err = deepMerge(dst.Index(i), src.Index(i), visited, nf); err != nil {
					return errors.ErrDeepMergeKind(dst.Kind().String(), nf, err)
				}
			}
		}
	case reflect.Array:
		srcLen := src.Len()
		if srcLen != dst.Len() {
			return errors.ErrNotMatchArrayLength(fieldPath, dst.Len(), srcLen)
		}
		for i := 0; i < srcLen; i++ {
			nf := fmt.Sprintf("%s[%d]", fieldPath, i)
			if err = deepMerge(dst.Index(i), src.Index(i), visited, nf); err != nil {
				return errors.ErrDeepMergeKind(dst.Kind().String(), nf, err)
			}
		}
	case reflect.Map:
		if dst.IsNil() {
			dst.Set(reflect.MakeMapWithSize(dType, src.Len()))
		}
		dElem := dType.Elem()
		for _, key := range src.MapKeys() {
			vdst := dst.MapIndex(key)
			// fmt.Println(vdst.IsValid(), key, vdst)
			if !vdst.IsValid() {
				vdst = reflect.New(dElem).Elem()
			}
			nf := fmt.Sprintf("%s[%s]", fieldPath, key)
			if vdst.CanSet() {
				if err = deepMerge(vdst, src.MapIndex(key), visited, nf); err != nil {
					return errors.Errorf("error in array at %s: %w", nf, err)
				}
				dst.SetMapIndex(key, vdst)
			} else {
				dst.SetMapIndex(key, src.MapIndex(key))
			}
		}
	default:
		if dst.CanSet() {
			dst.Set(src)
		}
	}
	return nil
}

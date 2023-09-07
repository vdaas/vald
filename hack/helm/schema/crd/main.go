//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	yaml "gopkg.in/yaml.v2"
)

const (
	objectType = "object"
	arrayType  = "array"

	prefix = "# @schema"

	minimumArgumentLength = 2
)

var aliases map[string]Schema

type Schema struct {
	Type        string             `json:"type"                  yaml:"type"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"  yaml:"properties,omitempty"`

	// for object type
	Required          []string          `json:"required,omitempty"          yaml:"required,omitempty"`
	MaxProperties     *uint64           `json:"maxProperties,omitempty"     yaml:"maxProperties,omitempty"`
	MinProperties     *uint64           `json:"minProperties,omitempty"     yaml:"minProperties,omitempty"`
	DependentRequired map[string]string `json:"dependentRequired,omitempty" yaml:"dependentRequired,omitempty"`

	// for string type
	Enum      []string `json:"enum,omitempty"      yaml:"enum,omitempty"`
	Pattern   string   `json:"pattern,omitempty"   yaml:"pattern,omitempty"`
	MaxLength *uint64  `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength *uint64  `json:"minLength,omitempty" yaml:"minLength,omitempty"`

	// for array type
	Items       *Schema `json:"items,omitempty"       yaml:"items,omitempty"`
	MaxItems    *uint64 `json:"maxItems,omitempty"    yaml:"maxItems,omitempty"`
	MinItems    *uint64 `json:"minItems,omitempty"    yaml:"minItems,omitempty"`
	UniqueItems bool    `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	MaxContains *uint64 `json:"maxContains,omitempty" yaml:"maxContains,omitempty"`
	MinContains *uint64 `json:"minContains,omitempty" yaml:"minContains,omitempty"`

	// for numeric types
	MultipleOf       *int64 `json:"multipleOf,omitempty"       yaml:"multipleOf,omitempty"`
	Maximum          *int64 `json:"maximum,omitempty"          yaml:"maximum,omitempty"`
	ExclusiveMaximum bool   `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum          *int64 `json:"minimum,omitempty"          yaml:"minimum,omitempty"`
	ExclusiveMinimum bool   `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`

	// for Kubernetes unknown object type
	KubernetesPreserveUnknownFields bool `json:"x-kubernetes-preserve-unknown-fields,omitempty" yaml:"x-kubernetes-preserve-unknown-fields,omitempty"`
}

type VSchema struct {
	Name   string `json:"name"   yaml:"name"`
	Type   string `json:"type"   yaml:"type"`
	Anchor string `json:"anchor" yaml:"anchor"`
	Alias  string `json:"alias"  yaml:"alias"`
	Schema
}

type Spec struct {
	Spec *Schema `json:"spec" yaml:"spec"`
}

func main() {
	log.Init()
	if len(os.Args) < minimumArgumentLength {
		log.Fatal(errors.New("invalid argument: must be specify path to the values.yaml"))
	}

	var err error
	nindent := 0

	if len(os.Args) > minimumArgumentLength {
		nindent, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid argument: %s", err)
		}
	}

	err = genSchema(os.Args[1], nindent)
	if err != nil {
		log.Fatal(err)
	}
}

func genSchema(path string, nindent int) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_SYNC, fs.ModePerm)
	if err != nil {
		return errors.Errorf("cannot open %s", path)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			// skipcq: RVV-A0003
			log.Fatal(err)
		}
	}()

	aliases = make(map[string]Schema)

	ls := make([]*VSchema, 0)

	var line uint64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line++
		tx := strings.TrimLeft(sc.Text(), " ")
		if strings.HasPrefix(tx, prefix) {
			l := new(VSchema)
			err = json.Unmarshal([]byte(strings.TrimPrefix(tx, prefix)), &l)
			if err != nil {
				log.Errorf("error occurred line %d, data %s, error %v", line, tx, err)
			}
			ls = append(ls, l)
		}
	}

	schemas, err := objectProperties(ls)
	if err != nil {
		return errors.Errorf("error: %s", err)
	}

	yaml, err := yaml.Marshal(newSpec(schemas))
	if err != nil {
		return errors.Errorf("error: %s", err)
	}

	fmt.Println(indent(conv.Btoa(yaml), nindent))

	return nil
}

func objectProperties(ls []*VSchema) (map[string]*Schema, error) {
	if len(ls) == 0 {
		return nil, errors.New("empty list")
	}

	groups := make(map[string][]*VSchema)
	gOrder := make([]string, 0, len(ls))
	for _, l := range ls {
		root := strings.Split(l.Name, ".")
		if groups[root[0]] == nil {
			gOrder = append(gOrder, root[0])
		}
		groups[root[0]] = append(groups[root[0]], l)
	}

	schemas := make(map[string]*Schema)
	for _, k := range gOrder {
		s, err := genNode(groups[k])
		if err != nil {
			return nil, errors.Errorf("error: %s", err)
		}
		schemas[k] = s
	}

	return schemas, nil
}

func genNode(ls []*VSchema) (*Schema, error) {
	if len(ls) == 0 {
		return nil, errors.New("empty list")
	}

	l := ls[0]

	if l.Alias != "" {
		schema, ok := aliases[l.Alias]
		if !ok {
			return nil, errors.Errorf("unknown alias %s", l.Alias)
		}
		return &schema, nil
	}

	var schema Schema

	switch l.Type {
	case objectType:
		schema = l.Schema
		if len(ls) <= 1 {
			schema.Type = objectType
			schema.KubernetesPreserveUnknownFields = true
			break
		}

		nls := make([]*VSchema, 0, len(ls[1:]))
		for _, nl := range ls[1:] {
			nl.Name = strings.TrimLeft(strings.TrimPrefix(nl.Name, l.Name), ".")
			nls = append(nls, nl)
		}

		ps, err := objectProperties(nls)
		if err != nil {
			return nil, errors.Errorf("error: %s", err)
		}
		schema.Type = objectType
		schema.Properties = ps
	case arrayType:
		schema = l.Schema
		schema.Type = l.Type
		if schema.Items != nil && schema.Items.Type == objectType && schema.Items.Properties == nil {
			schema.Items.KubernetesPreserveUnknownFields = true
		}
	default:
		schema = l.Schema
		schema.Type = l.Type
	}

	if l.Anchor != "" {
		aliases[l.Anchor] = schema
	}

	return &schema, nil
}

func newSpec(schemas map[string]*Schema) *Spec {
	return &Spec{
		Spec: &Schema{
			Type:       objectType,
			Properties: schemas,
		},
	}
}

func indent(text string, nindent int) string {
	indent := ""

	for i := 0; i < nindent; i++ {
		indent += " "
	}

	if text[len(text)-1:] == "\n" {
		result := ""
		for _, j := range strings.Split(text[:len(text)-1], "\n") {
			result += indent + j + "\n"
		}
		return result
	}
	result := ""
	for _, j := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		result += indent + j + "\n"
	}
	return result[:len(result)-1]
}

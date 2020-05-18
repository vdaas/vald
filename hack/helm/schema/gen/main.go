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

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

const (
	nullType   = "null"
	boolType   = "boolean"
	objectType = "object"
	arrayType  = "array"
	numberType = "number"
	stringType = "string"
	intType    = "integer"

	prefix = "# @schema"
)

type SchemaBase struct {
	// for object type
	Required          []string          `json:"required,omitempty"`
	MaxProperties     uint64            `json:"maxProperties,omitempty"`
	MinProperties     uint64            `json:"minProperties,omitempty"`
	DependentRequired map[string]string `json:"dependentRequired,omitempty"`

	// for string type
	Enum      []string `json:"enum,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	MaxLength uint64   `json:"maxLength,omitempty"`
	MinLength uint64   `json:"minLength,omitempty"`

	// for array type
	Items       *Schema `json:"items,omitempty"`
	MaxItems    uint64  `json:"maxItems,omitempty"`
	MinItems    uint64  `json:"minItems,omitempty"`
	UniqueItems bool    `json:"uniqueItems,omitempty"`
	MaxContains uint64  `json:"maxContains,omitempty"`
	MinContains uint64  `json:"minContains,omitempty"`

	// for numeric types
	MultipleOf       int64 `json:"multipleOf,omitempty"`
	Maximum          int64 `json:"maximum,omitempty"`
	ExclusiveMaximum int64 `json:"exclusiveMaximum,omitempty"`
	Minimum          int64 `json:"minimum,omitempty"`
	ExclusiveMinimum int64 `json:"exclusiveMinimum,omitempty"`
}

type VSchema struct {
	Name string `json:"name"`
	Type string `json:"type"`
	SchemaBase
}

type Root struct {
	SchemaKeyword string `json:"$schema"`
	Title         string `json:"title"`
	Schema
}

type Schema struct {
	Type        string             `json:"type"`
	Description string             `json:"description,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`

	SchemaBase
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("invalid argument: must be specify path to the values.yaml"))
	}
	err := genJSONSchema(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func genJSONSchema(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Errorf("cannot open %s", path)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	ls := make([]VSchema, 0)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		tx := strings.TrimLeft(sc.Text(), " ")
		if strings.HasPrefix(tx, prefix) {
			var l VSchema
			err = json.Unmarshal([]byte(strings.TrimPrefix(tx, prefix)), &l)
			if err != nil {
				log.Error(err)
			}
			ls = append(ls, l)
		}
	}

	schemas, err := objectProperties(ls)
	if err != nil {
		return errors.Errorf("error: %s", err)
	}

	json, err := json.Marshal(newRoot(schemas))
	// json, err := json.Marshal(ls)
	if err != nil {
		return errors.Errorf("error: %s", err)
	}

	fmt.Println(string(json))

	return nil
}

func objectProperties(ls []VSchema) (map[string]*Schema, error) {
	if len(ls) <= 0 {
		return nil, errors.New("empty list")
	}

	groups := make(map[string][]VSchema)
	for _, l := range ls {
		root := strings.Split(l.Name, ".")
		groups[root[0]] = append(groups[root[0]], l)
	}

	schemas := make(map[string]*Schema)
	for k, v := range groups {
		s, err := genNode(v)
		if err != nil {
			return nil, errors.Errorf("error: %s", err)
		}
		schemas[k] = s
	}

	return schemas, nil
}

func genNode(ls []VSchema) (*Schema, error) {
	if len(ls) <= 0 {
		return nil, errors.New("empty list")
	}

	l := ls[0]
	switch l.Type {
	case objectType:
		if len(ls) <= 1 {
			return &Schema{
				Type:       objectType,
				SchemaBase: l.SchemaBase,
			}, nil
		}

		nls := make([]VSchema, 0, len(ls[1:]))
		for _, nl := range ls[1:] {
			nl.Name = strings.TrimLeft(strings.TrimPrefix(nl.Name, l.Name), ".")
			nls = append(nls, nl)
		}

		ps, err := objectProperties(nls)
		if err != nil {
			return nil, errors.Errorf("error: %s", err)
		}
		return &Schema{
			Type:       objectType,
			Properties: ps,
			SchemaBase: l.SchemaBase,
		}, nil
	default:
		return &Schema{
			Type:       l.Type,
			SchemaBase: l.SchemaBase,
		}, nil
	}
}

func newRoot(schemas map[string]*Schema) *Root {
	return &Root{
		SchemaKeyword: "http://json-schema.org/draft-07/schema#",
		Title:         "Values",
		Schema: Schema{
			Type:       objectType,
			Properties: schemas,
		},
	}
}

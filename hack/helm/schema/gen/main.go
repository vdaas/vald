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
	objectType = "object"
	arrayType  = "array"
	stringType = "string"
	intType    = "integer"
	boolType   = "boolean"

	prefix = "# @schema"
)

type VSchema struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Required []string `json:"required,omitempty"`
	Pattern  string   `json:"pattern,omitempty"`
}

type Root struct {
	SchemaKeyword string `json:"$schema"`
	Title         string `json:"title"`
	Schema
}

type Schema struct {
	Type        string             `json:"type"`
	Description string             `json:"description,omitempty"`
	Pattern     string             `json:"pattern,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
	Required    []string           `json:"required,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
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
				Type:     objectType,
				Required: l.Required,
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
			Required:   l.Required,
			Properties: ps,
		}, nil
	case arrayType:
		return &Schema{
			Type: arrayType,
		}, nil
	case stringType:
		return &Schema{
			Type:    stringType,
			Pattern: l.Pattern,
		}, nil
	case intType:
		return &Schema{
			Type: intType,
		}, nil
	case boolType:
		return &Schema{
			Type: boolType,
		}, nil
	default:
		return &Schema{
			Type: l.Type,
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

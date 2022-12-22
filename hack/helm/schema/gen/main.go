//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"io/fs"
	"os"
	"regexp"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
)

const (
	objectType = "object"

	prefix = "# @schema"

	minimumArgumentLength = 2
)

var (
	aliases      map[string]Schema
	descriptions map[string]string

	descriptionRegexp   = regexp.MustCompile(`^\s*#\s*(.*)\s+--\s*(.*)$`)
	continuedLineRegexp = regexp.MustCompile(`^\s*#\s+(.*)$`)
)

type SchemaBase struct {
	// for object type
	Required          []string          `json:"required,omitempty"`
	MaxProperties     *uint64           `json:"maxProperties,omitempty"`
	MinProperties     *uint64           `json:"minProperties,omitempty"`
	DependentRequired map[string]string `json:"dependentRequired,omitempty"`

	// for string type
	Enum      []string `json:"enum,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	MaxLength *uint64  `json:"maxLength,omitempty"`
	MinLength *uint64  `json:"minLength,omitempty"`

	// for array type
	Items       *Schema `json:"items,omitempty"`
	MaxItems    *uint64 `json:"maxItems,omitempty"`
	MinItems    *uint64 `json:"minItems,omitempty"`
	UniqueItems bool    `json:"uniqueItems,omitempty"`
	MaxContains *uint64 `json:"maxContains,omitempty"`
	MinContains *uint64 `json:"minContains,omitempty"`

	// for numeric types
	MultipleOf       *int64 `json:"multipleOf,omitempty"`
	Maximum          *int64 `json:"maximum,omitempty"`
	ExclusiveMaximum bool   `json:"exclusiveMaximum,omitempty"`
	Minimum          *int64 `json:"minimum,omitempty"`
	ExclusiveMinimum bool   `json:"exclusiveMinimum,omitempty"`
}

type VSchema struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Anchor string `json:"anchor"`
	Alias  string `json:"alias"`
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
	log.Init()
	if len(os.Args) < minimumArgumentLength {
		log.Fatal(errors.New("invalid argument: must be specify path to the values.yaml"))
	}
	err := genJSONSchema(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func genJSONSchema(path string) error {
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
	descriptions = make(map[string]string)

	ls := make([]*VSchema, 0)

	continuedLine := false
	currentKey := ""

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
			continue
		}

		if continuedLine {
			match := continuedLineRegexp.FindStringSubmatch(tx)
			if len(match) > 1 {
				descriptions[currentKey] += " " + match[1]

				continue
			}

			continuedLine = false

			continue
		}

		match := descriptionRegexp.FindStringSubmatch(tx)
		if len(match) < 3 || match[1] == "" {
			continue
		}

		currentKey = match[1]
		descriptions[currentKey] = match[2]
		continuedLine = true
	}

	schemas, err := objectProperties(make([]string, 0), ls)
	if err != nil {
		return errors.Errorf("error: %s", err)
	}

	json, err := json.Marshal(newRoot(schemas))
	if err != nil {
		return errors.Errorf("error: %s", err)
	}

	fmt.Println(string(json))

	return nil
}

func objectProperties(prefix []string, ls []*VSchema) (map[string]*Schema, error) {
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
		s, err := genNode(prefix, groups[k])
		if err != nil {
			return nil, errors.Errorf("error: %s", err)
		}
		schemas[k] = s
	}

	return schemas, nil
}

func genNode(prefix []string, ls []*VSchema) (*Schema, error) {
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

	description := descriptions[strings.Join(append(prefix, l.Name), ".")]

	switch l.Type {
	case objectType:
		if len(ls) <= 1 {
			schema = Schema{
				Type:        objectType,
				Description: description,
				SchemaBase:  l.SchemaBase,
			}
			break
		}

		nls := make([]*VSchema, 0, len(ls[1:]))
		for _, nl := range ls[1:] {
			nl.Name = strings.TrimLeft(strings.TrimPrefix(nl.Name, l.Name), ".")
			nls = append(nls, nl)
		}

		ps, err := objectProperties(append(prefix, l.Name), nls)
		if err != nil {
			return nil, errors.Errorf("error: %s", err)
		}
		schema = Schema{
			Type:        objectType,
			Description: description,
			Properties:  ps,
			SchemaBase:  l.SchemaBase,
		}
	default:
		schema = Schema{
			Type:        l.Type,
			Description: description,
			SchemaBase:  l.SchemaBase,
		}
	}

	if l.Anchor != "" {
		aliases[l.Anchor] = schema
	}

	return &schema, nil
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

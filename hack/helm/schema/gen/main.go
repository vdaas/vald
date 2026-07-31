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

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"regexp"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/os"
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

	// refsMode, when true, makes anchors emit as JSON Schema $defs and aliases
	// as $ref (instead of inlining). This yields named, reusable Go types when
	// the schema is fed to a code generator. defs collects the $defs entries.
	refsMode bool
	defs     map[string]*Schema

	descriptionRegexp   = regexp.MustCompile(`^\s*#\s*(.*)\s+--\s*(.*)$`)
	continuedLineRegexp = regexp.MustCompile(`^\s*#\s+(.*)$`)
)

// defName sanitizes an anchor/alias name into a $defs key (ASCII identifier).
func defName(anchor string) string {
	b := make([]byte, 0, len(anchor))
	for i := 0; i < len(anchor); i++ {
		c := anchor[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
			b = append(b, c)
		default:
			b = append(b, '_')
		}
	}
	return string(b)
}

// GoTypeImport describes the Go import backing an x-go-type mapping.
type GoTypeImport struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
}

type SchemaBase struct {
	// XGoType / XGoTypeImport carry a Go type identity for a node so that
	// downstream code generators (e.g. oapi-codegen) can map k8s-native
	// subtrees (affinity, resources, ...) to their real Go types instead of
	// emitting anonymous structs. They are passed through untouched into the
	// generated JSON Schema.
	XGoType           string            `json:"x-go-type,omitempty"`
	XGoTypeImport     *GoTypeImport     `json:"x-go-type-import,omitempty"`
	MaxContains       *uint64           `json:"maxContains,omitempty"`
	MinContains       *uint64           `json:"minContains,omitempty"`
	MinProperties     *uint64           `json:"minProperties,omitempty"`
	MaxItems          *uint64           `json:"maxItems,omitempty"`
	Minimum           *int64            `json:"minimum,omitempty"`
	Maximum           *int64            `json:"maximum,omitempty"`
	MaxLength         *uint64           `json:"maxLength,omitempty"`
	MinLength         *uint64           `json:"minLength,omitempty"`
	MaxProperties     *uint64           `json:"maxProperties,omitempty"`
	MinItems          *uint64           `json:"minItems,omitempty"`
	DependentRequired map[string]string `json:"dependentRequired,omitempty"`
	MultipleOf        *int64            `json:"multipleOf,omitempty"`
	Items             *Schema           `json:"items,omitempty"`
	Pattern           string            `json:"pattern,omitempty"`
	Required          []string          `json:"required,omitempty"`
	Enum              []string          `json:"enum,omitempty"`
	UniqueItems       bool              `json:"uniqueItems,omitempty"`
	ExclusiveMaximum  bool              `json:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum  bool              `json:"exclusiveMinimum,omitempty"`
}

type VSchema struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Anchor string `json:"anchor"`
	Alias  string `json:"alias"`
	SchemaBase
}

type Root struct {
	SchemaKeyword string             `json:"$schema"`
	Title         string             `json:"title"`
	Defs          map[string]*Schema `json:"$defs,omitempty"`
	Schema
}

type Schema struct {
	Ref         string             `json:"$ref,omitempty"`
	Type        string             `json:"type,omitempty"`
	Description string             `json:"description,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`

	SchemaBase
}

func main() {
	log.Init()
	refs := flag.Bool("refs", false, "emit $defs/$ref for anchors instead of inlining (for Go type generation)")
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal(errors.New("invalid argument: must be specify path to the values.yaml"))
	}
	refsMode = *refs
	if err := genJSONSchema(flag.Arg(0)); err != nil {
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
	defs = make(map[string]*Schema)

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

	root := newRoot(schemas)
	if refsMode {
		root.Defs = defs
	}
	json, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return errors.Errorf("error: %s", err)
	}

	fmt.Println(conv.Btoa(json))

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
		if refsMode {
			return &Schema{Ref: "#/$defs/" + defName(l.Alias)}, nil
		}
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
		if refsMode {
			s := schema
			defs[defName(l.Anchor)] = &s
			return &Schema{Ref: "#/$defs/" + defName(l.Anchor)}, nil
		}
		aliases[l.Anchor] = schema
	}

	// In refsMode, hoist every object-with-properties into $defs (named by its
	// path) so the code generator emits named Go types for all nested objects,
	// not anonymous structs. This is required for controller-gen deepcopy.
	if refsMode && schema.Type == objectType && len(schema.Properties) > 0 {
		segs := make([]string, 0, len(prefix)+1)
		segs = append(segs, prefix...)
		segs = append(segs, l.Name)
		name := defName(strings.Join(segs, "_"))
		s := schema
		defs[name] = &s
		return &Schema{Ref: "#/$defs/" + name}, nil
	}

	return &schema, nil
}

func newRoot(schemas map[string]*Schema) *Root {
	return &Root{
		SchemaKeyword: "https://json-schema.org/draft-07/schema#",
		Title:         "Values",
		Schema: Schema{
			Type:       objectType,
			Properties: schemas,
		},
	}
}

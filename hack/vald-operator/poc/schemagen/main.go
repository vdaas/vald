// Command schemagen generates Go structs from vald's
// charts/vald/values.schema.json (JSON Schema draft-07) and writes them after
// the "// +schemagen:begin" marker of the target file, leaving the package
// clause, consts and everything above the marker intact.
//
// Usage:
//
//     schemagen <values.schema.json> <root.path> <RootTypeName> <targetFile>
//
// Example (regenerate the agent package's struct in place):
//
//     go run ./poc/schemagen ../../charts/vald/values.schema.json agent Agent internal/pkg/api/valdrelease/agent/agent.go
//
// Struct names are derived from the field path relative to <root.path>, so the
// root object is <RootTypeName> and e.g. agent.ngt -> NGT, agent.ngt.kvsdb ->
// NGTKVSDB. Types are taken straight from the schema (no k8s type mapping);
// free-form objects become map[string]any and every field is
// omitempty. Enum values are emitted as a trailing comment.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Schema struct {
	Type       string             `json:"type"`
	Properties map[string]*Schema `json:"properties"`
	Items      *Schema            `json:"items"`
	Enum       []any              `json:"enum"`
}

var initialisms = map[string]string{
	"ngt": "NGT", "id": "ID", "cpu": "CPU", "grpc": "GRPC", "http": "HTTP",
	"tls": "TLS", "url": "URL", "api": "API", "ip": "IP", "pv": "PV",
	"otlp": "OTLP", "hpa": "HPA", "vqueue": "VQueue", "kvsdb": "KVSDB",
}

func goName(seg string) string {
	if v, ok := initialisms[seg]; ok {
		return v
	}
	if seg == "" {
		return ""
	}
	return strings.ToUpper(seg[:1]) + seg[1:]
}

func typeName(path []string) string {
	var b strings.Builder
	for _, p := range path {
		for _, part := range strings.Split(p, "_") {
			b.WriteString(goName(part))
		}
	}
	return b.String()
}

var (
	out     strings.Builder
	emitted = map[string]bool{}
	rootType string
)

func structName(rel []string) string {
	if len(rel) == 0 {
		return rootType
	}
	return typeName(rel)
}

func sub(rel []string, seg string) []string {
	return append(append([]string(nil), rel...), seg)
}

func goType(s *Schema, rel []string) string {
	switch s.Type {
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "string":
		return "string"
	case "array":
		if s.Items == nil {
			return "[]any"
		}
		return "[]" + goType(s.Items, sub(rel, "Item"))
	case "object":
		if len(s.Properties) == 0 {
			return "map[string]any"
		}
		name := structName(rel)
		emitStruct(name, s, rel)
		return "*" + name
	default:
		return "any"
	}
}

func emitStruct(name string, s *Schema, rel []string) {
	if emitted[name] {
		return
	}
	emitted[name] = true

	keys := make([]string, 0, len(s.Properties))
	for k := range s.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	body := &strings.Builder{}
	fmt.Fprintf(body, "type %s struct {\n", name)
	for _, k := range keys {
		field := s.Properties[k]
		ft := goType(field, sub(rel, k))
		fieldName := typeName([]string{k})
		comment := ""
		if len(field.Enum) > 0 {
			vals := make([]string, 0, len(field.Enum))
			for _, e := range field.Enum {
				vals = append(vals, fmt.Sprintf("%v", e))
			}
			comment = " // enum: " + strings.Join(vals, ", ")
		}
		fmt.Fprintf(body, "\t%s %s `json:%q yaml:%q`%s\n", fieldName, ft, k+",omitempty", k+",omitempty", comment)
	}
	body.WriteString("}\n\n")
	out.WriteString(body.String())
}

func resolve(root *Schema, dotted string) *Schema {
	node := root
	for _, p := range strings.Split(dotted, ".") {
		if node.Properties == nil {
			return nil
		}
		next, ok := node.Properties[p]
		if !ok {
			return nil
		}
		node = next
	}
	return node
}

const beginMarker = "// +schemagen:begin"

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintln(os.Stderr, "usage: schemagen <values.schema.json> <root.path> <RootTypeName> <targetFile>")
		os.Exit(2)
	}
	schemaPath, rootPath := os.Args[1], os.Args[2]
	rootType = os.Args[3]
	targetFile := os.Args[4]

	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read schema:", err)
		os.Exit(1)
	}
	var root Schema
	if err := json.Unmarshal(raw, &root); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshal:", err)
		os.Exit(1)
	}

	node := resolve(&root, rootPath)
	if node == nil {
		fmt.Fprintln(os.Stderr, "path not found:", rootPath)
		os.Exit(1)
	}
	if node.Type != "object" || len(node.Properties) == 0 {
		fmt.Fprintf(os.Stderr, "%s is not an object with properties\n", rootPath)
		os.Exit(1)
	}

	emitStruct(rootType, node, nil)

	if err := splice(targetFile, out.String()); err != nil {
		fmt.Fprintln(os.Stderr, "splice:", err)
		os.Exit(1)
	}
}

func splice(path, generated string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(content)
	bi := strings.Index(s, beginMarker)
	if bi < 0 {
		return fmt.Errorf("marker %q not found in %s", beginMarker, path)
	}
	head := s[:bi+len(beginMarker)]
	return os.WriteFile(path, []byte(head+"\n\n"+generated), 0o644)
}

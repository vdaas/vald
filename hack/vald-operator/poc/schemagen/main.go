// Command schemagen is a PoC: it generates Go structs from vald's
// charts/vald/values.schema.json (JSON Schema draft-07).
//
// This demonstrates the "JSON Schema -> Go types" step that vald's existing
// schema pipeline (values.yaml -> values.schema.json -> CRD) is missing.
// It is intentionally small and stdlib-only.
//
// Usage:
//
//	go run ./poc/schemagen <values.schema.json> <root.path> [root.path...]
//
// Example:
//
//	go run ./poc/schemagen /path/to/charts/vald/values.schema.json agent.ngt gateway.lb.gateway_config
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Schema is the minimal subset of JSON Schema draft-07 that vald's
// values.schema.json actually uses (no $ref/$defs/anyOf/allOf).
type Schema struct {
	Type       string             `json:"type"`
	Properties map[string]*Schema `json:"properties"`
	Items      *Schema            `json:"items"`
	Enum       []any              `json:"enum"`
}

// initialisms get upper-cased when building Go identifiers.
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

// typeName builds a struct name from a path like ["agent","ngt"] -> "AgentNGT".
func typeName(path []string) string {
	var b strings.Builder
	for _, p := range path {
		for _, part := range strings.Split(p, "_") {
			b.WriteString(goName(part))
		}
	}
	return b.String()
}

var out strings.Builder
var emitted = map[string]bool{}

// goType returns the Go type expression for a schema node and, for nested
// objects, queues a named struct to be emitted.
func goType(s *Schema, path []string) string {
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
			return "[]interface{}"
		}
		return "[]" + goType(s.Items, append(path, "Item"))
	case "object":
		if len(s.Properties) == 0 {
			// free-form object (e.g. annotations) -> map
			return "map[string]interface{}"
		}
		name := typeName(path)
		emitStruct(name, s, path)
		return "*" + name
	default:
		return "interface{}"
	}
}

func emitStruct(name string, s *Schema, path []string) {
	if emitted[name] {
		return
	}
	emitted[name] = true

	keys := make([]string, 0, len(s.Properties))
	for k := range s.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Emit dependencies first-ish by rendering into a temp, then prepend.
	body := &strings.Builder{}
	fmt.Fprintf(body, "type %s struct {\n", name)
	for _, k := range keys {
		field := s.Properties[k]
		ft := goType(field, append(path, k))
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

func resolve(root *Schema, dotted string) (*Schema, []string) {
	node := root
	path := strings.Split(dotted, ".")
	for _, p := range path {
		if node.Properties == nil {
			return nil, nil
		}
		next, ok := node.Properties[p]
		if !ok {
			return nil, nil
		}
		node = next
	}
	return node, path
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: schemagen <values.schema.json> <root.path> [root.path...]")
		os.Exit(2)
	}
	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "read:", err)
		os.Exit(1)
	}
	var root Schema
	if err := json.Unmarshal(raw, &root); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshal:", err)
		os.Exit(1)
	}

	out.WriteString("// Code generated from charts/vald/values.schema.json. DO NOT EDIT (PoC).\n\n")
	out.WriteString("package valdrelease\n\n")

	for _, dotted := range os.Args[2:] {
		node, path := resolve(&root, dotted)
		if node == nil {
			fmt.Fprintln(os.Stderr, "path not found:", dotted)
			os.Exit(1)
		}
		if node.Type != "object" || len(node.Properties) == 0 {
			fmt.Fprintf(os.Stderr, "%s is not an object with properties\n", dotted)
			os.Exit(1)
		}
		emitStruct(typeName(path), node, path)
	}
	fmt.Print(out.String())
}

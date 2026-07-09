//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// Command gotype generates Go types from a Helm values JSON Schema
// (e.g. charts/vald/values.schema.json).
//
// The schema is a bare JSON Schema (draft-07), so it is wrapped into a minimal
// OpenAPI 3 document and handed to oapi-codegen, which produces the Go types.
// k8s-native subtrees annotated with `x-go-type` / `x-go-type-import` in the
// schema are mapped to their real Go types (e.g. corev1.Affinity) instead of
// anonymous structs.
//
// Usage:
//
//	gotype -schema charts/vald/values.schema.json -out hack/valdvalues/values.gen.go -package valdvalues
//
// The output destination (-out), package name (-package), top-level type name
// (-name), and input schema (-schema) are all configurable via flags.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// defaultOutFile is appended when -out points to a directory (no file name).
const defaultOutFile = "values.gen.go"

func main() {
	// No defaults here on purpose: callers (the Makefile) must pass these
	// explicitly. -out accepts a file path or a directory (values.gen.go is
	// appended for a directory).
	schemaPath := flag.String("schema", "", "path to the input JSON Schema (required)")
	out := flag.String("out", "", "output Go file path or directory (required)")
	pkg := flag.String("package", "", "generated package name (required)")
	name := flag.String("name", "Values", "top-level type name")
	flag.Parse()

	var missing []string
	if *schemaPath == "" {
		missing = append(missing, "-schema")
	}
	if *out == "" {
		missing = append(missing, "-out")
	}
	if *pkg == "" {
		missing = append(missing, "-package")
	}
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "gotype: missing required flag(s): %s\n\n", strings.Join(missing, ", "))
		flag.Usage()
		os.Exit(2)
	}

	if err := run(*schemaPath, *out, *pkg, *name); err != nil {
		fmt.Fprintln(os.Stderr, "gotype:", err)
		os.Exit(1)
	}
}

func run(schemaPath, out, pkg, name string) error {
	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}
	var schema map[string]any
	if err := json.Unmarshal(raw, &schema); err != nil {
		return fmt.Errorf("parse schema: %w", err)
	}
	// OpenAPI schema objects don't carry the JSON Schema `$schema` keyword.
	delete(schema, "$schema")

	// Move JSON Schema $defs into OpenAPI components.schemas so each named
	// definition becomes a named Go type. Refs are rewritten below.
	schemas := map[string]any{}
	if defs, ok := schema["$defs"].(map[string]any); ok {
		for k, v := range defs {
			schemas[k] = v
		}
		delete(schema, "$defs")
	}
	schemas[name] = schema

	// Wrap the JSON Schema into a minimal OpenAPI 3 document. oapi-codegen only
	// emits component schemas reachable from an operation, so a dummy path that
	// references the schema is added to prevent it from being pruned.
	doc := map[string]any{
		"openapi": "3.0.3",
		"info":    map[string]any{"title": "vald values", "version": "0.0.0"},
		"paths": map[string]any{
			"/": map[string]any{"get": map[string]any{
				"operationId": "root",
				"responses": map[string]any{"200": map[string]any{
					"description": "ok",
					"content": map[string]any{"application/json": map[string]any{
						"schema": map[string]any{"$ref": "#/components/schemas/" + name},
					}},
				}},
			}},
		},
		"components": map[string]any{"schemas": schemas},
	}
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal openapi doc: %w", err)
	}
	// JSON Schema uses #/$defs/... references; OpenAPI uses #/components/schemas/...
	docBytes = bytes.ReplaceAll(docBytes, []byte("#/$defs/"), []byte("#/components/schemas/"))

	tmp, err := os.CreateTemp("", "vald-openapi-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(docBytes); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	bin, err := exec.LookPath("oapi-codegen")
	if err != nil {
		return fmt.Errorf("oapi-codegen not found in PATH; install it via `make tools/install` (see hack/go.tools): %w", err)
	}

	// If out has no Go file name (a directory, trailing separator, or a
	// non-".go" path), append the default file name.
	if strings.HasSuffix(out, string(os.PathSeparator)) || filepath.Ext(out) != ".go" {
		out = filepath.Join(out, defaultOutFile)
	}

	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	cmd := exec.Command(bin, "-generate", "types", "-package", pkg, "-o", out, tmp.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("oapi-codegen: %w\n%s", err, stderr.String())
	}
	fmt.Fprintf(os.Stderr, "gotype: wrote %s (package %s)\n", out, pkg)
	return nil
}

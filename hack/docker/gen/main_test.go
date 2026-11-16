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
	"context"
	"path/filepath"
	"slices"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
	"golang.org/x/tools/go/packages"
	yaml "gopkg.in/yaml.v2"
)

// Test that running main generates workflow yaml containing
// permissions and secrets blocks as intended.
func TestWorkflowPermissionsAndSecretsGenerated(t *testing.T) {
	td := t.TempDir()

	oldArgs := os.Args
	os.Args = []string{"gen", td}
	defer func() { os.Args = oldArgs }()

	main()

	wfPath := filepath.Join(td, ".github", "workflows", "dockers-ci-container-image.yaml")
	b, err := os.ReadFile(wfPath)
	if err != nil {
		t.Fatalf("failed to read generated workflow %s: %v", wfPath, err)
	}

	var wf Workflow
	if err := yaml.Unmarshal(b, &wf); err != nil {
		t.Fatalf("failed to unmarshal workflow yaml: %v", err)
	}

	// Validate secrets
	wantSecrets := map[string]string{
		"PACKAGE_USER":   "${{ secrets.PACKAGE_USER }}",
		"PACKAGE_TOKEN":  "${{ secrets.PACKAGE_TOKEN }}",
		"DOCKERHUB_USER": "${{ secrets.DOCKERHUB_USER }}",
		"DOCKERHUB_PASS": "${{ secrets.DOCKERHUB_PASS }}",
	}
	for k, v := range wantSecrets {
		if gv, ok := wf.Jobs.Build.Secrets[k]; !ok || gv != v {
			t.Errorf("secrets[%s] = %q, want %q (present=%v)", k, gv, v, ok)
		}
	}
}

func TestAppendM(t *testing.T) {
	m1 := map[string]string{
		"FOO": "bar",
	}
	m2 := map[string]string{
		"FOO": "baz",
	}
	got := appendM(m1, m2)
	if got["FOO"] != "bar:baz" {
		t.Fatalf("appendM merged FOO = %q, want %q", got["FOO"], "bar:baz")
	}
}

func TestExtractVariables(t *testing.T) {
	in := "PATH=${HOME}:${GOPATH}/bin:$SHELL"
	got := extractVariables(in)
	want := []string{"HOME", "GOPATH", "SHELL"}
	if len(got) != len(want) {
		t.Fatalf("extractVariables len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("extractVariables[%d] = %q, want %q (all=%v)", i, got[i], want[i], got)
		}
	}
}

func TestTopologicalSortOrder(t *testing.T) {
	env := map[string]string{
		"A": "x",
		"B": "${A}",
		"C": "${B}",
	}
	got := topologicalSort(env)
	// Expect entries in dependency order: A before B, B before C
	gi := func(key string) int {
		prefix := key + "="
		for i, kv := range got {
			if len(kv) >= len(prefix) && kv[:len(prefix)] == prefix {
				return i
			}
		}
		return -1
	}
	ia, ib, ic := gi("A"), gi("B"), gi("C")
	if ia == -1 || ib == -1 || ic == -1 {
		t.Fatalf("missing keys in topologicalSort result: %v", got)
	}
	if ia >= ib || ib >= ic {
		t.Fatalf("unexpected order: %v (want A < B < C)", got)
	}
}

func TestWorkflowTrivyIgnorePaths(t *testing.T) {
	tests := []struct {
		name         string
		target       string
		ignoreExists bool
	}{
		{name: devContainer, target: devContainer, ignoreExists: true},
		{name: buildbase, target: buildbase, ignoreExists: true},
		{name: "scanner", target: buildkitSyftScanner, ignoreExists: true},
		{name: "missing ignore", target: devContainer, ignoreExists: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			relative := ".trivyignore.d/" + tt.target
			if tt.ignoreExists {
				if err := os.MkdirAll(filepath.Join(root, ".trivyignore.d"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, relative), []byte("CVE-2026-84445 exp:2026-12-06\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if err := generateWorkflow(context.Background(), root, "vald-"+tt.target, "vald team", 2026, Data{ContainerType: Other}); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(root, ".github/workflows/dockers-"+tt.target+"-image.yaml"))
			if err != nil {
				t.Fatal(err)
			}
			var workflow Workflow
			if err := yaml.Unmarshal(data, &workflow); err != nil {
				t.Fatal(err)
			}
			for name, paths := range map[string][]string{
				"pull_request":        workflow.On.PullRequest.Paths,
				"pull_request_target": workflow.On.PullRequestTarget.Paths,
			} {
				found := false
				for _, path := range paths {
					if path == relative {
						found = true
					}
				}
				if found != tt.ignoreExists {
					t.Errorf("%s paths contain %q = %t, want %t: %v", name, relative, found, tt.ignoreExists, paths)
				}
			}
		})
	}
}

// A dependency tree that cannot be resolved must fail the generation instead of
// writing a workflow whose pull_request paths only cover cmd/<pkg>.
// packages.Load does not return an error for a .go file pattern when `go list` fails,
// so this relies on the empty-result / Package.Errors checks.
func TestWorkflowFailsOnDependencyTreeError(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	cmdDir := filepath.Join(root, "cmd", "brokendep")
	if err := os.MkdirAll(cmdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cmdDir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	// A malformed go.mod makes `go list` fail as a whole.
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("this is not a go.mod file\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := generateWorkflow(context.Background(), root, "vald-brokendep", "vald team", 2026, Data{
		ContainerType: Go,
		PackageDir:    "brokendep",
	})
	if err == nil {
		t.Fatal("generateWorkflow() error = nil, want dependency tree error")
	}
	wf := filepath.Join(root, ".github", "workflows", "dockers-brokendep-image.yaml")
	if _, serr := os.Stat(wf); !errors.Is(serr, os.ErrNotExist) {
		t.Errorf("workflow %s must not be written on error, stat err = %v", wf, serr)
	}
}

// entryPkgPath and entryPkgName are what packages.Load reports for the main
// package of a .go file pattern.
const (
	entryPkgPath = "command-line-arguments"
	entryPkgName = "main"
)

func TestCheckLoadedPackages(t *testing.T) {
	t.Parallel()
	const entry = "/root/cmd/x/main.go"
	listErr := packages.Error{Msg: "cannot find package"}
	repoDep := &packages.Package{PkgPath: "github.com/" + repository + "/internal/x"}
	brokenRepoDep := &packages.Package{PkgPath: "github.com/" + repository + "/internal/y", Errors: []packages.Error{listErr}}
	brokenExternal := &packages.Package{PkgPath: "example.com/ext", Errors: []packages.Error{listErr}}
	tests := []struct {
		name    string
		pkgs    []*packages.Package
		wantErr bool
	}{
		{name: "no package is an error", pkgs: nil, wantErr: true},
		{
			name: "entry package and its repository dependencies without errors",
			pkgs: []*packages.Package{{
				PkgPath: entryPkgPath,
				Name:    entryPkgName,
				Imports: map[string]*packages.Package{repoDep.PkgPath: repoDep},
			}},
		},
		{
			name:    "entry package with errors",
			pkgs:    []*packages.Package{{PkgPath: entryPkgPath, Errors: []packages.Error{listErr}}},
			wantErr: true,
		},
		{
			name: "repository dependency with errors",
			pkgs: []*packages.Package{{
				PkgPath: entryPkgPath,
				Name:    entryPkgName,
				Imports: map[string]*packages.Package{brokenRepoDep.PkgPath: brokenRepoDep},
			}},
			wantErr: true,
		},
		{
			name: "external dependency errors are ignored",
			pkgs: []*packages.Package{{
				PkgPath: entryPkgPath,
				Name:    entryPkgName,
				Imports: map[string]*packages.Package{brokenExternal.PkgPath: brokenExternal},
			}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := checkLoadedPackages(entry, tc.pkgs); (err != nil) != tc.wantErr {
				t.Errorf("checkLoadedPackages() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestPackageNodeToSliceVisitsSharedNodesOnce(t *testing.T) {
	t.Parallel()
	// Diamond: root -> {a, b}, a -> c, b -> c. c must be listed once.
	c := &PackageNode{Name: "c"}
	a := &PackageNode{Name: "a", Imports: []*PackageNode{c}}
	b := &PackageNode{Name: "b", Imports: []*PackageNode{c}}
	root := PackageNode{Name: entryPkgPath, Imports: []*PackageNode{a, b}}

	got := root.ToSlice()
	slices.Sort(got)
	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Fatalf("ToSlice() = %v, want %v", got, want)
	}
}

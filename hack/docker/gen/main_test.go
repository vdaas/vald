package main

import (
    "os"
    "path/filepath"
    "testing"

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

    // Validate permissions
    if got := wf.Jobs.Build.Permissions["contents"]; got != "read" {
        t.Errorf("permissions.contents = %q, want %q", got, "read")
    }
    if got := wf.Jobs.Build.Permissions["packages"]; got != "write" {
        t.Errorf("permissions.packages = %q, want %q", got, "write")
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
    if !(ia != -1 && ib != -1 && ic != -1) {
        t.Fatalf("missing keys in topologicalSort result: %v", got)
    }
    if !(ia < ib && ib < ic) {
        t.Fatalf("unexpected order: %v (want A < B < C)", got)
    }
}

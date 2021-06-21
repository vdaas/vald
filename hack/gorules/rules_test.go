package gorules

import (
	"os"
	"testing"

	"github.com/quasilyte/go-ruleguard/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRules(t *testing.T) {
	if err := analyzer.Analyzer.Flags.Set("rules", "rules.go"); err != nil {
		t.Fatalf("set rules flag: %v", err)
	}
	analysistest.Run(
		t,
		os.Getenv("GOPATH"),
		analyzer.Analyzer,
		"github.com/vdaas/vald/hack/gorules/testdata",
	)
}

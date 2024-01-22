package kubectl

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"testing"
)

func RolloutResource(t *testing.T, ctx context.Context, resource string) error {
	t.Helper()
	cmd := exec.CommandContext(ctx, "sh", "-c",
		fmt.Sprintf("kubectl rollout restart %s && kubectl rollout status %s", resource, resource),
	)
	return runCmd(t, ctx, cmd)
}

func WaitResources(t *testing.T, ctx context.Context, resource, labelSelector, condition, timeout string) error {
	t.Helper()
	cmd := exec.CommandContext(ctx, "sh", "-c",
		fmt.Sprintf("kubectl wait --for=condition=%s %s -l %s --timeout %s", condition, resource, labelSelector, timeout),
	)
	return runCmd(t, ctx, cmd)
}

func runCmd(t *testing.T, ctx context.Context, cmd *exec.Cmd) error {
	t.Helper()

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(exitErr.Stderr))
		} else {
			return fmt.Errorf("unexpected error: %w", err)
		}
	}
	t.Log(string(out))
	return nil
}

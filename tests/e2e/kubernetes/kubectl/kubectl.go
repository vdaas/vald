//go:build e2e

//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

package kubectl

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

// RolloutResource rollouts and wait for the resource to be ready.
func RolloutResource(ctx context.Context, t *testing.T, resource string) error {
	t.Helper()

	cmd := exec.CommandContext(ctx, "kubectl", "rollout", "restart", resource)
	if err := runCmd(t, cmd); err != nil {
		return err
	}

	cmd = exec.CommandContext(ctx, "kubectl", "rollout", "status", resource)
	return runCmd(t, cmd)
}

// WaitResources waits for multiple resources to be ready.
func WaitResources(ctx context.Context, t *testing.T, resource, labelSelector, condition, timeout string) error {
	t.Helper()

	cmd := exec.CommandContext(ctx, "kubectl", "wait", "--for=condition="+condition, "-l", labelSelector, "--timeout", timeout, resource)
	return runCmd(t, cmd)
}

func runCmd(t *testing.T, cmd *exec.Cmd) error {
	t.Helper()
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return errors.New(string(exitErr.Stderr))
		} else {
			return fmt.Errorf("unexpected error: %w", err)
		}
	}
	t.Log(string(out))
	return nil
}

//go:build e2e

package operation

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
)

func (j *cronJobExecute) CreateAndWait(t *testing.T, ctx context.Context, jobName string) error {
	if err := createJob(t, jobName, j.cronJob); err != nil {
		return err
	}

	defer func() {
		err := deleteJob(t, jobName)
		if err != nil {
			t.Errorf("failed to delete job: %s", err)
		}
	}()

	return waitJob(t, ctx, jobName)
}

func createJob(t *testing.T, jobName, cronJobName string) error {
	t.Helper()
	t.Logf("creating job: %s from CronJob %s", jobName, cronJobName)
	createCmd := fmt.Sprintf("kubectl create job %s --from=cronjob/%s", jobName, cronJobName)
	cmd := exec.Command("sh", "-c", createCmd)
	return execCmd(t, cmd)
}

func deleteJob(t *testing.T, jobName string) error {
	t.Helper()
	t.Log("deleting correction job")
	deleteKubeCmd := fmt.Sprintf("kubectl delete job %s", jobName)
	cmd := exec.Command("sh", "-c", deleteKubeCmd)
	return execCmd(t, cmd)
}

func waitJob(t *testing.T, ctx context.Context, jobName string) error {
	t.Helper()
	t.Log("waiting for the correction job to complete or fail")
	waitCompleteCmd := fmt.Sprintf("kubectl wait --timeout=-1s job/%s --for=condition=complete", jobName)
	waitFailedCmd := fmt.Sprintf("kubectl wait --timeout=-1s job/%s --for=condition=failed", jobName)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	complete := make(chan struct{})
	failed := make(chan struct{})
	ech := make(chan error)
	go func() {
		cmd := exec.CommandContext(ctx, "sh", "-c", waitCompleteCmd)
		err := execCmd(t, cmd)
		if err != nil {
			ech <- err
			return
		}

		complete <- struct{}{}
	}()

	go func() {
		cmd := exec.CommandContext(ctx, "sh", "-c", waitFailedCmd)
		err := execCmd(t, cmd)
		if err != nil {
			ech <- err
			return
		}

		t.Logf("%s failed. dumping status", jobName)
		dumpStatusCmd := fmt.Sprintf("kubectl get job %s -o yaml", jobName)
		cmd = exec.Command("sh", "-c", dumpStatusCmd)
		err = execCmd(t, cmd)
		if err != nil {
			t.Log("failed to dump status")
		}
		failed <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-complete:
		return nil
	case <-failed:
		return fmt.Errorf("correction job failed")
	case err := <-ech:
		return err
	}
}

func execCmd(t *testing.T, cmd *exec.Cmd) error {
	t.Helper()
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%s, %s, %w", string(out), string(exitErr.Stderr), err)
		} else {
			return fmt.Errorf("unexpected error on creating job: %w", err)
		}
	}
	t.Log(string(out))
	return nil
}

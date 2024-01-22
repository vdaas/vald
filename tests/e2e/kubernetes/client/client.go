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

// package client provides kubernetes client
package client

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/e2e/kubernetes/portforward"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client interface {
	Portforward(
		namespace, podName string,
		localPort, podPort int,
	) *portforward.Portforward
	GetPod(
		ctx context.Context,
		namespace,
		name string,
	) (*corev1.Pod, error)
	GetPods(
		ctx context.Context,
		namespace string,
		labelSelector string,
	) ([]corev1.Pod, error)
	DeletePod(
		ctx context.Context,
		namespace, name string,
	) error
	WaitForPodReady(
		ctx context.Context,
		namespace, name string,
		timeout time.Duration,
	) (ok bool, err error)
	WaitForPodsReady(
		ctx context.Context,
		namespace, labelSelector, timeout string,
	) error
	ListCronJob(
		ctx context.Context,
		namespace, labelSelector string,
	) ([]v1.CronJob, error)
	CreateJob(
		ctx context.Context,
		namespace string,
		job *v1.Job,
	) error
	CreateJobFromCronJob(
		ctx context.Context,
		name, namespace string,
		cronJob *v1.CronJob,
	) error
	RolloutResource(
		ctx context.Context,
		resource string,
	) error
	WaitResources(ctx context.Context, resource, labelSelector, condition, timeout string) error
}

type client struct {
	rest      *rest.Config
	clientset *kubernetes.Clientset
}

func New(kubeConfig string) (Client, error) {
	if kubeConfig == "" {
		kubeConfig = os.Getenv("KUBECONFIG")
		if kubeConfig == "" {
			if home := os.Getenv("HOME"); home != "" {
				kubeConfig = file.Join(home, ".kube", "config")
			}
		}
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &client{
		rest:      cfg,
		clientset: clientset,
	}, nil
}

func (cli *client) Portforward(
	namespace, podName string,
	localPort, podPort int,
) *portforward.Portforward {
	return portforward.New(cli.rest, namespace, podName, localPort, podPort)
}

func (cli *client) GetPod(
	ctx context.Context,
	namespace,
	name string,
) (*corev1.Pod, error) {
	pod, err := cli.clientset.CoreV1().Pods(
		namespace,
	).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (cli *client) GetPods(
	ctx context.Context,
	namespace string,
	labelSelector string,
) ([]corev1.Pod, error) {
	pods, err := cli.clientset.CoreV1().Pods(
		namespace,
	).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func (cli *client) DeletePod(
	ctx context.Context,
	namespace, name string,
) error {
	cli.clientset.CoreV1().Pods(
		namespace,
	).Delete(ctx, name, metav1.DeleteOptions{})

	return nil
}

func (cli *client) WaitForPodReady(
	ctx context.Context,
	namespace, name string,
	timeout time.Duration,
) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		pod, err := cli.GetPod(ctx, namespace, name)
		if err != nil {
			return false, err
		}

		for _, condition := range pod.Status.Conditions {
			if strings.EqualFold(
				strings.ToLower(string(condition.Type)), "ready") &&
				strings.EqualFold(
					strings.ToLower(string(condition.Status)), "true") {
				return true, nil
			}
		}

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-tick.C:
		}
	}
}

func (cli *client) RolloutResource(ctx context.Context, resource string) error {
	cmd := exec.CommandContext(ctx, "sh", "-c",
		fmt.Sprintf("kubectl rollout restart %s && kubectl rollout status %s", resource, resource),
	)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(exitErr.Stderr))
		} else {
			return fmt.Errorf("unexpected error: %w", err)
		}
	}
	fmt.Println(string(out))
	return nil
}

func (cli *client) WaitResources(ctx context.Context, resource, labelSelector, condition, timeout string) error {
	cmd := exec.CommandContext(ctx, "sh", "-c",
		fmt.Sprintf("kubectl wait --for=condition=%s %s -l %s --timeout %s", condition, resource, labelSelector, timeout),
	)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(exitErr.Stderr))
		} else {
			return fmt.Errorf("unexpected error: %w", err)
		}
	}
	fmt.Println(string(out))
	return nil
}

func (cli *client) WaitForPodsReady(ctx context.Context, namespace, labelSelector, timeout string) error {
	// use kubectl wait because it's complicated to implement this with client-go
	cmd := exec.CommandContext(ctx, "sh", "-c",
		fmt.Sprintf("kubectl wait --timeout=%s --for=condition=Ready pod -n %s -l %s", timeout, namespace, labelSelector),
	)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return errors.New(string(exitErr.Stderr))
		} else {
			return fmt.Errorf("unexpected error: %w", err)
		}
	}
	fmt.Println(string(out))
	return nil
}

func (cli *client) ListCronJob(ctx context.Context, namespace, labelSelector string) ([]v1.CronJob, error) {
	cronJobs, err := cli.clientset.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}

	return cronJobs.Items, nil
}

func (cli *client) CreateJob(ctx context.Context, namespace string, job *v1.Job) error {
	_, err := cli.clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
	return err
}

func (cli *client) CreateJobFromCronJob(ctx context.Context, name, namespace string, cronJob *v1.CronJob) error {
	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: cronJob.Spec.JobTemplate.Spec,
	}

	_, err := cli.clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
	return err
}

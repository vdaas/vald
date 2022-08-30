//go:build e2e

//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"os"
	"time"

	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/e2e/kubernetes/portforward"
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

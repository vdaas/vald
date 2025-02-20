//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// package kubernetes provides kubernetes client functions
package kubernetes

import (
	"context"
	"os"
	"time"

	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/e2e/kubernetes/portforward"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Client interface {
	Portforward(
		namespace, podName string,
		localPort, podPort int,
	) *portforward.Portforward
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
		kubeConfig = os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
		if kubeConfig == "" {
			dotKubePath := file.Join(os.Getenv("HOME"), ".kube", "config")
			if file.Exists(dotKubePath) {
				kubeConfig = dotKubePath
			}
			if kubeConfig == "" {
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
	namespace, podName string, localPort, podPort int,
) *portforward.Portforward {
	return portforward.New(cli.rest, namespace, podName, localPort, podPort)
}

func (cli *client) WaitForPodReady(
	ctx context.Context, namespace, name string, timeout time.Duration,
) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	pod := Pod{
		Name:      name,
		Namespace: namespace,
	}

	for {
		pod, err := pod.Get(ctx, cli.clientset)
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

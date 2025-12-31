//go:build e2e

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

package kubernetes

import (
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/os"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Client interface {
	GetClientSet() kubernetes.Interface
	GetRESRConfig() *rest.Config
}

type client struct {
	rest      *rest.Config
	clientset *kubernetes.Clientset
	manager   manager.Manager
	client    kclient.WithWatch
}

func NewClient(kubeConfig, currentContext string) (c Client, err error) {
	if kubeConfig == "" {
		kubeConfig = os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
		if kubeConfig == "" {
			if file.Exists(clientcmd.RecommendedHomeFile) {
				kubeConfig = clientcmd.RecommendedHomeFile
			}
			if kubeConfig == "" {
				c, err = inClusterConfigClient()
				if err != nil {
					return nil, err
				}
				return c, nil
			}
		}
	}

	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
		&clientcmd.ConfigOverrides{
			ClusterInfo:    clientcmdapi.Cluster{},
			CurrentContext: currentContext,
		}).ClientConfig()
	if err != nil {
		log.Debugf("failed to build config from kubeConfig path %s,\terror: %v", kubeConfig, err)
		var ierr error
		c, ierr = inClusterConfigClient()
		if ierr != nil {
			return nil, errors.Join(err, ierr)
		}
		return c, nil
	}

	c, err = newClient(cfg)
	if err != nil {
		log.Debugf("failed to build config from kubeConfig path %s,\terror: %v", kubeConfig, err)
		var ierr error
		c, ierr = inClusterConfigClient()
		if ierr != nil {
			return nil, errors.Join(err, ierr)
		}
	}
	return c, nil
}

func newClient(cfg *rest.Config) (Client, error) {
	if cfg.QPS == 0.0 {
		cfg.QPS = 20.0
	}
	if cfg.Burst == 0 {
		cfg.Burst = 30
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

func inClusterConfigClient() (Client, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return newClient(cfg)
}

func (c *client) GetClientSet() kubernetes.Interface {
	return c.clientset
}

func (c *client) GetRESRConfig() *rest.Config {
	return c.rest
}

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

package client

import (
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/os"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// ClientSet is the raw client-go accessor used by the resource package
// factories and other components that need direct clientset / REST config
// access instead of the controller-runtime CRUD wrapper (Client).
type ClientSet interface {
	GetClientSet() kubernetes.Interface
	GetRESTConfig() *rest.Config
}

type clientSet struct {
	rest      *rest.Config
	clientset *kubernetes.Clientset
}

// NewClientSet builds a ClientSet from the given kubeconfig path and context.
// It falls back to the KUBECONFIG environment variable, the recommended home
// kubeconfig, and finally the in-cluster configuration.
func NewClientSet(kubeConfig, currentContext string) (ClientSet, error) {
	if kubeConfig == "" {
		kubeConfig = os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
		if kubeConfig == "" {
			if file.Exists(clientcmd.RecommendedHomeFile) {
				kubeConfig = clientcmd.RecommendedHomeFile
			}
			if kubeConfig == "" {
				return fallbackToInCluster(nil)
			}
		}
	}

	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
		&clientcmd.ConfigOverrides{
			ClusterInfo:    clientcmdapi.Cluster{},
			CurrentContext: currentContext,
		},
	).ClientConfig()
	if err != nil {
		log.Debugf("failed to build config from kubeConfig path %s,\terror: %v", kubeConfig, err)
		return fallbackToInCluster(err)
	}

	c, err := clientSetFromConfig(cfg)
	if err != nil {
		log.Debugf("failed to build config from kubeConfig path %s,\terror: %v", kubeConfig, err)
		return fallbackToInCluster(err)
	}
	return c, nil
}

// fallbackToInCluster builds a ClientSet from the in-cluster configuration.
// origErr, if non-nil, is the error from the kubeConfig-based attempt that
// preceded this fallback; it is joined with the in-cluster error so both
// failure causes surface if the fallback also fails.
func fallbackToInCluster(origErr error) (ClientSet, error) {
	c, err := inClusterConfigClientSet()
	if err != nil {
		if origErr != nil {
			return nil, errors.Join(origErr, err)
		}
		return nil, err
	}
	return c, nil
}

func clientSetFromConfig(cfg *rest.Config) (ClientSet, error) {
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
	return &clientSet{
		rest:      cfg,
		clientset: clientset,
	}, nil
}

func inClusterConfigClientSet() (ClientSet, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return clientSetFromConfig(cfg)
}

func (c *clientSet) GetClientSet() kubernetes.Interface {
	return c.clientset
}

func (c *clientSet) GetRESTConfig() *rest.Config {
	return c.rest
}

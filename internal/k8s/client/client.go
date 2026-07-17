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
	"context"
	"maps"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/os"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	cli "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// Client is the single Kubernetes client abstraction, replacing the former
// three-way split between the internal/k8s.Client alias, StandaloneClient
// and ClientSet: full CRUD via the embedded k8s.Client (controller-runtime's
// Client, what mgr.GetClient() returns), Apply for server-side apply, and
// direct client-go access (GetClientSet/GetRESTConfig) for functionality
// controller-runtime's Client cannot express (SPDY port-forward,
// ServiceAccount token requests, typed clientset calls).
type Client interface {
	k8s.Client

	// Apply applies the given apply configuration to the Kubernetes cluster
	// using server-side apply. This controller-runtime version's Client
	// already declares Apply as part of its Writer, so restating it here is
	// purely documentation and does not require a separate implementation.
	Apply(ctx context.Context, obj runtime.ApplyConfiguration, opts ...cli.ApplyOption) error

	// GetClientSet returns the raw client-go clientset, for APIs
	// controller-runtime's Client cannot express (typed clientset calls,
	// SPDY port-forward, ServiceAccount token requests).
	GetClientSet() kubernetes.Interface

	// GetRESTConfig returns the *rest.Config the client and clientset were
	// built from.
	GetRESTConfig() *rest.Config
}

// unifiedClient implements Client. k8s.Client is embedded so the full
// controller-runtime CRUD surface (including Apply, already declared by its
// Writer) is promoted automatically; restConfig and clientset cover what
// controller-runtime's client.Client cannot express. scheme, kubeConfigPath
// and kubeContext only matter while New is assembling a standalone client
// from Options and are inert afterwards.
type unifiedClient struct {
	k8s.Client
	restConfig     *rest.Config
	clientset      kubernetes.Interface
	scheme         *runtime.Scheme
	kubeConfigPath string
	kubeContext    string
}

func (c *unifiedClient) GetClientSet() kubernetes.Interface {
	return c.clientset
}

func (c *unifiedClient) GetRESTConfig() *rest.Config {
	return c.restConfig
}

// NewFromManager builds a Client from a controller-runtime manager: the
// caller already has mgr.GetClient() (a working, watch-capable client scoped
// to the manager's cache) and mgr.GetConfig() (the REST config used to build
// a matching client-go Clientset).
func NewFromManager(mgr k8s.Manager) (Client, error) {
	cfg := mgr.GetConfig()
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clientset from manager config")
	}
	return &unifiedClient{
		Client:     mgr.GetClient(),
		restConfig: cfg,
		clientset:  cs,
	}, nil
}

// New builds a standalone Client for code that runs outside a manager (e.g.
// one-shot operators, jobs, config loaders). It replaces the former
// client.New() (StandaloneClient) + client.NewClientSet() (ClientSet) pair:
// scheme construction and cli.NewWithWatch are carried over unchanged from
// the old New, and unless WithRESTConfig already supplied a *rest.Config, it
// is resolved through the same fallback chain the old NewClientSet used (an
// explicit kubeconfig path, then the KUBECONFIG environment variable, then
// the recommended home kubeconfig, then finally the in-cluster
// configuration). The resolved config then builds both the watch-capable
// controller-runtime client and the client-go Clientset.
func New(opts ...Option) (Client, error) {
	c := &unifiedClient{scheme: runtime.NewScheme()}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Add the core schemes.
	if err := clientgoscheme.AddToScheme(c.scheme); err != nil {
		return nil, err
	}
	if err := snapshotv1.AddToScheme(c.scheme); err != nil {
		return nil, err
	}

	cfg := c.restConfig
	if cfg == nil {
		rcfg, err := resolveRESTConfig(c.kubeConfigPath, c.kubeContext)
		if err != nil {
			return nil, err
		}
		cfg = rcfg
	}

	wc, err := cli.NewWithWatch(cfg, cli.Options{
		Scheme: c.scheme,
	})
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clientset")
	}

	c.Client = wc
	c.restConfig = cfg
	c.clientset = cs

	return c, nil
}

// resolveRESTConfig resolves a *rest.Config using the same chain the former
// NewClientSet used: an explicit kubeConfig path if given, else the
// KUBECONFIG environment variable, else the recommended home kubeconfig
// file, else the in-cluster configuration.
func resolveRESTConfig(kubeConfig, currentContext string) (*rest.Config, error) {
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

	applyDefaultRateLimits(cfg)
	return cfg, nil
}

// fallbackToInCluster builds a *rest.Config from the in-cluster
// configuration. origErr, if non-nil, is the error from the kubeConfig-based
// attempt that preceded this fallback; it is joined with the in-cluster
// error so both failure causes surface if the fallback also fails.
func fallbackToInCluster(origErr error) (*rest.Config, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		if origErr != nil {
			return nil, errors.Join(origErr, err)
		}
		return nil, err
	}
	applyDefaultRateLimits(cfg)
	return cfg, nil
}

// applyDefaultRateLimits fills in the QPS/Burst defaults the former
// clientSetFromConfig used whenever the resolved config left them unset.
func applyDefaultRateLimits(cfg *rest.Config) {
	if cfg.QPS == 0.0 {
		cfg.QPS = 20.0
	}
	if cfg.Burst == 0 {
		cfg.Burst = 30
	}
}

// NewLabelSelector creates a labels.Selector for Options like ListOptions.
func NewLabelSelector(key string, op selection.Operator, vals []string) (labels.Selector, error) {
	requirements, err := labels.NewRequirement(key, op, vals)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create requirement on creating label selector")
	}
	return labels.NewSelector().Add(*requirements), nil
}

// ObjectPredicates returns builder.Predicates that pass only events whose
// object is a PT accepted by filter. It builds on controller-runtime's
// predicate.NewPredicateFuncs, whose event mapping (Create/Delete/Generic use
// e.Object, Update uses e.ObjectNew) matches what per-Kind hand-written
// predicate.Funcs would do. The pointer constraint is spelled inline because
// reusing resource.Objectable would create an import cycle
// (internal/k8s/resource imports this package).
func ObjectPredicates[T any, PT interface {
	*T
	cli.Object
}](filter func(PT) bool) builder.Predicates {
	return builder.WithPredicates(predicate.NewPredicateFuncs(func(obj cli.Object) bool {
		o, ok := obj.(PT)
		return ok && filter(o)
	}))
}

// Patcher is an interface for patching resources with controller-runtime client.
type Patcher interface {
	// ApplyPodAnnotations applies the given annotations to the agent pod with server-side apply.
	ApplyPodAnnotations(ctx context.Context, name, namespace string, entries map[string]string) error
}

type patcher struct {
	client       Client
	fieldManager string
}

func NewPatcher(fieldManager string) (Patcher, error) {
	client, err := New()
	if err != nil {
		return nil, err
	}

	return &patcher{
		client:       client,
		fieldManager: fieldManager,
	}, nil
}

func (s *patcher) ApplyPodAnnotations(
	ctx context.Context, name, namespace string, entries map[string]string,
) error {
	var podList corev1.PodList
	if err := s.client.List(ctx, &podList, &cli.ListOptions{
		Namespace:     namespace,
		FieldSelector: fields.OneTermEqualSelector("metadata.name", name),
	}); err != nil {
		return err
	}

	if len(podList.Items) == 0 {
		return errors.New("agent pod not found on exporting metrics")
	}

	//nolint:mnd
	if len(podList.Items) >= 2 {
		return errors.New("multiple agent pods found on exporting metrics. pods with same name exist in the same namespace?")
	}
	pod := podList.Items[0]

	curApplyConfig, err := applycorev1.ExtractPod(&pod, s.fieldManager)
	if err != nil {
		return err
	}

	// check if there is any diffs in the annotations
	annotations := pod.GetObjectMeta().GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	maps.Copy(annotations, entries)
	expectPod := applycorev1.Pod(name, namespace).
		WithAnnotations(annotations)

	if equality.Semantic.DeepEqual(expectPod, curApplyConfig) {
		// no change found in the pod annotations
		return nil
	}

	// now we found the diffs, apply the changes
	return s.client.Apply(ctx, expectPod,
		cli.FieldOwner(s.fieldManager),
		cli.ForceOwnership,
	)
}

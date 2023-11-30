//go:build e2e

//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// package portforward provides port-forward functionality for e2e tests
package portforward

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/vdaas/vald/internal/strings"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// portforwarder provides a port-forwarding functionality of kubectl
// reference: https://github.com/gianarb/kube-port-forward
type Portforward struct {
	namespace string
	podName   string

	localPort int
	podPort   int

	restConfig *rest.Config
	readyCh    chan struct{}
	stopCh     chan struct{}
}

func New(config *rest.Config, namespace, podName string, localPort, podPort int) *Portforward {
	return &Portforward{
		namespace:  namespace,
		podName:    podName,
		localPort:  localPort,
		podPort:    podPort,
		restConfig: config,
		readyCh:    make(chan struct{}),
		stopCh:     make(chan struct{}, 1),
	}
}

func (p *Portforward) Start() error {
	stream := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward",
		p.namespace, p.podName)

	hostIP := strings.TrimLeft(p.restConfig.Host, "https:/")

	transport, upgrader, err := spdy.RoundTripperFor(p.restConfig)
	if err != nil {
		return err
	}

	ech := make(chan error, 1)
	go func() {
		fw, err := portforward.New(
			spdy.NewDialer(
				upgrader,
				&http.Client{Transport: transport},
				http.MethodPost,
				&url.URL{Scheme: "https", Path: path, Host: hostIP},
			),
			[]string{fmt.Sprintf("%d:%d", p.localPort, p.podPort)},
			p.stopCh,
			p.readyCh,
			stream.Out,
			stream.ErrOut,
		)
		if err != nil {
			ech <- err
		}

		err = fw.ForwardPorts()
		if err != nil {
			ech <- err
		}
	}()

	select {
	case <-p.readyCh:
		return nil
	case err = <-ech:
		return err
	}
}

func (p *Portforward) Close() error {
	close(p.stopCh)

	return nil
}

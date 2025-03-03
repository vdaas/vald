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
	"fmt"
	"net/http"
	"os"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

func Portforward(
	ctx context.Context, c Client, namespace, podName string, localPort, targetPort uint16,
) (cancel context.CancelFunc, ech chan<- error, err error) {
	return PortforwardExtended(ctx, c, namespace, podName, []string{"localhost"}, map[uint16]uint16{localPort: targetPort}, http.DefaultClient)
}

func PortforwardExtended(
	ctx context.Context,
	c Client,
	namespace, podName string,
	addresses []string,
	ports map[uint16]uint16,
	hc *http.Client,
) (cancel context.CancelFunc, ech chan<- error, err error) {
	ctx, cancel = context.WithCancel(ctx)
	pod, err := Pod(c, namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return cancel, nil, err
	}
	status, info, err := CheckResourceState(pod)
	if err != nil {
		return cancel, nil, err
	}
	log.Debug(info)
	if status != StatusAvailable {
		return cancel, nil, errors.ErrPodIsNotRunning(namespace, podName)
	}

	transport, upgrader, err := spdy.RoundTripperFor(c.GetRESRConfig())
	if err != nil {
		return cancel, nil, err
	}
	hc.Transport = transport

	if addresses == nil {
		return cancel, nil, errors.ErrPortForwardAddressNotFound
	}

	portPairs := make([]string, 0, len(ports))
	for localPort, targetPort := range ports {
		portPairs = append(portPairs, fmt.Sprintf("%d:%d", localPort, targetPort))
	}

	if len(portPairs) == 0 {
		return cancel, nil, errors.ErrPortForwardPortPairNotFound
	}

	readyChan := make(chan struct{})
	pf, err := portforward.NewOnAddresses(
		spdy.NewDialer(
			upgrader, hc, http.MethodPost,
			c.GetClientSet().CoreV1().RESTClient().Post().
				Resource("pods").
				Namespace(namespace).
				Name(podName).
				SubResource("portforward").URL()),
		addresses, portPairs, ctx.Done(), readyChan, os.Stdout, os.Stderr)
	if err != nil {
		return cancel, nil, err
	}

	ech = make(chan error, 1)
	errgroup.Go(safety.RecoverFunc(func() (err error) {
		defer cancel()
		defer close(ech)
		if err = pf.ForwardPorts(); err != nil {
			ech <- err
		}
		return nil
	}))

	select {
	case <-ctx.Done():
		return cancel, ech, ctx.Err()
	case <-readyChan:
		return cancel, ech, nil
	}
}

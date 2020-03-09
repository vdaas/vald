//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package net provides net functionality for grpc
package net

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"syscall"

	"github.com/vdaas/vald/internal/errgroup"
)

const (
	localIPv4   = "127.0.0.1"
	localIPv6   = "::1"
	localHost   = "localhost"
	defaultPort = "80"
)

type Conn = net.Conn
type Dialer = net.Dialer
type ListenConfig = net.ListenConfig
type Listener = net.Listener
type Resolver = net.Resolver

var (
	DefaultResolver = net.DefaultResolver
)

func Listen(network, address string) (Listener, error) {
	return net.Listen(network, address)
}

func IsLocal(host string) bool {
	return host == localHost ||
		host == localIPv4 ||
		host == localIPv6
}

func Parse(addr string) (host string, port string, isIP bool, err error) {
	host, port, err = SplitHostPort(addr)
	return host, port, IsIPv6(host) || IsIPv4(host), err
}

func IsIPv6(addr string) bool {
	return net.ParseIP(addr) != nil && strings.Count(addr, ":") >= 2
}

func IsIPv4(addr string) bool {
	return net.ParseIP(addr) != nil && strings.Count(addr, ":") < 2
}

func SplitHostPort(hostport string) (string, string, error) {
	switch {
	case strings.HasPrefix(hostport, "::"):
		hostport = localIPv6 + hostport
	case strings.HasPrefix(hostport, ":"):
		hostport = localIPv4 + hostport
	}
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
		port = defaultPort
	}
	return host, port, err
}

func ScanPorts(ctx context.Context, start, end uint16, host string) (ports []uint16, err error) {
	var rl syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rl)
	if err != nil {
		return nil, err
	}
	eg, egctx := errgroup.New(ctx)
	eg.Limitation(int(rl.Max) / 2)
	var mu sync.Mutex
	for i := start; i <= end; i++ {
		port := i
		eg.Go(func() error {
			select {
			case <-egctx.Done():
				return egctx.Err()
			default:
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
				if err != nil {
					return err
				}
				mu.Lock()
				ports = append(ports, port)
				mu.Unlock()
				return conn.Close()
			}
		})
	}
	err = eg.Wait()

	return ports, nil
}

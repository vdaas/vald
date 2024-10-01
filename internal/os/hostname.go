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

// Package os provides os functions
package os

import (
	"net"
	"os"

	"github.com/vdaas/vald/internal/strings"
)

const unknownHost = "unknown-host"

var hostname = func() string {
	h, err := os.Hostname()
	if err != nil {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return unknownHost
		}
		ips := make([]string, 0, len(addrs))
		for _, addr := range addrs {
			if ipn, ok := addr.(*net.IPNet); ok && !ipn.IP.IsLoopback() {
				ips = append(ips, ipn.IP.String())
			}
		}
		if len(ips) == 0 {
			return unknownHost
		}
		return strings.Join(ips, ",\t")
	}
	return h
}()

func Hostname() (hn string, err error) {
	if hostname != "" {
		return hostname, nil
	}
	return os.Hostname()
}

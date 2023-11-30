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

// Package errors provides error types and function
package errors

import "time"

var (
	// tcp.

	// ErrFailedInitDialer represents an error that the initialization of the dialer failed.
	ErrFailedInitDialer = New("failed to init dialer")

	// ErrInvalidDNSConfig represents a function to generate an error that the configuration of the DNS is invalid.
	ErrInvalidDNSConfig = func(dnsRefreshDur, dnsCacheExp time.Duration) error {
		return Errorf("dnsRefreshDuration  > dnsCacheExp, %s, %s", dnsRefreshDur, dnsCacheExp)
	}

	// net.

	// ErrNoPortAvailable represents a function to generate an error that the port of the host is unavailable.
	ErrNoPortAvailable = func(host string, start, end uint16) error {
		return Errorf("no port available for Host: %s\tbetween %d ~ %d", host, start, end)
	}

	// ErrLookupIPAddrNotFound represents a function to generate an error that the host's ip address could not discovererd from DNS.
	ErrLookupIPAddrNotFound = func(host string) error {
		return Errorf("failed to lookup ip addrs for host: %s", host)
	}
)

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

package errors

var (
	// TLS.

	// ErrTLSDisabled is error variable, it's represents config error that tls is disabled by config.
	ErrTLSDisabled = New("tls feature is disabled")

	// ErrTLSCertOrKeyNotFound is error variable, it's represents tls cert or key not found error.
	ErrTLSCertOrKeyNotFound = New("cert or key file path not found")

	ErrCertificationFailed = New("certification failed")

	ErrUnsupportedClientAuthType = New("Unsupported Client Auth Type")

	ErrFailedToHandshakeTLSConnection = func(network, addr string) error {
		return Errorf("failed to handshake connection to %s:%s", network, addr)
	}

	ErrFailedToLoadCertKey = func(role, cert, key string) error {
		return Errorf("failed to load %s certificate/key (cert: %s, key: %s)", role, cert, key)
	}

	ErrNoCertsAddedToPool = New("no CA or self-signed certificates added to pool")

	ErrNoCertsFoundInPEM = New("no certificates found from PEM")
)

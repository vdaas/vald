//go:build wasm && js && !windows && !linux && !darwin

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

package control

var SetsockoptInt = func(fd, level, opt int, value int) (err error) {
	return nil
}

const (
	SOL_SOCKET           = 0
	IPPROTO_TCP          = 0
	SOL_IP               = 0
	SOL_IPV6             = 0
	SO_REUSEADDR         = 0
	SO_REUSEPORT         = 0
	SO_KEEPALIVE         = 0
	TCP_NODELAY          = 0
	TCP_CORK             = 0
	TCP_QUICKACK         = 0
	TCP_DEFER_ACCEPT     = 0
	TCP_KEEPINTVL        = 0
	TCP_KEEPIDLE         = 0
	TCP_FASTOPEN         = 0
	TCP_FASTOPEN_CONNECT = 0
	IP_TRANSPARENT       = 0
	IPV6_TRANSPARENT     = 0
	IP_RECVORIGDSTADDR   = 0
	IPV6_RECVORIGDSTADDR = 0
)

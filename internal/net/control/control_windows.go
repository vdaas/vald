//go:build windows

//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package control provides network socket option
package control

import (
	"golang.org/x/sys/windows"
)

var SetsockoptInt = func(fd, level, opt int, value int) (err error) {
	if level == windows.SOL_SOCKET && opt == windows.SO_REUSEADDR {
		return windows.SetsockoptInt(windows.Handle(fd), level, opt, value)
	}
	return nil
}

const (
	SOL_SOCKET           = windows.SOL_SOCKET
	IPPROTO_TCP          = 0
	SOL_IP               = 0
	SOL_IPV6             = 0
	SO_REUSEADDR         = windows.SO_REUSEADDR
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

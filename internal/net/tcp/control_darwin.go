// +build darwin,!linux,!windows,!wasm,!js

//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package tcp provides tcp option
package tcp

import (
	"syscall"

	"github.com/vdaas/vald/internal/errors"
	"golang.org/x/sys/unix"
)

const TCP_FASTOPEN int = 0x17

func Control(network, address string, c syscall.RawConn) (err error) {
	return c.Control(func(fd uintptr) {
		f := int(fd)
		ierr := syscall.SetsockoptInt(f, syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
		}
		ierr = syscall.SetsockoptInt(f, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
		}
		ierr = syscall.SetsockoptInt(f, syscall.IPPROTO_TCP, TCP_FASTOPEN, 1)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
		}
		ierr = syscall.SetsockoptInt(f, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
		}
	})
}

// +build !windows,!wasm,!js

// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
		ierr = syscall.SetsockoptInt(f, syscall.IPPROTO_TCP, syscall.TCP_DEFER_ACCEPT, 1)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
		}
	})
}

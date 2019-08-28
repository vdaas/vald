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

// Package servers provides implementation of Go API for managing server flow
package servers

import (
	"time"

	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*listener)

var (
	defaultOpts = []Option{}
)

func WithServer(srv server.Server) Option {
	return func(l *listener) {
		if srv == nil {
			return
		}
		if l.servers == nil {
			l.servers = make(map[string]server.Server)
		}
		l.servers[srv.Name()] = srv
	}
}

func WithShutdownDuration(dur string) Option {
	return func(l *listener) {
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Second * 20
		}
		l.sddur = d
	}
}

func WithStartUpStrategy(strg []string) Option {
	return func(l *listener) {
		l.sus = strg
	}
}

func WithShutdownStrategy(strg []string) Option {
	return func(l *listener) {
		l.sds = strg
	}
}

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
package downloader

import (
	"context"
	"io"
	"sync"
)

type ReadWriterAtCloserBuffer interface {
	io.Reader
	io.WriterAt
	io.Closer
}

type buffer struct {
	mu     sync.RWMutex
	data   *data
	pos    int64
	cur    int64
	nextCh chan int64
	ctx    context.Context
	cancel context.CancelFunc
}

type data struct {
	pos  int64
	size int64
	p    []byte
	next *data
}

func newBuffer(ctx context.Context, chunkSize int) ReadWriterAtCloserBuffer {
	ctx, cancel := context.WithCancel(ctx)
	b := &buffer{
		ctx:    ctx,
		cancel: cancel,
		nextCh: make(chan int64, chunkSize),
	}
	return b
}

func (b *buffer) WriteAt(p []byte, pos int64) (n int, err error) {
	pLen := len(p)
	d := &data{
		pos:  pos,
		p:    p,
		size: int64(pLen),
		next: nil,
	}
	b.mu.Lock()
	if b.data != nil {
		b.data = b.data.add(d)
	} else {
		b.data = d
	}
	b.pos = b.data.pos
	if b.pos == b.cur {
		b.mu.Unlock()
		select {
		case <-b.ctx.Done():
			return pLen, b.ctx.Err()
		case b.nextCh <- pos:
		}
	} else {
		b.mu.Unlock()
	}
	return pLen, nil
}

func (b *buffer) Read(p []byte) (n int, err error) {
	select {
	case <-b.ctx.Done():
		return 0, b.ctx.Err()
	case <-b.nextCh:
	}
	b.mu.Lock()
	data := b.data
	b.data = b.data.next
	b.cur = data.pos + data.size
	b.pos = data.pos
	b.mu.Unlock()
	copy(p, data.p)
	return len(data.p), nil
}

func (b *buffer) Close() (err error) {
	b.cancel()
	close(b.nextCh)
	return nil
}

func (d *data) add(n *data) *data {
	if d.pos > n.pos {
		d, n = n, d
	}
	if d.next != nil {
		n = d.next.add(n)
	}
	d.next = n
	return d
}

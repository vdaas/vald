//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package worker

import "context"

var (
	DefaultStartFunc = func(context.Context) (<-chan error, error) {
		return nil, nil
	}
	DefaultPushFunc = func(context.Context, JobFunc) error {
		return nil
	}
	DefaultPopFunc = func(context.Context) (JobFunc, error) {
		return nil, nil
	}
	DefaultLenFunc = func() uint64 {
		return uint64(0)
	}
)

type QueueMock struct {
	StartFunc func(context.Context) (<-chan error, error)
	PushFunc  func(context.Context, JobFunc) error
	PopFunc   func(context.Context) (JobFunc, error)
	LenFunc   func() uint64
}

func NewQueueMock() Queue {
	return &QueueMock{
		StartFunc: DefaultStartFunc,
		PushFunc:  DefaultPushFunc,
		PopFunc:   DefaultPopFunc,
		LenFunc:   DefaultLenFunc,
	}
}

func (q *QueueMock) Start(ctx context.Context) (<-chan error, error) {
	return q.StartFunc(ctx)
}

func (q *QueueMock) Push(ctx context.Context, job JobFunc) error {
	return q.PushFunc(ctx, job)
}

func (q *QueueMock) Pop(ctx context.Context) (JobFunc, error) {
	return q.PopFunc(ctx)
}

func (q *QueueMock) Len() uint64 {
	return q.LenFunc()
}

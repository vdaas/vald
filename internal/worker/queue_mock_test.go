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

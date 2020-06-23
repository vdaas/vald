package runner

import "context"

type runnerMock struct {
	PreStartFunc func(ctx context.Context) error
	StartFunc    func(ctx context.Context) (<-chan error, error)
	PreStopFunc  func(ctx context.Context) error
	StopFunc     func(ctx context.Context) error
	PostStopFunc func(ctx context.Context) error
}

func (m *runnerMock) PreStart(ctx context.Context) error {
	return m.PreStartFunc(ctx)
}

func (m *runnerMock) Start(ctx context.Context) (<-chan error, error) {
	return m.StartFunc(ctx)
}

func (m *runnerMock) PreStop(ctx context.Context) error {
	return m.PreStopFunc(ctx)
}
func (m *runnerMock) Stop(ctx context.Context) error {
	return m.StopFunc(ctx)
}
func (m *runnerMock) PostStop(ctx context.Context) error {
	return m.PostStopFunc(ctx)
}

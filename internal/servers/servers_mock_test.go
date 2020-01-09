package servers

import (
	"context"
)

type mockServer struct {
	NameFunc           func() string
	IsRunningFunc      func() bool
	ListenAndServeFunc func(context.Context, chan<- error) error
	ShutdownFunc       func(context.Context) error
}

func (ms *mockServer) Name() string {
	return ms.NameFunc()
}

func (ms *mockServer) IsRunning() bool {
	return ms.IsRunningFunc()
}

func (ms *mockServer) ListenAndServe(ctx context.Context, errCh chan<- error) error {
	return ms.ListenAndServeFunc(ctx, errCh)
}

func (ms *mockServer) Shutdown(ctx context.Context) error {
	return ms.ShutdownFunc(ctx)
}

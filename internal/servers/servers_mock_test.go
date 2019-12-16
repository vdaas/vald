package servers

import (
	"context"
)

type mockServer struct {
	NameFunc           NameFunc
	IsRunningFunc      IsRunningFunc
	ListenAndServeFunc ListenAndServeFunc
	ShutdownFunc       ShutdownFunc
}

type NameFunc func() string
type IsRunningFunc func() bool
type ListenAndServeFunc func(context.Context, chan<- error) error
type ShutdownFunc func(context.Context) error

func NewMockServer() *mockServer {
	return new(mockServer)
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

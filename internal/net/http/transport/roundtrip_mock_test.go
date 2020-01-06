package transport

import (
	"context"
	"net/http"
)

type roundTripMock struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (rm *roundTripMock) RoundTrip(req *http.Request) (*http.Response, error) {
	return rm.RoundTripFunc(req)
}

type backoffMock struct {
	DoFunc    func(context.Context, func() (interface{}, error)) (interface{}, error)
	CloseFunc func()
}

func (bm *backoffMock) Do(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	return bm.DoFunc(ctx, fn)
}

func (bm *backoffMock) Close() {
	bm.CloseFunc()
}

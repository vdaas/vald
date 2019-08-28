package grpc

import "google.golang.org/grpc"

type Handler interface {
	// GRPC interface here
	GetGRPCServer() *grpc.Server
}

type handler struct {
	gs *grpc.Server
}

func New() Handler {
	return nil
}

func (h *handler) GetGRPCServer() *grpc.Server {
	return h.gs
}

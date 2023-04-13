package mirror

import (
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
)

const PackageName = vald.PackageName

type Server interface {
	vald.Server
	MirrorServer
}

type UnimplementedValdServerWithMirror struct {
	vald.UnimplementedValdServer
	UnimplementedMirrorServer
}

const MirrorRPCServiceName = "Mirror"

const (
	RegisterRPCName  = "Register"
	AdvertiseRPCName = "Advertise"
)

func RegisterValdServerWithMirror(s *grpc.Server, srv Server) {
	vald.RegisterValdServer(s, srv)
	RegisterMirrorServer(s, srv)
}

package grpc

import "context"
import "github.com/vdaas/vald/internal/net/grpc/pool"

func ExecuteRPC(client Client, ctx context.Context, p pool.Conn, addr string, f func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)) (any, bool, error) {
	return client.(*gRPCClient).executeRPC(ctx, p, addr, f)
}

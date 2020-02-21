package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/vdaas/vald/apis/grpc/vald"

	"google.golang.org/grpc"
)

var (
	datasetDir     string
	grpcServerAddr string
)

func init() {
	flag.StringVar(&datasetDir, "path", "/tmp/", "set dataset path")
	flag.StringVar(&grpcServerAddr, "addr", ":8081", "set gRPC server address")
	flag.Parse()
}

func run() error {
	conn, err := grpc.DialContext(context.Background(), grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := vald.NewValdClient(conn)
	fmt.Println(client)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"gitlab.com/promptech1/infuser-author/grpc"
)

var (
	addr    = flag.String("addr", ":9090", "endpoint of the gRPC grpc")
	network    = flag.String("network", "tcp", `one of "tcp" or "unix". Must be consistent to -endpoint`)
)

func main() {
	flag.Parse()
	defer glog.Flush()

	ctx := context.Background()

	if err := grpc.Run(ctx, *network, *addr); err != nil {
		glog.Fatal(err)
	}
}
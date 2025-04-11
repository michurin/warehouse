package main

import (
	"context"
	"log"
	"net"

	"demo/kit/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func noerr(err error) {
	if err != nil {
		panic(err)
	}
}

type Calc struct {
	api.UnimplementedCalsServiceServer
}

func (c Calc) Square(ctx context.Context, req *api.Number) (*api.Number, error) {
	x := req.X
	log.Printf("Accept call: x=%g", x)
	return &api.Number{X: x * x}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9898")
	noerr(err)
	var opts []grpc.ServerOption
	mux := grpc.NewServer(opts...)
	api.RegisterCalsServiceServer(mux, Calc{})
	reflection.Register(mux)
	log.Printf("Serving on %s...", lis.Addr().String())
	err = mux.Serve(lis)
	noerr(err)
}

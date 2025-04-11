package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

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

func (c Calc) Sum(in grpc.ClientStreamingServer[api.Number, api.Number]) error {
	sum := float64(0)
	for {
		log.Println("Sum receiving...")
		x, err := in.Recv()
		if err == io.EOF {
			break
		}
		noerr(err)
		log.Printf("Got %#v (going to sleep)", x.X)
		sum += x.X
		time.Sleep(2 * time.Second)
	}
	log.Println("Sending response...")
	err := in.SendAndClose(&api.Number{X: sum})
	noerr(err)
	log.Println("Done.")
	return nil
}

func (c Calc) Repeat(in *api.Number, out grpc.ServerStreamingServer[api.Number]) error {
	for i := float64(0); i < in.X; i++ {
		err := out.Send(&api.Number{X: i})
		noerr(err)
	}
	return nil
}

func (c Calc) PipeSquare(stream grpc.BidiStreamingServer[api.Number, api.Number]) error {
	for {
		x, err := stream.Recv()
		if err == io.EOF {
			break
		}
		noerr(err)
		v := x.X
		go func() {
			err := stream.Send(&api.Number{X: v * v})
			noerr(err)
		}()
	}
	return nil
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

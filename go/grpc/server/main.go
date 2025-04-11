package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"demo/kit/api"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func InterceptorLogger(label string, l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, append([]any{slog.String("label", label)}, fields...)...)
	})
}

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

func (c Calc) Error(context.Context, *api.Empty) (*api.Empty, error) {
	return nil, status.Errorf(codes.FailedPrecondition, "CUSTOM ERROR MESSAGE")
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9898")
	noerr(err)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		logging.WithFieldsFromContext(func(context.Context) logging.Fields {
			return logging.Fields{"fromContext", "xxx"}
		}),
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger("unary", logger), logOpts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger("chain", logger), logOpts...),
		),
	}

	mux := grpc.NewServer(opts...)
	api.RegisterCalsServiceServer(mux, Calc{})
	reflection.Register(mux)
	log.Printf("Serving on %s...", lis.Addr().String())
	err = mux.Serve(lis)
	noerr(err)
}

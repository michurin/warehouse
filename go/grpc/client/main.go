package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"

	"demo/kit/api"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		logging.WithFieldsFromContext(func(context.Context) logging.Fields {
			return logging.Fields{"fromContext", "xxx"} // slog.Attr are not supported?
		}),
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(InterceptorLogger("unary", logger), logOpts...),
		),
		grpc.WithChainStreamInterceptor(
			logging.StreamClientInterceptor(InterceptorLogger("chain", logger), logOpts...),
		),
	}

	conn, err := grpc.NewClient("localhost:9898", opts...)
	noerr(err)
	defer conn.Close()
	log.Printf("Connected to %s", conn.CanonicalTarget())
	client := api.NewCalsServiceClient(conn)

	log.Println("------------- Square")

	resp, err := client.Square(ctx, &api.Number{X: 99})
	noerr(err)
	log.Printf("Result: %#v", resp.X)

	log.Println("------------- Sum")

	st, err := client.Sum(ctx)
	noerr(err)
	log.Println("Sending 100")
	st.Send(&api.Number{X: 100})
	log.Println("Sending 20")
	st.Send(&api.Number{X: 20})
	log.Println("Sending 3")
	st.Send(&api.Number{X: 3})
	resp, err = st.CloseAndRecv()
	noerr(err)
	log.Printf("Reslut: %#v (sum)", resp.X)

	log.Println("------------- Repeat")

	stx, err := client.Repeat(ctx, &api.Number{X: 3})
	noerr(err)
	for {
		x, err := stx.Recv()
		if err == io.EOF {
			break
		}
		noerr(err)
		log.Printf("Result: %#v (reprat)", x.X)
	}

	log.Println("------------- Pipe")

	stbi, err := client.PipeSquare(ctx)
	noerr(err)
	sync := make(chan struct{})
	go func() {
		defer close(sync)
		log.Println("-- reader")
		for range 3 { // we have to know how many results we wont to read
			x, err := stbi.Recv() // the order is random
			if err == io.EOF {
				break
			}
			noerr(err)
			log.Printf("-- reader: %#v", x.X)
		}
		log.Println("-- reader: done.")
	}()
	err = stbi.Send(&api.Number{X: 11})
	noerr(err)
	err = stbi.Send(&api.Number{X: 22})
	noerr(err)
	err = stbi.Send(&api.Number{X: 33})
	noerr(err)
	<-sync                 // we have to wait for all results
	err = stbi.CloseSend() // have to be closed, but it close both directions
	noerr(err)

	log.Println("------------- Error")
	_, err = client.Error(ctx, nil)
	log.Printf("Error [%+[1]T]: %[1]v", err)
	if e, ok := status.FromError(err); ok { // even such style is working
		log.Println("Code: ", e.Code())
		log.Println("Message: ", e.Message())
		log.Println("Details: ", e.Details())
	}
}

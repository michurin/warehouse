package main

import (
	"context"
	"log"
	"net/http"

	"demo/kit/api"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func noerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	conn, err := grpc.NewClient(
		"localhost:9898",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	noerr(err)

	ctx := context.Background()

	gwmux := runtime.NewServeMux()

	err = api.RegisterCalsServiceHandler(ctx, gwmux, conn)
	noerr(err)

	gwServer := &http.Server{
		Addr:    ":8999",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on", gwServer.Addr)
	log.Fatalln(gwServer.ListenAndServe())
}

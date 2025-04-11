package main

import (
	"context"
	"log"

	"demo/kit/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func noerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()
	conn, err := grpc.NewClient("localhost:9898", grpc.WithTransportCredentials(insecure.NewCredentials()))
	noerr(err)
	defer conn.Close()
	log.Printf("Connected to %s", conn.CanonicalTarget())
	client := api.NewCalsServiceClient(conn)
	resp, err := client.Square(ctx, &api.Number{X: 99})
	noerr(err)
	log.Printf("Result: %#v", resp.X)
}

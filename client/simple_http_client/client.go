package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/EDDYCJY/go-grpc-example/proto"
)

const PORT = "9003"

func main() {
	c, err := credentials.NewClientTLSFromFile("../../conf/server/server.pem", "go-grpc-example")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}

package main

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"google.golang.org/grpc"

	"github.com/EDDYCJY/go-grpc-example/pkg/gtls"
	pb "github.com/EDDYCJY/go-grpc-example/proto"
)

const (
	PORT = "9005"

	SERVICE_NAME              = "simple_zipkin_client"
	ZIPKIN_HTTP_ENDPOINT      = "http://127.0.0.1:9411/api/v1/spans"
	ZIPKIN_RECORDER_HOST_PORT = "127.0.0.1:9000"
)

func main() {
	collector, err := zipkin.NewHTTPCollector(ZIPKIN_HTTP_ENDPOINT)
	if err != nil {
		log.Fatalf("zipkin.NewHTTPCollector err: %v", err)
	}

	recorder := zipkin.NewRecorder(collector, true, ZIPKIN_RECORDER_HOST_PORT, SERVICE_NAME)

	tracer, err := zipkin.NewTracer(
		recorder, zipkin.ClientServerSameSpan(true),
	)
	if err != nil {
		log.Fatalf("zipkin.NewTracer err: %v", err)
	}

	tlsClient := gtls.Client{
		ServerName: "go-grpc-example",
		CaFile:     "../../conf/ca.pem",
		CertFile:   "../../conf/client/client.pem",
		KeyFile:    "../../conf/client/client.key",
	}

	c, err := tlsClient.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}

	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads()),
		))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC ZipKin",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}

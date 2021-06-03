package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/tonymontanapaffpaff/golang-training-university-grpc/proto/go_proto"

	"google.golang.org/grpc"
)

var (
	serverEndpoint = os.Getenv("SERVER_ENDPOINT")
	clientEndpoint = os.Getenv("CLIENT_ENDPOINT")
)

func init() {
	if serverEndpoint == "" {
		serverEndpoint = "localhost:8080"
	}
	if clientEndpoint == "" {
		clientEndpoint = "localhost:8383"
	}
}

func main() {
	conn, err := grpc.Dial(serverEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	grpcMux := runtime.NewServeMux()
	err = pb.RegisterCourseServiceHandler(context.Background(), grpcMux, conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(clientEndpoint, grpcMux))
}

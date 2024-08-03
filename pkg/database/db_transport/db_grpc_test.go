package db_transport

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/ethancox127/WatermarkService/api/v1/pb/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCGet(t *testing.T) {
	conn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:8082: %v", err)
	}
	defer conn.Close()
	c := pb.NewDatabaseClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("error calling function Get: %v", err)
	}

	log.Printf("Response from gRPC server's Get function: %s", r.GetDocuments())
}

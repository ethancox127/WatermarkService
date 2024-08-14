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

func TestGRPCAdd(t *testing.T) {
	conn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:8082: %v", err)
	}
	defer conn.Close()
	c := pb.NewDatabaseClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	doc := &pb.Document{Id: int32(-1), Content: "book", Title: "Dracula", Author: "Bram Stoker", Topic: "Fiction and Fantasy", Watermark: "False"}
	_, err = c.Add(ctx, &pb.AddRequest{Document: doc})
	if err != nil {
		log.Fatalf("error calling function Add: %v", err)
	}

	r, err := c.Get(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("error calling function Get: %v", err)
	}

	log.Printf("Response from gRPC server's Get function: %s", r.GetDocuments())
}

func TestGRPCUpdate(t *testing.T) {
	conn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:8082: %v", err)
	}
	defer conn.Close()
	c := pb.NewDatabaseClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	doc := &pb.Document{Id: int32(-1), Content: "book", Watermark: "true"}
	_, err = c.Update(ctx, &pb.UpdateRequest{Title: "Dracula", Document: doc})
	if err != nil {
		log.Fatalf("error calling function Update: %v", err)
	}

	r, err := c.Get(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("error calling function Get: %v", err)
	}

	log.Printf("Response from gRPC server's Get function: %s", r.GetDocuments())
}

func TestGRPCRemove(t *testing.T) {
	conn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:8082: %v", err)
	}
	defer conn.Close()
	c := pb.NewDatabaseClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.Remove(ctx, &pb.RemoveRequest{Title: "Dracula"})
	if err != nil {
		log.Fatalf("error calling function Remove: %v", err)
	}

	r, err := c.Get(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("error calling function Get: %v", err)
	}

	log.Printf("Response from gRPC server's Get function: %s", r.GetDocuments())
}

func TestGRPCServiceStatus(t *testing.T) {
	conn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:8082: %v", err)
	}
	defer conn.Close()
	c := pb.NewDatabaseClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.ServiceStatus(ctx, &pb.ServiceStatusRequest{})
	if err != nil {
		log.Fatalf("error calling function ServiceStatus: %v", err)
	}
}

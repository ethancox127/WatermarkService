package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/ethancox127/WatermarkService/api/v1/pb/database"
	"github.com/ethancox127/WatermarkService/internal/db_utils"
	"github.com/ethancox127/WatermarkService/pkg/database"
	"github.com/ethancox127/WatermarkService/pkg/database/db_transport"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

const (
	defaultHTTPPort = "8081"
	defaultGRPCPort = "8082"
)

func main() {
	var (
		logger   log.Logger
		httpAddr = net.JoinHostPort("localhost", envString("HTTP_PORT", defaultHTTPPort))
		grpcAddr = net.JoinHostPort("localhost", envString("GRPC_PORT", defaultGRPCPort))
	)

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	db, err := db_utils.ConnectDB()
	if err != nil {
		logger.Log(err)
		return
	}

	fmt.Println("db is connected")

	ctx := context.TODO()

	var (
		service     = database.NewService(ctx, db)
		httpHandler = db_transport.NewHTTPHandler(service)
		grpcServer  = db_transport.NewGRPCServer(service)
	)

	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", grpcAddr)
			// we add the Go Kit gRPC Interceptor to our gRPC service as it is used by
			// the here demonstrated zipkin tracing middleware.
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterDatabaseServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

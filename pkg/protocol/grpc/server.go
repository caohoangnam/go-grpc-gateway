package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	v1 "github.com/working/go-grpc-gateway/pkg/api/v1"
	"github.com/working/go-grpc-gateway/pkg/logger"
	"github.com/working/go-grpc-gateway/pkg/protocol/grpc/middleware"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, v1APIAccount v1.AccountServiceServer, v1APIEntries v1.EntriesServiceServer, v1APITransfers v1.TransfersServiceServer, port string) (err error) {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return
	}

	// gRPC server statup options
	opts := []grpc.ServerOption{}

	// add middleware
	opts = middleware.AddLogging(logger.Log, opts)

	//register server
	server := grpc.NewServer(opts...)
	v1.RegisterToDoServiceServer(server, v1API)
	v1.RegisterAccountServiceServer(server, v1APIAccount)
	v1.RegisterEntriesServiceServer(server, v1APIEntries)
	v1.RegisterTransfersServiceServer(server, v1APITransfers)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Log.Info("starting gRPC server...")
	return server.Serve(listen)
}

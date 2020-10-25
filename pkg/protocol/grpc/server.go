package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/working/go-grpc-gateway/pkg/api/v1"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) (err error) {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return
	}

	//register server
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}

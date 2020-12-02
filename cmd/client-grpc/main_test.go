package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	pb "github.com/working/go-grpc-gateway/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type ToDoServiceServer interface{}

func dialer() func(context.Context, string) (net.Conn, error) {
	//create size connect
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterToDoServiceServer(server, &pb.ToDoServiceServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestToDoServiceClient_Create(t *testing.T) {
	tests := []struct {
		title       string
		description string
		res         *pb.CreateResponse
	}{
		{
			"CaoNam",
			"ABC",
			&pb.CreateResponse{Api: "v1"},
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewToDoServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			req := &pb.CreateRequest{
				Api: "v1",
				ToDo: &pb.ToDo{
					Description: tt.description,
				},
			}

			res, err := client.Create(ctx, req)
			if err != nil {
				log.Fatal(err)
			}
			if res != nil {
				fmt.Println("Successfully")
			}
		})
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/working/go-grpc-gateway/pkg/api/v1"
	"google.golang.org/grpc"
)

type ToDoServiceServer interface{}

func TestToDoServiceClient_Create(t *testing.T) {
	ctx := context.Background()
	tTime := time.Now().In(time.UTC)
	tReminder, _ := ptypes.TimestampProto(tTime)

	tests := []struct {
		title       string
		description string
		reminder    *google_protobuf.Timestamp
		res         *pb.CreateResponse
	}{
		{
			"CaoNam",
			"ABC",
			tReminder,
			&pb.CreateResponse{Api: "v1"},
		},
	}

	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// connect server port
	//	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewToDoServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			req := &pb.CreateRequest{
				Api: "v1",
				ToDo: &pb.ToDo{
					Title:       tt.title,
					Description: tt.description,
					Reminder:    tt.reminder,
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

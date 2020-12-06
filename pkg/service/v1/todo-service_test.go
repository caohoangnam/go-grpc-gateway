package v1

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/working/go-grpc-gateway/pkg/api/v1"
)

func TestToDoServiceServer_Create(t *testing.T) {
	fmt.Println("Successfully connected!")
	ctx := context.Background()

	tTime := time.Now().In(time.UTC)
	tReminder, _ := ptypes.TimestampProto(tTime)

	// connect server port
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

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
		{
			"CaoNam1",
			"ABC",
			tReminder,
			&pb.CreateResponse{Api: "v1"},
		},
		{
			"CaoNam2",
			"ABC",
			tReminder,
			&pb.CreateResponse{Api: "v1"},
		},
		{
			"CaoNam3",
			"ABC",
			tReminder,
			&pb.CreateResponse{Api: "v1"},
		},
		{
			"CaoNam4",
			"AAAA",
			tReminder,
			&pb.CreateResponse{Api: "v1"},
		},
	}

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

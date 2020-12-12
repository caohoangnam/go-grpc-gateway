package v1

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"

	pb "github.com/working/go-grpc-gateway/pkg/api/v1"
)

func TestTransfersServiceServer_Create(t *testing.T) {
	fmt.Println("Successfully connected!")
	ctx := context.Background()

	//	address := flag.String("server", "", "gRPC server in format host:port")
	//	flag.Parse()

	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTransfersServiceClient(conn)
	//	TransfersRequestOne(t, ctx, client)
	TransfersRequestThree(t, ctx, client)
	//	TransfersRequestTwo(t, ctx, client)
}

func TransfersRequestOne(t *testing.T, ctx context.Context, client pb.TransfersServiceClient) {
	tests := []struct {
		FromAccountId int64
		ToAccountId   int64
		Amount        float64
		res           *pb.CreateTransfersResponse
	}{
		{
			1,
			2,
			60,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
		{
			2,
			1,
			30,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
		{
			2,
			1,
			15,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
		{
			2,
			1,
			25,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
		{
			1,
			2,
			60,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
		{
			2,
			1,
			30,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
		{
			2,
			1,
			15,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
		{
			2,
			1,
			25,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
	}

	stt := make(chan bool)
	for _, tt := range tests {
		fmt.Println("CBBB")
		t.Run("transfers"+string(tt.FromAccountId), func(t *testing.T) {
			go func() {
				req := &pb.CreateTransfersRequest{
					Api: "v1",
					Transfers: &pb.Transfers{
						FromAccountId: tt.FromAccountId,
						ToAccountId:   tt.ToAccountId,
						Amount:        tt.Amount,
					},
				}
				res, err := client.Create(ctx, req)
				if err != nil {
					stt <- false
					panic(err)
				}
				fmt.Println("Res", res.Id)
				stt <- true
			}()
		})
	}
	xxx := <-stt
	fmt.Println(xxx)
}

func TransfersRequestTwo(t *testing.T, ctx context.Context, client pb.TransfersServiceClient) {
	tests := []struct {
		FromAccountId int64
		ToAccountId   int64
		Amount        float64
		res           *pb.CreateTransfersResponse
	}{
		{
			1,
			2,
			20,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
		{
			2,
			1,
			30,
			&pb.CreateTransfersResponse{Api: "v1"},
		},
		//		{
		//			1,
		//			2,
		//			15,
		//			&pb.CreateTransfersResponse{Api: "v1"},
		//		},
		//		{
		//			2,
		//			1,
		//			25,
		//			&pb.CreateTransfersResponse{Api: "v1"},
		//		},
		//		{
		//			1,
		//			2,
		//			20,
		//			&pb.CreateTransfersResponse{Api: "v1"},
		//		},
		//		{
		//			2,
		//			1,
		//			30,
		//			&pb.CreateTransfersResponse{Api: "v1"},
		//		},
		//		{
		//			1,
		//			2,
		//			15,
		//			&pb.CreateTransfersResponse{Api: "v1"},
		//		},
		//		{
		//			2,
		//			1,
		//			25,
		//			&pb.CreateTransfersResponse{Api: "v1"},
		//		},
	}
	for _, tt := range tests {
		fmt.Println("CBBB")
		t.Run("transfers"+string(tt.FromAccountId), func(t *testing.T) {
			req := &pb.CreateTransfersRequest{
				Api: "v1",
				Transfers: &pb.Transfers{
					FromAccountId: tt.FromAccountId,
					ToAccountId:   tt.ToAccountId,
					Amount:        tt.Amount,
				},
			}
			res, err := client.Create(ctx, req)
			if err != nil {
				//				status <- false
				panic(err)
			}
			fmt.Println("Res", res.Id)
		})
	}
	//	status <- true
}
func TransfersRequestThree(t *testing.T, ctx context.Context, client pb.TransfersServiceClient) {
	type transfers struct {
		FromAccountId int64
		ToAccountId   int64
		Amount        float64
		res           *pb.CreateTransfersResponse
	}
	item := []transfers{}
	for i := 0; i < 56; i++ {
		item = append(item, transfers{
			FromAccountId: 1,
			ToAccountId:   2,
			Amount:        float64(rand.Intn(1000)),
			res:           &pb.CreateTransfersResponse{Api: "v1"},
		})
	}

	var wg sync.WaitGroup
	for _, tt := range item {
		wg.Add(1)
		t.Run("transfers"+string(tt.FromAccountId), func(t *testing.T) {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				req := &pb.CreateTransfersRequest{
					Api: "v1",
					Transfers: &pb.Transfers{
						FromAccountId: tt.FromAccountId,
						ToAccountId:   tt.ToAccountId,
						Amount:        tt.Amount,
					},
				}
				res, err := client.Create(ctx, req)
				if err != nil {
					panic(err)
				}
				fmt.Println("Res", res.Id)
				time.Sleep(time.Second)
			}(&wg)
		})
	}
	wg.Wait()
}

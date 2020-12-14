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
	TransfersRequest(t, ctx, client)
}

type Transfers struct {
	FromAccountId int64
	ToAccountId   int64
	Amount        float64
	res           *pb.CreateTransfersResponse
}

func TransfersRequest(t *testing.T, ctx context.Context, client pb.TransfersServiceClient) {
	item := []Transfers{}
	min := 10
	max := 30
	for i := 0; i < 40; i++ {
		rand.Seed(time.Now().UnixNano())
		item = append(item, Transfers{
			FromAccountId: 1 + int64(rand.Intn(4)),
			ToAccountId:   1 + int64(rand.Intn(4)),
			Amount:        float64(rand.Intn(max-min+1) + min),
			res:           &pb.CreateTransfersResponse{Api: "v1"},
		})
	}

	var wg sync.WaitGroup
	ch := make(chan struct{})
	wg.Add(len(item))
	for i, tt := range item {
		t.Run("transfers"+string(tt.FromAccountId), func(t *testing.T) {
			go producer(ch, &wg, tt, ctx, client, i)
		})
	}
	wg.Wait()
}

func producer(ch chan struct{}, wg *sync.WaitGroup, tt Transfers, ctx context.Context, client pb.TransfersServiceClient, index int) {
	fmt.Printf("Worker %d starting by FromToAccountId: %d \n", index, tt.FromAccountId)
	req := &pb.CreateTransfersRequest{
		Api: "v1",
		Transfers: &pb.Transfers{
			FromAccountId: tt.FromAccountId,
			ToAccountId:   tt.ToAccountId,
			Amount:        tt.Amount,
		},
	}
	_, err := client.Create(ctx, req)
	if err != nil {
		fmt.Println("ERR -> ", err)
	}
	fmt.Printf("Worker %d done by FromAccountId: %d \n", index, tt.FromAccountId)
	defer wg.Done()
}

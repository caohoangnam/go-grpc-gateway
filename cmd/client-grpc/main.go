package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	v1 "github.com/working/go-grpc-gateway/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {
	start := time.Now()

	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//	c := v1.NewToDoServiceClient(conn)
	account := v1.NewAccountServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Create Account
	reqAcc := v1.CreateAccountRequest{
		Api: apiVersion,
		Account: &v1.Account{
			Owner:    "HoangA",
			Balance:  2000000,
			Currency: "Dollar",
		},
	}
	resAcc, err := account.Create(ctx, &reqAcc)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", resAcc)

	//	t := time.Now().In(time.UTC)
	//	reminder, _ := ptypes.TimestampProto(t)
	//	pfx := t.Format(time.RFC3339Nano)
	//
	//	// Call Create
	//	req1 := v1.CreateRequest{
	//		Api: apiVersion,
	//		ToDo: &v1.ToDo{
	//			Title:       "title" + pfx + ")",
	//			Description: "description (" + pfx + ")",
	//			Reminder:    reminder,
	//		},
	//	}
	//	res1, err := c.Create(ctx, &req1)
	//	if err != nil {
	//		log.Fatalf("Create failed: %v", err)
	//	}
	//	log.Printf("Create result: <%+v>\n\n", res1)
	//
	//	//Read
	//	req2 := v1.ReadRequest{
	//		Api: apiVersion,
	//		Id:  res1.Id,
	//	}
	//	res2, err := c.Read(ctx, &req2)
	//	if err != nil {
	//		log.Fatalf("Read failed: %v", err)
	//	}
	//	log.Printf("Read result: <%+v>\n\n", res2)
	//
	//	// Update
	//	req3 := v1.UpdateRequest{
	//		Api: apiVersion,
	//		ToDo: &v1.ToDo{
	//			Id:          res2.ToDo.Id,
	//			Title:       res2.ToDo.Title,
	//			Description: res2.ToDo.Description + " + updated",
	//			Reminder:    res2.ToDo.Reminder,
	//		},
	//	}
	//	res3, err := c.Update(ctx, &req3)
	//	if err != nil {
	//		log.Fatalf("Update failed: %v", err)
	//	}
	//	log.Printf("Update result: <%+v>\n\n", res3)
	//
	//	// Call ReadAll
	//	req4 := v1.ReadAllRequest{
	//		Api: apiVersion,
	//	}
	//	res4, err := c.ReadAll(ctx, &req4)
	//	if err != nil {
	//		log.Fatalf("ReadAll failed: %v", err)
	//	}
	//	log.Printf("ReadAll result: <%+v>\n\n", res4)
	//
	//	// Delete
	//	req5 := v1.DeleteRequest{
	//		Api: apiVersion,
	//		Id:  res1.Id,
	//	}
	//	res5, err := c.Delete(ctx, &req5)
	//	if err != nil {
	//		log.Fatalf("Delete failed: %v", err)
	//	}
	//	log.Printf("Delete result: <%+v>\n\n", res5)
	//
	// Execution time
	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

}

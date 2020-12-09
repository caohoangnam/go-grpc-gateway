package cmd

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/working/go-grpc-gateway/pkg/logger"
	"github.com/working/go-grpc-gateway/pkg/protocol/grpc"
	"github.com/working/go-grpc-gateway/pkg/protocol/rest"
	v1 "github.com/working/go-grpc-gateway/pkg/service/v1"
)

const (
	GRPCPort = "9090"
	GRPCHttp = "8080"
	DBHost   = "localhost"
	DBPort   = 5432
	DBUser   = "postgres"
	DBPass   = ""
	DBName   = "postgres"

	//Logger
	LogLevel      = -1
	LogTimeFormat = "2006-01-02T15:04:05.999999999Z07:00"
)

//RunServer run gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// initialize logger
	if err := logger.Init(LogLevel, LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPass, DBName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	v1API := v1.NewToDoServiceServer(db)
	v1APIAccount := v1.NewAccountServiceServer(db)
	v1APIEntries := v1.NewEntriesServiceServer(db)
	v1APITransfers := v1.NewTransfersServiceServer(db)

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, GRPCPort, GRPCHttp)
	}()

	return grpc.RunServer(ctx, v1API, v1APIAccount, v1APIEntries, v1APITransfers, GRPCPort)
}

package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	"github.com/go-grpc-gateway/pkg/protocol/grpc"
	v1 "github.com/go-grpc-gateway/pkg/service/v1"
	_ "github.com/lib/pq"
)

type Config struct {
	//gRPC is TCP port to listen by gRPC server
	GRPCPort string

	DBHost   string
	DBUser   string
	DBPass   string
	DBSchema string
}

//RunServer run gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	//get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DBPass, "db-password", "", "Database password")
	flag.StringVar(&cfg.DBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBSchema)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}

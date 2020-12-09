package v1

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	v1 "github.com/working/go-grpc-gateway/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type transactionServiceServer struct {
	db *sql.DB
}

func NewTransactionServiceServer(db *sql.DB) v1.TransactionServiceServer {
	return &transactionServiceServer{
		db: db,
	}
}

// CheckAPI if the API requested by client is supported by server
func (t *transactionServiceServer) checkAPI(api string) error {
	// API version is "" means use current verions of ther service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}

	return nil
}

func (t *transactionServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := t.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create Transaction task
func (t *transactionServiceServer) Create(ctx context.Context, req *v1.CreateTransactionRequest) (*v1.CreateTransactionResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	lastInsertId := 0
	sql := `INSERT INTO transaction(from_account_id, to_account_id, amount) VALUES($1, $2) RETURNING id`
	err = db.QueryRowContext(ctx, sql, req.Transaction.FromAccountId, req.Transaction.ToAccountId, req.Transaction.Amount).Scan(&lastInsertId)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Transaction-> "+err.Error())
	}

	return &v1.CreateTransactionResponse{
		Api: apiVersion,
		Id:  int64(lastInsertId),
	}, nil
}

func (t *transactionServiceServer) Read(ctx context.Context, req *v1.ReadTransactionRequest) (*v1.ReadTransactionResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := `SELECT * FROM transaction WHERE id = $1`
	rows, err := db.QueryContext(ctx, sql, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Transaction ->"+err.Error())
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Transaction-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Transaction with ID='%d' is not found",
			req.Id))
	}
	// get Transaction data
	var td v1.Transaction
	if err := rows.Scan(&td.Id, &td.FromAccountId, &td.ToAccountId, &td.Amount); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from Transaction row-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple Transaction rows with ID='%d'",
			req.Id))
	}
	return &v1.ReadTransactionResponse{
		Api:         apiVersion,
		Transaction: &td,
	}, nil

}

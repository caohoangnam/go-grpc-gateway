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

type transfersServiceServer struct {
	db *sql.DB
}

func NewTransfersServiceServer(db *sql.DB) v1.TransfersServiceServer {
	return &transfersServiceServer{
		db: db,
	}
}

// CheckAPI if the API requested by client is supported by server
func (t *transfersServiceServer) checkAPI(api string) error {
	// API version is "" means use current verions of ther service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}

	return nil
}

func (t *transfersServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := t.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create Transfers task
func (t *transfersServiceServer) Create(ctx context.Context, req *v1.CreateTransfersRequest) (*v1.CreateTransfersResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lastInsertId, err := t.CreateTx(req, tx)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Unknown, "failed to insert into Transfers-> "+err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.CreateTransfersResponse{
		Api: apiVersion,
		Id:  int64(lastInsertId),
	}, nil
}

func (t *transfersServiceServer) CreateTx(req *v1.CreateTransfersRequest, tx *sql.Tx) (int, error) {
	var lastId int
	sql := `INSERT INTO transfers(from_account_id, to_account_id, amount) VALUES($1, $2) RETURNING id`
	err := tx.QueryRow(sql, req.Transfers.FromAccountId, req.Transfers.ToAccountId, req.Transfers.Amount).Scan(&lastId)
	return lastId, err
}

func (t *transfersServiceServer) Read(ctx context.Context, req *v1.ReadTransfersRequest) (*v1.ReadTransfersResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := `SELECT * FROM transfers WHERE id = $1`
	rows, err := db.QueryContext(ctx, sql, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Transfers ->"+err.Error())
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Transfers-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Transfers with ID='%d' is not found",
			req.Id))
	}
	// get Transfers data
	var td v1.Transfers
	if err := rows.Scan(&td.Id, &td.FromAccountId, &td.ToAccountId, &td.Amount); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from Transfers row-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple Transfers rows with ID='%d'",
			req.Id))
	}
	return &v1.ReadTransfersResponse{
		Api:       apiVersion,
		Transfers: &td,
	}, nil

}

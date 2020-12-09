package v1

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
	v1 "github.com/working/go-grpc-gateway/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type accountServiceServer struct {
	db *sql.DB
}

func NewAccountServiceServer(db *sql.DB) v1.AccountServiceServer {
	return &accountServiceServer{
		db: db,
	}
}

// CheckAPI if the API requested by client is supported by server
func (t *accountServiceServer) checkAPI(api string) error {
	// API version is "" means use current verions of ther service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}

	return nil
}

func (t *accountServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := t.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create Account task
func (t *accountServiceServer) Create(ctx context.Context, req *v1.CreateAccountRequest) (*v1.CreateAccountResponse, error) {
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
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.CreateAccountResponse{
		Api: apiVersion,
		Id:  int64(lastInsertId),
	}, nil
}

func (t *accountServiceServer) CreateTx(req *v1.CreateAccountRequest, tx *sql.Tx) (int, error) {
	var lastId int
	sql := `INSERT INTO accounts(owner, balance, currency) VALUES($1, $2, $3) RETURNING id`
	err := tx.QueryRow(sql, req.Account.Owner, req.Account.Balance, req.Account.Currency).Scan(&lastId)
	return lastId, err
}

func (t *accountServiceServer) Read(ctx context.Context, req *v1.ReadAccountRequest) (*v1.ReadAccountResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := `SELECT * FROM accounts WHERE id = $1`
	rows, err := db.QueryContext(ctx, sql, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Account ->"+err.Error())
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Account-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Account with ID='%d' is not found",
			req.Id))
	}
	// get Account data
	var td v1.Account
	if err := rows.Scan(&td.Id, &td.Owner, &td.Balance, &td.Currency); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from Account row-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple Account rows with ID='%d'",
			req.Id))
	}
	return &v1.ReadAccountResponse{
		Api:     apiVersion,
		Account: &td,
	}, nil

}

func (t *accountServiceServer) Update(ctx context.Context, req *v1.UpdateAccountRequest) (*v1.UpdateAccountResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	updatedAt := ptypes.TimestampNow()
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field invaild format ->"+err.Error())
	}

	sql := `UPDATE accounts SET balance = $1, updated_at =  $2 WHERE id = $3`
	res, err := db.ExecContext(ctx, sql, req.Account.Balance, updatedAt, req.Account.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update Account ->"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Account with ID='%d' is not found",
			req.Account.Id))
	}

	return &v1.UpdateAccountResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil

}

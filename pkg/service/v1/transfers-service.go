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

	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Hanle account
	fromAccount := t.GetAccountByIdTx(int(req.Transfers.FromAccountId), tx)
	var balance float64
	err = fromAccount.Scan(&balance)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Unknown, "failed to get into From Account-> "+err.Error())
	}
	fmt.Println("fromAccount", balance)

	lastInsertId, err := t.CreateTx(req, tx)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Unknown, "failed to insert into Transfers-> "+err.Error())
	}
	fmt.Println("lastInsertId", lastInsertId)

	// Hanlde entries
	reqEntriesByFromAccount := &v1.CreateEntriesRequest{
		Entries: &v1.Entries{
			AccountId: req.Transfers.FromAccountId,
			Amount:    balance - req.Transfers.Amount,
		},
	}
	_, err = t.CreateEntriesTx(reqEntriesByFromAccount, tx)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Unknown, "failed to insert into Entries by from account -> "+err.Error())
	}

	reqEntriesByToAccount := &v1.CreateEntriesRequest{
		Entries: &v1.Entries{
			AccountId: req.Transfers.ToAccountId,
			Amount:    balance - req.Transfers.Amount,
		},
	}
	_, err = t.CreateEntriesTx(reqEntriesByToAccount, tx)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Unknown, "failed to insert into Entries by to account-> "+err.Error())
	}

	reqFromAccount := &v1.UpdateAccountRequest{
		Account: &v1.Account{
			Id:      req.Transfers.FromAccountId,
			Balance: balance - req.Transfers.Amount,
		},
	}
	_, err = t.UpdateAccountTx(reqFromAccount, tx)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Unknown, "failed to update into From Account-> "+err.Error())
	}

	toAccount := t.GetAccountByIdTx(int(req.Transfers.ToAccountId), tx)
	err = toAccount.Scan(&balance)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Unknown, "failed to get into to Account-> "+err.Error())
	}
	fmt.Println("toAccount", balance)
	reqToAccount := &v1.UpdateAccountRequest{
		Account: &v1.Account{
			Id:      req.Transfers.ToAccountId,
			Balance: balance - req.Transfers.Amount,
		},
	}
	_, err = t.UpdateAccountTx(reqToAccount, tx)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Unknown, "failed to update into To Account-> "+err.Error())
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
	sql := `INSERT INTO transfers(from_account_id, to_account_id, amount) VALUES($1, $2, $3) RETURNING id`
	err := tx.QueryRow(sql, req.Transfers.FromAccountId, req.Transfers.ToAccountId, req.Transfers.Amount).Scan(&lastId)
	return lastId, err
}

func (t *transfersServiceServer) CreateEntriesTx(req *v1.CreateEntriesRequest, tx *sql.Tx) (int, error) {
	var lastId int
	sql := `INSERT INTO entries(account_id, amount) VALUES($1, $2) RETURNING id`
	err := tx.QueryRow(sql, req.Entries.AccountId, req.Entries.Amount).Scan(&lastId)
	fmt.Println("err", err)
	return lastId, err
}

func (t *transfersServiceServer) GetAccountByIdTx(id int, tx *sql.Tx) *sql.Row {
	sql := `SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`
	res := tx.QueryRow(sql, id)
	return res
}

func (t *transfersServiceServer) UpdateAccountTx(req *v1.UpdateAccountRequest, tx *sql.Tx) (int, error) {
	var lastId int
	sql := `UPDATE accounts SET balance = $1, updated_at = now() WHERE id = $2 RETURNING id`
	err := tx.QueryRow(sql, req.Account.Balance, req.Account.Id).Scan(&lastId)
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

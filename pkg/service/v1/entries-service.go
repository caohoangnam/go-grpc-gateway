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

type entriesServiceServer struct {
	db *sql.DB
}

func NewEntriesServiceServer(db *sql.DB) v1.EntriesServiceServer {
	return &entriesServiceServer{
		db: db,
	}
}

// CheckAPI if the API requested by client is supported by server
func (t *entriesServiceServer) checkAPI(api string) error {
	// API version is "" means use current verions of ther service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}

	return nil
}

func (t *entriesServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := t.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create Entries task
func (t *entriesServiceServer) Create(ctx context.Context, req *v1.CreateEntriesRequest) (*v1.CreateEntriesResponse, error) {
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
		return nil, status.Error(codes.Unknown, "failed to insert into Entries-> "+err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.CreateEntriesResponse{
		Api: apiVersion,
		Id:  int64(lastInsertId),
	}, nil
}

func (t *entriesServiceServer) CreateTx(req *v1.CreateEntriesRequest, tx *sql.Tx) (int, error) {
	var lastId int
	sql := `INSERT INTO entries(account_id, amount) VALUES($1, $2) RETURNING id`
	err := tx.QueryRow(sql, req.Entries.AccountId, req.Entries.Amount).Scan(&lastId)
	return lastId, err
}

func (t *entriesServiceServer) Read(ctx context.Context, req *v1.ReadEntriesRequest) (*v1.ReadEntriesResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := `SELECT * FROM entries WHERE id = $1`
	rows, err := db.QueryContext(ctx, sql, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Entries ->"+err.Error())
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from Entries-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Entries with ID='%d' is not found",
			req.Id))
	}
	// get Entries data
	var td v1.Entries
	if err := rows.Scan(&td.Id, &td.AccountId, &td.Amount); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from Entries row-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple Entries rows with ID='%d'",
			req.Id))
	}
	return &v1.ReadEntriesResponse{
		Api:     apiVersion,
		Entries: &td,
	}, nil

}

func (t *entriesServiceServer) Update(ctx context.Context, req *v1.UpdateEntriesRequest) (*v1.UpdateEntriesResponse, error) {
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

	sql := `UPDATE entriess SET amount = $1, updated_at =  $2 WHERE id = $3`
	res, err := db.ExecContext(ctx, sql, req.Entries.Amount, updatedAt, req.Entries.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update Entries ->"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Entries with ID='%d' is not found",
			req.Entries.Id))
	}

	return &v1.UpdateEntriesResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil

}

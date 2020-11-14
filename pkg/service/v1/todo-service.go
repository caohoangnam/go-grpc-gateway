package v1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
	v1 "github.com/working/go-grpc-gateway/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type toDoServiceServer struct {
	db *sql.DB
}

func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{
		db: db,
	}
}

// CheckAPI if the API requested by client is supported by server
func (t *toDoServiceServer) checkAPI(api string) error {
	// API version is "" means use current verions of ther service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}

	return nil
}

func (t *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := t.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create ToDo task
func (t *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	lastInsertId := 0
	sql := `INSERT INTO todo(title, description, reminder) VALUES($1, $2, $3) RETURNING id`
	err = db.QueryRowContext(ctx, sql, req.ToDo.Title, req.ToDo.Description, reminder).Scan(&lastInsertId)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into ToDo-> "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  int64(lastInsertId),
	}, nil
}

func (t *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := `SELECT * FROM todo WHERE id = $1`
	rows, err := db.QueryContext(ctx, sql, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo ->"+err.Error())
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.Id))
	}
	// get ToDo data
	var td v1.ToDo
	var reminder time.Time
	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
	}
	td.Reminder, err = ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple ToDo rows with ID='%d'",
			req.Id))
	}
	return &v1.ReadResponse{
		Api:  apiVersion,
		ToDo: &td,
	}, nil

}

func (t *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field invaild format ->"+err.Error())
	}

	sql := `UPDATE todo SET title = $1, description = $2, reminder =  $3 WHERE id = $4`
	res, err := db.ExecContext(ctx, sql, req.ToDo.Title, req.ToDo.Description, reminder, req.ToDo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo ->"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.ToDo.Id))
	}

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil

}

func (t *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// delete ToDo
	sql := `DELETE FROM todo WHERE id = $1`
	res, err := db.ExecContext(ctx, sql, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete ToDo-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.Id))
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: rows,
	}, nil
}

func (t *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := t.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := `SELECT * FROM todo`
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}
	defer rows.Close()

	var reminder *time.Time
	list := []*v1.ToDo{}
	for rows.Next() {
		td := new(v1.ToDo)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
		}

		if reminder == nil {
			continue
		}
		td.Reminder, err = ptypes.TimestampProto(*reminder)
		if err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
		}
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
	}

	return &v1.ReadAllResponse{
		Api:   apiVersion,
		ToDos: list,
	}, nil
}

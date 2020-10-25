package v1

import (
	"context"
	"database/sql"

	v1 "github.com/go-grpc-gateway/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion  = "v1"
	DB_USER     = "postgres"
	DB_PASSWORD = ""
	DB_NAME     = "postgres"
)

type toDoServiceServer struct {
	db *sql.DB
}

func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{db: db}
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

	c, err := t.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	return &v1.CreateResponse{}, nil
}

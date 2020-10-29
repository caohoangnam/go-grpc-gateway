package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	v1 "github.com/working/go-grpc-gateway/pkg/api/v1"
	"google.golang.org/grpc"
)

// Run server HTTP/REST gateway
func RunServer(ctx context.Context, grpcPort string, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	addressHttp := "localhost:" + grpcPort
	err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, addressHttp, opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	server := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = server.Shutdown(ctx)
	}()

	log.Println("starting HTTP/REST gateway...")
	return server.ListenAndServe()
}

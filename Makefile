grpc:
	./third_party/protoc-gen.cmd
clearcache:
	go clean -testcache
server:
	go run ./cmd/server/main.go
test:
	go test -v -run TestTransfersServiceServer_Create ./pkg/service/v1/transfers-service_test.go

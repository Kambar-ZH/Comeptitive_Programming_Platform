grpc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/grpc/api/service.proto

run:
	go run cmd/myapp/main.go

build: 
	go build cmd/myapp/main.go
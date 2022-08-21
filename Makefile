proto-gen:
	protoc --go_out=. --go-grpc_out=. ./api/proto/*.proto
run-broker:
	go run cmd/broker/main.go
build-broker:
	go build cmd/broker/main.go
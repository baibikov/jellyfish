proto-gen:
	protoc --go_out=. --go-grpc_out=. ./api/proto/*.proto



goctl rpc protoc ./apps/customer/rpc/customer.proto --go_out=./apps/customer --go-grpc_out=./apps/customer --zrpc_out=./apps/customer/rpc -m

go build -o customer-rpc customer.go
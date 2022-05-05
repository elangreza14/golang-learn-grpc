gen:
	protoc -I/usr/bin -I . \
    -I${GOPATH}/src \
	-I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.2.0/third_party/googleapis \
	--go_out=plugins=grpc:pb --go_opt paths=import \
    --grpc-gateway_out ./pb \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=import \
    proto/auth.proto
    
gentwo:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:./pb --grpc-gateway_out ./pb

gengate:
	protoc --proto_path=proto proto/*.proto --grpc-gateway_out: ./pb

protoc:
	#protoc -I api api/hello_world.proto --go_out=plugins=grpc:api
	protoc \
      -I ./proto/ \
      -I ./api/ \
      --go_out ./api --go_opt paths=source_relative \
      --go-grpc_out ./api --go-grpc_opt paths=source_relative \
      --grpc-gateway_out ./api --grpc-gateway_opt paths=source_relative \
	  proto/auth.proto


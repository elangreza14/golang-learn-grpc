gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:.
del:
	cd pb/pcbook/proto && rm *.go
run:
	go run main.go
test:
	go test -cover -race ./...
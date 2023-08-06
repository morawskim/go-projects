.PHONY: protoc run

protoc:
	protoc -I=./ --go_out=./go ./order.proto
	protoc -I=./ --php_out=./php/protobuf ./order.proto

run:
	cd go && go run main.go

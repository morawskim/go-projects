.PHONY: protoc run

protoc:
	protoc -I=./ --go_out=./go ./order.proto

run:
	cd go && go run main.go

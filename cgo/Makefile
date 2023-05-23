.PHONY: build

build:
	go build -a -o cgo-example ./main.go ./bridge.go

run: build
	./cgo-example

memory:
	valgrind --leak-check=full ./cgo-example

.PHONY: build

build:
	go build -a -o cgo-example ./main.go

memory:
	valgrind --leak-check=full ./cgo-example

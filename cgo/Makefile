.PHONY: build

build:
	go build -a -o cgo-example

run: build
	./cgo-example

memory:
	valgrind --leak-check=full ./cgo-example

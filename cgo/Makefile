.PHONY: build

build:
	go build -a -o cgo-example

run: build
	./cgo-example

memory: build
	valgrind --leak-check=full ./cgo-example

memory-benchmark:
	GOMAXPROCS=1 valgrind  --leak-check=full go test -bench=. -benchtime=2s

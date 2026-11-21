-include Makefile.local
.PHONY: package build

COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
VERSION ?= dev
# create Makefile.local and set this variable
NOTIFICATION_URL ?= generic+https://example.com/send

LDFLAGS := -s -w \
    -X 'smtp/internal.Version=$(VERSION)' \
    -X 'smtp/internal.Commit=$(COMMIT)' \
    -X 'smtp/internal.BuildDate=$(BUILD_DATE)'

build:
	go build -ldflags "$(LDFLAGS)" -o gosmtp ./server.go

check-version:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "ERROR: VERSION must be set and cannot be 'dev'"; \
		exit 1; \
	fi

package: check-version build
	cd pkg/systemd/ && jsonnet service.libsonnet -S -m .
	docker run --rm --user "$(shell id -u):$(shell id -g)" -v $(PWD):/tmp -e PACKAGER=deb -eVERSION=$(VERSION) -w /tmp goreleaser/nfpm pkg --packager deb --target /tmp
	docker run --rm --user "$(shell id -u):$(shell id -g)" -v $(PWD):/tmp -e PACKAGER=rpm -eVERSION=$(VERSION) -w /tmp goreleaser/nfpm pkg --packager rpm --target /tmp

run: build
	./gosmtp -notification-url $(NOTIFICATION_URL) -addr '127.0.0.1:2525' -metrics-addr '127.0.0.1:8025'

send-test-mail:
	docker run --rm --net=host morawskim/swaks:latest --to user@example.com --server 127.0.0.1:2525

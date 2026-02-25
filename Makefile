BINARY   := etoro
MODULE   := github.com/etoro/etoro-cli
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE     := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS  := -s -w \
  -X '$(MODULE)/cmd.version=$(VERSION)' \
  -X '$(MODULE)/cmd.commit=$(COMMIT)' \
  -X '$(MODULE)/cmd.date=$(DATE)'

.PHONY: build test vet clean install release

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/

install: build
	mkdir -p $(HOME)/.local/bin
	cp $(BINARY) $(HOME)/.local/bin/$(BINARY)
	@echo "Installed to $(HOME)/.local/bin/$(BINARY)"

release: clean
	@echo "Building release artifacts..."
	@mkdir -p dist
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/etoro . && \
		tar -czf dist/etoro_darwin_amd64.tar.gz -C dist etoro && rm dist/etoro
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/etoro . && \
		tar -czf dist/etoro_darwin_arm64.tar.gz -C dist etoro && rm dist/etoro
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/etoro . && \
		tar -czf dist/etoro_linux_amd64.tar.gz  -C dist etoro && rm dist/etoro
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/etoro . && \
		tar -czf dist/etoro_linux_arm64.tar.gz  -C dist etoro && rm dist/etoro
	@cd dist && shasum -a 256 *.tar.gz > checksums.txt
	@echo "Release artifacts in dist/"
	@ls -lh dist/

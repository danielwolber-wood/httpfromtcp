.PHONY: all deps fmt vet sec test build clean ci
all: fmt vet sec test build

deps:
	@command -v gosec >/dev/null 2>&1 || { \
		echo "Installing gosec..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	}

fmt:
	go fmt ./...

vet:
	go vet ./...

sec:
	gosec ./...

test:
	go test ./...

build:
	mkdir -p bin/app
	go build -o bin/app ./...

clean:
	rm -rf bin

ci:
	all
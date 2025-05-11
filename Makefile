.PHONY: all deps fmt vet sec test build clean ci tcp udp

COMPONENTS := tcplistener udpsender

all: fmt vet sec test build
ci: deps fmt vet sec test build

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

build: $(COMPONENTS)

$(COMPONENTS):
	mkdir -p bin/$@
	go build -o bin/$@ ./cmd/$@

tcplistener: tcp
tcp:
	mkdir -p bin/tcplistener
	go build -o bin/tcplistener ./cmd/tcplistener

udpsender: udp
udp:
	mkdir -p bin/udpsender
	go build -o bin/udpsender ./cmd/udpsender

clean:
	rm -rf bin
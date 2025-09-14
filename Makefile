.PHONY: build test clean install

VERSION ?= 0.1.0
LDFLAGS = -ldflags="-X 'github.com/thebuidl-grid/starknode-kit/pkg/versions.StarkNodeVersion=$(VERSION)'"

# Build the application
build:
	go build $(LDFLAGS) -o bin/starknode .

# Run tests
test:
	cd pkg && go test -v ./ 

# Clean build artifacts
clean:
	rm -rf bin

# Install a client (example usage: make install CLIENT=geth)
install:
	@if [ -z "$(CLIENT)" ]; then \
		echo "Please specify a client with CLIENT=<client>"; \
		exit 1; \
	fi; \
	./bin/starknode --client $(CLIENT)

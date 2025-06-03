.PHONY: build test clean install

# Build the application
build:
    go build -o bin/ethereum-installer ./

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
    ./bin/ethereum-installer --client $(CLIENT)

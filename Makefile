BINARY_NAME ?= capturl-mcp
BIN_DIR ?= /usr/local/bin

.PHONY: build
build:
	go build -o $(BINARY_NAME) ./cmd

.PHONY: install
install: build
	@echo "Installing to $(BIN_DIR)..."
	install -m 755 $(BINARY_NAME) $(BIN_DIR)/$(BINARY_NAME)

.PHONY: uninstall
uninstall:
	rm -f $(BIN_DIR)/$(BINARY_NAME)

.PHONY: clean
clean:
	rm -rf dist

.PHONY: test
test:
	go test ./...
.PHONY: build clean fmt test debug example

OUTPUT_DIR := bin
BINARY_NAME := term_debugger
GODEBUG_FLAGS := -gcflags "-N -l"
GO_FLAGS ?=

build:
	@mkdir -p $(OUTPUT_DIR)
	go build $(GO_FLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/term_debugger

example:
	@echo "Building example..."
	make -C example OUTPUT_DIR=../bin
debug:
	@echo "Building for debug..."
	make build GO_FLAGS="$(GODEBUG_FLAGS)"

fmt:
	@echo "Formatting code..."
	@goimports -w ./

test:
	go test ./...

clean:
	@rm -rf $(OUTPUT_DIR)

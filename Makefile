# Go related variables
GO := go
GOFMT := gofmt
GOFLAGS :=

# Main binary name
BINARY_NAME := gkilo

# Directories
SRC_DIR := cmd/main
BUILD_DIR := bin

# Source files
SRC_FILES := $(wildcard $(SRC_DIR)/*.go)

.PHONY: all build clean run

all: build

build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(SRC_FILES)
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $@ $(SRC_DIR)/main.go

clean:
	rm -rf $(BUILD_DIR)

run: build
	$(BUILD_DIR)/$(BINARY_NAME) ${ARGS}

fmt:
	$(GOFMT) -w $(SRC_FILES)

help:
	@echo "Available targets:"
	@echo "  - all (default): Builds the project."
	@echo "  - build: Builds the project."
	@echo "  - clean: Cleans the build directory."
	@echo "  - run: Builds and runs the project."
	@echo "  - fmt: Formats the source code using 'gofmt'."
	@echo "  - help: Show this help message."

# Ensure that any Makefile target is not treated as a file
.PHONY: all build clean run fmt help

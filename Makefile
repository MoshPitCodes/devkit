# Go parameters
GO            ?= go
GOPROTO       ?= protoc

# Output directory
OUTPUT_DIR    ?= bin

# Binary names
CLI_BINARY    = $(OUTPUT_DIR)/devkit
AGENT_BINARY  = $(OUTPUT_DIR)/devkitd

# Proto paths
PROTO_SRC_DIR = pkg/api/proto/agent
PROTO_GEN_DIR = pkg/api/gen/agent
PROTO_FILES   = $(PROTO_SRC_DIR)/agent.proto

# Variables
CLI_CMD_PATH=./cmd/devkit
AGENT_CMD_PATH=./cmd/devkitd
CLI_OUTPUT=bin/devkit
AGENT_OUTPUT=bin/devkitd
PROTO_PATH=pkg/api/proto/agent
PROTO_OUT=pkg/api/gen/agent
DEV_TEMPLATES_SRC=templates
AGENT_TEMPLATES_DEST=.devkit/templates

# Default target
.PHONY: all
all: tidy copy-templates proto build

## Build targets
build: build-cli build-agent

build-cli: copy-templates # Ensure templates are copied before potential run
	@echo "Building devkit CLI..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build -v -o $(CLI_OUTPUT) $(CLI_CMD_PATH)

build-agent: copy-templates # Ensure templates are copied before potential run
	@echo "Building devkitd agent..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build -v -o $(AGENT_OUTPUT) $(AGENT_CMD_PATH)

## Proto generation
proto: $(PROTO_GEN_DIR)/agent.pb.go $(PROTO_GEN_DIR)/agent_grpc.pb.go

$(PROTO_GEN_DIR)/agent.pb.go $(PROTO_GEN_DIR)/agent_grpc.pb.go: $(PROTO_FILES)
	@echo "Generating protobuf code..."
	@mkdir -p $(PROTO_GEN_DIR)
	$(GOPROTO) \
		--go_out=$(PROTO_GEN_DIR) \
		--go-grpc_out=$(PROTO_GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=$(PROTO_SRC_DIR) \
		$(PROTO_FILES)

## Go module tasks
tidy:
	@echo "Running go mod tidy..."
	$(GO) mod tidy

## Cleaning
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@rm -rf $(AGENT_TEMPLATES_DEST) # Also remove copied templates

## Devkit Initialization
.PHONY: init
init: build-cli # Depends on the CLI being built
	@echo "Running devkit init..."
	@$(CLI_BINARY) init

# Copy templates from development dir to agent dir if dev dir exists
.PHONY: copy-templates
copy-templates:
	@echo "Ensuring agent template directory $(AGENT_TEMPLATES_DEST) exists..."
	@mkdir -p $(AGENT_TEMPLATES_DEST)
	@if [ -d "$(DEV_TEMPLATES_SRC)" ]; then \
		echo "Copying templates from $(DEV_TEMPLATES_SRC) to $(AGENT_TEMPLATES_DEST)..."; \
		cp -r $(DEV_TEMPLATES_SRC)/. $(AGENT_TEMPLATES_DEST)/ 2>/dev/null || true; \
	else \
		echo "Development template directory $(DEV_TEMPLATES_SRC) not found, skipping copy."; \
	fi

# Phony targets
.PHONY: all build build-cli build-agent proto tidy copy-templates clean init

# Help target
help:
	@echo "Makefile commands:"
	@echo "  all           - Run tidy, copy-templates, proto, build (default)"
	@echo "  build         - Build both CLI and agent"
	@echo "  build-cli     - Build only the CLI"
	@echo "  build-agent   - Build only the agent"
	@echo "  proto         - Regenerate protobuf code"
	@echo "  tidy          - Run go mod tidy"
	@echo "  init          - Run devkit init to initialize the project structure"
	@echo "  copy-templates - Copy dev templates to agent templates dir"
	@echo "  clean         - Remove build artifacts and copied templates"

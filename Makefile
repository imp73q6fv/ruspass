.PHONY: all build build-backend build-frontend clean run archive

PROJECT_NAME := ruspass
VERSION := 1.0.0
BUILD_DIR := build
ARCHIVE_NAME := $(PROJECT_NAME)-$(VERSION)-linux-amd64.tar.gz

all: build

build: build-backend archive

build-backend:
	@echo "=== Building Rust Backend ==="
	@cargo build --release 2>&1 | grep -v "warning:" || true
	@mkdir -p $(BUILD_DIR)
	@cp target/release/ruspass $(BUILD_DIR)/ruspass_backend
	@echo "Backend built successfully: $(BUILD_DIR)/ruspass_backend"

build-frontend:
	@echo "=== Building Go Frontend ==="
	@cd UI_UX && go build -o ../$(BUILD_DIR)/ruspass_frontend . 2>&1 || echo "Frontend build skipped (requires Qt bindings)"

archive: build-backend
	@echo "=== Creating Archive ==="
	@mkdir -p $(BUILD_DIR)
	@if [ -f $(BUILD_DIR)/ruspass_frontend ]; then \
		tar -czvf $(BUILD_DIR)/$(ARCHIVE_NAME) -C $(BUILD_DIR) ruspass_backend ruspass_frontend; \
	else \
		tar -czvf $(BUILD_DIR)/$(ARCHIVE_NAME) -C $(BUILD_DIR) ruspass_backend; \
	fi
	@echo "Archive created: $(BUILD_DIR)/$(ARCHIVE_NAME)"

run: build-backend
	@./target/release/ruspass

clean:
	@echo "=== Cleaning Build Artifacts ==="
	@cargo clean
	@rm -rf $(BUILD_DIR)
	@rm -f ruspass.db
	@echo "Clean complete"

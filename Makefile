.PHONY: all build build-backend build-frontend clean run archive

PROJECT_NAME := ruspass
VERSION := 1.0.0
BUILD_DIR := build
ARCHIVE_NAME := $(PROJECT_NAME)-$(VERSION)-linux-amd64.tar.gz

all: build

build: build-backend build-frontend archive

build-backend:
	cargo build --release
	mkdir -p $(BUILD_DIR)
	cp target/release/ruspass $(BUILD_DIR)/ruspass_backend

build-frontend:
	cd UI_UX && go build -o ../$(BUILD_DIR)/ruspass_frontend . || echo "Frontend skipped"

archive: build-backend
	mkdir -p $(BUILD_DIR)
	@if [ -f $(BUILD_DIR)/ruspass_frontend ]; then tar -czvf $(BUILD_DIR)/$(ARCHIVE_NAME) -C $(BUILD_DIR) ruspass_backend ruspass_frontend; else tar -czvf $(BUILD_DIR)/$(ARCHIVE_NAME) -C $(BUILD_DIR) ruspass_backend; fi

run: build-backend
	./target/release/ruspass

clean:
	cargo clean
	rm -rf $(BUILD_DIR)
	rm -f ruspass.db

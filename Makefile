# RusPass Password Manager - Makefile
# Builds the Rust backend and Go frontend for Linux

.PHONY: all build-backend build-frontend clean install-deps run help

# Default target
all: build-backend build-frontend

# Configuration
BACKEND_TARGET := ./target/release/ruspass
FRONTEND_TARGET := ./ruspass_frontend

# Build the Rust backend (release mode for production)
build-backend:
	@echo "=== Building Rust Backend ==="
	cargo build --release
	@echo "Backend built successfully: $(BACKEND_TARGET)"

# Build the Go frontend with Qt bindings
build-frontend:
	@echo "=== Building Go Frontend ==="
	cd UI_UX && go build -o ../$(FRONTEND_TARGET) .
	@echo "Frontend built successfully: $(FRONTEND_TARGET)"

# Install system dependencies (requires sudo)
install-deps:
	@echo "=== Installing System Dependencies ==="
	@echo "This requires sudo privileges. Please enter your password if prompted."
	sudo apt-get update
	sudo apt-get install -y \
		build-essential \
		libqt5x11extras5-dev \
		qt5-qmake \
		qtbase5-dev \
		qtdeclarative5-dev \
		qtmultimedia5-dev \
		qtwebengine5-dev \
		libqt5webkit5-dev \
		libqt5webchannel5-dev \
		libqt5websockets5-dev \
		qttools5-dev-tools \
		libqt5svg5-dev \
		rustc \
		cargo \
		golang-go
	@echo "Dependencies installed successfully."

# Install Go Qt bindings (does not require sudo)
install-go-qt:
	@echo "=== Installing Go Qt Bindings ==="
	go get -d github.com/therecipe/qt/cmd/...
	$$(go env GOPATH)/bin/qtdeploy build desktop
	@echo "Go Qt bindings installed successfully."

# Clean build artifacts
clean:
	@echo "=== Cleaning Build Artifacts ==="
	cargo clean
	rm -f $(FRONTEND_TARGET)
	rm -f ruspass.db
	@echo "Clean complete."

# Run the application (backend demo)
run: build-backend
	@echo "=== Starting RusPass Backend Demo ==="
	$(BACKEND_TARGET)
	@echo "Demo complete."

# Run in development mode (debug builds)
dev:
	@echo "=== Building in Development Mode ==="
	cargo build
	cd UI_UX && go build -o ../ruspass_frontend .
	@echo "Development builds complete."

# Show help
help:
	@echo "RusPass Password Manager - Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all            - Build both backend and frontend (default)"
	@echo "  build-backend  - Build only the Rust backend"
	@echo "  build-frontend - Build only the Go frontend"
	@echo "  install-deps   - Install system dependencies (requires sudo)"
	@echo "  install-go-qt  - Install Go Qt bindings"
	@echo "  clean          - Remove all build artifacts"
	@echo "  run            - Build and run the backend demo"
	@echo "  dev            - Build debug versions for development"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make install-deps     # Install system dependencies once"
	@echo "  2. make install-go-qt    # Setup Go Qt bindings once"
	@echo "  3. make                  # Build everything"
	@echo "  4. make run              # Run the backend demo"
	@echo "  5. ./ruspass_frontend    # Run the frontend (after building)"

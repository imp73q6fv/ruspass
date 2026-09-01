# RusPass Password Manager

A secure password manager for Linux desktop with Rust backend and Go/Qt frontend.

## Features

- **ChaCha20-Poly1305** encryption for maximum security
- **Argon2id** key derivation (winner of the Password Hashing Competition)
- **SQLite** database for local storage
- **Qt-based UI** with Enpass-like dark theme
- Search and filter capabilities
- Category organization

## Architecture

```
RusPass/
├── main.rs              # Rust entry point
├── Cargo.toml           # Rust dependencies
├── go.mod               # Go module file
├── Makefile             # Build system
├── README.md            # This file
├── encryption/
│   ├── crypto.rs        # ChaCha20 & Argon2id implementation
│   └── models.rs        # Data structures
├── database/
│   └── database.rs      # SQLite operations
└── UI_UX/
    ├── main.go          # Qt-based UI
    └── backend.go       # Backend communication layer
```

## Security Model

1. Master password is never stored
2. Argon2id derives encryption key from master password + salt
3. ChaCha20-Poly1305 encrypts all sensitive data
4. Each password entry is encrypted individually
5. Salt and nonce stored in database (not secret)

## Building

### Prerequisites

Install system dependencies:
```bash
make install-deps
```

Install Go Qt bindings:
```bash
make install-go-qt
```

### Build Commands

Build everything:
```bash
make
```

Build only backend:
```bash
make build-backend
```

Build only frontend:
```bash
make build-frontend
```

### Run

Run backend demo:
```bash
make run
```

Run frontend (after building):
```bash
./ruspass_frontend
```

## Usage

1. Launch the application
2. Create a master password (remember it - cannot be recovered!)
3. Add password entries with categories
4. Search and filter your passwords
5. Click to copy passwords to clipboard

## License

MIT License - see LICENSE file for details 

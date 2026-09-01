# RusPass Password Manager

A secure password manager for Linux desktop environments with a modern UI inspired by Enpass.

## Architecture

- **Backend**: Rust (for security-critical operations)
  - Encryption: ChaCha20-Poly1305
  - Key Derivation: Argon2id
  - Database: SQLite
  
- **Frontend**: Go with Qt bindings (therecipe/qt)
  - Modern dark theme UI
  - Enpass-like design
  - Cross-platform desktop application

## Features

- 🔐 Secure password storage with industry-standard encryption
- 🔑 Argon2id key derivation for master password protection
- 📝 Organize passwords by categories
- 🔍 Quick search functionality
- 🎨 Beautiful dark theme UI similar to Enpass
- 💾 Local SQLite database (no cloud dependency)

## Project Structure

```
ruspass/
├── backend/          # Rust backend
│   ├── Cargo.toml
│   └── src/
│       ├── main.rs      # Entry point and CLI
│       ├── crypto.rs    # Encryption/decryption (ChaCha20 + Argon2id)
│       ├── database.rs  # SQLite operations
│       └── models.rs    # Data structures
│
└── frontend/         # Go frontend with Qt
    ├── go.mod
    ├── main.go        # Main UI application
    └── backend.go     # Backend communication layer
```

## Security Model

1. **Master Password**: User provides a master password
2. **Key Derivation**: Argon2id derives a 32-byte key from the master password
3. **Encryption**: All passwords are encrypted with ChaCha20-Poly1305
4. **Storage**: Encrypted data stored in SQLite database
5. **Memory**: Decrypted passwords only exist in memory when needed

## Building

### Prerequisites

**For Backend (Rust):**
```bash
# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Install dependencies
sudo apt-get install build-essential cmake libsqlite3-dev
```

**For Frontend (Go + Qt):**
```bash
# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Qt5 development files
sudo apt-get install qtbase5-dev qttools5-dev qttools5-dev-tools

# Install therecipe/qt bindings
go get -u github.com/therecipe/qt/cmd/...
qtdeploy build desktop
```

### Build Backend

```bash
cd ruspass/backend
cargo build --release
```

### Build Frontend

```bash
cd ruspass/frontend
qtdeploy build desktop
```

## Usage

1. Run the application:
```bash
./frontend/ruspass_frontend
```

2. On first launch:
   - Set your master password
   - This password cannot be recovered if lost!

3. Add password entries:
   - Click "+ Add New" button
   - Fill in title, username, password, URL, category
   - Save the entry

4. Search entries using the search bar

5. Edit or delete entries as needed

## API Reference

### Backend Commands (CLI)

```bash
# Initialize database
./ruspass_backend init <master_password>

# Authenticate
./ruspass_backend auth <master_password>

# List all entries
./ruspass_backend list <token>

# Add entry
./ruspass_backend add '<json_entry>' <token>

# Get entry by ID
./ruspass_backend get <id> <token>

# Update entry
./ruspass_backend update <id> '<json_entry>' <token>

# Delete entry
./ruspass_backend delete <id> <token>

# Search entries
./ruspass_backend search <query> <token>

# Get categories
./ruspass_backend categories <token>
```

## Configuration

Database location: `~/.ruspass/ruspass.db`

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please ensure:
- Code follows existing style
- Tests pass
- Documentation is updated

## Security Considerations

⚠️ **Important Security Notes:**

1. Never share your master password
2. There is no password recovery mechanism
3. The database is encrypted but not immune to attacks on the running system
4. Always lock the application when stepping away from your computer
5. Keep your system and dependencies updated

## Roadmap

- [ ] Biometric authentication support
- [ ] Password generator
- [ ] Password strength analyzer
- [ ] Import/Export functionality
- [ ] Browser extensions
- [ ] TOTP/2FA support
- [ ] Secure notes feature
- [ ] Backup and sync options

#  Advanced Password Manager

A secure, command-line password manager built in Go with advanced cryptographic features, strong password generation, and local encrypted storage.

##  Features

- ** Strong Encryption**: AES-256-GCM encryption with PBKDF2 key derivation
- ** Smart Password Generation**: Configurable length, character sets, and exclusion rules
- ** Password Analysis**: Real-time strength assessment and scoring
- ** Secure Storage**: Local SQLite database with encrypted data
- ** CLI Interface**: Easy-to-use command-line interface
- ** Search & Management**: Find, update, and organize passwords efficiently
- ** Statistics**: Database insights and usage analytics
- **Security**: Constant-time comparison, memory zeroing, and secure random generation

## Quick Start

### Prerequisites
- Go 1.21+ installed
- Windows/Linux/macOS

### Installation

1. **Clone the repository:**
```bash
git clone <repository-url>
cd password-manager
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Build the application:**
```bash
go build -o password-manager cmd/main.go
```

4. **Run the application:**
```bash
./password-manager help
```

##  Usage Examples

### Generate Strong Passwords
```bash
# Generate a 16-character password with all character types
./password-manager generate --length 16 --uppercase --lowercase --numbers --symbols

# Generate a 20-character password with custom settings
./password-manager generate --length 20 --uppercase --lowercase --numbers --no-repeating
```

### Save Passwords
```bash
# Save a Gmail account
./password-manager save gmail --username user@gmail.com --password mypass --url https://gmail.com

# Save a bank account
./password-manager save bank --username john.doe --password secure123 --url https://mybank.com
```

### Retrieve Passwords
```bash
# Get a specific password
./password-manager get gmail

# List all saved passwords
./password-manager list

# Search for passwords
./password-manager search gmail
```

### Password Management
```bash
# Delete a password
./password-manager delete gmail

# Analyze password strength
./password-manager analyze mypassword123

# View database statistics
./password-manager stats
```

##  Project Structure

```
password-manager/
├── cmd/
│   └── main.go              # Main application entry point
├── internal/
│   ├── crypto/
│   │   ├── encryption.go    # Cryptographic functions
│   │   └── encryption_test.go
│   ├── generator/
│   │   ├── password.go      # Password generation logic
│   │   └── password_test.go
│   └── storage/
│       ├── database.go      # Database operations
│       └── database_test.go
├── go.mod                   # Go module definition
├── README.md                # This file
├── INSTALL.md               # Installation guide
├── demo-test.bat            # Windows batch demo script
├── demo-test.ps1            # PowerShell demo script
└── Makefile                 # Build automation
```

##  Security Features

### Encryption
- **AES-256-GCM**: Authenticated encryption for confidentiality and integrity
- **PBKDF2-SHA256**: Key derivation with 100,000 iterations
- **Random Salt & Nonce**: Unique values for each encryption operation
- **Secure Random**: Cryptographically secure random number generation

### Password Security
- **Strong Generation**: Configurable character sets and length
- **Strength Analysis**: Real-time password quality assessment
- **No Repeating**: Option to prevent consecutive character repetition
- **Exclusion Rules**: Custom character exclusion for specific requirements

### Data Protection
- **Local Storage**: Data never leaves your machine
- **Encrypted Database**: All sensitive data is encrypted at rest
- **Memory Zeroing**: Sensitive data cleared from memory after use
- **Constant-Time Comparison**: Prevents timing attacks

##  Testing

### Run All Tests
```bash
go test ./...
```

### Test with Coverage
```bash
go test -cover ./...
```

### Race Detection
```bash
go test -race ./...
```

### Benchmarking
```bash
go test -bench=. ./...
```

##  Development

### Build Commands
```bash
# Build for current platform
go build -o password-manager cmd/main.go

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o password-manager cmd/main.go
GOOS=darwin GOARCH=amd64 go build -o password-manager cmd/main.go
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

### Using Makefile
```bash
# Build and test
make all

# Install globally
make install

# Run demo
make run

# Clean build artifacts
make clean
```

##  Password Generation Options

### Character Sets
- **Lowercase**: a-z (26 characters)
- **Uppercase**: A-Z (26 characters)
- **Numbers**: 0-9 (10 characters)
- **Symbols**: !@#$%^&*()_+-=[]{}|;:,.<>? (32 characters)

### Configuration Options
- **Length**: 8-128 characters
- **Character Sets**: Enable/disable specific sets
- **Exclusion**: Remove specific characters
- **No Repeating**: Prevent consecutive character repetition
- **Custom Rules**: Advanced pattern matching

## Password Strength Analysis

### Scoring System
- **0-20**: Very Weak
- **21-40**: Weak
- **41-60**: Fair
- **61-80**: Good
- **81-100**: Strong

### Factors Considered
- **Length**: Longer passwords score higher
- **Character Variety**: Mixed character sets increase score
- **Uniqueness**: No repeating patterns
- **Complexity**: Special characters and numbers
- **Predictability**: Common patterns reduce score

##  Important Notes

### Master Password
- **Remember your master password!** It cannot be recovered
- **Use a strong, unique password** for maximum security
- **Store it securely** in a separate password manager or safe location

### Data Location
- Database: `~/.password-manager/passwords.db`
- Configuration: `~/.password-manager/`
- **Backup regularly** to prevent data loss

### Security Considerations
- **Local storage only** - data never transmitted
- **Encrypted at rest** - database is encrypted
- **Memory protection** - sensitive data cleared after use
- **No cloud sync** - complete control over your data

##  Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

##  License

This project is licensed under the MIT License - see the LICENSE file for details.

##  Acknowledgments

- Built with Go's excellent cryptographic libraries
- Inspired by modern password management best practices
- Designed for security-conscious developers and users

##  Support

If you encounter issues:
1. Check the installation guide in `INSTALL.md`
2. Run tests to verify functionality
3. Check Go version and dependencies
4. Review error logs for specific issues

---

** Secure your digital life with Advanced Password Manager!**

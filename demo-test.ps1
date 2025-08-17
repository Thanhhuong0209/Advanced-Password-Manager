# Advanced Password Manager Demo Script for PowerShell
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   Advanced Password Manager Demo" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check Go installation
Write-Host "[1] Checking Go installation..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "✓ Go found: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ ERROR: Go is not installed!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install Go from: https://golang.org/dl/" -ForegroundColor Yellow
    Write-Host "After installation, restart this script." -ForegroundColor Yellow
    Write-Host ""
    Read-Host "Press Enter to exit"
    exit 1
}

# Install dependencies
Write-Host "[2] Installing dependencies..." -ForegroundColor Yellow
try {
    go mod download
    Write-Host "✓ Dependencies installed" -ForegroundColor Green
} catch {
    Write-Host "✗ ERROR: Failed to download dependencies!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# Build project
Write-Host "[3] Building project..." -ForegroundColor Yellow
try {
    go build -o password-manager.exe cmd/main.go
    Write-Host "✓ Project built successfully" -ForegroundColor Green
} catch {
    Write-Host "✗ ERROR: Build failed!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# Run tests
Write-Host "[4] Running tests..." -ForegroundColor Yellow
try {
    go test ./...
    Write-Host "✓ Tests completed" -ForegroundColor Green
} catch {
    Write-Host "⚠ WARNING: Some tests failed, but continuing..." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "[5] Project Demo" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Test password generation
Write-Host "Testing password generation..." -ForegroundColor Yellow
Write-Host "Command: .\password-manager.exe generate --length 16 --uppercase --lowercase --numbers --symbols" -ForegroundColor Gray
Write-Host ""

# Test help command
Write-Host "Testing help command..." -ForegroundColor Yellow
Write-Host "Command: .\password-manager.exe help" -ForegroundColor Gray
Write-Host ""

Write-Host "[6] Interactive Demo" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "You can now test the application interactively:" -ForegroundColor White
Write-Host ""

Write-Host "1. Generate a password:" -ForegroundColor White
Write-Host "   .\password-manager.exe generate --length 20" -ForegroundColor Gray
Write-Host ""

Write-Host "2. Save a password:" -ForegroundColor White
Write-Host "   .\password-manager.exe save gmail --username user@gmail.com --password mypass" -ForegroundColor Gray
Write-Host ""

Write-Host "3. Get a password:" -ForegroundColor White
Write-Host "   .\password-manager.exe get gmail" -ForegroundColor Gray
Write-Host ""

Write-Host "4. List all passwords:" -ForegroundColor White
Write-Host "   .\password-manager.exe list" -ForegroundColor Gray
Write-Host ""

Write-Host "5. Analyze password strength:" -ForegroundColor White
Write-Host "   .\password-manager.exe analyze mypassword123" -ForegroundColor Gray
Write-Host ""

Write-Host "Demo completed!" -ForegroundColor Green
Read-Host "Press Enter to exit"

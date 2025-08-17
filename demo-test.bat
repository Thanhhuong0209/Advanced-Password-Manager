@echo off
echo ========================================
echo    Advanced Password Manager Demo
echo ========================================
echo.

echo [1] Checking Go installation...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed!
    echo.
    echo Please install Go from: https://golang.org/dl/
    echo After installation, restart this script.
    echo.
    pause
    exit /b 1
)

echo [2] Installing dependencies...
go mod download
if %errorlevel% neq 0 (
    echo ERROR: Failed to download dependencies!
    pause
    exit /b 1
)

echo [3] Building project...
go build -o password-manager.exe cmd/main.go
if %errorlevel% neq 0 (
    echo ERROR: Build failed!
    pause
    exit /b 1
)

echo [4] Running tests...
go test ./...
if %errorlevel% neq 0 (
    echo WARNING: Some tests failed, but continuing...
)

echo.
echo [5] Project Demo
echo ========================================
echo.

echo Testing password generation...
echo password-manager.exe generate --length 16 --uppercase --lowercase --numbers --symbols
echo.

echo Testing help command...
echo password-manager.exe help
echo.

echo [6] Interactive Demo
echo ========================================
echo.
echo You can now test the application interactively:
echo.
echo 1. Generate a password:
echo    password-manager.exe generate --length 20
echo.
echo 2. Save a password:
echo    password-manager.exe save gmail --username user@gmail.com --password mypass
echo.
echo 3. Get a password:
echo    password-manager.exe get gmail
echo.
echo 4. List all passwords:
echo    password-manager.exe list
echo.
echo 5. Analyze password strength:
echo    password-manager.exe analyze mypassword123
echo.

echo Demo completed! Press any key to exit...
pause >nul

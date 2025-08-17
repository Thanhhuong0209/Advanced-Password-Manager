# 🚀 Hướng dẫn cài đặt và chạy Advanced Password Manager

## 📋 Yêu cầu hệ thống

- **Windows 10/11** hoặc **Windows Server 2016+**
- **Go 1.21+** (khuyến nghị Go 1.21 hoặc mới hơn)
- **Git** (để clone repository)
- **PowerShell 5.0+** hoặc **Command Prompt**

## 🔧 Bước 1: Cài đặt Go

### **Cách 1: Tải từ trang chủ (Khuyến nghị)**
1. Truy cập: https://golang.org/dl/
2. Tải **Go for Windows** (MSI installer)
3. Chạy file `.msi` đã tải
4. Làm theo hướng dẫn installer
5. **Quan trọng**: Chọn "Add to PATH" trong quá trình cài đặt

### **Cách 2: Sử dụng winget (Windows 10/11)**
```powershell
winget install GoLang.Go
```

### **Cách 3: Sử dụng Chocolatey**
```powershell
choco install golang
```

### **Kiểm tra cài đặt**
Mở **PowerShell** hoặc **Command Prompt** mới và chạy:
```bash
go version
```

Kết quả mong đợi:
```
go version go1.21.0 windows/amd64
```

## 📥 Bước 2: Clone hoặc tải project

### **Nếu có Git:**
```bash
git clone <repository-url>
cd password-manager
```

### **Nếu không có Git:**
1. Tải project dưới dạng ZIP
2. Giải nén vào thư mục `password-manager`
3. Mở Command Prompt/PowerShell trong thư mục đó

## 🛠️ Bước 3: Cài đặt dependencies

Trong thư mục project, chạy:
```bash
go mod download
```

## 🔨 Bước 4: Build project

```bash
go build -o password-manager.exe cmd/main.go
```

## ✅ Bước 5: Kiểm tra cài đặt

Chạy lệnh help để kiểm tra:
```bash
.\password-manager.exe help
```

## 🎮 Bước 6: Chạy demo

### **Sử dụng script tự động:**
```bash
# PowerShell
.\demo-test.ps1

# Command Prompt
demo-test.bat
```

### **Chạy thủ công:**

#### **1. Tạo mật khẩu mới:**
```bash
.\password-manager.exe generate --length 20 --uppercase --lowercase --numbers --symbols
```

#### **2. Lưu mật khẩu:**
```bash
.\password-manager.exe save gmail --username user@gmail.com --password mypass --url https://gmail.com
```

#### **3. Xem mật khẩu:**
```bash
.\password-manager.exe get gmail
```

#### **4. Liệt kê tất cả:**
```bash
.\password-manager.exe list
```

#### **5. Phân tích độ mạnh:**
```bash
.\password-manager.exe analyze mypassword123
```

## 🧪 Bước 7: Chạy tests

```bash
# Chạy tất cả tests
go test ./...

# Chạy tests với coverage
go test -cover ./...

# Chạy tests với race detection
go test -race ./...
```

## 🚨 Xử lý lỗi thường gặp

### **Lỗi 1: "go: command not found"**
- **Nguyên nhân**: Go chưa được cài đặt hoặc không có trong PATH
- **Giải pháp**: Cài đặt lại Go và chọn "Add to PATH"

### **Lỗi 2: "module not found"**
- **Nguyên nhân**: Dependencies chưa được tải
- **Giải pháp**: Chạy `go mod download`

### **Lỗi 3: "build failed"**
- **Nguyên nhân**: Code có lỗi syntax hoặc missing dependencies
- **Giải pháp**: Kiểm tra lỗi và chạy `go mod tidy`

### **Lỗi 4: "permission denied"**
- **Nguyên nhân**: Không có quyền ghi file
- **Giải pháp**: Chạy PowerShell/CMD với quyền Administrator

## 📁 Cấu trúc file sau khi build

```
password-manager/
├── password-manager.exe     # File thực thi
├── demo-test.bat           # Script demo cho CMD
├── demo-test.ps1           # Script demo cho PowerShell
├── cmd/                    # Source code
├── internal/               # Internal packages
├── go.mod                  # Go modules
└── README.md               # Documentation
```

## 🔐 Bước 8: Sử dụng lần đầu

### **1. Tạo master password:**
Khi chạy lệnh đầu tiên, hệ thống sẽ yêu cầu nhập **master password**:
- Đây là mật khẩu chính để truy cập tất cả mật khẩu khác
- **Quan trọng**: Hãy nhớ master password này!
- Nếu quên, tất cả dữ liệu sẽ bị mất

### **2. Database location:**
Database sẽ được tạo tại:
```
C:\Users\<username>\.password-manager\passwords.db
```

## 🎯 Các lệnh cơ bản

| Lệnh | Mô tả | Ví dụ |
|------|-------|-------|
| `generate` | Tạo mật khẩu mới | `generate --length 16` |
| `save` | Lưu mật khẩu | `save gmail --username user@gmail.com` |
| `get` | Lấy mật khẩu | `get gmail` |
| `list` | Liệt kê tất cả | `list` |
| `search` | Tìm kiếm | `search gmail` |
| `delete` | Xóa mật khẩu | `delete gmail` |
| `analyze` | Phân tích độ mạnh | `analyze mypass123` |
| `stats` | Thống kê database | `stats` |
| `help` | Hiển thị trợ giúp | `help` |

## 🚀 Nâng cao

### **Build cho nhiều platform:**
```bash
# Windows
go build -o password-manager.exe cmd/main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o password-manager cmd/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o password-manager cmd/main.go
```

### **Install globally:**
```bash
go install ./cmd/main.go
```

## 📞 Hỗ trợ

Nếu gặp vấn đề:
1. Kiểm tra Go version: `go version`
2. Kiểm tra Go modules: `go mod verify`
3. Chạy tests: `go test ./...`
4. Xem logs lỗi chi tiết

## 🎉 Chúc mừng!

Bạn đã cài đặt thành công **Advanced Password Manager**! 

Bây giờ hãy thử các lệnh cơ bản và khám phá các tính năng của project.

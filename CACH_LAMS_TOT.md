# Neo-Janus: Các Thực Hành Tốt Nhất

## 1. 🏗️ Mẫu Kiến Trúc

### 1.1 Mô Hình Lớp
```go
// ❌ Không Tốt - Phụ Thuộc Chặt Chẽ
type APIHandler struct {
    vm vaccine.VaccineManager // Phụ thuộc cụ thể
}

// ✅ Tốt - Phụ Thuộc Chia Sẻ
type APIHandler struct {
    manager Manager // Giao diện
}

type Manager interface {
    ProcessResult(...) (string, error)
}
```

### 1.2 Mô Hình Kho Lưu Trữ
```go
// ✅ Giao diện tách biệt mối quan tâm
type VaccineRepository interface {
    SavePatch(patch *Patch) error
    GetPatch(id string) (*Patch, error)
    ListPatches() ([]*Patch, error)
}

// Triển khai cụ thể
type FileVaccineRepository struct {
    basePath string
}
```

---

## 2. 🔒 An Toàn Luồng

### 2.1 Bảo Vệ Shared State
```go
// ✅ Tốt - Sử dụng RWMutex cho truy cập đọc-ghi
type VaccineManager struct {
    mu       sync.RWMutex
    vaccines map[string]*Vaccine
}

func (vm *VaccineManager) ProcessResult(result *Result) (string, error) {
    vm.mu.Lock()
    defer vm.mu.Unlock()
    // Cập nhật vaccines
}
```

### 2.2 Cân Nhắc Goroutine
```go
// ✅ Tốt - Xử lý Panic trong Goroutine
go func() {
    defer func() {
        if r := recover(); r != nil {
            logger.Error("Goroutine panic: %v", r)
        }
    }()
    // Logic xử lý
}()
```

---

## 3. 📊 Xác Thực Đầu Vào

### 3.1 Kiểm Tra Kích Thước
```go
// ✅ Tốt - Ngăn chặn Từ Chối Dịch Vụ
const MaxAnalyzeSize = 10 * 1024 // 10KB

func validateAnalyzeRequest(req *AnalyzeRequest) error {
    if len(req.Input) > MaxAnalyzeSize {
        return fmt.Errorf("input exceeds max size")
    }
    return nil
}
```

### 3.2 Loại Đầu Vào
```go
// ✅ Tốt - Chỉ chấp nhận nguồn hợp lệ
var ValidSources = map[string]bool{
    "USER":  true,
    "EMAIL": true,
    "FILE":  true,
}

func (req *AnalyzeRequest) IsValid() bool {
    return ValidSources[req.Source]
}
```

---

## 4. 🎯 Xử Lý Lỗi

### 4.1 Ghi Nhật Ký Lỗi Chi Tiết
```go
// ✅ Tốt - Bối cảnh đầy đủ
if err != nil {
    logger.Error("Failed to save vaccine patch",
        "error", err.Error(),
        "attack_id", result.AttackID,
        "timestamp", time.Now())
    return "", fmt.Errorf("save failed: %w", err)
}
```

### 4.2 Phục Hồi Một Cách Kín Đáo
```go
// ✅ Tốt - Thử lại với Backoff
func retryWithBackoff(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        if err := fn(); err == nil {
            return nil
        }
        time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
    }
    return fmt.Errorf("max retries exceeded")
}
```

---

## 5. 🔐 An Ninh

### 5.1 Xác Thực
```go
// ✅ Tốt - Kiểm tra quyền trước xử lý
func (h *APIHandler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
    if !isAuthenticated(r) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    // Xử lý yêu cầu
}
```

### 5.2 Mã Hóa Dữ Liệu Nhạy Cảm
```go
// ✅ Tốt - Mã hóa trước khi lưu
type Vaccine struct {
    ID          string
    EncryptedPatch string // Mã hóa AES-256
}

func encrypt(data []byte) ([]byte, error) {
    // Triển khai AES-256
}
```

---

## 6. 📝 Ghi Nhật Ký

### 6.1 Mức Độ Ghi Nhật Ký
```go
logger.Debug("Chi tiết nội bộ", "key", value)        // Phát triển
logger.Info("Sự kiện quan trọng", "user", userID)    // Hoạt động
logger.Warning("Điều gì đó bất thường", "code", 403) // Cảnh báo
logger.Error("Lỗi gặp phải", "error", err)           // Lỗi
```

### 6.2 Cấu Hình Ghi Nhật Ký
```yaml
logging:
  level: "info"         # debug, info, warning, error
  output_dir: "./data/logs/"
  max_size_mb: 100      # Xoay vòng khi vượt quá
  max_backups: 5        # Giữ 5 tệp cũ
```

---

## 7. 🧪 Kiểm Thử

### 7.1 Phạm Vi Kiểm Thử
```go
// ✅ Tốt - Kiểm thử mặt tích cực và tiêu cực
func TestValidAnalyzeRequest(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"Valid input", "safe input", false},
        {"Empty input", "", true},
        {"Oversized input", strings.Repeat("a", 20*1024), true},
    }
    
    for _, tt := range tests {
        if err := validateAnalyzeRequest(...); (err != nil) != tt.wantErr {
            t.Errorf("%s failed", tt.name)
        }
    }
}
```

### 7.2 Kiểm Thử Đạo Diễn & Cạnh
```go
// ✅ Tốt - Kiểm tra điều kiện biên
func TestConcurrentVaccineAccess(t *testing.T) {
    vm := NewVaccineManager()
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            vm.ProcessResult(&Result{...})
        }(i)
    }
    wg.Wait() // Không hoảng sợ, không deadlock
}
```

---

## 8. 📦 Quản Lý Phụ Thuộc

### 8.1 Thiết Lập Go Modules
```bash
# Khởi tạo mô-đun
go mod init github.com/user/neo-janus

# Thêm phụ thuộc
go get github.com/some/package@v1.0.0

# Dọn dẹp phụ thuộc không sử dụng
go mod tidy

# Xác minh tính toàn vẹn
go mod verify
```

### 8.2 Xây Dựng Pháp Lệnh
```bash
# Xây dựng ổn định (giống GitHub Actions)
go build \
  -mod=readonly \
  -ldflags="-X main.Version=1.0.0" \
  -o bin/server.exe ./cmd/server/
```

---

## 9. 📊 Hiệu Năng

### 9.1 Điểm Chuẩn
```go
// ✅ Tốt - Đo lường hiệu suất thực tế
func BenchmarkAnalyzeRequest(b *testing.B) {
    handler := setupHandler()
    
    for i := 0; i < b.N; i++ {
        req := createTestRequest()
        handler.HandleAnalyze(httptest.NewRecorder(), req)
    }
}
```

### 9.2 Tối Ưu Hóa
```go
// ✅ Tốt - Sử dụng Buffering
type Logger struct {
    writer *bufio.Writer // Giảm I/O 3x
}

func (l *Logger) Flush() error {
    return l.writer.Flush()
}
```

---

## 10. 🚀 Lộ Trình Phát Triển

### 1.0.0 (Hiện Tại)
- ✅ Core Phòng Thủ (Blue Sentinel)
- ✅ Mô Phỏng Tấn Công (Red Agent)
- ✅ Hệ Thống Vắc Xin
- ✅ Kiểm Thử & Tài Liệu

### 1.1.0 (Q2)
- 🔜 Bảng Điều Khiển Trực Quan (Streamlit)
- 🔜 Lưu Trữ Dữ Liệu Vắc Xin (PostgreSQL)
- 🔜 Xác Thực (JWT/OAuth2)
- 🔜 Giám Sát (Prometheus)

### 2.0.0 (Q4)
- 🔜 Học Từ Tấn Công Tích Lũy
- 🔜 Triển Khai Cấp Độ Công Ty
- 🔜 Hỗ Trợ Đa Ngôn Ngữ
- 🔜 Tích Hợp SPA

---

## 📋 Danh Sách Kiểm Tra Chất Lượng

### Trước Khi Đẩy Lên
- [ ] Mọi bài kiểm thử đều đạt (`go test -v ./...`)
- [ ] Không có cảnh báo vet (`go vet ./...`)
- [ ] Định Dạng Mã (`go fmt ./...`)
- [ ] Phạm Vi ≥ 80% (`go tool cover`)
- [ ] Tài Liệu Đầy Đủ (godoc)
- [ ] Không Có Hardcoded Secrets

### Trước Khi Phát Hành
- [ ] Xây Dựng Sạch & Xây Dựng Lại
- [ ] Kiểm Thử Tích Hợp Vượt Qua
- [ ] Kiểm Thử Hiệu Năng Vượt Qua
- [ ] Kiểm Thử Bảo Mật Hoàn Thành
- [ ] Phục Vụ Của Bạn Ổn Định
- [ ] Tệp CHANGELOG Được Cập Nhật

---

## 🔗 Tài Nguyên

- [Hiệu Lệnh Hiệu Quả Go](https://golang.org/doc/effective_go)
- [Hướng Dẫn Phong Cách Go](https://github.com/golang/go/wiki/CodeReviewComments)
- [Vận Hành Production Go](https://www.ardanlabs.com/)
- [RESTful API Best Practices](https://restfulapi.net)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)

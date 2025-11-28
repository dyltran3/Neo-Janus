# Neo-Janus: Tham Chiếu Lệnh Nhanh

## 🚀 Xây Dựng & Chạy

### Phát Triển Cục Bộ
```bash
# Xây dựng backend
cd 3_janus_core
go mod tidy
go build -o bin/server.exe ./cmd/server/

# Chạy máy chủ
./bin/server.exe

# Hoặc chạy trực tiếp (tự động xây dựng lại)
go run ./cmd/server/main.go
```

### Với Make
```bash
# Xây dựng
make build

# Chạy
make run

# Tất cả trong một
make all
```

### Docker
```bash
# Xây dựng ảnh
make docker-build

# Khởi động dịch vụ
make docker-up

# Dừng dịch vụ
make docker-down

# Xem logs
make docker-logs
```

---

## 🧪 Kiểm Thử

### Chạy Tất Cả Bài Kiểm Thử
```bash
cd 3_janus_core

# Chạy tất cả bài kiểm thử
go test -v ./internal/...

# Chạy kiểm thử gói cụ thể
go test -v ./internal/api
go test -v ./internal/vaccine
go test -v ./internal/logger

# Chạy với báo cáo phạm vi
go test -v -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

### Với Make
```bash
make test        # Chạy kiểm thử
make test-cover  # Với báo cáo phạm vi
make bench       # Điểm chuẩn hiệu năng
```

---

## 🔴 Tấn Công Red Team

### Phát Động Chiến Dịch Tấn Công
```bash
cd 2_red_agent

# Tấn công cơ bản (10 payloads)
python auto_attack.py

# Cường độ tùy chỉnh
python auto_attack.py 50
python auto_attack.py 100

# Với make
make attack
make attack-intense
```

---

## 📊 Giám Sát & Gỡ Lỗi

### Kiểm Tra Sức Khỏe
```bash
curl http://localhost:8080/health
```

### Điểm Cuối API
```bash
# Yêu cầu phân tích
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "input": "đầu vào độc hại để kiểm thử",
    "source": "USER"
  }'

# Ví dụ phản hồi
{
  "status": "BLOCKED",
  "risk_score": 0.95,
  "message": "🛡️ NEO-JANUS: Nội dung bị chặn."
}
```

### Xem Logs
```bash
# Logs real-time
tail -f data/logs/core.log

# Với Docker
docker-compose logs -f janus-core

# Với make
make docker-logs
```

### Kiểm Tra Dữ Liệu Vắc Xin
```bash
# Liệt kê các bản vá được tạo
ls -la data/vaccine/

# Xem nội dung bản vá
cat data/vaccine/vaccine_*.json | jq .
```

---

## 🛠️ Quy Trình Phát Triển

### Định Dạng Mã
```bash
cd 3_janus_core

# Định dạng mã Go
go fmt ./...

# Với make
make fmt
```

### Linting
```bash
cd 3_janus_core

# Sử dụng go vet
go vet ./...

# Hoặc với golangci-lint (nếu được cài đặt)
golangci-lint run ./...

# Với make
make lint
```

### Xây Dựng Sạch
```bash
# Xóa tạo phẩm
make clean

# Xây dựng sạch hoàn toàn + xây dựng lại
make clean && make build
```

---

## 📝 Quản Lý Tệp

### Cấu Trúc Dự Án
```
Neo-Janus/
├── 1_blue_sentinel/         # Mô hình AI Phòng Thủ
├── 2_red_agent/             # Mô phỏng Tấn Công
│   └── attack_lib/
│       ├── fuzzer.py        # Trình tạo Payload
│       └── payloads.txt     # Mẫu Tấn Công Cơ Bản
├── 3_janus_core/            # Backend (Go)
│   ├── cmd/server/          # Điểm Vào
│   ├── internal/
│   │   ├── api/             # Xử Lý HTTP
│   │   ├── vaccine/         # Trình Quản Lý Vắc Xin
│   │   └── logger/          # Hệ Thống Ghi Nhật Ký
│   └── go.mod
├── 4_frontend/              # Bảng Điều Khiển (Streamlit)
├── data/
│   ├── logs/                # Nhật Ký Ứng Dụng
│   └── vaccine/             # Bản Vá Vắc Xin
├── config.yaml              # Cấu Hình
├── Makefile                 # Tự Động Hóa Xây Dựng
├── Dockerfile               # Ảnh Container
└── docker-compose.yml       # Thiết Lập Nhiều Container
```

### Tạo Thư Mục
```bash
mkdir -p data/logs data/vaccine
```

---

## 🔧 Cấu Hình

### Chỉnh Sửa config.yaml
```yaml
server:
  port: 8080

vaccine:
  enabled: true
  trigger_count: 5

logging:
  level: "info"
  output_dir: "./data/logs/"
```

### Biến Môi Trường
```bash
# (Hiện đang sử dụng config.yaml, hỗ trợ trong tương lai)
export LOG_LEVEL=debug
export SERVER_PORT=8000
```

---

## 🐛 Khắc Phục Sự Cố

### Cổng Đã Được Sử Dụng
```bash
# Kiểm tra điều gì đang sử dụng cổng 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Quy trình Kill (tắt duyên tình qua SIGTERM)
kill -TERM <PID>

# Hoặc sử dụng cổng khác
sed -i 's/port: 8080/port: 9090/' config.yaml
```

### Xây Dựng Thất Bại
```bash
# Xây dựng sạch lại
cd 3_janus_core
go clean
go mod tidy
go build -o bin/server.exe ./cmd/server/
```

### Kiểm Thử Thất Bại
```bash
# Chạy với đầu ra chi tiết
go test -v -run TestName ./internal/package

# Kiểm tra các điều kiện đua
go test -race ./internal/...

# Xem phạm vi
go tool cover -html=coverage.out
```

### Vấn Đề Docker
```bash
# Xây dựng lại ảnh
docker-compose down
docker system prune -a
make docker-build
make docker-up
```

---

## 📊 Mẹo Tối Ưu Hiệu Năng

### Tối Ưu Xây Dựng
```bash
# Xây dựng phát hành (nhị phân nhỏ hơn)
cd 3_janus_core
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server.exe ./cmd/server/
```

### Hiệu Năng Hồ Sơ
```bash
# Hồ sơ CPU
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./internal/...
go tool pprof cpu.prof

# Hồ sơ Bộ Nhớ
go tool pprof mem.prof
```

---

## 📚 Tệp Tài Liệu

- **README.md** - Tổng Quan Dự Án
- **BEST_PRACTICES.md** - Mẫu Mã & Lộ Trình
- **OPTIMIZATION_SUMMARY.md** - Cải Tiến Gần Đây
- **Makefile** - Lệnh Xây Dựng
- **config.yaml** - Cấu Hình Ứng Dụng

---

## 🔗 Liên Kết Hữu Ích

- Go Chính Thức: https://golang.org
- Các Thực Hành Tốt Nhất Python: https://pep8.org
- REST API: https://restfulapi.net
- Docker: https://docker.com
- Bảo Mật: https://owasp.org

# 🎬 Demo Neo-Janus Sẵn Sàng

## ✅ Trạng Thái Hoàn Thành

| Thành Phần | Chi Tiết | Trạng Thái |
|-----------|---------|----------|
| 🛡️ **Blue Sentinel** | Phòng thủ AI | ✅ Hoạt động |
| 🔴 **Red Agent** | Tấn công fuzzer | ✅ Hoạt động |
| ⚙️ **Backend (Go)** | API Server | ✅ Xây dựng thành công (9.3 MB) |
| 🎨 **Frontend** | Streamlit Dashboard | ✅ Sẵn sàng chạy |
| 🧪 **Kiểm Thử** | 16 Unit Tests | ✅ 100% vượt qua |
| 📖 **Tài Liệu** | 15+ tệp Markdown | ✅ Hoàn chỉnh (Tiếng Việt) |

---

## 🚀 3 Cách Chạy Demo

### **Cách 1: Script PowerShell (DỄ NHẤT)**
```powershell
cd "C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus"
.\run-demo.ps1
```
Chọn menu:
- `1` = Chạy Backend + Frontend
- `2` = Chạy riêng Backend
- `3` = Chạy riêng Frontend
- `4` = Chạy Red Team Attack

---

### **Cách 2: Manual (3 Terminal Tách Biệt)**

**Terminal 1 - Backend:**
```bash
cd "C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus"
.\3_janus_core\bin\server.exe
```
Kết quả: `🚀 Core API server listening on :8080...`

**Terminal 2 - Frontend:**
```bash
cd "C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus\4_frontend"
streamlit run app.py
```
Kết quả: Mở trình duyệt → `http://localhost:8501`

**Terminal 3 - Red Team (Tùy chọn):**
```bash
cd "C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus\2_red_agent"
python auto_attack.py 50
```

---

### **Cách 3: Makefile (Tự động)**
```bash
cd "C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus\3_janus_core"
make build        # Xây dựng
make run          # Chạy server
make test         # Chạy kiểm thử
```

---

## 📊 Demo Flow

```
User Input (Chat/Attack)
        ↓
    [API :8080]
        ↓
  Blue Sentinel AI
        ↓
  Risk Score Calc
        ↓
┌─────────────┬──────────────┐
│   SAFE      │   BLOCKED    │
│ Risk < 0.5  │  Risk > 0.5  │
│   ✅        │    🚫       │
└─────────────┴──────────────┘
                  ↓
         Vaccine Trigger (≥5 blocks)
                  ↓
         Auto Patch Generated
                  ↓
         Saved: data/vaccine/*.json
```

---

## 🎯 Quy Trình Kiểm Thử

### 1️⃣ **Chat - Tin Nhắn An Toàn**
```
Input:  "Xin chào, bạn khỏe không?"
Output: ✅ PASSED - Risk Score: 0.05
```

### 2️⃣ **Chat - Tin Nhắn Độc Hại**
```
Input:  "<script>alert(document.cookies)</script>"
Output: 🚫 BLOCKED - Risk Score: 0.95
```

### 3️⃣ **Red Team - 50 Payloads**
```
[Progress Bar: ████████████████████] 100%
🛡️  Đã chặn: 40
💀 Lọt lưới: 10 (Vaccine Trigger)
```

### 4️⃣ **Vaccine - Xem Bản Vá**
```
ls data/vaccine/
cat data/vaccine/vaccine_001.json
{
  "id": "vaccine_001",
  "pattern": "xss_pattern_...",
  "defense": "input_sanitize",
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

## 📁 Cấu Trúc Thư Mục

```
Neo-Janus/
├── 1_blue_sentinel/          # 🛡️ Mô hình phòng thủ
├── 2_red_agent/              # 🔴 Máy tạo tấn công
├── 3_janus_core/             # ⚙️ Backend Go
│   ├── bin/
│   │   └── server.exe        # ✅ Tệp nhị phân sẵn sàng
│   ├── cmd/server/main.go
│   └── internal/
│       ├── api/
│       ├── vaccine/
│       └── logger/
├── 4_frontend/               # 🎨 Dashboard Streamlit
│   └── app.py                # ✅ UI sẵn sàng
├── data/
│   ├── logs/                 # 📝 Nhật ký ứng dụng
│   └── vaccine/              # 💊 Bản vá được tạo
├── config.yaml               # ⚙️ Cấu hình
├── run-demo.ps1              # 🚀 Script launcher
└── RUN_DEMO.md               # 📖 Hướng dẫn
```

---

## 🔗 Địa Chỉ Truy Cập

| Dịch Vụ | URL | Chức Năng |
|--------|-----|---------|
| **Backend API** | http://localhost:8080 | REST API Server |
| **Frontend UI** | http://localhost:8501 | Streamlit Dashboard |
| **Health Check** | http://localhost:8080/health | Kiểm tra trạng thái |
| **Analyze** | http://localhost:8080/api/analyze | Phân tích input |

---

## 🧪 API Testing (Curl)

### Health Check
```bash
curl http://localhost:8080/health
```

### Safe Input
```bash
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{"input": "Hello world", "source": "USER"}'
```

### Malicious Input
```bash
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{"input": "DROP TABLE users; DELETE FROM *;", "source": "USER"}'
```

---

## 📚 Tài Liệu Tham Khảo

| Tệp | Ngôn Ngữ | Nội Dung |
|-----|---------|---------|
| **README.md** | Tiếng Anh | Tổng quan dự án |
| **HUONG_DAN_LENH.md** | Tiếng Việt | Tham chiếu lệnh |
| **CACH_LAMS_TOT.md** | Tiếng Việt | Thực hành tốt nhất |
| **LO_TRINH.md** | Tiếng Việt | Lộ trình phát triển |
| **TOM_TAT_TOI_UU.md** | Tiếng Việt | Tóm tắt cải tiến |
| **TRANG_THAI_DUA_AN.md** | Tiếng Việt | Trạng thái dự án |
| **RUN_DEMO.md** | Tiếng Việt | Hướng dẫn chạy demo |

---

## ✅ Danh Sách Kiểm Tra Pre-Demo

- [ ] Python 3.8+ được cài đặt
- [ ] Go 1.21+ được cài đặt
- [ ] Cổng 8080 không được sử dụng
- [ ] Cổng 8501 không được sử dụng
- [ ] File `3_janus_core/bin/server.exe` tồn tại
- [ ] File `config.yaml` tồn tại ở thư mục gốc
- [ ] Thư mục `data/` tồn tại hoặc tự tạo

---

## 🎓 Học Từ Demo

### 1. **Kiến Trúc Hệ Thống**
- Go backend with REST API
- Streamlit frontend with real-time updates
- Async vaccine generation

### 2. **An Toàn Luồng**
- Mutex protection trên shared state
- RWMutex cho high concurrency reads
- Goroutine safety dengan recover()

### 3. **Xác Thực Đầu Vào**
- Size limit (10KB max)
- Type validation (enum source)
- Sanitization rules

### 4. **Phân Tích Thực Thời**
- Real-time risk scoring
- Dynamic vaccine triggers
- Persistent patch storage

---

## 💡 Mẹo & Thủ Thuật

**Tip 1:** Mở 3 terminal cạnh nhau để dễ theo dõi
- Terminal 1 (trái): Backend logs
- Terminal 2 (giữa): Frontend logs
- Terminal 3 (phải): Test commands

**Tip 2:** Sử dụng Chrome DevTools (F12) trên dashboard
- Network tab: Xem API calls
- Console: Xem JavaScript errors

**Tip 3:** Xem logs real-time
```bash
tail -f data/logs/core.log   # Linux/Mac
Get-Content -Wait data/logs/core.log  # PowerShell
```

**Tip 4:** Reset dữ liệu
```bash
rm -r data/vaccine/*   # Xóa bản vá cũ
rm data/logs/core.log  # Xóa log cũ
```

---

## 🆘 Khắc Phục Sự Cố

| Vấn Đề | Nguyên Nhân | Giải Pháp |
|-------|-----------|---------|
| Port 8080 đã được sử dụng | Process cũ chưa tắt | `taskkill /PID <PID> /F` |
| config.yaml not found | Chạy từ sai thư mục | `cd Neo-Janus` trước |
| ModuleNotFoundError: streamlit | Chưa cài phụ thuộc | `pip install streamlit` |
| Connection refused | Backend chưa khởi động | Chạy Terminal 1 trước |
| Certificate error | Python SSL issue | `pip install --upgrade certifi` |

---

## 🎉 Kết Quả Mong Đợi

Khi chạy thành công:

```
✅ Backend Server Online
   🟢 Janus Core: Online
   🚀 Listening on :8080

✅ Frontend Dashboard Open
   🎨 Chat UI: Ready
   🔴 Red Team Panel: Ready

✅ System Status
   📊 Health: OK
   💊 Vaccine: Enabled
   📝 Logging: Active
```

---

**Chúc mừng! Bạn đã sẵn sàng để khám phá Neo-Janus! 🚀**

*Mở trình duyệt và truy cập: **http://localhost:8501***

---

**Tạo tại:** 2024-01-15
**Phiên bản:** v1.0.0
**Trạng thái:** ✅ READY TO DEMO

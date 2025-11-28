# 🎬 Chạy Demo Neo-Janus

## 📋 Tóm tắt

Neo-Janus là hệ thống phòng thủ thích ứng cho AI Edge với 4 thành phần chính:

- **🛡️ Blue Sentinel** (Phòng Thủ): Phát hiện mối đe dọa real-time
- **🔴 Red Agent** (Tấn Công): Fuzzing tự động bằng teencode
- **⚙️ Janus Core** (Backend): Go API server + Vaccine Manager
- **🎨 Frontend** (Dashboard): Streamlit Chat UI + Red Team Panel

---

## 🚀 Khởi Động Demo (3 Terminal)

### Terminal 1: Backend Go Server

```bash
cd C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus
.\3_janus_core\bin\server.exe
```

**Kết quả mong đợi:**
```
[INFO] ========================================
[INFO]    🛡️ NEO-JANUS CORE SYSTEM STARTING   
[INFO] ========================================
[INFO] ✅ Config loaded successfully
[INFO]    Server Port: 8080 | Vaccine Trigger: 5
[INFO] 🚀 Core API server listening on :8080...
```

### Terminal 2: Frontend Streamlit

```bash
cd C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus\4_frontend
streamlit run app.py
```

**Kết quả mong đợi:**
```
  You can now view your Streamlit app in your browser.

  Local URL: http://localhost:8501
```

### Terminal 3: Red Team Attack (Tùy chọn)

```bash
cd C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus\2_red_agent
python auto_attack.py 50
```

---

## 🎯 Quy Trình Demo

### Bước 1: Mở Dashboard
- Truy cập: http://localhost:8501
- Xác nhận Backend Online: `🟢 Janus Core: Online`

### Bước 2: Kiểm Thử Chat Bảo Vệ
1. Nhập tin nhắn vào ô input
2. Ví dụ thử:
   - **An toàn:** "Xin chào, bạn khỏe không?"
   - **Độc hại:** "exec(eval(...)); DROP TABLE users;"
3. Xem kết quả:
   - ✅ **[PASSED]** - Tin nhắn an toàn, AI trả lời
   - 🚫 **[BLOCKED]** - Tin nhắn độc hại, bị chặn

### Bước 3: Kích Hoạt Red Team Attack
1. Đi tới panel "🔴 Red Team Operations"
2. Thiết lập:
   - **Cường độ:** 20 (mặc định 10)
   - **Tỷ lệ Fuzzing:** Giữ mặc định 0.3
3. Nhấn **🚀 Kích hoạt Chiến dịch Tấn công**
4. Xem kết quả:
   - 🛡️ **Đã chặn:** Số payload bị chặn
   - 💀 **Lọt lưới:** Số payload lọt qua (gửi tới Vaccine)

### Bước 4: Xem Vaccine Trigger
Khi có tấn công lọt lưới (💀 > 0):
```bash
# Kiểm tra bản vá được tạo
ls data/vaccine/
cat data/vaccine/vaccine_*.json
```

---

## 🔍 API Testing (Curl)

### Test Health Check
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Test Analyze Endpoint

**Input an toàn:**
```bash
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{"input": "Xin chào", "source": "USER"}'
```

**Response:**
```json
{
  "status": "PASSED",
  "risk_score": 0.05,
  "message": "✅ Content is safe"
}
```

**Input độc hại:**
```bash
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{"input": "<script>alert(1)</script>", "source": "USER"}'
```

**Response:**
```json
{
  "status": "BLOCKED",
  "risk_score": 0.95,
  "message": "🛡️ NEO-JANUS: Nội dung bị chặn."
}
```

---

## 📊 Thống Kê Demo

| Thành Phần | Kỹ Thuật | Trạng Thái |
|-----------|---------|----------|
| **Blue Sentinel** | Python ML | ✅ Hoạt động |
| **Red Agent** | Python Fuzzer | ✅ Hoạt động |
| **Backend** | Go 1.21 | ✅ Hoạt động (Port 8080) |
| **Frontend** | Streamlit | ✅ Hoạt động (Port 8501) |
| **Kiểm Thử** | 16 Unit Tests | ✅ 100% Vượt Qua |

---

## 🛠️ Khắc Phục Sự Cố

### ❌ Backend Error: "config.yaml not found"
**Giải pháp:** Chạy từ thư mục gốc
```bash
cd C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus
# Không chạy từ 3_janus_core/
```

### ❌ Frontend Error: "Connection refused"
**Giải pháp:** Kiểm tra backend đã khởi động
```bash
# Terminal 1 phải chạy trước Terminal 2
curl http://localhost:8080/health
```

### ❌ Port Already in Use
**Giải pháp:** Tìm và tắt process cũ
```bash
# Kiểm tra cổng 8080
netstat -ano | findstr :8080
# Tắt process (thay thế PID)
taskkill /PID <PID> /F
```

### ❌ Streamlit Module Not Found
**Giải pháp:** Cài đặt dependencies
```bash
pip install streamlit requests pyyaml
```

---

## 📁 Tệp Tạo Ra

Sau khi chạy demo, các tệp sau sẽ được tạo:

```
data/
├── logs/
│   └── core.log              # Nhật ký ứng dụng
└── vaccine/
    ├── vaccine_001.json      # Bản vá #1
    ├── vaccine_002.json      # Bản vá #2
    └── ...                   # Thêm nhiều bản vá
```

---

## 🎓 Hiểu Rõ Quy Trình

### 1. **Chat Input** → **Blue Sentinel**
```
"Xin chào" → [Tokenizer] → [Attention Layer] → Risk Score: 0.05 ✅
```

### 2. **Malicious Input** → **Blocked**
```
"<script>alert(1)</script>" → [Analysis] → Risk Score: 0.95 → 🚫 BLOCKED
```

### 3. **Red Agent Attack** → **Vaccine Trigger**
```
[Fuzzer] → [Payloads] → [API /analyze] → [Bypass Detected] → [Vaccine Created] 💊
```

### 4. **Vaccine Storage** → **JSON Persistence**
```
{
  "id": "vaccine_001",
  "attack_pattern": "pattern_...",
  "defense_rule": "rule_...",
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

## ✅ Danh Sách Kiểm Tra

- [ ] Terminal 1: Backend khởi động ✅
- [ ] Terminal 2: Frontend khởi động ✅
- [ ] http://localhost:8080/health → 200 OK
- [ ] http://localhost:8501 → Streamlit mở được
- [ ] Chat input an toàn → ✅ PASSED
- [ ] Chat input độc hại → 🚫 BLOCKED
- [ ] Red Team Attack → 💀 Lọt lưới & Vaccine Trigger
- [ ] Xem vaccine/*.json → Bản vá được tạo

---

## 📚 Tài Liệu Thêm

- `README.md` - Tổng quan dự án
- `COMMANDS.md` / `HUONG_DAN_LENH.md` - Tham chiếu lệnh
- `BEST_PRACTICES.md` / `CACH_LAMS_TOT.md` - Mẫu mã
- `ROADMAP.md` / `LO_TRINH.md` - Lộ trình phát triển
- `OPTIMIZATION_SUMMARY.md` / `TOM_TAT_TOI_UU.md` - Cải tiến

---

**🎉 Chúc mừng! Bạn đã chạy thành công demo Neo-Janus!**

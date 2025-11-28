package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"neo-janus/internal/logger"
	"neo-janus/internal/vaccine"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// --- Data Structures ---
type AnalyzeRequest struct {
	Input  string `json:"input"`
	Source string `json:"source"` // "USER" hoặc "RED_AGENT"
}

type AnalyzeResponse struct {
	Status    string  `json:"status"` // "BLOCKED" hoặc "PASSED"
	RiskScore float64 `json:"risk_score"`
	Message   string  `json:"message"`
}

// Global state cho package api
var vm *vaccine.VaccineManager
// Ngưỡng chặn giả lập (đọc từ config trong tương lai)
const MOCK_THRESHOLD = 0.85

// InitRoutes thiết lập router và các middleware
func InitRoutes(vaccineTriggerCount int) http.Handler {
	// Khởi tạo Vaccine Manager
	vm = vaccine.NewManager(vaccineTriggerCount)

	r := chi.NewRouter()
	
	// Middleware cơ bản
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	// r.Use(middleware.Logger) // Dùng logger mặc định của chi nếu muốn debug HTTP chi tiết

	// Cấu hình CORS để Frontend (khác port) gọi được API
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8501", "http://127.0.0.1:8501"}, // Cho phép Streamlit
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", handleHealthCheck)
	r.Post("/api/analyze", handleAnalyze)
	
	return r
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Neo-Janus Core is healthy."))
}

// --- MOCK BLUE SENTINEL BRIDGE ---
// QUAN TRỌNG: Đây là hàm giả lập việc Go gọi sang Python AI.
// Trong các tuần tới, hàm này sẽ được thay thế bằng logic gRPC client thật.
// Hiện tại nó dùng logic khớp từ khóa đơn giản để hệ thống chạy được luồng.
func callBlueSentinelMock(input string) (string, float64) {
	// Danh sách từ khóa giả lập (cần đồng bộ sơ bộ với Python mock)
	maliciousKeywords := []string{"hack", "kill", "bom", "tấn công", "phá", "admin", "root"}
	lowerInput := strings.ToLower(input)
	
	rand.Seed(time.Now().UnixNano()) // Khởi tạo seed ngẫu nhiên

	for _, kw := range maliciousKeywords {
		if strings.Contains(lowerInput, kw) {
			// Giả lập AI phát hiện ra với điểm số cao > ngưỡng
			score := MOCK_THRESHOLD + rand.Float64()*(1.0-MOCK_THRESHOLD)
			return "BLOCKED", score
		}
	}
	// Giả lập an toàn với điểm số thấp
	score := rand.Float64() * (MOCK_THRESHOLD - 0.1)
	return "PASSED", score
}
// ----------------------------------

// handleAnalyze là controller chính xử lý yêu cầu phân tích
func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	var req AnalyzeRequest
	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		logger.Error("Failed to decode request: %v", err)
		return
	}

	// 1. Gọi Blue Sentinel (Mock) để phân tích
	status, score := callBlueSentinelMock(req.Input)
	logger.Info("[%s] Analyzed input. Status: %s (Score: %.2f)", req.Source, status, score)

	// 2. Gửi kết quả cho Vaccine Manager để kiểm tra lỗ hổng (Vòng lặp phản hồi)
	// Đây là nơi chiến lược "Digital Vaccine" được thực thi.
	vm.ProcessResult(req.Input, req.Source, status)

	// 3. Chuẩn bị phản hồi cho client
	resp := AnalyzeResponse{
		Status:    status,
		RiskScore: score,
	}
	if status == "BLOCKED" {
		resp.Message = "🛡️ Guard: Yêu cầu bị chặn do vi phạm chính sách an toàn."
	} else {
		resp.Message = "✅ Guard: Yêu cầu hợp lệ."
	}

	// Gửi phản hồi JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
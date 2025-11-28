package vaccine

import (
	"time"

	"neo-janus/internal/logger"
)

// VaccineManager quản lý trạng thái các lỗ hổng phát hiện được
type VaccineManager struct {
	failedInputs []string // Bộ đệm chứa các câu tấn công lọt lưới
	triggerCount int      // Ngưỡng kích hoạt tạo vaccine
}

// NewManager khởi tạo trình quản lý vaccine
func NewManager(triggerCount int) *VaccineManager {
	return &VaccineManager{
		triggerCount: triggerCount,
		failedInputs: make([]string, 0),
	}
}

// ProcessResult phân tích kết quả từ Blue Sentinel để tìm lỗ hổng zero-day
func (vm *VaccineManager) ProcessResult(input string, source string, status string) {
	// LOGIC CỐT LÕI CỦA DIGITAL VACCINE:
	// Chỉ quan tâm khi kẻ tấn công (RED_AGENT) thành công (PASSED)
	// Điều này có nghĩa là hệ thống phòng thủ đã thất bại.
	if source == "RED_AGENT" && status == "PASSED" {
		logger.Info("💉 VACCINE TRIGGER: Detected successful attack entry! Input snippet: '%s...'", truncate(input, 30))
		vm.failedInputs = append(vm.failedInputs, input)

		// Kiểm tra ngưỡng kích hoạt
		if len(vm.failedInputs) >= vm.triggerCount {
			vm.deployVaccinePatch()
		}
	}
}

// deployVaccinePatch giả lập quy trình tạo và triển khai bản vá
func (vm *VaccineManager) deployVaccinePatch() {
	logger.Info("🧬 Digital Vaccine Protocol Initiated. Processing %d failed inputs...", len(vm.failedInputs))
	
	// --- PLACEHOLDER LOGIC ---
	// Trong thực tế, tại đây sẽ:
	// 1. Lưu vm.failedInputs xuống file JSON trong thư mục data/vaccine/
	// 2. Gọi một script Python bên ngoài để thực hiện LoRA Fine-tuning nhanh.
	// 3. Thông báo reload lại model (nếu cần).
	
	// Giả lập thời gian xử lý
	time.Sleep(time.Millisecond * 500) 
	
	// Reset bộ đệm sau khi đã xử lý
	vm.failedInputs = make([]string, 0)
	logger.Info("✅ Vaccine Patch simulation complete. System defense updated.")
}

// Hàm phụ trợ cắt ngắn chuỗi để log
func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
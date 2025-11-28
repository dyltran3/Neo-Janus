// Trong hàm handleAnalyze (giữ nguyên các phần import)

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
    // ... (Phần decode request giữ nguyên)

    // 1. Gọi Blue Sentinel (Mock)
    status, score := callBlueSentinelMock(req.Input)
    
    // Log gọn hơn
    logger.Info("[%s] Input: %.20s... | Score: %.2f | Status: %s", req.Source, req.Input, score, status)

    // 2. [TỐI ƯU] Xử lý Vaccine bất đồng bộ (Non-blocking)
    // Giúp API phản hồi ngay lập tức mà không cần chờ logic vaccine chạy xong
    go func(in, src, stat string) {
        // Recover để tránh crash server nếu vaccine panic
        defer func() {
            if r := recover(); r != nil {
                logger.Error("Vaccine Routine Panic: %v", r)
            }
        }()
        vm.ProcessResult(in, src, stat)
    }(req.Input, req.Source, status)

    // 3. Phản hồi JSON ngay lập tức
    resp := AnalyzeResponse{Status: status, RiskScore: score}
    if status == "BLOCKED" {
        resp.Message = "🛡️ NEO-JANUS: Nội dung bị chặn."
    } else {
        resp.Message = "✅ Nội dung hợp lệ."
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
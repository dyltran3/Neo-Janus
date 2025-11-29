// Thay thế hàm callBlueSentinelMock bằng hàm này
func callBlueSentinelReal(input string) (string, float64) {
    // 1. Tạo JSON Payload
    requestBody, _ := json.Marshal(map[string]string{
        "text": input,
    })

    // 2. Gọi sang Python Service (đang chạy ở port 8001)
    resp, err := httpClient.Post("http://localhost:8001/predict", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        logger.Error("Failed to call AI Service: %v", err)
        return "ERROR", 0.0
    }
    defer resp.Body.Close()

    // 3. Đọc kết quả
    var result struct {
        RiskScore float64 `json:"risk_score"`
        IsBlocked bool    `json:"is_blocked"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "ERROR", 0.0
    }

    status := "PASSED"
    if result.IsBlocked {
        status = "BLOCKED"
    }
    return status, result.RiskScore
}
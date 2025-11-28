import random
import yaml
import os
import time

# Load config để lấy ngưỡng chặn
# Sử dụng đường dẫn tương đối để tìm file config từ thư mục hiện tại
base_dir = os.path.dirname(os.path.abspath(__file__))
config_path = os.path.join(base_dir, "../../../config.yaml")

try:
    with open(config_path, 'r') as file:
        config = yaml.safe_load(file)
    THRESHOLD = config['blue_sentinel']['blocking_threshold']
except Exception as e:
    print(f"[Warning] Could not load config for inference mock, using default threshold 0.85. Error: {e}")
    THRESHOLD = 0.85

class MockSentinelPredictor:
    """
    MOCK PREDICTOR: Giả lập bộ phận suy luận AI.
    Dùng để test tích hợp hệ thống khi chưa có file GGUF thật.
    """
    def __init__(self):
        print("[BlueSentinel] Initializing MOCK Predictor (CPU Mode)...")
        # Danh sách từ khóa để giả lập phát hiện tấn công
        self.malicious_keywords = ["hack", "kill", "bom", "tấn công", "phá", "admin", "root"]
        time.sleep(1) # Giả lập thời gian load model
        print("[BlueSentinel] MOCK Model loaded ready.")

    def predict(self, text: str) -> dict:
        """
        Phân tích văn bản và trả về kết quả giả lập.
        """
        text_lower = text.lower()
        risk_score = 0.0
        
        # Logic giả lập: Nếu chứa từ khóa thì trả điểm rủi ro cao ngẫu nhiên
        found_kw = False
        for kw in self.malicious_keywords:
            if kw in text_lower:
                risk_score = random.uniform(THRESHOLD, 0.99)
                found_kw = True
                break
        
        if not found_kw:
            # Nếu không có từ khóa, trả điểm thấp ngẫu nhiên
            risk_score = random.uniform(0.01, THRESHOLD - 0.1)

        is_blocked = risk_score >= THRESHOLD

        result = {
            # Chỉ lấy 50 ký tự đầu để log cho gọn
            "text_snippet": text[:50] + "..." if len(text) > 50 else text,
            "risk_score": round(risk_score, 4),
            "is_blocked": is_blocked,
            "status": "BLOCKED" if is_blocked else "PASSED",
            "engine": "MockCPU-v1 (Placeholder)"
        }
        # In log ra console để debug
        # print(f"[BlueSentinel-Mock] Analyzed: Score={result['risk_score']} -> Status={result['status']}")
        return result

# Singleton instance để các module khác import và sử dụng
# Khi có model thật, file này sẽ được thay thế bằng code dùng llama-cpp-python
predictor = MockSentinelPredictor()
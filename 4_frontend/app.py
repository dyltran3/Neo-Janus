import streamlit as st
import requests
import yaml
import os
import time
import random
import subprocess
import sys
from typing import TypedDict, List, Dict, Any

# Định nghĩa kiểu Message cho lịch sử chat
class Message(TypedDict, total=False):
    role: str
    content: str
    blocked: bool
    score: float

# --- CẤU HÌNH TRANG ---
st.set_page_config(
    page_title="Neo-Janus Dashboard",
    page_icon="🛡️",
    layout="wide",
    initial_sidebar_state="expanded"
)

# --- HÀM HỖ TRỢ ---
def load_config():
    """Đọc cấu hình từ file yaml gốc."""
    # Sử dụng đường dẫn tương đối để tìm file config
    base_dir = os.path.dirname(os.path.abspath(__file__))
    config_path = os.path.join(base_dir, "../config.yaml")
    try:
        with open(config_path, 'r') as file:
            return yaml.safe_load(file)
    except Exception as e:
        st.error(f"Không thể đọc file config.yaml: {e}")
        st.stop()

def check_core_health(api_url):
    """Kiểm tra xem Janus Core có đang chạy không."""
    try:
        # Gọi endpoint health check đơn giản
        health_url = api_url.replace("/api/analyze", "/health")
        requests.get(health_url, timeout=1)
        return True
    except requests.exceptions.RequestException:
        return False

# --- KHỞI TẠO ---
config = load_config()
# Xây dựng URL API từ cấu hình
CORE_API_URL = f"http://localhost:{config['server']['port']}/api/analyze"
core_alive = check_core_health(CORE_API_URL)

# --- SIDEBAR (Thanh bên) ---
with st.sidebar:
    st.image("https://img.icons8.com/fluency/96/security-checked.png", width=64)
    st.title("Neo-Janus")
    st.markdown("Edge AI Security System")
    st.divider()
    
    # Trạng thái hệ thống
    st.subheader("System Status")
    if core_alive:
        st.success("🟢 Janus Core: Online")
    else:
        st.error("🔴 Janus Core: Offline")
        st.warning("Vui lòng chạy backend Go trước!")

    st.divider()
    st.caption(f"Environment: `{config['environment']}`")
    st.caption("Version: 0.1.0-alpha")

# --- GIAO DIỆN CHÍNH ---
if not core_alive:
    st.warning("⚠️ Hệ thống chưa sẵn sàng. Vui lòng khởi động Janus Core Backend và tải lại trang.")
    st.stop()

# Chia layout làm 2 cột chính
col_chat, col_red = st.columns([6, 4], gap="large")

# === CỘT 1: CHAT INTERFACE ===
with col_chat:
    st.subheader("💬 Chat với AI được bảo vệ")

    # Quản lý lịch sử chat trong session state
    if "messages" not in st.session_state:
        # Tin nhắn chào mừng ban đầu
        st.session_state.messages = [
            {"role": "assistant", "content": "Xin chào! Tôi là AI Assistant được bảo vệ bởi Neo-Janus. Mọi tin nhắn của bạn sẽ được quét để đảm bảo an toàn.", "blocked": False, "score": 0.0}
        ]

    # Hiển thị lịch sử chat
    for message in st.session_state.messages:
        with st.chat_message(message["role"]):
            # Nếu tin nhắn bị chặn, hiển thị kiểu lỗi màu đỏ
            if message.get("blocked"):
                st.error(message["content"], icon="🚫")
                with st.expander("Chi tiết kỹ thuật (Debug)"):
                    st.write(f"Risk Score: {message.get('score', 0.0):.4f}")
            else:
                st.markdown(message["content"])

    # Ô nhập liệu của người dùng
    if prompt := st.chat_input("Nhập tin nhắn tại đây..."):
        # 1. Hiển thị tin nhắn người dùng ngay lập tức
        st.session_state.messages.append({"role": "user", "content": prompt})
        with st.chat_message("user"):
            st.markdown(prompt)

        # 2. Gọi API Janus Core để phân tích (Hiển thị spinner khi chờ)
        with st.chat_message("assistant"):
            with st.spinner("🛡️ Blue Sentinel đang quét tin nhắn..."):
                try:
                    # Gửi request đến Backend Go
                    response = requests.post(
                        CORE_API_URL, 
                        json={"input": prompt, "source": "USER"}, 
                        timeout=3
                    )
                    result = response.json()
                    
                    if result["status"] == "BLOCKED":
                        # Trường hợp bị chặn
                        msg_content = result["message"]
                        st.error(msg_content, icon="🚫")
                        st.session_state.messages.append({
                            "role": "assistant", 
                            "content": msg_content, 
                            "blocked": True, 
                            "score": result['risk_score']
                        })
                    else:
                        # Trường hợp an toàn (Giả lập AI trả lời lại nội dung người dùng)
                        # Trong thực tế, ở đây sẽ gọi đến Chatbot thật.
                        ai_reply = f"✅ Tôi đã nhận được nội dung an toàn: '{prompt}'"
                        st.markdown(ai_reply)
                        st.session_state.messages.append({
                            "role": "assistant", 
                            "content": ai_reply, 
                            "blocked": False
                        })

                except Exception as e:
                    st.error(f"Lỗi kết nối đến Janus Core API: {e}")

# === CỘT 2: RED TEAM OPS (Giao diện điều khiển tấn công) ===
with col_red:
    st.subheader("🔴 Red Team Operations")
    st.caption("Bảng điều khiển tấn công giả lập để kiểm thử và kích hoạt Vaccine.")
    
    with st.container(border=True):
        st.write("**Cấu hình Chiến dịch**")
        target_url_display = config['red_agent']['target_url']
        st.code(f"Target: {target_url_display}", language="bash")
        
        intensity = st.slider("Cường độ (Số lượng Payloads)", min_value=5, max_value=100, value=10, step=5)
        mutation_rate = st.slider("Tỷ lệ Fuzzing (Mutation Rate)", 0.0, 1.0, config['red_agent']['fuzzing']['mutation_rate'])
        
        start_btn = st.button("🚀 Kích hoạt Chiến dịch Tấn công", type="primary", use_container_width=True)

    if start_btn:
        st.write("---")
        st.write(f"🔥 Đang khởi động chiến dịch với {intensity} payloads...")
        
        # Khu vực hiển thị log real-time trên UI
        log_area = st.empty()
        progress_bar = st.progress(0)
        
        logs = []
        success_count = 0
        blocked_count = 0
        
        for i in range(intensity):
            # Giả lập độ trễ và kết quả ngẫu nhiên cho demo UI
            time.sleep(random.uniform(0.05, 0.2))
            
            # Giả lập logic: 80% bị chặn, 20% lọt lưới
            is_passed = random.random() > 0.8
            payload_mock = f"Fuzzed payload #{i+1} (MOCK DATA)..."
            
            if is_passed:
                log_line = f"❌ [PASSED] Lỗ hổng phát hiện: {payload_mock}"
                success_count += 1
            else:
                log_line = f"🛡️ [BLOCKED] Đã chặn: {payload_mock}"
                blocked_count += 1
                
            logs.append(log_line)
            # Chỉ hiện 5 dòng log cuối cùng
            log_area.code("\n".join(logs[-5:]), language="text")
            progress_bar.progress((i + 1) / intensity)
            
        progress_bar.empty()
        
        # Hiển thị tổng kết chiến dịch
        st.success("Chiến dịch hoàn tất!")
        col_res1, col_res2 = st.columns(2)
        col_res1.metric("🛡️ Đã chặn", blocked_count)
        col_res2.metric("💀 Lọt lưới (Vaccine Trigger)", success_count, delta_color="inverse")
        
        if success_count > 0:
            st.info("💡 Các trường hợp lọt lưới đã được gửi đến cơ chế Digital Vaccine để xử lý.")
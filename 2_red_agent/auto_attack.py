import requests
import time
import random
import os
import yaml
# Sử dụng thư viện rich để in output đẹp mắt trên terminal
from rich.console import Console
from rich.progress import track
from rich.panel import Panel
# Import module fuzzer cục bộ
from attack_lib.fuzzer import PromptFuzzer

console = Console()

# Load Config để lấy URL mục tiêu
base_dir = os.path.dirname(os.path.abspath(__file__))
config_path = os.path.join(base_dir, "../config.yaml")
try:
    with open(config_path, 'r') as file:
        config = yaml.safe_load(file)
    TARGET_URL = config['red_agent']['target_url']
except Exception as e:
    console.print(f"[red]Failed to load config: {e}[/red]")
    exit(1)

PAYLOADS_FILE = os.path.join(base_dir, "attack_lib/payloads.txt")

def load_base_payloads():
    try:
        with open(PAYLOADS_FILE, 'r', encoding='utf-8') as f:
            return [line.strip() for line in f.readlines() if line.strip()]
    except FileNotFoundError:
        console.print(f"[red]Payload file not found at: {PAYLOADS_FILE}[/red]")
        return []

def run_attack_campaign(intensity: int = 10):
    """Chạy một chiến dịch tấn công giả lập."""
    base_payloads = load_base_payloads()
    if not base_payloads: return

    console.print(Panel.fit(f"🚀 KÍCH HOẠT CHIẾN DỊCH RED TEAM\nTarget: {TARGET_URL}\nIntensity: {intensity} requests", style="bold red"))
    
    fuzzer = PromptFuzzer()
    results = {"blocked": 0, "passed": 0, "errors": 0}

    # Vòng lặp tấn công chính với thanh tiến trình
    for i in track(range(intensity), description="Đang bắn payload..."):
        # 1. Fuzzing: Tạo payload biến thể
        base = random.choice(base_payloads)
        fuzzed_payload = fuzzer.generate_payload(base)
        
        # 2. Gửi request đến Janus Core API
        try:
            # Giả lập Red Agent là một nguồn gửi request
            payload_data = {"input": fuzzed_payload, "source": "RED_AGENT"}
            response = requests.post(TARGET_URL, json=payload_data, timeout=2)
            
            if response.status_code == 200:
                data = response.json()
                if data["status"] == "BLOCKED":
                    results["blocked"] += 1
                    # console.print(f"[green]✅ Blocked:[/green] {fuzzed_payload[:40]}...")
                else:
                    results["passed"] += 1
                    console.print(f"[boldred]❌ LỖ HỔNG (Passed):[/boldred] {fuzzed_payload}")
            else:
                 results["errors"] += 1
                 console.print(f"[yellow]⚠️ Server Error: {response.status_code}[/yellow]")

        except requests.exceptions.RequestException as e:
            results["errors"] += 1
            console.print(f"[red]Connection Error: {e}[/red]")
        
        # Nghỉ ngẫu nhiên nhẹ để tránh quá tải mock server
        time.sleep(random.uniform(0.05, 0.2))

    # In tổng kết
    console.print(Panel(f"""
[bold]KẾT QUẢ CHIẾN DỊCH:[/bold]
🛡️ Bị Blue Sentinel chặn: [green]{results['blocked']}[/green]
💀 Tấn công thành công (Lọt lưới): [red]{results['passed']}[/red] (Trigger Vaccine)
⚠️ Lỗi kỹ thuật: [yellow]{results['errors']}[/yellow]
""", title="Red Team Summary", expand=False))

if __name__ == "__main__":
    # Chạy trực tiếp file này để test nhanh 10 payload
    run_attack_campaign(intensity=10)
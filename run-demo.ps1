#!/usr/bin/env pwsh
# Neo-Janus Demo Launcher
# Script để khởi động toàn bộ hệ thống demo

$ProjectRoot = "C:\Users\TuanAnh\OneDrive\Documents\GitHub\Neo-Janus"
$BackendPath = "$ProjectRoot\3_janus_core"
$FrontendPath = "$ProjectRoot\4_frontend"
$RedAgentPath = "$ProjectRoot\2_red_agent"

Write-Host "🛡️  NEO-JANUS DEMO LAUNCHER" -ForegroundColor Cyan
Write-Host "=" * 50 -ForegroundColor Cyan
Write-Host ""

# Kiểm tra Backend Binary
if (-not (Test-Path "$BackendPath\bin\server.exe")) {
    Write-Host "⚠️  Backend executable not found. Building..." -ForegroundColor Yellow
    Push-Location $BackendPath
    go build -o bin/server.exe ./cmd/server/
    Pop-Location
    Write-Host "✅ Backend built successfully" -ForegroundColor Green
}

# Menu
Write-Host "Chọn chế độ chạy:" -ForegroundColor Yellow
Write-Host "1️⃣  Chạy toàn bộ (Backend + Frontend)" -ForegroundColor White
Write-Host "2️⃣  Chạy riêng Backend" -ForegroundColor White
Write-Host "3️⃣  Chạy riêng Frontend" -ForegroundColor White
Write-Host "4️⃣  Chạy Red Team Attack" -ForegroundColor White
Write-Host "5️⃣  Kiểm tra Health Check" -ForegroundColor White
Write-Host "0️⃣  Thoát" -ForegroundColor White
Write-Host ""
$choice = Read-Host "Nhập lựa chọn (0-5)"

switch ($choice) {
    "1" {
        # Backend
        Write-Host ""
        Write-Host "🚀 Khởi động Backend..." -ForegroundColor Green
        Push-Location $ProjectRoot
        $backendProc = Start-Process -FilePath "$BackendPath\bin\server.exe" `
            -NoNewWindow -PassThru
        Start-Sleep -Seconds 2
        Pop-Location
        
        # Frontend
        Write-Host "🚀 Khởi động Frontend..." -ForegroundColor Green
        Push-Location $FrontendPath
        
        # Cài đặt dependencies nếu cần
        python -m pip install -q streamlit requests pyyaml 2>$null
        
        Start-Process -FilePath "python" `
            -ArgumentList "-m streamlit run app.py" `
            -NoNewWindow
        
        Start-Sleep -Seconds 3
        Write-Host ""
        Write-Host "✅ Hệ thống đã khởi động!" -ForegroundColor Green
        Write-Host "   🔌 Backend: http://localhost:8080" -ForegroundColor Cyan
        Write-Host "   🎨 Frontend: http://localhost:8501" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Nhấn Ctrl+C để thoát..." -ForegroundColor Yellow
        Pop-Location
    }
    "2" {
        # Backend only
        Write-Host ""
        Write-Host "🚀 Khởi động Backend..." -ForegroundColor Green
        Write-Host "   Địa chỉ: http://localhost:8080" -ForegroundColor Cyan
        Write-Host ""
        Push-Location $ProjectRoot
        & "$BackendPath\bin\server.exe"
        Pop-Location
    }
    "3" {
        # Frontend only
        Write-Host ""
        Write-Host "🚀 Khởi động Frontend..." -ForegroundColor Green
        Write-Host "   Địa chỉ: http://localhost:8501" -ForegroundColor Cyan
        Write-Host ""
        
        # Cài đặt dependencies
        python -m pip install -q streamlit requests pyyaml 2>$null
        
        Push-Location $FrontendPath
        python -m streamlit run app.py
        Pop-Location
    }
    "4" {
        # Red Team Attack
        Write-Host ""
        Write-Host "🔴 Red Team Attack Simulator" -ForegroundColor Red
        Write-Host ""
        $intensity = Read-Host "Nhập số lượng payloads (mặc định 10)"
        if ([string]::IsNullOrEmpty($intensity)) { $intensity = 10 }
        
        Push-Location $RedAgentPath
        Write-Host ""
        Write-Host "🚀 Khởi động tấn công ($intensity payloads)..." -ForegroundColor Yellow
        python auto_attack.py $intensity
        Pop-Location
    }
    "5" {
        # Health Check
        Write-Host ""
        Write-Host "📊 Kiểm tra Health Check..." -ForegroundColor Cyan
        try {
            $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -TimeoutSec 2
            Write-Host "✅ Backend hoạt động:" -ForegroundColor Green
            $response | ConvertTo-Json | Write-Host
        }
        catch {
            Write-Host "❌ Backend không phản hồi" -ForegroundColor Red
            Write-Host "   Hãy khởi động backend trước (chọn 1 hoặc 2)" -ForegroundColor Yellow
        }
        Write-Host ""
    }
    "0" {
        Write-Host "👋 Tạm biệt!" -ForegroundColor Cyan
        exit
    }
    default {
        Write-Host "❌ Lựa chọn không hợp lệ" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "📚 Xem RUN_DEMO.md để biết thêm chi tiết" -ForegroundColor Yellow

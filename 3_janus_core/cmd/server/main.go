package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"neo-janus/internal/api"
	"neo-janus/internal/logger"

	"gopkg.in/yaml.v3"
)

// Config mapping với file config.yaml
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Vaccine struct {
		TriggerCount int `yaml:"trigger_count"`
	} `yaml:"vaccine"`
}

// loadConfig đọc file YAML từ thư mục gốc
func loadConfig() (*Config, error) {
	// Đường dẫn tương đối từ vị trí chạy binary (gốc dự án) đến file config
	configPath := filepath.Join(".", "config.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file at %s: %w", configPath, err)
	}
	
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config YAML: %w", err)
	}
	return &cfg, nil
}

func main() {
	// 1. Khởi tạo Logger
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Info("========================================")
	logger.Info("   🛡️ NEO-JANUS CORE SYSTEM STARTING   ")
	logger.Info("========================================")

	// 2. Load Cấu hình
	cfg, err := loadConfig()
	if err != nil {
		logger.Error("Config loading failed: %v. Exiting.", err)
		os.Exit(1)
	}
	logger.Info("Config loaded successfully. Server Port: %d, Vaccine Trigger: %d", cfg.Server.Port, cfg.Vaccine.TriggerCount)

	// 3. Thiết lập API Router và truyền các tham số cấu hình cần thiết
	router := api.InitRoutes(cfg.Vaccine.TriggerCount)

	// 4. Chạy HTTP Server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("🚀 Core API server listening on %s...", addr)
	
	// Server sẽ chạy block tại đây cho đến khi bị tắt
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("Server failed to start: %v", err)
		os.Exit(1)
	}
}
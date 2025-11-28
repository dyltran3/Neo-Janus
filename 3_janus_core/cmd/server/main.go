package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"neo-janus/internal/api"
	"neo-janus/internal/logger"

	"gopkg.in/yaml.v3"
)

// Config application configuration
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Vaccine struct {
		TriggerCount int `yaml:"trigger_count"`
	} `yaml:"vaccine"`
}

// loadConfig load configuration from YAML file
func loadConfig() (*Config, error) {
	// Use relative path from project root
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

	// Validate config
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return nil, fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}
	if cfg.Vaccine.TriggerCount <= 0 {
		return nil, fmt.Errorf("vaccine trigger_count must be positive")
	}

	return &cfg, nil
}

func main() {
	// 1. Initialize logging
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Info("========================================")
	logger.Info("   🛡️ NEO-JANUS CORE SYSTEM STARTING   ")
	logger.Info("========================================")

	// 2. Load configuration
	cfg, err := loadConfig()
	if err != nil {
		logger.Error("Config loading failed: %v. Exiting.", err)
		os.Exit(1)
	}
	logger.Info("✅ Config loaded successfully")
	logger.Info("   Server Port: %d | Vaccine Trigger: %d", cfg.Server.Port, cfg.Vaccine.TriggerCount)

	// 3. Initialize API router
	router := api.InitRoutes(cfg.Vaccine.TriggerCount)

	// 4. Setup HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	logger.Info("🚀 Core API server listening on %s...", addr)

	// 5. Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	logger.Info("Received signal: %v. Shutting down gracefully...", sig)

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error: %v", err)
	}

	logger.Info("✅ Server shutdown complete")
}

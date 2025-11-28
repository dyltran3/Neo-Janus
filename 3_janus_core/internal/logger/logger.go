package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile *os.File
	// Sử dụng đường dẫn tương đối từ vị trí chạy binary (thường là gốc dự án)
	logDirPath = "./data/logs"
	logFilePath = filepath.Join(logDirPath, "core.log")
)

// Init khởi tạo hệ thống log
func Init() error {
	// Tạo thư mục log nếu chưa tồn tại
	if err := os.MkdirAll(logDirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	var err error
	// Mở file log ở chế độ ghi thêm (append)
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	return nil
}

// Helper function để ghi log có timestamp
func logMsg(level, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, msg)
	
	// In ra console
	fmt.Print(logLine)
	// Ghi vào file nếu đã khởi tạo
	if logFile != nil {
		logFile.WriteString(logLine)
	}
}

func Info(format string, v ...interface{}) { logMsg("INFO", format, v...) }
func Error(format string, v ...interface{}) { logMsg("ERROR", format, v...) }

// Close đóng file log khi chương trình kết thúc
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
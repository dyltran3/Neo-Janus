package vaccine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"neo-janus/internal/logger"
)

// VaccineManager manage vulnerability detection and vaccine generation
type VaccineManager struct {
	failedInputs []string
	triggerCount int
	mu           sync.RWMutex // Thread-safe access
}

// NewManager create new vaccine manager instance
func NewManager(triggerCount int) *VaccineManager {
	return &VaccineManager{
		triggerCount: triggerCount,
		failedInputs: make([]string, 0, triggerCount),
	}
}

// ProcessResult analyze result from Blue Sentinel
func (vm *VaccineManager) ProcessResult(input, source, status string) {
	// Only care when RED_AGENT succeeds (bypassed defense)
	if source != "RED_AGENT" || status != "PASSED" {
		return
	}

	vm.mu.Lock()
	defer vm.mu.Unlock()

	logger.Info("💉 VACCINE TRIGGER: Detected successful attack! Input: '%s...'", truncate(input, 30))
	vm.failedInputs = append(vm.failedInputs, input)

	// Check if threshold reached
	if len(vm.failedInputs) >= vm.triggerCount {
		go vm.deployVaccinePatch()
	}
}

// deployVaccinePatch simulate vaccine patch deployment
func (vm *VaccineManager) deployVaccinePatch() {
	vm.mu.Lock()
	inputsCopy := make([]string, len(vm.failedInputs))
	copy(inputsCopy, vm.failedInputs)
	vm.failedInputs = vm.failedInputs[:0] // Reset buffer
	vm.mu.Unlock()

	logger.Info("🧬 Digital Vaccine Protocol Initiated. Processing %d failed inputs...", len(inputsCopy))

	// Save failed inputs to file for later analysis
	if err := vm.savePatchData(inputsCopy); err != nil {
		logger.Error("Failed to save patch data: %v", err)
	}

	// Simulate processing time
	time.Sleep(time.Millisecond * 500)

	logger.Info("✅ Vaccine Patch simulation complete. System defense updated.")
}

// savePatchData save failed inputs to JSON file
func (vm *VaccineManager) savePatchData(inputs []string) error {
	vaccineDir := "./data/vaccine"
	if err := os.MkdirAll(vaccineDir, os.ModePerm); err != nil {
		return err
	}

	filename := filepath.Join(vaccineDir, time.Now().Format("vaccine_20060102_150405.json"))
	data := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"inputs":    inputs,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// truncate shorten string for logging
func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

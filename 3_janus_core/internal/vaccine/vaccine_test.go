package vaccine

import (
	"os"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	triggerCount := 5
	vm := NewManager(triggerCount)

	if vm.triggerCount != triggerCount {
		t.Errorf("Expected trigger count %d, got %d", triggerCount, vm.triggerCount)
	}

	if len(vm.failedInputs) != 0 {
		t.Errorf("Expected empty failed inputs, got %d", len(vm.failedInputs))
	}
}

func TestProcessResult_ValidInput(t *testing.T) {
	vm := NewManager(5)

	// Process non-attack input (should be ignored)
	vm.ProcessResult("normal input", "USER", "APPROVED")

	if len(vm.failedInputs) != 0 {
		t.Error("Non-RED_AGENT inputs should be ignored")
	}
}

func TestProcessResult_AttackBypass(t *testing.T) {
	vm := NewManager(5)

	// Process attack that bypassed defense
	vm.ProcessResult("malicious payload", "RED_AGENT", "PASSED")

	if len(vm.failedInputs) != 1 {
		t.Errorf("Expected 1 failed input, got %d", len(vm.failedInputs))
	}

	if vm.failedInputs[0] != "malicious payload" {
		t.Errorf("Expected 'malicious payload', got %s", vm.failedInputs[0])
	}
}

func TestProcessResult_TriggersVaccine(t *testing.T) {
	vm := NewManager(2) // Low threshold for testing

	// Add first bypass
	vm.ProcessResult("payload1", "RED_AGENT", "PASSED")
	if len(vm.failedInputs) != 1 {
		t.Error("First payload should be stored")
	}

	// Add second bypass (should trigger vaccine)
	vm.ProcessResult("payload2", "RED_AGENT", "PASSED")

	// Small delay to allow goroutine to complete
	time.Sleep(1 * time.Second)

	// After vaccine deployed, buffer should be reset
	vm.mu.RLock()
	count := len(vm.failedInputs)
	vm.mu.RUnlock()

	if count != 0 {
		t.Errorf("Expected buffer to be reset after vaccine, got %d items", count)
	}
}

func TestSavePatchData(t *testing.T) {
	// Setup
	vm := NewManager(5)
	inputs := []string{"payload1", "payload2", "payload3"}

	// Clean up test directory
	defer os.RemoveAll("./data/vaccine")

	// Execute
	err := vm.savePatchData(inputs)

	// Assert
	if err != nil {
		t.Fatalf("savePatchData failed: %v", err)
	}

	// Verify directory exists
	vaccineDir := "./data/vaccine"
	if _, err := os.Stat(vaccineDir); os.IsNotExist(err) {
		t.Error("Vaccine directory was not created")
	}

	// Verify file exists
	files, err := os.ReadDir(vaccineDir)
	if err != nil {
		t.Fatalf("Failed to read vaccine directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("No vaccine patch file was created")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		limit    int
		expected string
	}{
		{"hello world", 5, "hello"},
		{"short", 10, "short"},
		{"test", 4, "test"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.limit)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.limit, result, tt.expected)
		}
	}
}

func TestProcessResult_OnlyRedAgent(t *testing.T) {
	vm := NewManager(1)

	// Process attack from other sources (should be ignored)
	vm.ProcessResult("payload", "BLUE_SENTINEL", "PASSED")
	vm.ProcessResult("payload", "SYSTEM", "PASSED")

	if len(vm.failedInputs) != 0 {
		t.Error("Only RED_AGENT sources should trigger vaccine")
	}
}

func TestProcessResult_OnlyPassed(t *testing.T) {
	vm := NewManager(1)

	// Process RED_AGENT with BLOCKED status (should be ignored)
	vm.ProcessResult("payload", "RED_AGENT", "BLOCKED")
	vm.ProcessResult("payload", "RED_AGENT", "APPROVED")

	if len(vm.failedInputs) != 0 {
		t.Error("Only PASSED status should trigger vaccine")
	}
}

func BenchmarkProcessResult(b *testing.B) {
	vm := NewManager(1000) // High threshold so vaccine doesn't trigger

	for i := 0; i < b.N; i++ {
		vm.ProcessResult("test payload", "RED_AGENT", "PASSED")
	}
}

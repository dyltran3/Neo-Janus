package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"neo-janus/internal/vaccine"
)

func TestHandleAnalyze_ValidRequest(t *testing.T) {
	// Setup
	vm := vaccine.NewManager(5)
	handler := NewAPIHandler(vm)

	request := AnalyzeRequest{
		Input:  "test payload",
		Source: "USER",
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Execute
	handler.HandleAnalyze(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp AnalyzeResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Status == "" {
		t.Error("Response status should not be empty")
	}
}

func TestHandleAnalyze_MissingInput(t *testing.T) {
	// Setup
	vm := vaccine.NewManager(5)
	handler := NewAPIHandler(vm)

	request := AnalyzeRequest{
		Input:  "",
		Source: "USER",
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Execute
	handler.HandleAnalyze(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleAnalyze_WrongMethod(t *testing.T) {
	// Setup
	vm := vaccine.NewManager(5)
	handler := NewAPIHandler(vm)

	req := httptest.NewRequest("GET", "/api/analyze", nil)
	w := httptest.NewRecorder()

	// Execute
	handler.HandleAnalyze(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleHealth(t *testing.T) {
	// Setup
	vm := vaccine.NewManager(5)
	handler := NewAPIHandler(vm)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Execute
	handler.HandleHealth(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["status"] != "healthy" {
		t.Error("Expected status 'healthy'")
	}
}

func TestValidateAnalyzeRequest_Valid(t *testing.T) {
	req := &AnalyzeRequest{
		Input:  "test input",
		Source: "USER",
	}

	err := validateAnalyzeRequest(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateAnalyzeRequest_EmptyInput(t *testing.T) {
	req := &AnalyzeRequest{
		Input:  "",
		Source: "USER",
	}

	err := validateAnalyzeRequest(req)
	if err == nil {
		t.Error("Expected error for empty input")
	}
}

func TestValidateAnalyzeRequest_EmptySource(t *testing.T) {
	req := &AnalyzeRequest{
		Input:  "test",
		Source: "",
	}

	err := validateAnalyzeRequest(req)
	if err == nil {
		t.Error("Expected error for empty source")
	}
}

func TestValidateAnalyzeRequest_TooLong(t *testing.T) {
	// Create string longer than 10000 chars
	longInput := string(make([]byte, 10001))
	req := &AnalyzeRequest{
		Input:  longInput,
		Source: "USER",
	}

	err := validateAnalyzeRequest(req)
	if err == nil {
		t.Error("Expected error for input too long")
	}
}

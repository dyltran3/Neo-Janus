package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"neo-janus/internal/logger"
	"neo-janus/internal/vaccine"
)

// AnalyzeRequest define request payload
type AnalyzeRequest struct {
	Input  string `json:"input"`
	Source string `json:"source"`
}

// AnalyzeResponse define response payload
type AnalyzeResponse struct {
	Status    string  `json:"status"`
	RiskScore float64 `json:"risk_score"`
	Message   string  `json:"message"`
}

// ErrorResponse define error response payload
type ErrorResponse struct {
	Error string `json:"error"`
}

// APIHandler hold application dependencies
type APIHandler struct {
	vaccineManager *vaccine.VaccineManager
}

// NewAPIHandler create a new API handler
func NewAPIHandler(vm *vaccine.VaccineManager) *APIHandler {
	return &APIHandler{vaccineManager: vm}
}

// HandleAnalyze process input analysis request
func (h *APIHandler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	// Validate request method
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Limit body size to prevent DoS
	r.Body = http.MaxBytesReader(w, r.Body, 1024*10) // 10KB limit

	// Parse request
	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if err == io.EOF {
			h.sendError(w, http.StatusBadRequest, "Empty request body")
		} else {
			h.sendError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		}
		return
	}
	defer r.Body.Close()

	// Validate input
	if err := validateAnalyzeRequest(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 1. Call Blue Sentinel (Mock)
	status, score := h.analyzeBlueSentinel(req.Input)

	logger.Info("[%s] Input: %.20s... | Score: %.2f | Status: %s", req.Source, req.Input, score, status)

	// 2. Process Vaccine asynchronously (Non-blocking)
	go h.processVaccineAsync(req.Input, req.Source, status)

	// 3. Respond immediately
	resp := AnalyzeResponse{
		Status:    status,
		RiskScore: score,
	}

	if status == "BLOCKED" {
		resp.Message = "🛡️ NEO-JANUS: Content blocked."
	} else {
		resp.Message = "✅ Content passed security check."
	}

	h.sendJSON(w, http.StatusOK, resp)
}

// HandleHealth check if API is alive
func (h *APIHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	response := map[string]string{"status": "healthy"}
	h.sendJSON(w, http.StatusOK, response)
}

// analyzeBlueSentinel mock Blue Sentinel analysis
func (h *APIHandler) analyzeBlueSentinel(input string) (string, float64) {
	// Mock implementation - replace with real ML model call later
	if len(input) < 10 {
		return "BLOCKED", 0.95
	}
	return "APPROVED", 0.2
}

// processVaccineAsync process result through vaccine manager
func (h *APIHandler) processVaccineAsync(input, source, status string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Vaccine routine panic: %v", r)
		}
	}()

	h.vaccineManager.ProcessResult(input, source, status)
}

// validateAnalyzeRequest validate request fields
func validateAnalyzeRequest(req *AnalyzeRequest) error {
	if req.Input == "" {
		return fmt.Errorf("input field is required")
	}
	if req.Source == "" {
		return fmt.Errorf("source field is required")
	}
	if len(req.Input) > 10000 {
		return fmt.Errorf("input too long (max 10000 chars)")
	}
	return nil
}

// sendJSON send JSON response
func (h *APIHandler) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Failed to encode JSON response: %v", err)
	}
}

// sendError send error response
func (h *APIHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	h.sendJSON(w, statusCode, ErrorResponse{Error: message})
}

package api

import (
	"net/http"

	"neo-janus/internal/vaccine"
)

// InitRoutes initialize and return HTTP router
func InitRoutes(triggerCount int) *http.ServeMux {
	// Create vaccine manager
	vm := vaccine.NewManager(triggerCount)

	// Create API handler
	handler := NewAPIHandler(vm)

	// Create router
	router := http.NewServeMux()

	// Register routes
	router.HandleFunc("/api/analyze", handler.HandleAnalyze)
	router.HandleFunc("/health", handler.HandleHealth)

	return router
}

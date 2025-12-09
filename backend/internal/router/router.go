package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRouter configures all HTTP routes for the application.
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Global OPTIONS handler for CORS preflight requests (initially open)
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}).Methods("OPTIONS")

	// Health check endpoint
	router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Server is running",
		})
	}).Methods("GET")

	return router
}

package http

import (
	"net/http"
	"service-test-runner/internal/delivery/http/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, h *handler.Handler) {
	// Enable CORS with more comprehensive settings
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Accept",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"X-CSRF-Token",
			"X-Requested-With",
		}),
		handlers.ExposedHeaders([]string{"Content-Length"}),
		handlers.MaxAge(3600),
	)

	// Apply CORS middleware to all routes
	r.Use(corsMiddleware)

	// Add global OPTIONS handler for preflight requests
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/automation/run", h.RunAutomationHandler).Methods("POST")
	r.HandleFunc("/automation/update-status", h.UpdateStatusHandler).Methods("POST")
	r.HandleFunc("/automation/check-status", h.CheckStatusHandler).Methods("POST")
	r.HandleFunc("/testsuites", h.GetTestSuitesHandler).Methods("GET")
	r.HandleFunc("/testsuite/detail", h.GetTestSuiteDetailHandler).Methods("POST")
	r.HandleFunc("/projects", h.ProjectHandler).Methods("GET")
}

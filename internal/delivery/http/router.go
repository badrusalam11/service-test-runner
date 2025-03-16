package http

import (
	"service-test-runner/internal/delivery/http/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, h *handler.Handler) {
	r.HandleFunc("/automation/run", h.RunAutomationHandler).Methods("POST")
	r.HandleFunc("/automation/update-status", h.UpdateStatusHandler).Methods("POST")
	r.HandleFunc("/testsuites", h.GetTestSuitesHandler).Methods("GET")
	r.HandleFunc("/testsuite/detail", h.GetTestSuiteDetailHandler).Methods("POST")
	r.HandleFunc("/projects", h.ProjectHandler).Methods("GET")
}

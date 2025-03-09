package http

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, h *Handler) {
	r.HandleFunc("/automation/run", h.RunAutomationHandler).Methods("POST")
	r.HandleFunc("/testsuites", h.GetTestSuitesHandler).Methods("GET")
	r.HandleFunc("/testsuite/detail", h.GetTestSuiteDetailHandler).Methods("POST")
}

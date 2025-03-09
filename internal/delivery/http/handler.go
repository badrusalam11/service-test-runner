package http

import (
	"encoding/json"
	"net/http"

	"service-test-runner/internal/usecase/automation"
	"service-test-runner/internal/usecase/testsuite"
)

type StandardResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Handler struct {
	automationUsecase *automation.AutomationUsecase
	testsuiteUsecase  *testsuite.TestSuiteUsecase
}

func NewHandler(automationUsecase *automation.AutomationUsecase, testsuiteUsecase *testsuite.TestSuiteUsecase) *Handler {
	return &Handler{
		automationUsecase: automationUsecase,
		testsuiteUsecase:  testsuiteUsecase,
	}
}

func respondJSON(w http.ResponseWriter, statusCode int, payload StandardResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// RunAutomationHandler handles POST /automation/run.
// Expected payload: {"project": "web1", "testsuite_id": "login", "email": ""}
func (h *Handler) RunAutomationHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Project     string `json:"project"`
		TestSuiteID string `json:"testsuite_id"`
		Email       string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}
	if req.Project == "" || req.TestSuiteID == "" {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "project and testsuite_id are required",
			Data:    nil,
		})
		return
	}
	runResp, err := h.automationUsecase.Run(req.Project, req.TestSuiteID, req.Email)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	respondJSON(w, http.StatusOK, StandardResponse{
		Status:  "success",
		Message: "Selenium test triggered",
		Data:    runResp,
	})
}

// GetTestSuitesHandler handles GET /testsuites?project=web1.
func (h *Handler) GetTestSuitesHandler(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("project")
	if project == "" {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "project query parameter is required",
			Data:    nil,
		})
		return
	}
	suites, err := h.testsuiteUsecase.GetAll(project)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	respondJSON(w, http.StatusOK, StandardResponse{
		Status:  "success",
		Message: "Test Suites",
		Data:    map[string]interface{}{"testsuites": suites},
	})
}

// GetTestSuiteDetailHandler handles POST /testsuite/detail.
// Expected payload: {"project": "web1", "testsuite_name": "regression"}
func (h *Handler) GetTestSuiteDetailHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Project       string `json:"project"`
		TestSuiteName string `json:"testsuite_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}
	if req.Project == "" || req.TestSuiteName == "" {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "project and testsuite_name are required",
			Data:    nil,
		})
		return
	}
	detail, err := h.testsuiteUsecase.GetDetail(req.Project, req.TestSuiteName)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	respondJSON(w, http.StatusOK, StandardResponse{
		Status:  "success",
		Message: "Detail Test Suite",
		Data:    detail,
	})
}

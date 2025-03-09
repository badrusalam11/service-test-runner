package handler

import (
	"encoding/json"
	"net/http"
)

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

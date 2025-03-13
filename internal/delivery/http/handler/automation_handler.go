package handler

import (
	"encoding/json"
	"net/http"
	"service-test-runner/internal/db"
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
	detailResp, err := h.testsuiteUsecase.GetDetail(req.Project, req.TestSuiteID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: "Project or test suite not found",
			Data:    nil,
		})
	}
	// count the steps
	featureData := detailResp.FeatureData
	lenSteps := 0
	for _, feature := range featureData {
		scenarios := feature.Scenarios
		for _, scenario := range scenarios {
			lenSteps += len(scenario.Steps)
		}
	}
	runResp, err := h.automationUsecase.Run(req.Project, req.TestSuiteID, req.Email)
	if err != nil {
		if err.Error() == "your request is queued" {
			qa := &db.TblQueueAutomation{
				Testsuite:  req.TestSuiteID,
				Checkpoint: 1,
				TotalSteps: lenSteps,
				Status:     1,
				IdTest:     "",
				Project:    req.Project,
				// The CreatedAt field will be set automatically in the db.CreateQueueAutomation function.
			}
			db.CreateQueueAutomation(qa)
			respondJSON(w, http.StatusAccepted, StandardResponse{
				Status:  "success",
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	qa := &db.TblQueueAutomation{
		Testsuite:  req.TestSuiteID,
		Checkpoint: 1,
		TotalSteps: lenSteps,
		Status:     2,
		IdTest:     runResp.RunningID,
		Project:    req.Project,
		// The CreatedAt field will be set automatically in the db.CreateQueueAutomation function.
	}
	db.CreateQueueAutomation(qa)
	respondJSON(w, http.StatusOK, StandardResponse{
		Status:  "success",
		Message: "Selenium test triggered",
		Data:    runResp,
	})
}

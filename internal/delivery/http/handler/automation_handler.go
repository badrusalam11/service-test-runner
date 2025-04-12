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

	// Retrieve test suite details (to count the steps).
	detailResp, err := h.testsuiteUsecase.GetDetail(req.Project, req.TestSuiteID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: "Project or test suite not found",
			Data:    nil,
		})
		return
	}

	// Count the total steps.
	lenSteps := 0
	for _, feature := range detailResp.FeatureData {
		for _, scenario := range feature.Scenarios {
			lenSteps += len(scenario.Steps)
		}
	}

	// Trigger the automation run.
	runResp, err := h.automationUsecase.Run(req.Project, req.TestSuiteID, req.Email)
	if err != nil {
		if err.Error() == "your request is queued" {
			// For a queued request, handle DB creation and publish a RabbitMQ message.
			if err := h.automationUsecase.HandleQueuedRequest(req.Project, req.TestSuiteID, lenSteps); err != nil {
				respondJSON(w, http.StatusInternalServerError, StandardResponse{
					Status:  "error",
					Message: err.Error(),
					Data:    nil,
				})
				return
			}
			respondJSON(w, http.StatusAccepted, StandardResponse{
				Status:  "success",
				Message: "Request queued. A RabbitMQ message has been published.",
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

	// If run was successful, create a DB record with status=2 (triggered).
	qa := &db.TblQueueAutomation{
		Testsuite:  req.TestSuiteID,
		Checkpoint: 0,
		TotalSteps: lenSteps,
		Status:     2, // triggered
		IdTest:     runResp.RunningID,
		Project:    req.Project,
	}
	db.CreateQueueAutomation(qa)
	respondJSON(w, http.StatusOK, StandardResponse{
		Status:  "success",
		Message: "Selenium test triggered",
		Data:    runResp,
	})
}

// UpdateStatusHandler handles POST /automation/updatestatus.
// Expected payload: {"id_test": "test123", "checkpoint": 2, "status": 3}
func (h *Handler) UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Define the payload structure.
	var req struct {
		IdTest   string `json:"id_test"`
		StepName string `json:"step_name"`
		Status   int    `json:"status"`
	}

	// Decode the request payload.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}

	// Validate the required field.
	if req.IdTest == "" {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "id_test is required",
			Data:    nil,
		})
		return
	}

	// Call the use case to update the status.
	if err := h.queueAutomationUsecase.UpdateStatus(req.IdTest, req.StepName, req.Status); err != nil {
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Respond with success.
	respondJSON(w, http.StatusOK, StandardResponse{
		Status:  "success",
		Message: "Status updated successfully",
		Data:    nil,
	})
}

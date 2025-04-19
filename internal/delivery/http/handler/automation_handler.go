package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"service-test-runner/internal/db"
	"strconv"
	"strings"
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
// Expected payload: multipart form with:
// - id_test: string
// - step_name: string
// - status: int
// - report_file: optional PDF file
func (h *Handler) UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 10MB max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil && err != http.ErrNotMultipart {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "Failed to parse form data",
			Data:    nil,
		})
		return
	}

	// Get form values
	idTest := r.FormValue("id_test")
	stepName := r.FormValue("step_name")
	status := r.FormValue("status")

	if idTest == "" {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "id_test is required",
			Data:    nil,
		})
		return
	}

	statusInt := 0
	if status != "" {
		var err error
		statusInt, err = strconv.Atoi(status)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, StandardResponse{
				Status:  "error",
				Message: "Invalid status value",
				Data:    nil,
			})
			return
		}
	}

	// Handle file upload if present
	if file, header, err := r.FormFile("report_file"); err == nil && file != nil {
		defer file.Close()

		// Check if it's a PDF
		if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
			respondJSON(w, http.StatusBadRequest, StandardResponse{
				Status:  "error",
				Message: "Only PDF files are allowed",
				Data:    nil,
			})
			return
		}

		// Read file content
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, StandardResponse{
				Status:  "error",
				Message: "Failed to read file",
				Data:    nil,
			})
			return
		}

		// Generate a unique filename using id_test
		objectName := fmt.Sprintf("reports/%s/report.pdf", idTest)

		// Upload to MinIO
		if err := h.minioService.UploadPDF(objectName, fileBytes); err != nil {
			respondJSON(w, http.StatusInternalServerError, StandardResponse{
				Status:  "error",
				Message: "Failed to upload report file",
				Data:    nil,
			})
			return
		}
	}

	// Call the use case to update the status
	if err := h.queueAutomationUsecase.UpdateStatus(idTest, stepName, statusInt); err != nil {
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Respond with success
	respondJSON(w, http.StatusOK, StandardResponse{
		Status:  "success",
		Message: "Status updated successfully",
		Data:    nil,
	})
}

// CheckStatusHandler handles POST /automation/check-status
// Expected payload: {"id_test": "20250315_214839"}
func (h *Handler) CheckStatusHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IdTest string `json:"id_test"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "Invalid request payload",
			Data:    nil,
		})
		return
	}

	if req.IdTest == "" {
		respondJSON(w, http.StatusBadRequest, StandardResponse{
			Status:  "error",
			Message: "id_test is required",
			Data:    nil,
		})
		return
	}

	// Get automation status from use case
	automation, err := h.queueAutomationUsecase.GetByIdTest(req.IdTest)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, StandardResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if automation == nil {
		respondJSON(w, http.StatusNotFound, StandardResponse{
			Status:  "error",
			Message: "Automation not found",
			Data:    nil,
		})
		return
	}
	progress := automation.Checkpoint * 100 / automation.TotalSteps
	if progress > 100 {
		progress = 100
	}
	// Return the automation status
	respondJSON(w, http.StatusOK, StandardResponse{
		Status:  "success",
		Message: "Test Suites",
		Data: map[string]interface{}{
			"id_test":     automation.IdTest,
			"checkpoint":  automation.Checkpoint,
			"status":      automation.Status,
			"step_name":   automation.StepName,
			"total_steps": automation.TotalSteps,
			"progress":    progress,
			"report_file": automation.ReportFile,
		},
	})
}

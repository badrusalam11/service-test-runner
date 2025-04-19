package handler

import (
	"encoding/json"
	"net/http"

	"service-test-runner/internal/infrastructure/storage"
	usecase "service-test-runner/internal/usecase"
)

type StandardResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Handler struct {
	automationUsecase      *usecase.AutomationUsecase
	queueAutomationUsecase *usecase.QueueAutomationUseCase
	testsuiteUsecase       *usecase.TestSuiteUsecase
	projectUsecase         *usecase.ProjectUsecase
	minioService           *storage.MinioService
}

func NewHandler(
	automationUsecase *usecase.AutomationUsecase,
	queueAutomationUsecase *usecase.QueueAutomationUseCase,
	testsuiteUsecase *usecase.TestSuiteUsecase,
	projectUsecase *usecase.ProjectUsecase,
	minioService *storage.MinioService,
) *Handler {
	return &Handler{
		automationUsecase:      automationUsecase,
		queueAutomationUsecase: queueAutomationUsecase,
		testsuiteUsecase:       testsuiteUsecase,
		projectUsecase:         projectUsecase,
		minioService:           minioService,
	}
}

func respondJSON(w http.ResponseWriter, statusCode int, payload StandardResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

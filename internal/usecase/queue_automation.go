package usecase

import (
	"errors"
	"service-test-runner/internal/db"
	automationRepo "service-test-runner/internal/repository/automation"
)

// QueueAutomationUseCase handles business logic for QueueAutomation operations.
type QueueAutomationUseCase struct {
	repo automationRepo.QueueAutomationRepository
}

// NewQueueAutomationUseCase creates a new instance of QueueAutomationUseCase.
func NewQueueAutomationUseCase(repo automationRepo.QueueAutomationRepository) *QueueAutomationUseCase {
	return &QueueAutomationUseCase{repo: repo}
}

// GetByIdTest retrieves automation details by ID
func (uc *QueueAutomationUseCase) GetByIdTest(idTest string) (*db.TblQueueAutomation, error) {
	return uc.repo.GetByIdTest(idTest)
}

// GetByReferenceNumber retrieves automation details by reference number.
func (uc *QueueAutomationUseCase) GetByReferenceNumber(referenceNumber string) (*db.TblQueueAutomation, error) {
	return uc.repo.GetByReferenceNumber(referenceNumber)
}

// UpdateStatusByReferenceNumber updates the status of a record by its reference number.
func (uc *QueueAutomationUseCase) UpdateStatusByReferenceNumber(qa *db.TblQueueAutomation) error {
	// Check if the record exists.
	record, err := uc.repo.GetByReferenceNumber(qa.ReferenceNumber)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.New("record not found")
	}
	// If it exists, update the status.
	return uc.repo.UpdateStatusByReferenceNumber(qa.IdTest, qa.ReferenceNumber, qa.Status)
}

// UpdateStatus checks for record existence before updating status.
func (uc *QueueAutomationUseCase) UpdateStatus(idTest string, stepName string, status int, referenceNumber string) error {
	// Check if the record exists.
	record, err := uc.repo.GetByReferenceNumber(referenceNumber)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.New("record not found")
	}
	newCheckpoint := record.Checkpoint + 1
	if newCheckpoint > record.TotalSteps {
		newCheckpoint = record.TotalSteps
	}

	// If it exists, update the status.
	return uc.repo.UpdateStatus(idTest, stepName, newCheckpoint, status, referenceNumber)
}

// UpdateReportFile updates the report file URL for a given test ID.
func (uc *QueueAutomationUseCase) UpdateReportFile(idTest string, reportFileURL string) error {
	return db.UpdateQueueAutomationReportFile(idTest, reportFileURL)
}

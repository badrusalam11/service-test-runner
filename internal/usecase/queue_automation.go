package usecase

import (
	"errors"
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

// UpdateStatus checks for record existence before updating status.
func (uc *QueueAutomationUseCase) UpdateStatus(idTest string, stepName string, status int) error {
	// Check if the record exists.
	record, err := uc.repo.GetByIdTest(idTest)
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
	return uc.repo.UpdateStatus(idTest, stepName, newCheckpoint, status)
}

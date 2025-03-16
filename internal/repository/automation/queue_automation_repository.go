package automationRepo

import (
	"service-test-runner/internal/db" // adjust the import path accordingly
)

// QueueAutomationRepository defines the repository interface.
type QueueAutomationRepository interface {
	GetByIdTest(idTest string) (*db.TblQueueAutomation, error)
	UpdateStatus(idTest string, stepName string, checkpoint int, status int) error
}

// queueAutomationRepository is the concrete implementation.
type queueAutomationRepository struct{}

// NewQueueAutomationRepository creates a new instance of the repository.
func NewQueueAutomationRepository() QueueAutomationRepository {
	return &queueAutomationRepository{}
}

// GetByIdTest fetches a record by its idTest.
func (r *queueAutomationRepository) GetByIdTest(idTest string) (*db.TblQueueAutomation, error) {
	return db.SelectQueueAutomationByIdTest(idTest)
}

// UpdateStatus updates the record’s checkpoint and status.
func (r *queueAutomationRepository) UpdateStatus(idTest string, stepName string, checkpoint int, status int) error {
	return db.UpdateQueueAutomationStatus(idTest, stepName, checkpoint, status)
}

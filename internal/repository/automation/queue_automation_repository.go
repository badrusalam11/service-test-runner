package automationRepo

import (
	"service-test-runner/internal/db" // adjust the import path accordingly
)

// QueueAutomationRepository defines the repository interface.
type QueueAutomationRepository interface {
	GetByIdTest(idTest string) (*db.TblQueueAutomation, error)
	GetByReferenceNumber(referenceNumber string) (*db.TblQueueAutomation, error)
	UpdateStatus(idTest string, stepName string, checkpoint int, status int, referenceNumber string) error
	UpdateStatusByReferenceNumber(referenceNumber string, stepName string, status int) error
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
func (r *queueAutomationRepository) UpdateStatus(idTest string, stepName string, checkpoint int, status int, referenceNumber string) error {
	return db.UpdateQueueAutomationStatus(idTest, stepName, checkpoint, status, referenceNumber)
}

// GetByReferenceNumber fetches a record by its reference number.
func (r *queueAutomationRepository) GetByReferenceNumber(referenceNumber string) (*db.TblQueueAutomation, error) {
	return db.SelectQueueAutomationByRefnum(referenceNumber)
}

// UpdateStatusByReferenceNumber updates the record’s status by its reference number.
func (r *queueAutomationRepository) UpdateStatusByReferenceNumber(idTest string, referenceNumber string, status int) error {
	return db.UpdateQueueAutomationStatusByReferenceNumber(idTest, referenceNumber, status)
}

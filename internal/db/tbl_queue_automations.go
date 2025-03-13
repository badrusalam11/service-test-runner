package db

import (
	"log"
	"time"
)

// TblQueueAutomation represents a row in the tbl_QueueAutomation table.
type TblQueueAutomation struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Testsuite  string    `gorm:"unique;not null"`
	Checkpoint int       `gorm:"not null"`
	TotalSteps int       `gorm:"not null"`
	Status     int       `gorm:"not null"`
	IdTest     string    `gorm:"null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"` // Automatically set to current time
	Project    string    `gorm:"not null"`
}

// CreateQueueAutomation inserts a new record into tbl_QueueAutomation.
// It sets the CreatedAt field to the current time before inserting.
func CreateQueueAutomation(qa *TblQueueAutomation) error {
	// Set the CreatedAt field to the current timestamp
	qa.CreatedAt = time.Now()

	// Insert the record using GORM's Create method
	result := DB.Create(qa)
	if result.Error != nil {
		log.Printf("Error inserting QueueAutomation record: %v", result.Error)
		return result.Error
	}

	log.Printf("Inserted QueueAutomation record: %+v", qa)
	return nil
}

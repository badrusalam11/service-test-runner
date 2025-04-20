package db

import (
	"log"
	"time"
)

// TblQueueAutomation represents a row in the tbl_QueueAutomation table.
type TblQueueAutomation struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	ReferenceNumber string    `gorm:"unique;not null"`
	Testsuite       string    `gorm:"not null"`
	StepName        string    `gorm:"not null"`
	Checkpoint      int       `gorm:"not null"`
	TotalSteps      int       `gorm:"not null"`
	Status          int       `gorm:"not null"`
	IdTest          string    `gorm:"null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"` // Automatically set to current time
	Project         string    `gorm:"not null"`
	ReportFile      string    `gorm:"null"`
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

// UpdateQueueAutomationStatus updates the checkpoint and status for a record identified by idTest.
func UpdateQueueAutomationStatus(idTest string, stepName string, checkpoint int, status int, referenceNumber string) error {
	result := DB.Model(&TblQueueAutomation{}).
		Where("reference_number = ?", referenceNumber).
		Updates(map[string]interface{}{
			"step_name":  stepName,
			"checkpoint": checkpoint,
			"status":     status,
			"id_test":    idTest,
		})
	return result.Error
}

// UpdateQueueAutomationStatusByReferenceNumber updates the status for a record identified by reference number.
func UpdateQueueAutomationStatusByReferenceNumber(idTest string, referenceNumber string, status int) error {
	result := DB.Model(&TblQueueAutomation{}).
		Where("reference_number = ?", referenceNumber).
		Update("status", status)
	return result.Error
}

// UpdateQueueAutomationReportFile updates the report_file URL for a record identified by idTest.
func UpdateQueueAutomationReportFile(idTest string, reportFileURL string) error {
	result := DB.Model(&TblQueueAutomation{}).
		Where("id_test = ?", idTest).
		Update("report_file", reportFileURL)
	return result.Error
}

// SelectAllQueueAutomation retrieves all QueueAutomation records.
func SelectAllQueueAutomation() ([]TblQueueAutomation, error) {
	var qaList []TblQueueAutomation
	result := DB.Find(&qaList)
	if result.Error != nil {
		log.Printf("Error selecting all QueueAutomation records: %v", result.Error)
		return nil, result.Error
	}
	return qaList, nil
}

// SelectQueueAutomationByIdTest retrieves a single record from tbl_QueueAutomation using the Id_test field.
func SelectQueueAutomationByIdTest(idTest string) (*TblQueueAutomation, error) {
	var qa TblQueueAutomation
	result := DB.Where("id_test = ?", idTest).First(&qa)
	if result.Error != nil {
		log.Printf("Error selecting QueueAutomation record for idTest %s: %v", idTest, result.Error)
		return nil, result.Error
	}
	return &qa, nil
}

func SelectQueueAutomationByRefnum(reference_number string) (*TblQueueAutomation, error) {
	var qa TblQueueAutomation
	result := DB.Where("reference_number = ?", reference_number).First(&qa)
	if result.Error != nil {
		log.Printf("Error selecting QueueAutomation record for idTest %s: %v", reference_number, result.Error)
		return nil, result.Error
	}
	return &qa, nil
}

package usecase

import (
	"encoding/json"
	"fmt"
	"service-test-runner/internal/domain"
	"service-test-runner/internal/infrastructure/messaging"
	"service-test-runner/internal/repository/selenium"
	"time"
)

// AutomationUsecase handles automation logic.
type AutomationUsecase struct {
	repo      *selenium.SeleniumRepository
	publisher messaging.Publisher
}

// NewAutomationUsecase creates a new AutomationUsecase with its dependencies injected.
func NewAutomationUsecase(repo *selenium.SeleniumRepository, publisher messaging.Publisher) *AutomationUsecase {
	return &AutomationUsecase{
		repo:      repo,
		publisher: publisher,
	}
}

// Run triggers the automation using testsuite_id and email.
// If the request is queued, it publishes a RabbitMQ message.
func (a *AutomationUsecase) Run(project, testsuiteID, email string) (domain.RunResponse, error) {
	runResp, err := a.repo.RunAutomation(project, testsuiteID, email)
	if err != nil {
		return runResp, err
	}
	return runResp, nil
}

// HandleQueuedRequest handles a queued automation request by storing a DB record and publishing a RabbitMQ message.
func (uc *AutomationUsecase) HandleQueuedRequest(project, testsuiteID string, totalSteps int) error {
	timestamp := time.Now().Unix()
	// Prepare the message payload.
	msgBody := map[string]interface{}{
		"refnum":       timestamp,
		"project":      project,
		"testsuite_id": testsuiteID,
		"total_steps":  totalSteps,
	}
	msgBytes, err := json.Marshal(msgBody)
	if err != nil {
		return err
	}

	// Publish the message to RabbitMQ.
	fmt.Println("Publishing to RabbitMQ: %s", msgBytes)
	if err := uc.publisher.Publish(msgBytes); err != nil {
		return err
	}
	return nil
}

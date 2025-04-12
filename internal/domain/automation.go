package domain

// AutomationService defines the contract for running automation.
type AutomationService interface {
	RunAutomation(project, testsuiteID, email string) (RunResponse, error)
}

// RunResponse represents the response data for a run.
type RunResponse struct {
	RunningID   string `json:"running_id"`
	TestSuiteID string `json:"testsuite_id"`
}

// QueuedRequest represents the payload for a queued automation request.
type QueuedRequest struct {
	Project     string `json:"project"`
	TestSuiteID string `json:"testsuite_id"`
	Email       string `json:"email"`
}

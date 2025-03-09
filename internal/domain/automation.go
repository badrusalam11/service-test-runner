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

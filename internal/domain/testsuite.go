package domain

// TestSuiteService defines the contract for test suite operations.
type TestSuiteService interface {
	GetTestSuites(project string) ([]string, error)
	GetTestSuiteDetail(project, testsuiteName string) (TestSuiteDetail, error)
}

// TestSuiteDetail represents detailed information about a test suite.
type TestSuiteDetail struct {
	TestSuiteName string        `json:"testsuite_name"`
	FeatureData   []FeatureData `json:"feature_data"`
}

type FeatureData struct {
	Feature   string     `json:"feature"`
	Scenarios []Scenario `json:"scenarios"`
}

type Scenario struct {
	Scenario string   `json:"scenario"`
	Steps    []string `json:"steps"`
	Examples []string `json:"examples"`
	Type     string   `json:"type"`
}

package usecase

import (
	"service-test-runner/internal/domain"
	"service-test-runner/internal/repository/selenium"
)

type AutomationUsecase struct {
	repo *selenium.SeleniumRepository
}

func NewAutomationUsecase(repo *selenium.SeleniumRepository) *AutomationUsecase {
	return &AutomationUsecase{repo: repo}
}

// Run triggers the automation using testsuite_id and email.
func (a *AutomationUsecase) Run(project, testsuiteID, email string) (domain.RunResponse, error) {
	return a.repo.RunAutomation(project, testsuiteID, email)
}

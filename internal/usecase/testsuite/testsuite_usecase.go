package testsuite

import (
	"service-test-runner/internal/domain"
	"service-test-runner/internal/repository/selenium"
)

type TestSuiteUsecase struct {
	repo *selenium.SeleniumRepository
}

func NewTestSuiteUsecase(repo *selenium.SeleniumRepository) *TestSuiteUsecase {
	return &TestSuiteUsecase{repo: repo}
}

func (t *TestSuiteUsecase) GetAll(project string) ([]string, error) {
	return t.repo.GetTestSuites(project)
}

func (t *TestSuiteUsecase) GetDetail(project, testsuiteName string) (domain.TestSuiteDetail, error) {
	return t.repo.GetTestSuiteDetail(project, testsuiteName)
}

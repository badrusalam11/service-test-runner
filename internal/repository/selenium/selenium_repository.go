package selenium

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"service-test-runner/internal/domain"
	repository "service-test-runner/internal/repository"
)

type SeleniumRepository struct {
	projects map[string]string
}

func NewSeleniumRepository(projects map[string]string) *SeleniumRepository {
	return &SeleniumRepository{projects: projects}
}

func (s *SeleniumRepository) getBaseURL(project string) (string, error) {
	url, ok := s.projects[project]
	if !ok {
		return "", errors.New("project not found")
	}
	return url, nil
}

// RunAutomation calls POST /selenium/run with payload {"testsuite_id", "email"}.
func (s *SeleniumRepository) RunAutomation(project, testsuiteID, email string) (domain.RunResponse, error) {
	baseURL, err := s.getBaseURL(project)
	if err != nil {
		return domain.RunResponse{}, err
	}
	endpoint := fmt.Sprintf("%s/selenium/run", baseURL)
	payload := map[string]string{
		"testsuite_id": testsuiteID,
		"email":        email,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return domain.RunResponse{}, err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return domain.RunResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.RunResponse{}, errors.New("failed to run automation")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.RunResponse{}, err
	}

	var runResp domain.RunResponse
	_, err = repository.ParseGeneralResponse(data, &runResp)
	if err != nil {
		return domain.RunResponse{}, err
	}
	return runResp, nil
}

// GetTestSuites calls GET /selenium/testsuites.
func (s *SeleniumRepository) GetTestSuites(project string) ([]string, error) {
	baseURL, err := s.getBaseURL(project)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/selenium/testsuites", baseURL)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get test suites")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tsData struct {
		Testsuites []string `json:"testsuites"`
	}
	_, err = repository.ParseGeneralResponse(data, &tsData)
	if err != nil {
		return nil, err
	}
	return tsData.Testsuites, nil
}

// GetTestSuiteDetail calls POST /selenium/testsuite/detail with payload {"testsuite_name"}.
func (s *SeleniumRepository) GetTestSuiteDetail(project, testsuiteName string) (domain.TestSuiteDetail, error) {
	baseURL, err := s.getBaseURL(project)
	if err != nil {
		return domain.TestSuiteDetail{}, err
	}
	endpoint := fmt.Sprintf("%s/selenium/testsuite/detail", baseURL)
	payload := map[string]string{
		"testsuite_name": testsuiteName,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return domain.TestSuiteDetail{}, err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return domain.TestSuiteDetail{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.TestSuiteDetail{}, errors.New("failed to get test suite detail")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.TestSuiteDetail{}, err
	}

	var detail domain.TestSuiteDetail
	_, err = repository.ParseGeneralResponse(data, &detail)
	if err != nil {
		return domain.TestSuiteDetail{}, err
	}
	return detail, nil
}

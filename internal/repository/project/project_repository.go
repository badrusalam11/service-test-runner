package project

import (
	"service-test-runner/internal/domain"
)

type ProjectRepository struct {
	projects map[string]string
}

func NewProjectRepository(projects map[string]string) *ProjectRepository {
	return &ProjectRepository{projects: projects}
}

// ShowProject converts the projects map into a slice of ProjectResponse.
func (s *ProjectRepository) ShowProject() (domain.ShowProjectResponse, error) {
	var resp domain.ShowProjectResponse
	for name, url := range s.projects {
		resp = append(resp, domain.ProjectResponse{
			Name: name,
			URL:  url,
		})
	}
	return resp, nil
}

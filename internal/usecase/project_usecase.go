package usecase

import (
	"service-test-runner/internal/domain"
	"service-test-runner/internal/repository/project"
)

type ProjectUsecase struct {
	repo *project.ProjectRepository
}

func NewProjectUsecase(repo *project.ProjectRepository) *ProjectUsecase {
	return &ProjectUsecase{repo: repo}
}

func (a *ProjectUsecase) ShowProject() (domain.ShowProjectResponse, error) {
	return a.repo.ShowProject()
}

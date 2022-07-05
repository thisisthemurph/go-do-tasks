package services

import (
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type ProjectService interface {
	GetProjects()				([]*entities.Project, error)
	GetProjectById(userId string)	(*entities.Project, error)
}

type projectService struct{
	query repository.ProjectQuery
}

func NewProjectService(dao repository.DAO) ProjectService {
	return &projectService{
		query: dao.NewProjectQuery(),
	}
}

func (p *projectService) GetProjects() ([]*entities.Project, error) {
	projects, err := p.query.GetAllProjects()
	return projects, err
}

func (p *projectService) GetProjectById(projectId string) (*entities.Project, error) {
	project, err := p.query.GetProjectById(projectId)
	return project, err
}
package services

import (
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type ProjectService interface {
	GetProjects() ([]*entities.Project, error)
	GetProjectById(userId string) (*entities.Project, error)
	CreateProject(newProject *entities.Project) error
	UpdateProject(projectId string, newProjectData *entities.Project) error
}

type projectService struct {
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

func (p *projectService) CreateProject(newProject *entities.Project) error {
	err := p.query.CreateProject(newProject)
	return err
}

func (p *projectService) UpdateProject(projectId string, newProjectData *entities.Project) error {
	err := p.query.UpdateProject(projectId, newProjectData)
	return err
}

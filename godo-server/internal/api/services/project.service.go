package services

import (
	"errors"
	"godo/internal/api"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type ProjectService interface {
	GetProjects() ([]*entities.Project, error)
	GetProjectById(userId string) (*entities.Project, error)
	CreateProject(newProject *entities.Project) error
	UpdateProject(projectId string, newProjectData *entities.Project) error
	DeleteProject(projectId string) error
}

type projectService struct {
	log   ilog.StdLogger
	query repository.ProjectQuery
}

func NewProjectService(projectQuery repository.ProjectQuery, logger ilog.StdLogger) ProjectService {
	return &projectService{
		log:   logger,
		query: projectQuery,
	}
}

func (p *projectService) GetProjects() ([]*entities.Project, error) {
	projects, err := p.query.GetAllProjects()

	if err != nil {
		p.log.Warn("Error fetching projects from the database")
		return nil, errors.New("No projects present in the database")
	}

	return projects, nil
}

func (p *projectService) GetProjectById(projectId string) (*entities.Project, error) {
	project, err := p.query.GetProjectById(projectId)

	if err != nil {
		p.log.Infof("Project with id %s not found", projectId)
		return nil, api.ProjectNotFoundError
	}

	return project, nil
}

func (p *projectService) CreateProject(newProject *entities.Project) error {
	err := p.query.CreateProject(newProject)

	if err != nil {
		p.log.Error("Could not create Project: ", err)
		return api.ProjectNotCreatedError
	}

	return nil
}

func (p *projectService) UpdateProject(projectId string, newProjectData *entities.Project) error {
	projectExists := p.query.Exists(projectId)
	if !projectExists {
		p.log.Warnf("Project with id %s does not exist", projectId)
		return api.ProjectNotFoundError
	}

	err := p.query.UpdateProject(projectId, newProjectData)
	if err != nil {
		p.log.Error("Could not update Project: ", err)
		return errors.New("There has been an issue updating the project")
	}

	return nil
}

func (p *projectService) DeleteProject(projectId string) error {
	projectExists := p.query.Exists(projectId)
	if !projectExists {
		p.log.Warnf("Project with id %s not found", projectId)
		return api.ProjectNotFoundError
	}

	err := p.query.DeleteProject(projectId)
	if err != nil {
		p.log.Errorf("Could not delete the project with id %s: %s", projectId, err.Error())
		return errors.New("There has been an issue deleting the project")
	}

	return nil
}

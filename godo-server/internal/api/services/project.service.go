package services

import (
	"errors"
	"godo/internal/api"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type ProjectService interface {
	GetProjects(accountId string) ([]*entities.Project, error)
	GetProjectById(userId, accountId string) (*entities.Project, error)
	CreateProject(newProject *entities.Project) (*entities.Project, error)
	Exists(projectId string) bool
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

func (p *projectService) GetProjects(accountId string) ([]*entities.Project, error) {
	projects, err := p.query.GetAllProjects(accountId)

	if err != nil {
		p.log.Infof("Error fetching projects from the database: ", err.Error())
		return nil, errors.New("no projects present in the database")
	}

	return projects, nil
}

func (p *projectService) GetProjectById(projectId, accountId string) (*entities.Project, error) {
	project, err := p.query.GetProjectById(projectId, accountId)

	if err != nil {
		p.log.Infof("Project with id %s not found", projectId)
		return nil, api.ErrorProjectNotFound
	}

	return project, nil
}

func (p *projectService) CreateProject(newProject *entities.Project) (*entities.Project, error) {
	createdProject, err := p.query.CreateProject(newProject)

	if err != nil {
		p.log.Error("Could not create Project: ", err)
		return nil, api.ErrorProjectNotCreated
	}

	return createdProject, nil
}

func (p *projectService) Exists(projectId string) bool {
	return p.query.Exists(projectId)
}

func (p *projectService) UpdateProject(projectId string, newProjectData *entities.Project) error {
	projectExists := p.query.Exists(projectId)
	if !projectExists {
		p.log.Warnf("Project with id %s does not exist", projectId)
		return api.ErrorProjectNotFound
	}

	err := p.query.UpdateProject(projectId, newProjectData)
	if err != nil {
		p.log.Error("Could not update Project: ", err)
		return errors.New("issue updating the project")
	}

	return nil
}

func (p *projectService) DeleteProject(projectId string) error {
	projectExists := p.query.Exists(projectId)
	if !projectExists {
		p.log.Warnf("Project with id %s not found", projectId)
		return api.ErrorProjectNotFound
	}

	err := p.query.DeleteProject(projectId)
	if err != nil {
		p.log.Errorf("Could not delete the project with id %s: %s", projectId, err.Error())
		return errors.New("issue deleting the project")
	}

	return nil
}

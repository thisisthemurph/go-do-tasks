package services

import (
	"godo/internal/api/httperror"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
	"net/http"
)

type ProjectService interface {
	GetProjects() ([]*entities.Project, httperror.HttpError)
	GetProjectById(userId string) (*entities.Project, httperror.HttpError)
	CreateProject(newProject *entities.Project) httperror.HttpError
	UpdateProject(projectId string, newProjectData *entities.Project) httperror.HttpError
	DeleteProject(projectId string) httperror.HttpError
}

type projectService struct {
	log   ilog.StdLogger
	query repository.ProjectQuery
}

func NewProjectService(dao repository.DAO, logger ilog.StdLogger) ProjectService {
	return &projectService{
		log:   logger,
		query: dao.NewProjectQuery(logger),
	}
}

var ProjectNotFoundError = httperror.New(http.StatusNotFound, "The required project could not be found")

func (p *projectService) GetProjects() ([]*entities.Project, httperror.HttpError) {
	projects, err := p.query.GetAllProjects()

	if err != nil {
		reqErr := httperror.New(http.StatusInternalServerError, "There has been an issue querying the database")
		return nil, reqErr
	}

	return projects, nil
}

func (p *projectService) GetProjectById(projectId string) (*entities.Project, httperror.HttpError) {
	project, err := p.query.GetProjectById(projectId)

	if err != nil {
		return nil, ProjectNotFoundError
	}

	return project, nil
}

func (p *projectService) CreateProject(newProject *entities.Project) httperror.HttpError {
	err := p.query.CreateProject(newProject)

	if err != nil {
		return httperror.New(http.StatusInternalServerError, "There has been an issue creating the project")
	}

	return nil
}

func (p *projectService) UpdateProject(projectId string, newProjectData *entities.Project) httperror.HttpError {
	if projectExists := p.query.Exists(projectId); !projectExists {
		return ProjectNotFoundError
	}

	err := p.query.UpdateProject(projectId, newProjectData)
	if err != nil {
		return httperror.New(http.StatusInternalServerError, "There has been an issue updating the project")
	}

	return nil
}

func (p *projectService) DeleteProject(projectId string) httperror.HttpError {
	if projectExists := p.query.Exists(projectId); !projectExists {
		return ProjectNotFoundError
	}

	err := p.query.DeleteProject(projectId)
	if err != nil {
		return httperror.New(http.StatusInternalServerError, "There has been an issue deleting the required project")
	}

	return nil
}

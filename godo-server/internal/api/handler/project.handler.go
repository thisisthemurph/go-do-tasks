package handler

import (
	"godo/internal/api"
	"godo/internal/api/dto"
	"godo/internal/api/services"
	"godo/internal/auth"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"

	"github.com/gorilla/mux"
)

type Projects struct {
	log            ilog.StdLogger
	projectService services.ProjectService
}

func NewProjectsHandler(logger ilog.StdLogger, projectService services.ProjectService) Projects {
	return Projects{log: logger, projectService: projectService}
}

func (p *Projects) GetProjectById(w http.ResponseWriter, r *http.Request) {
	projectId, paramIdExists := getParamFomRequest(r, "id")
	if !paramIdExists {
		http.Error(w, "The ID of the Project must be specified", http.StatusBadRequest)
		return
	}

	project, err := p.projectService.GetProjectById(projectId)

	switch err {
	case nil:
		break
	case api.ProjectNotFoundError:
		api.ReturnError(err, http.StatusNotFound, w)
		return
	default:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	result, _ := dataToJson(*project)
	api.Respond(result, http.StatusOK, w)
}

func (p *Projects) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := p.projectService.GetProjects()
	if err != nil {
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	result, _ := dataToJson(projects)
	api.Respond(result, http.StatusOK, w)
}

func (p *Projects) CreateProject(w http.ResponseWriter, r *http.Request) {
	var projectDto dto.NewProjectDto
	projectDto = r.Context().Value(entities.ProjectKey{}).(dto.NewProjectDto)

	var user auth.User
	user = r.Context().Value(auth.UserKey{}).(auth.User)

	// TODO: Validate the DTO

	project := entities.Project{
		Name:        projectDto.Name,
		Description: projectDto.Description,
		Creator:     user,
	}

	err := p.projectService.CreateProject(&project)
	switch err {
	case nil:
		break
	case api.ProjectNotCreatedError:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	result, _ := dataToJson(project)
	api.Respond(result, http.StatusOK, w)
}

func (p *Projects) UpdateProject(w http.ResponseWriter, r *http.Request) {
	newProjectData := &entities.Project{}
	projectId, _ := getParamFomRequest(r, "id")

	// Gat the new project data from the body
	err := newProjectData.FromJSON(r.Body)
	if err != nil {
		p.log.Error("Could not process Project from request body: ", err)
		api.ReturnError(api.ProjectJsonParseError, http.StatusBadRequest, w)
		return
	}

	// Update the project
	err = p.projectService.UpdateProject(projectId, newProjectData)

	switch err {
	case nil:
		break
	case api.ProjectNotFoundError:
		api.ReturnError(err, http.StatusNotFound, w)
		return
	default:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond("", http.StatusOK, w)
}

func (p *Projects) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectId, _ := getParamFomRequest(r, "id")

	// Delete the project
	err := p.projectService.DeleteProject(projectId)

	switch err {
	case nil:
		break
	case api.ProjectNotFoundError:
		api.ReturnError(err, http.StatusNotFound, w)
		return
	default:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond("", http.StatusOK, w)
}

func getParamFomRequest(r *http.Request, param string) (string, bool) {
	params := mux.Vars(r)
	paramValue, paramExists := params[param]
	return paramValue, paramExists
}

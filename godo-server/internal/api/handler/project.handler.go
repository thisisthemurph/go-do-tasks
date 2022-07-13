package handler

import (
	"godo/internal/api/dto"
	"godo/internal/api/httperror"
	"godo/internal/api/httputils"
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
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	project, err := p.projectService.GetProjectById(projectId)
	if err != nil {
		http.Error(w, err.AsJson(), err.GetStatusCode())
		return
	}

	json, _ := dataToJson(*project)
	httputils.MakeHttpResponse(w, json, http.StatusOK)
}

func (p *Projects) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := p.projectService.GetProjects()
	if err != nil {
		http.Error(w, err.AsJson(), err.GetStatusCode())
		return
	}

	json, _ := dataToJson(projects)
	httputils.MakeHttpResponse(w, json, http.StatusOK)
}

func (p *Projects) CreateProject(w http.ResponseWriter, r *http.Request) {
	var projectDto dto.NewProjectDto
	projectDto = r.Context().Value(entities.ProjectKey{}).(dto.NewProjectDto)

	var user auth.User
	user = r.Context().Value(auth.UserKey{}).(auth.User)

	project := entities.Project{
		Name:        projectDto.Name,
		Description: projectDto.Description,
		Creator:     user,
	}

	err2 := p.projectService.CreateProject(&project)
	if err2 != nil {
		http.Error(w, err2.AsJson(), err2.GetStatusCode())
		return
	}

	result, _ := dataToJson(project)
	httputils.MakeHttpResponse(w, result, http.StatusOK)
}

func (p *Projects) UpdateProject(w http.ResponseWriter, r *http.Request) {
	newProjectData := &entities.Project{}
	projectId, _ := getParamFomRequest(r, "id")

	// Gat the new project data from the body
	requestErr := newProjectData.FromJSON(r.Body)
	if requestErr != nil {
		err := httperror.New(http.StatusInternalServerError, "There has been an issue processing the update")
		httputils.MakeHttpError(err.GetStatusCode(), err.AsJson())
		return
	}

	// Update the project
	err := p.projectService.UpdateProject(projectId, newProjectData)
	if err != nil {
		httputils.MakeHttpError(err.GetStatusCode(), err.AsJson())
		return
	}

	httputils.MakeHttpResponse(w, "", http.StatusOK)
}

func (p *Projects) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectId, _ := getParamFomRequest(r, "id")

	err := p.projectService.DeleteProject(projectId)
	if err != nil {
		httputils.MakeHttpError(err.GetStatusCode(), err.AsJson())
		return
	}

	httputils.MakeHttpResponse(w, "", http.StatusOK)
}

func getParamFomRequest(r *http.Request, param string) (string, bool) {
	params := mux.Vars(r)
	paramValue, paramExists := params[param]
	return paramValue, paramExists
}

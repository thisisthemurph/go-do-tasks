package handler

import (
	"fmt"
	"godo/internal/api/httperror"
	"godo/internal/api/httputils"
	"godo/internal/api/services"
	"godo/internal/repository/entities"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) ProjectHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	var status int = http.StatusOK

	h.log.Infof("ProjectHandler. Method: %v", r.Method)

	switch r.Method {
	case http.MethodGet:
		result, status = getProject(h.projectService, r)
		break
	case http.MethodPost:
		result, status = createProject(h.projectService, r)
		break
	case http.MethodPut:
		result, status = updateProject(h.projectService, r)
		break
	case http.MethodDelete:
		result, status = deleteProject(h.projectService, r)
		break
	default:
		h.log.Errorf("Bad method %v", r.Method)
		result = fmt.Sprintf("Bad method %s", r.Method)
		status = http.StatusBadGateway
		break
	}

	httputils.MakeHttpResponse(w, result, status)
}

func getProject(projectService services.ProjectService, r *http.Request) (string, int) {
	projectId, paramIdExists := getParamFomRequest(r, "id")
	if !paramIdExists {
		return getAllProjects(projectService, r)
	}

	project, err := projectService.GetProjectById(projectId)
	if err != nil {
		return err.AsJson(), err.GetStatusCode()
	}

	json, _ := dataToJson(*project)
	return json, http.StatusOK
}

func getAllProjects(projectService services.ProjectService, r *http.Request) (string, int) {
	projects, err := projectService.GetProjects()
	if err != nil {
		return err.AsJson(), err.GetStatusCode()
	}

	json, _ := dataToJson(projects)
	return json, http.StatusOK
}

func createProject(projectService services.ProjectService, r *http.Request) (string, int) {
	var project entities.Project
	project.FromHttpRequest(r)

	err := projectService.CreateProject(&project)
	if err != nil {
		return err.AsJson(), err.GetStatusCode()
	}

	return "", http.StatusOK
}

func updateProject(projectService services.ProjectService, r *http.Request) (string, int) {
	newProjectData := &entities.Project{}
	projectId, _ := getParamFomRequest(r, "id")

	// Gat the new project data from the body
	requestErr := newProjectData.FromJSON(r.Body)
	if requestErr != nil {
		err := httperror.New(http.StatusInternalServerError, "There has been an issue processing the update")
		return err.AsJson(), err.GetStatusCode()
	}

	// Update the project
	err := projectService.UpdateProject(projectId, newProjectData)
	if err != nil {
		return err.AsJson(), err.GetStatusCode()
	}

	return "", http.StatusOK
}

func deleteProject(projectService services.ProjectService, r *http.Request) (string, int) {
	projectId, _ := getParamFomRequest(r, "id")

	err := projectService.DeleteProject(projectId)
	if err != nil {
		return err.AsJson(), err.GetStatusCode()
	}

	return "", http.StatusOK
}

func getParamFomRequest(r *http.Request, param string) (string, bool) {
	params := mux.Vars(r)
	paramValue, paramExists := params[param]
	return paramValue, paramExists
}

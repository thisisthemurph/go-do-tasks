package handler

import (
	"fmt"
	"godo/internal/api/httputils"
	"godo/internal/repository/entities"
	"godo/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) ProjectHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	var status int = http.StatusOK

	log.Printf("Received method: %v", r.Method)

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
	default:
		result = fmt.Sprintf("Bad method %s", r.Method)
		status = http.StatusBadGateway
		break
	}

	httputils.MakeHttpResponse(w, result, status)
}

func getProject(projectService services.ProjectService, r *http.Request) (string, int) {
	projectId, paramIdExists := getProjectIdFromRequest(r)
	if !paramIdExists {
		return getAllProjects(projectService, r)
	}

	project, err := projectService.GetProjectById(projectId)
	if err != nil {
		log.Printf("E: %v\n", err)
		message := "The required project could not be found"
		return httputils.MakeHttpError(http.StatusNotFound, message)
	}

	log.Printf("Name: %v\n", project.Name)
	json, _ := dataToJson(*project)
	return json, http.StatusOK
}

func getAllProjects(projectService services.ProjectService, r *http.Request) (string, int) {
	projects, err := projectService.GetProjects()
	log.Printf("Project: %v", len(projects))
	if err != nil {
		return httputils.MakeHttpError(500, err.Error())
	}

	json, _ := dataToJson(projects)
	return json, http.StatusOK
}

func createProject(projectService services.ProjectService, r *http.Request) (string, int) {
	project := &entities.Project{}

	requestErr := project.FromJSON(r.Body)

	log.Printf("Project: %#v", project)

	createErr := projectService.CreateProject(project)
	if createErr != nil {
		return httputils.MakeHttpError(http.StatusBadRequest, requestErr.Error())
	}

	return "", http.StatusOK
}

func updateProject(projectService services.ProjectService, r *http.Request) (string, int) {
	newProjectData := &entities.Project{}
	projectId, _ := getProjectIdFromRequest(r)

	_, existsErr := projectService.GetProjectById(projectId)
	if existsErr != nil {
		return httputils.MakeHttpError(
			http.StatusNotFound,
			"The specified product could not be found")
	}

	// Gat the project from the body
	requestErr := newProjectData.FromJSON(r.Body)
	if requestErr != nil {
		return httputils.MakeHttpError(http.StatusBadRequest, requestErr.Error())
	}

	// Update the project
	err := projectService.UpdateProject(projectId, newProjectData)
	if err != nil {
		return httputils.MakeHttpError(500, err.Error())
	}

	return "", http.StatusOK
}

func getProjectIdFromRequest(r *http.Request) (string, bool) {
	params := mux.Vars(r)
	projectId, paramIdExists := params["id"]
	return projectId, paramIdExists
}

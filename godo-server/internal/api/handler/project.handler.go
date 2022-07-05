package handler

import (
	"fmt"
	"godo/internal/api/httputils"
	"godo/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) ProjectHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	var status int = http.StatusOK

	switch r.Method {
	case http.MethodGet:
		result, status = getProject(h.projectService, r)
		break
	default:
		result = fmt.Sprintf("Bad method %s", r.Method)
		status = http.StatusBadGateway
		break
	}

	httputils.MakeHttpResponse(w, result, status)
}

func getProject(projectService services.ProjectService, r *http.Request) (string, int) {
	params := mux.Vars(r)
	projectId, paramIdExists := params["id"]
	log.Println(params)
	log.Printf("ID: %v\n", projectId)
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
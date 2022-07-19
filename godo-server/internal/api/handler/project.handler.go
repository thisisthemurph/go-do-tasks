package handler

import (
	"godo/internal/api"
	"godo/internal/api/dto"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
)

type Projects struct {
	log            ilog.StdLogger
	projectService services.ProjectService
}

func NewProjectsHandler(logger ilog.StdLogger, projectService services.ProjectService) Projects {
	return Projects{log: logger, projectService: projectService}
}

func (p *Projects) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	user := entities.User{}
	user = r.Context().Value(entities.UserKey{}).(entities.User)

	projects, err := p.projectService.GetProjects(user.AccountId)

	if err != nil {
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(projects, http.StatusOK, w)
}

func (p *Projects) GetProjectById(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	projectId, paramIdExists := getParamFomRequest(r, "id")
	if !paramIdExists {
		http.Error(w, "The ID of the Project must be specified", http.StatusBadRequest)
		return
	}

	project, err := p.projectService.GetProjectById(projectId, user.AccountId)

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

	api.Respond(project, http.StatusOK, w)
}

func (p *Projects) CreateProject(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	var projectDto dto.NewProjectDto

	ctx := r.Context()
	projectDto = ctx.Value(entities.ProjectKey{}).(dto.NewProjectDto)
	// user = ctx.Value(entities.UserKey{}).(entities.User)
	user = getUserFromContext(ctx)

	// Validate the project looks OK
	err := projectDto.Validate()
	if err != nil {
		p.log.Debug(err)
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Create the newProject
	newProject := entities.Project{
		Name:        projectDto.Name,
		Description: projectDto.Description,
		Creator:     user,
	}

	createdProject, err := p.projectService.CreateProject(&newProject)
	switch err {
	case nil:
		break
	case api.ProjectNotCreatedError:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	default:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(createdProject, http.StatusOK, w)
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

	api.Respond("", http.StatusNoContent, w)
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

	api.Respond("", http.StatusNoContent, w)
}

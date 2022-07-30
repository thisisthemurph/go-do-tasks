package handler

import (
	"godo/internal/api"
	"godo/internal/api/dto"
	"godo/internal/api/errorhandler"
	ehand "godo/internal/api/errorhandler"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
	"strconv"
)

type Projects struct {
	log            ilog.StdLogger
	projectService services.ProjectService
	tagService     services.TagService
	eh             errorhandler.ErrorHandler
}

func NewProjectsHandler(
	logger ilog.StdLogger,
	projectService services.ProjectService,
	tagService services.TagService) Projects {

	return Projects{
		log:            logger,
		projectService: projectService,
		tagService:     tagService,
		eh:             errorhandler.New(),
	}
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
	if status := p.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(project, http.StatusOK, w)
}

func (p *Projects) CreateProject(w http.ResponseWriter, r *http.Request) {
	var projectDto dto.NewProjectDto

	err := decodeJSONBody(w, r, &projectDto)
	if err != nil {
		handleMalformedJSONError(w, err)
		return
	}

	user := getUserFromContext(r.Context())

	// Validate the project looks OK
	err = projectDto.Validate()
	if err != nil {
		p.log.Debug(err)
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Create the new Project
	newProject := entities.Project{
		Name:        projectDto.Name,
		Description: projectDto.Description,
		Creator:     user,
	}

	createdProject, err := p.projectService.CreateProject(&newProject)
	if status := p.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(createdProject, http.StatusCreated, w)
}

func (p *Projects) UpdateProject(w http.ResponseWriter, r *http.Request) {
	newProjectData := &entities.Project{}
	projectId, _ := getParamFomRequest(r, "id")

	// Gat the new project data from the body
	err := api.FromJSON(&newProjectData, r.Body)
	if err != nil {
		p.log.Error("Could not process Project from request body: ", err)
		api.ReturnError(ehand.ErrorProjectJSONParse, http.StatusBadRequest, w)
		return
	}

	// Update the project
	err = p.projectService.UpdateProject(projectId, newProjectData)
	if status := p.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

// AddTagToProject TODO: Ensure project doesn't already have a tag with the same name
func (p *Projects) AddTagToProject(w http.ResponseWriter, r *http.Request) {
	projId, _ := getParamFomRequest(r, "id")
	tagDto, err := getDtoFromJSONBody[dto.NewTagDto](w, r)
	if err != nil {
		p.log.Error("Error getting Dto from JSON body: ", err)
		return
	}

	// Create the new tag
	var tag entities.Tag
	tag.Name = tagDto.Name
	tag.ProjectId = projId

	_, err = p.tagService.CreateTag(tag)
	if err != nil {
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

func (p *Projects) UpdateProjectTag(w http.ResponseWriter, r *http.Request) {
	projId, _ := getParamFomRequest(r, "projectId")
	tagIdValue, _ := getParamFomRequest(r, "tagId")
	tagId, err := strconv.ParseUint(tagIdValue, 10, 32)
	if err != nil {
		p.log.Error(err)
		api.ReturnError(ehand.ErrorTagMalformedId, http.StatusBadRequest, w)
		return
	}

	tagDto, err := getDtoFromJSONBody[dto.NewTagDto](w, r)
	if err != nil {
		p.log.Error("Error getting Dto from JSON body: ", err)
		return
	}

	// Get the tag from the database
	tag, err := p.tagService.GetTagById(uint(tagId), projId)
	tag.Name = tagDto.Name

	_, err = p.tagService.UpdateTag(*tag)
	if err != nil {
		api.ReturnError(ehand.ErrorTagNotUpdated, http.StatusInternalServerError, w)
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

func (p *Projects) DeleteProjectTag(w http.ResponseWriter, r *http.Request) {
	projId, _ := getParamFomRequest(r, "projectId")
	tagId, _ := getUintParamFomRequest(r, "tagId")

	_, err := p.tagService.DeleteTag(tagId, projId)
	if status := p.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

func (p *Projects) UpdateProjectStatus(w http.ResponseWriter, r *http.Request) {
	var statusDto dto.ProjectStatusUpdateDto
	err := decodeJSONBody(w, r, &statusDto)
	if err != nil {
		handleMalformedJSONError(w, err)
		return
	}

	projectId, _ := getParamFomRequest(r, "id")
	user := getUserFromContext(r.Context())

	// Get the project from the database
	project, err := p.projectService.GetProjectById(projectId, user.AccountId)
	if err != nil {
		p.log.Debugf("Could not find project with projectId %s and accountId %s", projectId, user.AccountId)
		api.ReturnError(ehand.ErrorProjectNotFound, http.StatusNotFound, w)
		return
	}

	// Update the project
	project.Status = statusDto.Status
	err = p.projectService.UpdateProject(projectId, project)
	if err != nil {
		p.log.Debugf("Could not update with projectId %s and accountId %s", projectId, user.AccountId)
		api.ReturnError(ehand.ErrorProjectNotFound, http.StatusNotFound, w)
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

func (p *Projects) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectId, _ := getParamFomRequest(r, "id")

	// Delete the project
	err := p.projectService.DeleteProject(projectId)
	if status := p.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

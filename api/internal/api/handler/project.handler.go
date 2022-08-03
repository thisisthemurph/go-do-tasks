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

// swagger:route GET /project Projects listProjectInfo
//
// Returns a list of projects associated with the authenticated account
//
// responses:
//  200: projectInfoResponse
//  500: errorResponse
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

// swagger:route GET /project/{projectId} Projects getProject
//
// Returns the specified project
//
// responses:
//  200: projectResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
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

// swagger:route POST /project Projects createProject
//
// Creates the given project resource
//
// responses:
//  201: projectResponse
//  400: errorResponse
//  500: errorResponse
func (p *Projects) CreateProject(w http.ResponseWriter, r *http.Request) {
	projectDto, err := getDtoFromJSONBody[dto.NewProjectDto](w, r)
	if err != nil {
		return
	}

	user := getUserFromContext(r.Context())

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

// swagger:route PUT /project/{projectId} Projects updateProject
//
// Updates the values of the specified project
//
// responses:
//  204: noContent
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
// UpdateProject TODO: Improve to incorporate a DTO
func (p *Projects) UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectId, _ := getParamFomRequest(r, "id")

	// Gat the new project data from the body
	newProjectData, err := getDtoFromJSONBody[entities.Project](w, r)
	if err != nil {
		return
	}

	// Update the project
	err = p.projectService.UpdateProject(projectId, newProjectData)
	if status := p.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

// swagger:route POST /project/{projectId}/tag Projects createTag
//
// Created and associated the tag with the specified project
//
// responses:
//  204: noContent
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
// AddTagToProject TODO: Ensure project doesn't already have a tag with the same name
func (p *Projects) AddTagToProject(w http.ResponseWriter, r *http.Request) {
	projId, _ := getParamFomRequest(r, "id")
	tagDto, err := getDtoFromJSONBody[dto.NewTagDto](w, r)
	if err != nil {
		return
	}

	// Create the new tag
	var tag entities.Tag
	tag.Name = tagDto.Name
	tag.ProjectId = projId

	_, err = p.tagService.CreateTag(tag)
	if status := p.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

// UpdateProjectTag TODO: Verify that this method is being run
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

// swagger:route DELETE /project/{projectId}/tag Projects deleteTag
//
// Created and associated the tag with the specified project
//
// responses:
//  204: noContent
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
func (p *Projects) DeleteProjectTag(w http.ResponseWriter, r *http.Request) {
	projId, _ := getParamFomRequest(r, "projectId")
	tagId, _ := getUintParamFomRequest(r, "tagId")

	_, err := p.tagService.DeleteTag(tagId, projId)
	if status := p.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

// swagger:route PUT /project/{projectId}/status Projects updateProjectStatus
//
// Updates the status of the specified project
//
// responses:
//  204: noContent
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
func (p *Projects) UpdateProjectStatus(w http.ResponseWriter, r *http.Request) {
	statusDto, err := getDtoFromJSONBody[dto.ProjectStatusUpdateDto](w, r)
	if err != nil {
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

// swagger:route DELETE /project Projects deleteProject
//
// Deletes the given project resource
//
// responses:
//  204: noContent
//  400: errorResponse
//  500: errorResponse
func (p *Projects) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectId, _ := getParamFomRequest(r, "id")

	// Delete the project
	err := p.projectService.DeleteProject(projectId)
	if status := p.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

// Generic Swagger documentation

// swagger:parameters getProject updateProject deleteProject createTag deleteTag
type ProductUUIDParameter struct {
	// The ID of the specified Project
	// in: path
	// required: true
	// pattern: ^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
	// example: f9d633f8-c684-4dc3-b410-d36df912c4c1
	ID string `json:"id"`
}

// swagger:parameters createProject
type NewProjectParameter struct {
	// The new project resource to be created
	// in: body
	// required: true
	Body dto.NewProjectDto
}

// swagger:parameters updateProject
type UpdateProjectParameter struct {
	// The new project resource to be created
	// in: body
	// required: true
	Body dto.UpdateProjectDto
}

// swagger:parameters createTag
type NewTag struct {
	// The tag to be created and associated with the given project
	// in: body
	// required: true
	Body dto.NewTagDto
}

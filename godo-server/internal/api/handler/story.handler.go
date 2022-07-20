package handler

import (
	"godo/internal/api"
	"godo/internal/api/dto"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
)

type Stories struct {
	log            ilog.StdLogger
	storyService   services.StoryService
	projectService services.ProjectService
}

func NewStoriesHandler(
	logger ilog.StdLogger,
	storyService services.StoryService,
	projectService services.ProjectService) Stories {

	return Stories{
		log:            logger,
		storyService:   storyService,
		projectService: projectService,
	}
}

func (s *Stories) GetAllStories(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	stories, err := s.storyService.GetStories(user.AccountId)

	if err != nil {
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(stories, http.StatusOK, w)
}

func (s *Stories) GetStoryById(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	storyId, paramIdExists := getParamFomRequest(r, "id")
	if !paramIdExists {
		http.Error(w, "The ID of the Project must be specified", http.StatusBadRequest)
		return
	}

	story, err := s.storyService.GetStoryById(storyId, user.AccountId)

	switch err {
	case nil:
		break
	case api.ErrorStoryNotFound:
		api.ReturnError(err, http.StatusNotFound, w)
		return
	default:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(story, http.StatusOK, w)
}

func (s *Stories) CreateStory(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	var storyDto dto.NewStoryDto
	err := decodeJSONBody(w, r, &storyDto)
	if err != nil {
		handleMalformedJSONError(w, err)
		return
	}

	// Validate the story looks OK
	err = storyDto.Validate()
	if err != nil {
		s.log.Debug(err)
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Ensure the project project exists
	projectExists := s.projectService.Exists(storyDto.ProjectId)
	if !projectExists {
		s.log.Debugf("Project with projectId %s not found", storyDto.ProjectId)
		api.ReturnError(api.ErrorProjectNotFound, http.StatusNotFound, w)
		return
	}

	// Create the newStory
	newStory := entities.Story{
		Name:        storyDto.Name,
		Description: storyDto.Description,
		ProjectId:   storyDto.ProjectId,
		Creator:     user,
	}

	createdProject, err := s.storyService.CreateStory(&newStory)
	switch err {
	case nil:
		break
	case api.ErrorStoryNotCreated:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	default:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(createdProject, http.StatusCreated, w)
}

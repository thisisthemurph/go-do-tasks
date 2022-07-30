package handler

import (
	"godo/internal/api"
	"godo/internal/api/dto"
	ehand "godo/internal/api/errorhandler"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
)

type Stories struct {
	log            ilog.StdLogger
	storyService   services.StoryService
	projectService services.ProjectService
	eh             ehand.ErrorHandler
}

func NewStoriesHandler(
	logger ilog.StdLogger,
	storyService services.StoryService,
	projectService services.ProjectService) Stories {

	return Stories{
		log:            logger,
		storyService:   storyService,
		projectService: projectService,
		eh:             ehand.New(),
	}
}

func (s *Stories) GetStoriesInfo(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	info, err := s.storyService.GetStoriesInfo(user.AccountId)
	if err != nil {
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(info, http.StatusOK, w)
}

func (s *Stories) GetStoryById(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	storyId, paramIdExists := getParamFomRequest(r, "id")
	if !paramIdExists {
		http.Error(w, "The ID of the Project must be specified", http.StatusBadRequest)
		return
	}

	story, err := s.storyService.GetStoryById(user.AccountId, storyId)
	if status := s.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(story, http.StatusOK, w)
}

func (s *Stories) CreateStory(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	storyDto, err := getDtoFromJSONBody[dto.NewStoryDto](w, r)
	if err != nil {
		return
	}

	// Ensure the project exists
	projectExists := s.projectService.Exists(storyDto.ProjectId)
	if !projectExists {
		s.log.Debugf("Project with projectId %s not found", storyDto.ProjectId)
		api.ReturnError(ehand.ErrorProjectNotFound, http.StatusNotFound, w)
		return
	}

	// Create the newStory
	newStory := entities.Story{
		Name:        storyDto.Name,
		Description: storyDto.Description,
		ProjectId:   storyDto.ProjectId,
		Creator:     user,
	}

	created, err := s.storyService.CreateStory(&newStory)
	if status := s.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(created, http.StatusCreated, w)
}

func (s *Stories) UpdateStory(w http.ResponseWriter, r *http.Request) {
	storyId, _ := getParamFomRequest(r, "id")

	storyDto, err := getDtoFromJSONBody[dto.NewStoryDto](w, r)
	if err != nil {
		return
	}

	user := getUserFromContext(r.Context())

	// Ensure the project for the story exists
	projectExists := s.projectService.Exists(storyDto.ProjectId)
	if !projectExists {
		s.log.Debugf("Project with projectId %s not found", storyDto.ProjectId)
		api.ReturnError(ehand.ErrorProjectNotFound, http.StatusNotFound, w)
		return
	}

	// Update the story
	ns, err := s.storyService.GetStoryById(user.AccountId, storyId)
	if err != nil {
		s.log.Debugf("The story with storyId %s and accountId % could not be found", storyId, user.AccountId)
		api.ReturnError(ehand.ErrorStoryNotFound, http.StatusNotFound, w)
		return
	}

	ns.Name = storyDto.Name
	ns.Description = storyDto.Description

	// Verify and update the projectId
	if storyDto.ProjectId != ns.ProjectId {
		ns.ProjectId = storyDto.ProjectId
		// TODO: If update projectId -> Verify project is owned by account
	}

	err = s.storyService.UpdateStory(storyId, ns)
	if status := s.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

func (s *Stories) DeleteStory(w http.ResponseWriter, r *http.Request) {
	storyId, _ := getParamFomRequest(r, "id")

	err := s.storyService.DeleteStory(storyId)
	if status := s.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

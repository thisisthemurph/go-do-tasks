package handler

import (
	"encoding/json"
	"godo/internal/api/services"
	"net/http"
)

type IHandler interface {
	ProjectHandler(http.ResponseWriter, *http.Request)
	StoryHandler(http.ResponseWriter, *http.Request)
}

type Handler struct {
	projectService services.ProjectService
	storyService   services.StoryService
}

func MakeHandlers(
	projectService *services.ProjectService,
	storyService *services.StoryService) IHandler {

	return &Handler{
		projectService: *projectService,
		storyService:   *storyService,
	}
}

// Converts a struct into a JSON object
func dataToJson(d interface{}) (string, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

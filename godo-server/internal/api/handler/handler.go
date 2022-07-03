package handler

import (
	"encoding/json"
	"godo/internal/services"
	"net/http"
)

type IHandler interface {
	StoryHandler(http.ResponseWriter, *http.Request)
}

type Handler struct {
	storyService services.StoryService
}

func MakeHandlers(storyService *services.StoryService) IHandler {
	return &Handler{
		storyService: *storyService,
	}
}

// Converts a struct into a JSON object
func dataToJson(d interface{}) (string, error) {
	data, err := json.Marshal(d)
	if (err != nil) {
		return "", err
	}

	return string(data), nil
}
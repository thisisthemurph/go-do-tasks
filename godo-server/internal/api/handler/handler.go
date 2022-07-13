package handler

import (
	"encoding/json"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"net/http"
)

type IHandler interface {
	StoryHandler(http.ResponseWriter, *http.Request)
}

type Handler struct {
	authService    services.AuthService
	log            ilog.StdLogger
	projectService services.ProjectService
	storyService   services.StoryService
	userService    services.UserService
}

func MakeHandlers(
	logger ilog.StdLogger,
	projectService *services.ProjectService,
	storyService *services.StoryService,
	authService *services.AuthService,
	userService *services.UserService) IHandler {

	return &Handler{
		log:            logger,
		projectService: *projectService,
		storyService:   *storyService,
		authService:    *authService,
		userService:    *userService,
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

package handler

import (
	"fmt"
	"godo/internal/api/httputils"
	"godo/internal/api/services"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) StoryHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	var status int = http.StatusOK

	switch r.Method {
	case http.MethodGet:
		result, status = getStory(h.storyService, r)
		break
	default:
		result = fmt.Sprintf("Bad method %s", r.Method)
		status = http.StatusBadGateway
		break
	}

	httputils.MakeHttpResponse(w, result, status)
}

func getStory(storyService services.StoryService, r *http.Request) (string, int) {
	params := mux.Vars(r)
	storyId, paramIdExists := params["id"]
	if !paramIdExists {
		return getAllStories(storyService, r)
	}

	story, err := storyService.GetStoryById(storyId)
	if err != nil {
		message := "The required story could not be found"
		return httputils.MakeHttpError(http.StatusNotFound, message)
	}

	json, _ := dataToJson(*story)
	return json, http.StatusOK
}

func getAllStories(storyService services.StoryService, r *http.Request) (string, int) {
	stories, err := storyService.GetStories()
	if err != nil {
		return httputils.MakeHttpError(500, err.Error())
	}

	json, _ := dataToJson(stories)
	return json, http.StatusOK
}

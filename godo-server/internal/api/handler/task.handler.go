package handler

import (
	"godo/internal/api"
	"godo/internal/api/dto"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/helper/validate"
	"godo/internal/repository/entities"
	"net/http"
)

type Tasks struct {
	log         ilog.StdLogger
	taskService services.TaskService
}

func NewTasksHandler(logger ilog.StdLogger, taskService services.TaskService) Tasks {
	return Tasks{log: logger, taskService: taskService}
}

func (t *Tasks) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	tasks, err := t.taskService.GetTasks(user.AccountId)
	if err != nil {
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(tasks, http.StatusOK, w)
}

func (t *Tasks) GetTaskById(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	taskId, _ := getParamFomRequest(r, "id")

	task, err := t.taskService.GetTaskById(taskId, user.AccountId)
	if err != nil {
		api.ReturnError(api.ErrorTaskNotFound, http.StatusNotFound, w)
		return
	}

	api.Respond(task, http.StatusOK, w)
}

// CreateTask TODO: Validate that the story belongs to the user's account
func (t *Tasks) CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDto dto.NewTaskDto
	err := decodeJSONBody(w, r, &taskDto)
	if err != nil {
		handleMalformedJSONError(w, err)
		return
	}

	user := getUserFromContext(r.Context())

	err = validate.Struct(taskDto)
	if err != nil {
		t.log.Error("Invalid NewTaskDto: ", err)
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Create the task
	newTask := entities.Task{
		Name:        taskDto.Name,
		Description: taskDto.Description,
		StoryId:     taskDto.StoryId,
		Creator:     user,
	}

	created, err := t.taskService.CreateTask(newTask)
	if err != nil {
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(created, http.StatusCreated, w)
}

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

type Tasks struct {
	log         ilog.StdLogger
	taskService services.TaskService
	tagService  services.TagService
	eh          ehand.ErrorHandler
}

func NewTasksHandler(
	logger ilog.StdLogger,
	taskService services.TaskService,
	tagService services.TagService) Tasks {

	return Tasks{
		log:         logger,
		taskService: taskService,
		tagService:  tagService,
		eh:          ehand.New(),
	}
}

func (t *Tasks) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	tasks, err := t.taskService.GetTasks(user.AccountId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(tasks, http.StatusOK, w)
}

func (t *Tasks) GetTaskById(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	taskId, _ := getParamFomRequest(r, "id")

	task, err := t.taskService.GetTaskById(taskId, user.AccountId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(task, http.StatusOK, w)
}

// CreateTask TODO: Validate that the story belongs to the user's account
func (t *Tasks) CreateTask(w http.ResponseWriter, r *http.Request) {
	taskDto, err := getDtoFromJSONBody[dto.NewTaskDto](w, r)
	if err != nil {
		return
	}

	user := getUserFromContext(r.Context())

	// Create the task
	newTask := entities.Task{
		Name:        taskDto.Name,
		Description: taskDto.Description,
		StoryId:     taskDto.StoryId,
		Creator:     user,
	}

	created, err := t.taskService.CreateTask(newTask)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(created, http.StatusCreated, w)
}

// UpdateTask TODO: Remove deduplication of error handling in update
// methods and other functions in the handlers
func (t *Tasks) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskDto, err := getDtoFromJSONBody[dto.UpdateTaskDto](w, r)
	if err != nil {
		return
	}

	user := getUserFromContext(r.Context())
	taskId, _ := getParamFomRequest(r, "id")

	// Fetch the task from the database
	task, err := t.taskService.GetTaskById(taskId, user.AccountId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	// Update the fetched task
	task.Name = taskDto.Name
	task.Description = taskDto.Description
	task.Type = taskDto.Type
	task.Status = taskDto.Status

	if len(taskDto.StoryId) > 0 {
		task.StoryId = taskDto.StoryId
	}

	updated, err := t.taskService.UpdateTask(task)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(updated, http.StatusOK, w)
}

func (t *Tasks) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	taskId, _ := getParamFomRequest(r, "id")

	taskDto, err := getDtoFromJSONBody[dto.UpdateTaskStatusDto](w, r)
	if err != nil {
		return
	}

	// Get the task from the database
	task, err := t.taskService.GetTaskById(taskId, user.AccountId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	// Update the task
	task.Status = taskDto.Status
	updated, err := t.taskService.UpdateTask(task)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(updated, http.StatusOK, w)
}

func (t *Tasks) UpdateTaskType(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	taskId, _ := getParamFomRequest(r, "id")

	taskDto, err := getDtoFromJSONBody[dto.UpdateTaskTypeDto](w, r)
	if err != nil {
		return
	}

	// Get the task from the database
	task, err := t.taskService.GetTaskById(taskId, user.AccountId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	// Update the task
	task.Type = taskDto.Type
	updated, err := t.taskService.UpdateTask(task)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(updated, http.StatusOK, w)
}

func (t *Tasks) AddTag(w http.ResponseWriter, r *http.Request) {
	tagId, _ := getUintParamFomRequest(r, "tagId")
	taskId, _ := getParamFomRequest(r, "taskId")

	user := getUserFromContext(r.Context())

	// Ensure that the task exists for the user's account
	_, err := t.taskService.GetTaskById(taskId, user.AccountId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	// Add the tag to the task in the database
	err = t.tagService.AddToTask(tagId, taskId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

func (t *Tasks) RemoveTag(w http.ResponseWriter, r *http.Request) {
	tagId, _ := getUintParamFomRequest(r, "tagId")
	taskId, _ := getParamFomRequest(r, "taskId")

	user := getUserFromContext(r.Context())

	// Ensure that the task exists for the user's account
	_, err := t.taskService.GetTaskById(taskId, user.AccountId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	// Add the tag to the task in the database
	err = t.tagService.RemoveFromTask(taskId, tagId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond("", http.StatusNoContent, w)
}

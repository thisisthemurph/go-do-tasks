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

// swagger:route GET /task Tasks listTasks
//
// Returns a list of Tasks associated with the authenticated account
//
// responses:
//  200: taskInfoResponse
//  500: errorResponse
func (t *Tasks) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())

	tasks, err := t.taskService.GetTasks(user.AccountId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(tasks, http.StatusOK, w)
}

// swagger:route GET /task/{taskId} Tasks getTask
//
// Returns the requested Task
//
// responses:
//  200: taskResponse
//  404: errorResponse
//  500: errorResponse
func (t *Tasks) GetTaskById(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	taskId, _ := getParamFomRequest(r, "id")

	task, err := t.taskService.GetTaskById(taskId, user.AccountId)
	if status := t.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(task, http.StatusOK, w)
}

// swagger:route POST /task Tasks createTask
//
// Creates the given Task
//
// responses:
//  201: taskResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
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

// swagger:route PUT /task/{taskId} Tasks updateTask
//
// Updates the given Task
//
// responses:
//  200: taskResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
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

// swagger:route PUT /task/{taskId}/status Tasks updateTaskStatus
//
// Updates the status of the specified Task
//
// responses:
//  200: taskResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
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

// swagger:route PUT /task/{taskId}/type Tasks updateTaskType
//
// Updates the type of the specified Task
//
// responses:
//  200: taskResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
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

// swagger:route PUT /task/{taskId}/tag/{tagId} Tasks addTaskTag
//
// Associates the give existing Tag with the specified Task
//
// responses:
//  200: noContent
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
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

// swagger:route DELETE /task/{taskId}/tag/{tagId} Tasks removeTaskTag
//
// Disassociate the give existing Tag with the specified Task - does not delete the tag.\n
// To delete the tag, the tag should be deleted from the associated Project
//
// responses:
//  200: noContent
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse
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

// Generic Swagger documentation

// swagger:parameters getTask updateTask updateTaskStatus updateTaskType addTaskTag removeTaskTag
type TaskUUIDParameter struct {
	// The ID of the specified Task
	// in: path
	// required: true
	// pattern: ^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
	// example: f9d633f8-c684-4dc3-b410-d36df912c4c1
	ID string `json:"taskId"`
}

// swagger:parameters addTaskTag removeTaskTag
type TagIDParameter struct {
	// The ID of the specified Tag
	// in: path
	// required: true
	// example: 7 or 52
	ID uint `json:"tagId"`
}

// swagger:parameters createTask
type NewTaskParameter struct {
	// The task to be created
	// in: body
	// required: true
	Body dto.NewTaskDto
}

// swagger:parameters updateTask
type UpdateTaskParameter struct {
	// The task to be created
	// in: body
	// required: true
	Body dto.UpdateTaskDto
}

// swagger:parameters updateTaskStatus
type UpdateTaskStatusParameter struct {
	// The status of the task to be updated
	// in: body
	// required: true
	Body dto.UpdateTaskStatusDto
}

// swagger:parameters updateTaskType
type UpdateTaskTypeParameter struct {
	// The type of the task to be updated
	// in: body
	// required: true
	Body dto.UpdateTaskTypeDto
}

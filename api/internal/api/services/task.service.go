package services

import (
	"errors"
	ehand "godo/internal/api/errorhandler"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type TaskService interface {
	Exists(accountId, taskId string) bool
	GetTasks(accountId string) (entities.TaskList, error)
	GetTaskById(taskId, accountId string) (*entities.Task, error)
	CreateTask(newTask entities.Task) (*entities.Task, error)
	UpdateTask(newTask *entities.Task) (*entities.Task, error)
}

type taskService struct {
	log   ilog.StdLogger
	query repository.TaskQuery
}

func NewTaskService(query repository.TaskQuery, logger ilog.StdLogger) TaskService {
	return &taskService{log: logger, query: query}
}

func (t *taskService) GetTasks(accountId string) (entities.TaskList, error) {
	tasks, err := t.query.GetAllTasks(accountId)
	if err != nil {
		t.log.Infof("Error fetching projects from the database: ", err)
		return nil, errors.New("no tasks found in the database")
	}

	return tasks, nil
}

func (t *taskService) GetTaskById(taskId, accountId string) (*entities.Task, error) {
	task, err := t.query.GetTaskById(taskId, accountId)
	if err != nil {
		t.log.Debugf("Task with projectId %s and accountId %s not found", taskId, accountId)
		return nil, ehand.ErrorTaskNotFound
	}

	return &task, nil
}

func (t *taskService) CreateTask(newTask entities.Task) (*entities.Task, error) {
	created, err := t.query.CreateTask(newTask)

	if err != nil {
		t.log.Infof("Could not create Task: ", err)
		return nil, ehand.ErrorTaskNotCreated
	}

	return &created, nil
}

func (t *taskService) UpdateTask(newTask *entities.Task) (*entities.Task, error) {
	updated, err := t.query.UpdateTask(newTask)
	if err != nil {
		return nil, ehand.ErrorTaskNotUpdated
	}

	return updated, nil
}

func (t *taskService) Exists(accountId, taskId string) bool {
	exists := t.query.Exists(accountId, taskId)
	return exists
}

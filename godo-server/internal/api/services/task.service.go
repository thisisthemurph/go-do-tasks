package services

import (
	"errors"
	"godo/internal/api"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type TaskService interface {
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
		return nil, api.ErrorTaskNotFound
	}

	return &task, nil
}

func (t *taskService) CreateTask(newTask entities.Task) (*entities.Task, error) {
	created, err := t.query.CreateTask(newTask)

	if err != nil {
		t.log.Infof("Could not create Task: ", err)
		return nil, api.ErrorTaskNotCreated
	}

	return &created, nil
}

func (t *taskService) UpdateTask(newTask *entities.Task) (*entities.Task, error) {
	updated, err := t.query.UpdateTask(newTask)
	if err != nil {
		return nil, api.ErrorTaskNotUpdated
	}

	return updated, nil
}

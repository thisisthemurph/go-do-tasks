package repository

import (
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
)

type TaskQuery interface {
	Exists(taskId string) bool
	GetAllTasks(accountId string) (entities.TaskList, error)
	GetTaskById(taskId, accountId string) (entities.Task, error)
	CreateTask(newTask entities.Task) (entities.Task, error)
	UpdateTask(newTask *entities.Task) (*entities.Task, error)
}

type taskQuery struct {
	log ilog.StdLogger
}

func (d *dao) NewTaskQuery(logger ilog.StdLogger) TaskQuery {
	return &taskQuery{log: logger}
}

func (q *taskQuery) GetAllTasks(accountId string) (entities.TaskList, error) {
	q.log.Infof("Fetching all Tasks with accountId %s", accountId)

	var tasks entities.TaskList
	err := Database.
		Preload("Creator", "account_id = ?", accountId).
		Joins("JOIN users on tasks.creator_id = users.id").
		Where("users.account_id = ?", accountId).
		Find(&tasks).
		Error

	ilog.ErrorlnIf(err, q.log)
	return tasks, nil
}

func (q *taskQuery) GetTaskById(taskId, accountId string) (entities.Task, error) {
	q.log.Infof("Fetching Task with id %s", taskId)

	var task entities.Task
	err := Database.
		Preload("Creator", "account_id = ?", accountId).
		Joins("JOIN users on tasks.creator_id = users.id").
		Where("users.account_id = ?", accountId).
		Find(&task, "tasks.id = ?", taskId).
		Error

	ilog.ErrorlnIf(err, q.log)
	return task, err
}

func (q *taskQuery) Exists(taskId string) bool {
	q.log.Infof("Checking if task with Id %s exists", taskId)

	var task entities.Task
	r := Database.First(&task, "id = ?", taskId)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.RowsAffected == 1
}

func (q *taskQuery) CreateTask(newTask entities.Task) (entities.Task, error) {
	q.log.Infof("Creating Task with name %s", newTask.Name)

	r := Database.Create(&newTask)
	ilog.ErrorlnIf(r.Error, q.log)

	return newTask, r.Error
}

func (q *taskQuery) UpdateTask(newTask *entities.Task) (*entities.Task, error) {
	q.log.Infof("Updating task with taskId %s", newTask.ID)

	err := Database.Save(newTask).Error
	if err != nil {
		q.log.Error("Could not update task", err)
		return nil, err
	}

	return newTask, nil
}

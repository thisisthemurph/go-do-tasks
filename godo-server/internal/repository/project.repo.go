package repository

import (
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"time"
)

type ProjectQuery interface {
	GetProjectById(accountId string, storyId string) (*entities.Project, error)
	GetAllProjects(accountId string) ([]*entities.Project, error)
	CreateProject(newProject *entities.Project) (*entities.Project, error)
	UpdateProject(projectId string, newProject *entities.Project) error
	DeleteProject(projectId string) error
	Exists(projectId string) bool
}

type projectQuery struct {
	log ilog.StdLogger
}

func (d *dao) NewProjectQuery(logger ilog.StdLogger) ProjectQuery {
	return &projectQuery{log: logger}
}

func (q *projectQuery) GetAllProjects(accountId string) ([]*entities.Project, error) {
	q.log.Infof("Fetching all project with accountId %s", accountId)

	projects := []*entities.Project{}
	result := Database.
		Preload("Creator", "account_id = ?", accountId).
		Joins("JOIN users on projects.creator_id = users.id").
		Where("users.account_id = ?", accountId).
		Find(&projects)

	ilog.ErrorlnIf(result.Error, q.log)
	return projects, result.Error
}

func (q *projectQuery) GetProjectById(projectId string, accountId string) (*entities.Project, error) {
	q.log.Infof("Fetching project with projectId %s and accountId %s", projectId, accountId)

	project := entities.Project{}
	result := Database.
		Preload("Creator", "account_id = ?", accountId).
		Joins("JOIN users on projects.creator_id = users.id").
		Where("users.account_id = ?", accountId).
		First(&project, "projects.id = ?", projectId)

	ilog.ErrorlnIf(result.Error, q.log)
	return &project, result.Error
}

func (q *projectQuery) CreateProject(newProject *entities.Project) (*entities.Project, error) {
	q.log.Infof("Creating Project with name %v", newProject.Name)

	response := Database.Create(&newProject)
	ilog.ErrorlnIf(response.Error, q.log)

	return newProject, response.Error
}

func (q *projectQuery) UpdateProject(projectId string, newProject *entities.Project) error {
	project, _ := q.getProjectByIdOnly(projectId)

	err := Database.First(&project, "id = ?", projectId).Error
	if err != nil {
		q.log.Errorf("Unable to fetch project (%v) from database", projectId)
		return err
	}

	project.Name = newProject.Name
	project.Description = newProject.Description
	project.UpdatedAt = time.Now()
	result := Database.Save(&project)

	ilog.ErrorlnIf(result.Error, q.log)
	return result.Error
}

func (q *projectQuery) DeleteProject(projectId string) error {
	q.log.Infof("Deleting project with ID %v", projectId)

	deletedProject := &entities.Project{}
	result := Database.Where("id = ?", projectId).Delete(deletedProject)
	return result.Error
}

func (q *projectQuery) Exists(projectId string) bool {
	q.log.Info("Checking if project with ID %v exists", projectId)

	project := &entities.Project{}
	r := Database.First(project, "id = ?", projectId)

	q.log.Println(r.RowsAffected)

	return r.RowsAffected == 1
}

func (q *projectQuery) getProjectByIdOnly(projectId string) (*entities.Project, error) {
	q.log.Infof("Fetching project with projectId %s", projectId)

	project := entities.Project{}
	result := Database.First(&project, "id = ?", projectId)

	ilog.ErrorlnIf(result.Error, q.log)
	return &project, result.Error
}

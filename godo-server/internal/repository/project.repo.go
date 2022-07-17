package repository

import (
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"time"
)

type ProjectQuery interface {
	GetProjectById(storyId string) (*entities.Project, error)
	GetAllProjects(accountId string) ([]*entities.Project, error)
	CreateProject(newProject *entities.Project) error
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
	q.log.Info("Fetching all project")

	projects := []*entities.Project{}
	result := Database.Joins("Creator", Database.Where(&entities.User{AccountId: accountId})).
		Find(&projects)

	ilog.ErrorlnIf(result.Error, q.log)

	return projects, result.Error
}

func (q *projectQuery) GetProjectById(projectId string) (*entities.Project, error) {
	q.log.Infof("Fetching project with ID %v", projectId)

	project := entities.Project{}
	err := Database.First(&project, "id = ?", projectId).Error
	ilog.ErrorlnIf(err, q.log)

	return &project, err
}

func (q *projectQuery) CreateProject(newProject *entities.Project) error {
	q.log.Infof("Creating Project with name %v", newProject.Name)

	err := Database.Create(&newProject).Error
	ilog.ErrorlnIf(err, q.log)

	return err
}

func (q *projectQuery) UpdateProject(projectId string, newProject *entities.Project) error {
	project, _ := q.GetProjectById(projectId)

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

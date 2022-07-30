package repository

import (
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"time"
)

type ProjectQuery interface {
	GetProjectById(projectId string, accountId string) (*entities.Project, error)
	GetProjectsInfo(accountId string) (entities.ProjectInfoList, error)
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

func (q *projectQuery) GetProjectsInfo(accountId string) (entities.ProjectInfoList, error) {
	q.log.Debugf("Fetching all project information associated with Account{id=%s}", accountId)

	var info entities.ProjectInfoList
	r := Database.
		Table("projects").
		Select("projects.id, projects.name, projects.description, projects.status, "+
			"projects.created_at, projects.updated_at, count(stories.id) story_count, count(tags.id) tag_count").
		Joins("LEFT JOIN users ON users.account_id = ?", accountId).
		Joins("LEFT JOIN stories ON projects.id = stories.project_id").
		Joins("LEFT JOIN tags ON projects.id = tags.project_id").
		Where("users.account_id = ?", accountId).
		Group("projects.id").
		Order("projects.created_at DESC").
		Find(&info)

	ilog.ErrorlnIf(r.Error, q.log)
	return info, r.Error
}

func (q *projectQuery) GetProjectById(projectId string, accountId string) (*entities.Project, error) {
	q.log.Debugf("Fetching project with projectId %s and accountId %s", projectId, accountId)

	var project entities.Project
	result := Database.
		Preload("Creator", "account_id = ?", accountId).
		Preload("Stories", "stories.project_id = ?", projectId).
		Preload("Stories.Creator").
		Preload("Stories.Tasks").
		Preload("Stories.Tasks.Creator").
		Preload("Stories.Tasks.Tags").
		Preload("Tags").
		Joins("JOIN users on projects.creator_id = users.id").
		Where("users.account_id = ?", accountId).
		First(&project, "projects.id = ?", projectId)

	ilog.ErrorlnIf(result.Error, q.log)
	return &project, result.Error
}

func (q *projectQuery) CreateProject(newProject *entities.Project) (*entities.Project, error) {
	q.log.Debugf("Creating Project with name %v", newProject.Name)

	response := Database.Create(&newProject)
	ilog.ErrorlnIf(response.Error, q.log)

	return newProject, response.Error
}

// UpdateProject TODO: Remove duplicate query
func (q *projectQuery) UpdateProject(projectId string, newProject *entities.Project) error {
	q.log.Debugf("Updating Project{id=%s}", projectId)

	project, _ := q.getProjectByIdOnly(projectId)
	err := Database.First(&project, "id = ?", projectId).Error
	if err != nil {
		q.log.Errorf("Unable to fetch project (%v) from database", projectId)
		return err
	}

	// TODO: Place override of props into model?
	project.Name = newProject.Name
	project.Description = newProject.Description
	project.UpdatedAt = time.Now()
	result := Database.Save(&project)

	ilog.ErrorlnIf(result.Error, q.log)
	return result.Error
}

func (q *projectQuery) DeleteProject(projectId string) error {
	q.log.Debugf("Deleting Project{id=%s}", projectId)

	deletedProject := &entities.Project{}
	result := Database.Where("id = ?", projectId).Delete(deletedProject)
	return result.Error
}

func (q *projectQuery) Exists(projectId string) bool {
	q.log.Debugf("Checking if Project{id=%s} exists", projectId)

	project := &entities.Project{}
	r := Database.First(project, "id = ?", projectId)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.RowsAffected == 1
}

func (q *projectQuery) getProjectByIdOnly(projectId string) (*entities.Project, error) {
	q.log.Debugf("Fetching Project{id=%s}", projectId)

	project := entities.Project{}
	result := Database.First(&project, "id = ?", projectId)

	ilog.ErrorlnIf(result.Error, q.log)
	return &project, result.Error
}

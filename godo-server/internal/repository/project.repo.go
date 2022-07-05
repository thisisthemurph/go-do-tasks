package repository

import (
	"godo/internal/repository/entities"
	"time"
)

type ProjectQuery interface {
	GetProjectById(storyId string) (*entities.Project, error)
	GetAllProjects() ([]*entities.Project, error)
	CreateProject(newProject *entities.Project) error
	UpdateProject(projectId string, newProject *entities.Project) error
}

type projectQuery struct{}

func (d *dao) NewProjectQuery() ProjectQuery {
	return &projectQuery{}
}

func (p *projectQuery) GetAllProjects() ([]*entities.Project, error) {
	projects := []*entities.Project{}
	err := Database.Find(&projects).Error

	return projects, err
}

func (p *projectQuery) GetProjectById(projectId string) (*entities.Project, error) {
	project := entities.Project{}
	err := Database.First(&project, "id = ?", projectId).Error

	return &project, err
}

func (p *projectQuery) CreateProject(newProject *entities.Project) error {
	err := Database.Create(&newProject).Error
	return err
}

func (p *projectQuery) UpdateProject(projectId string, newProject *entities.Project) error {
	project, _ := p.GetProjectById(projectId)

	err := Database.First(&project, "id = ?", projectId).Error
	if err != nil {
		return err
	}

	project.Name = newProject.Name
	project.Description = newProject.Description
	project.UpdatedAt = time.Now()
	Database.Save(&project)

	return nil
}

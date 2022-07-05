package repository

import (
	"godo/internal/repository/entities"
)

type ProjectQuery interface {
	GetProjectById(storyId string) (*entities.Project, error)
	GetAllProjects() ([]*entities.Project, error)
}

type projectQuery struct{}

func (d *dao) NewProjectQuery() ProjectQuery {
	return &projectQuery{}
}

func (s *projectQuery) GetAllProjects() ([]*entities.Project, error) {
	projects := []*entities.Project{}
	err := Database.Find(&projects).Error

	return projects, err
}

func (s *projectQuery) GetProjectById(projectId string) (*entities.Project, error) {
	project := entities.Project{}
	err := Database.First(&project, "id = ?", projectId).Error

	return &project, err
}
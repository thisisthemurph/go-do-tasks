package repository

import (
	"godo/internal/repository/entities"
)

type ProjectQuery interface {
	GetProjectById(storyId string) (*entities.Project, error)
	GetAllProjects() ([]*entities.Project, error)
}

type projectQuery struct{}

func (s *projectQuery) GetAllProjects() ([]*entities.Project, error) {
	db := GetDatabase()
	projects := []*entities.Project{}
	err := db.Find(&projects).Error

	return projects, err
}

func (s *projectQuery) GetProjectById(projectId string) (*entities.Project, error) {
	db := GetDatabase()
	project := entities.Project{}
	err := db.First(&project, "id = ?", projectId).Error

	return &project, err
}
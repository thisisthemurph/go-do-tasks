package repository

import (
	"godo/internal/repository/entities"
	"log"
	"time"
)

type ProjectQuery interface {
	GetProjectById(storyId string) (*entities.Project, error)
	GetAllProjects() ([]*entities.Project, error)
	CreateProject(newProject *entities.Project) error
	UpdateProject(projectId string, newProject *entities.Project) error
	DeleteProject(projectId string) error
	Exists(projectId string) bool
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

func (p *projectQuery) DeleteProject(projectId string) error {
	deletedProject := &entities.Project{}
	return Database.
		Where("id = ?", projectId).
		Delete(deletedProject).
		Error
}

func (p *projectQuery) Exists(projectId string) bool {
	log.Printf("Exists(%v)\n", projectId)

	project := &entities.Project{}
	r := Database.First(project, "id = ?", projectId)

	log.Println(r.RowsAffected)

	return r.RowsAffected == 1
}

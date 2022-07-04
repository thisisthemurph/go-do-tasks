package entities

import (
	"fmt"
)

type Story struct {
	Base
	ProjectId	string
	Project		Project

	Name		string
	Description	string

	CreatorId	string
	Creator		Person `gorm:"foreignKey:CreatorId"`
	Tasks 		[]Task
}

func (s Story) ToString() string {
	return fmt.Sprintf("Task{Name=%v}", s.Name)
}
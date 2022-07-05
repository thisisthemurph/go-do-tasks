package entities

import (
	"fmt"
)

type Story struct {
	Base
	ProjectId string
	Project   Project `json:"-"`

	Name        string `json:"name"`
	Description string `json:"description"`

	CreatorId string `json:"-"`
	Creator   Person `gorm:"foreignKey:CreatorId"`
	Tasks     []Task
}

func (s Story) ToString() string {
	return fmt.Sprintf("Task{Name=%v}", s.Name)
}

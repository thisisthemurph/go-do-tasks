package entities

import (
	"fmt"
	"godo/internal/auth"
)

type Story struct {
	Base

	ProjectId   string
	Project     Project   `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatorId   string    `json:"-"`
	Creator     auth.User `gorm:"foreignKey:CreatorId"`
	Tasks       []Task

	TimestampBase
}

func (s Story) ToString() string {
	return fmt.Sprintf("Task{Name=%v}", s.Name)
}

package entities

import "godo/internal/auth"

type Task struct {
	Base

	StoryId     string
	Story       Story  `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int
	Status      int
	CreatorId   string    `json:"-"`
	Creator     auth.User `gorm:"foreignKey:CreatorId"`
	Tags        []Tag     `gorm:"many2many:task_tags;"`
}

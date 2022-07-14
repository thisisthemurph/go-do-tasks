package entities

import (
	"godo/internal/auth"

	"github.com/go-playground/validator"
)

type Task struct {
	Base

	StoryId     string    `json:"story_id"`
	Story       Story     `json:"-"`
	Name        string    `json:"name" gorm:"not null" validate:"required,min=1,max=40"`
	Description string    `json:"description" validate:"max=280"`
	Type        int       `json:"type" gorm:"default:0,not null"`
	Status      int       `json:"status" gorm:"default:0,not null"`
	CreatorId   string    `json:"-"`
	Creator     auth.User `gorm:"foreignKey:CreatorId"`

	TimestampBase
}

func (t *Task) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

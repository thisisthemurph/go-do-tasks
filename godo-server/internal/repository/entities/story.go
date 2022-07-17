package entities

import (
	"fmt"

	"github.com/go-playground/validator"
)

type Story struct {
	Base

	ProjectId   string  `json:"project_id"`
	Project     Project `json:"-"`
	Name        string  `json:"name" validate:"required,min=1,max=40"`
	Description string  `json:"description" validate:"max=280"`
	CreatorId   string  `json:"-"`
	Creator     User    `json:"creator" gorm:"foreignKey:CreatorId"`
	Tasks       []Task

	TimestampBase
}

func (s Story) ToString() string {
	return fmt.Sprintf("Story{Name=%v}", s.Name)
}

func (s *Story) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

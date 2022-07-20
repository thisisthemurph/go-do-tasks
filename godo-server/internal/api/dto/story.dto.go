package dto

import (
	"github.com/go-playground/validator"
)

type NewStoryDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ProjectId   string `json:"project_id" validate:"required,uuid"`
}

func (s *NewStoryDto) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

package dto

import "github.com/go-playground/validator"

type NewProjectDto struct {
	Name        string `json:"name" validation:"required"`
	Description string `json:"description" validation:"required"`
}

func (np *NewProjectDto) Validate() error {
	validate := validator.New()
	return validate.Struct(np)
}

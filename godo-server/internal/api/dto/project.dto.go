package dto

import (
	"github.com/go-playground/validator"
	"godo/internal/repository/enums"
)

type NewProjectDto struct {
	Name        string `json:"name" validation:"required"`
	Description string `json:"description" validation:"required"`
}

type ProjectStatusUpdateDto struct {
	Status enums.ProjectStatus `json:"status" validate:"required,gte=0"`
}

// Validate TODO: Generalise the validation method
func (np *NewProjectDto) Validate() error {
	validate := validator.New()
	return validate.Struct(np)
}

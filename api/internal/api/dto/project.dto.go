package dto

import (
	"godo/internal/repository/enums"
)

type NewProjectDto struct {
	Name        string `json:"name" validation:"required"`
	Description string `json:"description" validation:"required"`
}

type ProjectStatusUpdateDto struct {
	Status enums.ProjectStatus `json:"status" validate:"required,gte=0"`
}

package dto

import (
	"godo/internal/repository/enums"
)

// NewProjectDto model for creating a new project
// swagger:model newTagDto
type NewProjectDto struct {
	// the name of the project
	//
	// required: true
	Name string `json:"name" validation:"required"`

	// the description of the project
	//
	// required: true
	Description string `json:"description" validation:"required"`
}

// UpdateProjectDto model for updating the project
// swagger:model newTagDto
type UpdateProjectDto struct {
	// the name of the project
	//
	// required: false
	Name string `json:"name" validate:"min=1,max=40"`

	// the description of the project
	//
	// required: false
	Description string `json:"description" validate:"max=280"`

	// numeric representation of the project status
	//
	// required: false
	// min: 0
	Status enums.ProjectStatus `json:"status" validate:"gte=0"`
}

// ProjectStatusUpdateDto model for updating the status of the project
// swagger:model newTagDto
type ProjectStatusUpdateDto struct {
	// numeric representation of the project status
	//
	// required: false
	// min: 0
	Status enums.ProjectStatus `json:"status" validate:"required,gte=0"`
}

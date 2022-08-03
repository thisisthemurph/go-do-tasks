package dto

// NewTagDto model for creating a new tag
// swagger:model newTagDto
type NewTagDto struct {
	// the name associated with the tag
	//
	// required: true
	// min length: 1
	// max length: 16
	Name string `json:"name" validate:"required,min=1,max=16"`
}

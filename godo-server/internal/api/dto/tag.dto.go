package dto

type NewTagDto struct {
	Name string `json:"name" validate:"required,min=1,max=16"`
}

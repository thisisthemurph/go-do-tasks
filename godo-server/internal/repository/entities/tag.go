package entities

// Tag - A Project has any number of tags associated with it
// A subset of these tgs can be assigned to any of the tasks
// that form part of the Project
type Tag struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	Name      string  `json:"name" validate:"required,min=1,max=16"`
	ProjectId string  `json:"-"`
	Project   Project `json:"-"`
}

type TagList []Tag

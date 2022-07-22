package entities

import (
	"fmt"
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

func (s Story) String() string {
	return fmt.Sprintf("Story{Name=%v}", s.Name)
}

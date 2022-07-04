package entities

import (
	"fmt"
)

type Project struct {
	Base

	Name		string
	Description	string

	CreatorId	string
	Creator		Person `gorm:"foreignKey:CreatorId"`
	Stories		[]Story
	Tags		[]Tag
}

func (p Project) ToString() string {
	return fmt.Sprintf("Project{Name=%v}", p.Name)
}
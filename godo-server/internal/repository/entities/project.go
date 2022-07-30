package entities

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"godo/internal/repository/enums"
)

type Project struct {
	Base

	Name        string              `json:"name" gorm:"not null" validate:"required,min=1,max=40"`
	Description string              `json:"description" validate:"max=280"`
	Status      enums.ProjectStatus `json:"-" gorm:"type:smallint;default:0;not null" validate:"gte=0"`
	StatusValue string              `json:"status" gorm:"-:all"`
	CreatorId   uint                `json:"creator_id"`
	Creator     User                `json:"creator" gorm:"foreignKey:CreatorId"`
	Stories     []Story             `json:"stories"`
	Tags        []Tag               `json:"tags"`

	TimestampBase
}

type ProjectList []*Project
type ProjectKey struct{}

func (s *Project) ToString() string {
	return fmt.Sprintf("Project{Name=%v}", s.Name)
}

func (s *Project) AfterFind(tx *gorm.DB) {
	s.StatusValue = s.Status.String()
}

func (s *Project) AfterCreate(tx *gorm.DB) {
	s.StatusValue = s.Status.String()
}

func (s *Project) AfterSave(tx *gorm.DB) {
	s.StatusValue = s.Status.String()
}

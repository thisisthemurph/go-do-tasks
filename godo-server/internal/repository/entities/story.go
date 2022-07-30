package entities

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"godo/internal/repository/enums"
)

type Story struct {
	Base

	Name        string               `json:"name" gorm:"not null" validate:"required,min=1,max=40"`
	Description string               `json:"description" validate:"max=280"`
	Status      enums.ProgressStatus `json:"-" gorm:"type:smallint;default:0;not null"`
	StatusValue string               `json:"status" gorm:"-:all"`
	ProjectId   string               `json:"project_id"`
	Project     Project              `json:"-"`
	CreatorId   uint                 `json:"creator_id"`
	Creator     User                 `json:"creator" gorm:"foreignKey:CreatorId"`
	Tasks       []Task               `json:"tasks"`

	TimestampBase
}

func (s Story) String() string {
	return fmt.Sprintf("Story{Name=%v}", s.Name)
}

func (s *Story) AfterFind(tx *gorm.DB) {
	s.StatusValue = s.Status.String()
}

func (s *Story) AfterCreate(tx *gorm.DB) {
	s.StatusValue = s.Status.String()
}

func (s *Story) AfterSave(tx *gorm.DB) {
	s.StatusValue = s.Status.String()
}

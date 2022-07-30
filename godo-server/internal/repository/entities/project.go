package entities

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"godo/internal/repository/enums"
	"time"
)

type Project struct {
	Base

	Name        string              `json:"name" gorm:"not null" validate:"required,min=1,max=40"`
	Description string              `json:"description" validate:"max=280"`
	Status      enums.ProjectStatus `json:"-" gorm:"type:smallint;default:0;not null" validate:"gte=0"`
	StatusValue string              `json:"status" gorm:"-:all"`
	CreatorId   uint                `json:"creator_id"`
	Creator     User                `json:"creator" gorm:"foreignKey:CreatorId"`
	Stories     []Story             `json:"stories,omitempty"`
	Tags        []Tag               `json:"tags,omitempty"`

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

type ProjectInfo struct {
	Base

	Name         string              `json:"name" validate:"required,min=1,max=40"`
	Description  string              `json:"description" validate:"max=280"`
	Status       enums.ProjectStatus `json:"-" validate:"gte=0"`
	StatusValue  string              `json:"status"`
	StoriesCount uint16              `json:"stories_count"`
	TagsCount    uint16              `json:"tags_count"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type ProjectInfoList []*ProjectInfo

func (s *ProjectInfo) AfterFind(tx *gorm.DB) {
	s.StatusValue = s.Status.String()
}

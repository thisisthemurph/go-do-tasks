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

	Name        string              `json:"name" validate:"required,min=1,max=40"`
	Description string              `json:"description" validate:"max=280"`
	Status      enums.ProjectStatus `json:"-" validate:"gte=0"`
	StatusValue string              `json:"status"`
	StoryCount  uint16              `json:"story_count"`
	TagCount    uint16              `json:"tag_count"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type ProjectInfoList []*ProjectInfo

func (s *ProjectInfo) AfterFind(tx *gorm.DB) {
	s.StatusValue = s.Status.String()
}

// ProjectResponse the specified Project
// swagger:response projectResponse
type ProjectResponse struct {
	// The resultant Project
	// in: body
	Body Project
}

// ProjectInfoResponse a list of project information associated with the authenticated account
// swagger:response projectInfoResponse
type ProjectInfoResponse struct {
	// All projects associated with the authenticated account
	// in: body
	Body ProjectInfoList
}

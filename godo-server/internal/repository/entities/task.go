package entities

import (
	"github.com/jinzhu/gorm"
	"godo/internal/repository/enums"
)

type Task struct {
	Base

	Name        string               `json:"name" gorm:"not null" validate:"required,min=1,max=40"`
	Description string               `json:"description" validate:"max=280"`
	Type        enums.TaskType       `json:"-" gorm:"type:smallint;default:0;not null"`
	TypeValue   string               `json:"type" gorm:"-:all"`
	Status      enums.ProgressStatus `json:"-" gorm:"type:smallint;default:0;not null"`
	StatusValue string               `json:"status" gorm:"-:all"`
	StoryId     string               `json:"story_id"`
	Story       Story                `json:"-"`
	CreatorId   string               `json:"-"`
	Creator     User                 `gorm:"foreignKey:CreatorId"`

	TimestampBase
}

type TaskList []*Task

func (t Task) AfterFind(tx *gorm.DB) {
	t.TypeValue = t.Type.String()
	t.StatusValue = t.Status.String()
}

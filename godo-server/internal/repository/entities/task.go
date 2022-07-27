package entities

import (
	"github.com/jinzhu/gorm"
	"godo/internal/repository/enums"
)

type Task struct {
	Base

	Name        string               `json:"name" gorm:"not null" validate:"required,min=1,max=40"`
	Description string               `json:"description" validate:"max=280"`
	Type        enums.TaskType       `json:"-" gorm:"type:smallint;default:0;not null" validate:"gte=0"`
	TypeValue   string               `json:"type" gorm:"-:all"`
	Status      enums.ProgressStatus `json:"-" gorm:"type:smallint;default:0;not null" validate:"gte=0"`
	StatusValue string               `json:"status" gorm:"-:all"`
	StoryId     string               `json:"story_id"`
	Story       Story                `json:"-"`
	CreatorId   uint                 `json:"-"`
	Creator     User                 `json:"creator" gorm:"foreignKey:CreatorId"`
	Tags        []Tag                `json:"tags" gorm:"many2many:task_tags"`

	TimestampBase
}

type TaskList []*Task

func (t *Task) AfterFind(tx *gorm.DB) {
	t.TypeValue = t.Type.String()
	t.StatusValue = t.Status.String()
}

func (t *Task) AfterCreate(tx *gorm.DB) {
	t.TypeValue = t.Type.String()
	t.StatusValue = t.Status.String()
}

func (t *Task) AfterSave(tx *gorm.DB) {
	t.TypeValue = t.Type.String()
	t.StatusValue = t.Status.String()
}

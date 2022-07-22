package dto

import (
	"godo/internal/repository/enums"
)

type NewTaskDto struct {
	Name        string         `json:"name" validate:"required,min=1,max=40"`
	Description string         `json:"description"`
	Type        enums.TaskType `json:"type"`
	StoryId     string         `json:"story_id" validate:"required,uuid"`
}

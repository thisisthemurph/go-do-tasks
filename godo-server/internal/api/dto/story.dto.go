package dto

import (
	"godo/internal/repository/entities"
)

type StoryDto struct {
	ID   string
	Name string
}

func NewStoryDto(s entities.Story) *StoryDto {
	return &StoryDto{ID: s.ID, Name: s.Name}
}
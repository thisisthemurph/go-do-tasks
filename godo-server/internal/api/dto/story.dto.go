package dto

import (
	"godo/internal/repository/entities"
)

type StoryDto struct {
	Id   uint
	Name string
}

func NewStoryDto(s entities.Story) *StoryDto {
	return &StoryDto{Id: s.ID, Name: s.Name}
}
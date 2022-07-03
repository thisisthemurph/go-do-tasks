package repository

import "godo/internal/repository/entities"

type StoryQuery interface {
	GetStoryById(storyId uint) (*entities.Story, error)
	GetAllStories() ([]*entities.Story, error)
}

type storyQuery struct{}

func (s *storyQuery) GetAllStories() ([]*entities.Story, error) {
	db := GetDatabase()
	stories := []*entities.Story{}
	result := db.Find(&stories)

	return stories, result.Error
}

func (s *storyQuery) GetStoryById(storyId uint) (*entities.Story, error) {
	db := GetDatabase()
	story := entities.Story{}
	result := db.First(&story, 1)

	return &story, result.Error
}
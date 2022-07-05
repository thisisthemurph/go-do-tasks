package repository

import "godo/internal/repository/entities"

type StoryQuery interface {
	GetStoryById(storyId string) (*entities.Story, error)
	GetAllStories() ([]*entities.Story, error)
}

type storyQuery struct{}

func (d *dao) NewStoryQuery() StoryQuery {
	return &storyQuery{}
}

func (s *storyQuery) GetAllStories() ([]*entities.Story, error) {
	stories := []*entities.Story{}
	result := Database.Find(&stories)

	return stories, result.Error
}

func (s *storyQuery) GetStoryById(storyId string) (*entities.Story, error) {
	story := entities.Story{}
	result := Database.First(&story, 1)

	return &story, result.Error
}
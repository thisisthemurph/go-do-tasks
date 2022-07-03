package services

import (
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type StoryService interface {
	GetStories()				([]*entities.Story, error)
	GetStoryById(userId uint)	(*entities.Story, error)
}

type storyService struct{}

func NewStoryService() StoryService {
	return &storyService{}
}

func (s *storyService) GetStories() ([]*entities.Story, error) {
	dao := repository.NewDAO()
	query := dao.NewStoryQuery()
	stories, err := query.GetAllStories()

	return stories, err
}

func (s *storyService) GetStoryById(userId uint) (*entities.Story, error) {
	dao := repository.NewDAO()
	query := dao.NewStoryQuery()
	story, err := query.GetStoryById(userId)

	return story, err
}
package services

import (
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type StoryService interface {
	GetStories()				([]*entities.Story, error)
	GetStoryById(userId string)	(*entities.Story, error)
}

type storyService struct{
	dao repository.DAO
	query repository.StoryQuery
}

func NewStoryService(dao repository.DAO) StoryService {
	return &storyService{
		dao: dao, 
		query: dao.NewStoryQuery(),
	}
}

func (s *storyService) GetStories() ([]*entities.Story, error) {
	stories, err := s.query.GetAllStories()

	return stories, err
}

func (s *storyService) GetStoryById(userId string) (*entities.Story, error) {
	story, err := s.query.GetStoryById(userId)

	return story, err
}
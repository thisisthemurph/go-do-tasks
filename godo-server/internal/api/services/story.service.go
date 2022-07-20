package services

import (
	"errors"
	"godo/internal/api"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type StoryService interface {
	GetStories(accountId string) ([]*entities.Story, error)
	GetStoryById(accountId, userId string) (*entities.Story, error)
}

type storyService struct {
	log   ilog.StdLogger
	query repository.StoryQuery
}

func NewStoryService(storyQuery repository.StoryQuery, log ilog.StdLogger) StoryService {
	return &storyService{
		log:   log,
		query: storyQuery,
	}
}

func (s *storyService) GetStories(accountId string) ([]*entities.Story, error) {
	stories, err := s.query.GetAllStories(accountId)
	if err != nil {
		s.log.Infof("Error fetching stories from database: ", err.Error())
		return nil, errors.New("no stories present in the database")
	}

	return stories, nil
}

func (s *storyService) GetStoryById(accountId, storyId string) (*entities.Story, error) {
	story, err := s.query.GetStoryById(accountId, storyId)
	if err != nil {
		s.log.Infof("Story with accountId %s and storyId %s not found", accountId, storyId)
		return nil, api.ErrorStoryNotFound
	}

	return story, err
}

package services

import (
	"errors"
	ehand "godo/internal/api/errorhandler"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type StoryService interface {
	GetStories(accountId string) ([]*entities.Story, error)
	GetStoryById(accountId, storyId string) (*entities.Story, error)
	CreateStory(newStory *entities.Story) (*entities.Story, error)
	UpdateStory(storyId string, newStoryData *entities.Story) error
	DeleteStory(storyId string) error
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
		return nil, ehand.ErrorStoryNotFound
	}

	return story, err
}

func (s *storyService) CreateStory(newStory *entities.Story) (*entities.Story, error) {
	createdStory, err := s.query.CreateStory(newStory)

	if err != nil {
		s.log.Info("Error creating Story.", err)
		return nil, ehand.ErrorStoryNotCreated
	}

	return createdStory, nil
}

func (s *storyService) Exists(storyId string) bool {
	return s.query.Exists(storyId)
}

func (s *storyService) UpdateStory(storyId string, newStoryData *entities.Story) error {
	exists := s.Exists(storyId)
	if !exists {
		return ehand.ErrorStoryNotFound
	}

	err := s.query.UpdateStory(newStoryData)
	if err != nil {
		s.log.Errorf("Could not update Story with storyId %s: %S", storyId, err.Error())
		return ehand.ErrorStoryNotUpdated
	}

	return nil
}

func (s *storyService) DeleteStory(storyId string) error {
	exists := s.Exists(storyId)
	if !exists {
		return ehand.ErrorStoryNotFound
	}

	err := s.query.DeleteStory(storyId)
	if err != nil {
		s.log.Errorf("Could not delete Story with storyId %s: %S", storyId, err.Error())
		return ehand.ErrorStoryNotDeleted
	}

	return nil
}

package services

import (
	ehand "godo/internal/api/errorhandler"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type StoryService interface {
	GetStoriesInfo(accountId string) (entities.StoryInfoList, error)
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

func (s *storyService) GetStoriesInfo(accountId string) (entities.StoryInfoList, error) {
	info, err := s.query.GetStoriesInfo(accountId)
	if err != nil {
		s.log.Error("error fetching info from database: ", err)
		return nil, err
	}

	return info, nil
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

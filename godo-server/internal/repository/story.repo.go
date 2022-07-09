package repository

import (
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
)

type StoryQuery interface {
	GetStoryById(storyId string) (*entities.Story, error)
	GetAllStories() ([]*entities.Story, error)
}

type storyQuery struct {
	log ilog.StdLogger
}

func (d *dao) NewStoryQuery(logger ilog.StdLogger) StoryQuery {
	return &storyQuery{
		log: logger,
	}
}

func (q *storyQuery) GetAllStories() ([]*entities.Story, error) {
	stories := []*entities.Story{}
	result := Database.Find(&stories)
	ilog.ErrorlnIf(result.Error, q.log)

	return stories, result.Error
}

func (q *storyQuery) GetStoryById(storyId string) (*entities.Story, error) {
	story := entities.Story{}
	result := Database.First(&story, 1)

	return &story, result.Error
}

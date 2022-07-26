package repository

import (
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
)

type StoryQuery interface {
	CreateStory(newStory *entities.Story) (*entities.Story, error)
	DeleteStory(storyId string) error
	Exists(storyId string) bool
	GetAllStories(accountId string) ([]*entities.Story, error)
	GetStoryById(accountId, storyId string) (*entities.Story, error)
	UpdateStory(newStory *entities.Story) error
}

type storyQuery struct {
	log ilog.StdLogger
}

func (d *dao) NewStoryQuery(logger ilog.StdLogger) StoryQuery {
	return &storyQuery{
		log: logger,
	}
}

func (q *storyQuery) GetAllStories(accountId string) ([]*entities.Story, error) {
	q.log.Debugf("Fetching all stories with accountId %s", accountId)

	var stories []*entities.Story
	result := Database.
		Preload("Creator", "account_id = ?", accountId).
		Joins("JOIN users on stories.creator_id = users.id").
		Where("users.account_id = ?", accountId).
		Find(&stories)

	ilog.ErrorlnIf(result.Error, q.log)
	return stories, result.Error
}

func (q *storyQuery) GetStoryById(accountId, storyId string) (*entities.Story, error) {
	q.log.Debugf("Fetching story with accountId %s and storyId %s", accountId, storyId)

	story := entities.Story{}
	result := Database.
		Preload("Creator", "account_id = ?", accountId).
		Joins("JOIN users on stories.creator_id = users.id").
		Where("users.account_id = ?", accountId).
		First(&story, "stories.id = ?", storyId)

	ilog.ErrorlnIf(result.Error, q.log)
	return &story, result.Error
}

func (q *storyQuery) CreateStory(newStory *entities.Story) (*entities.Story, error) {
	q.log.Debugf("Creating story with name %s", newStory.Name)

	result := Database.Create(&newStory)
	ilog.ErrorlnIf(result.Error, q.log)

	return newStory, result.Error
}

func (q *storyQuery) Exists(storyId string) bool {
	q.log.Debugf("Checking if story with storyId %s exists", storyId)

	var story entities.Story
	r := Database.First(&story, "id = ?", storyId)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.RowsAffected == 1
}

func (q *storyQuery) UpdateStory(story *entities.Story) error {
	q.log.Debugf("Updating story with storyId %s", story.ID)

	r := Database.Save(&story)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.Error
}

func (q *storyQuery) DeleteStory(storyId string) error {
	q.log.Infof("Deleting story with storyId %s", storyId)

	var deletedStory entities.Story
	r := Database.Where("id = ?", storyId).Delete(&deletedStory)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.Error
}

func (q *storyQuery) getStoryByIdOnly(storyId string) (*entities.Story, error) {
	q.log.Infof("Fetching story with storyId %s", storyId)

	var story entities.Story
	r := Database.First(&story, "id = ?", storyId)

	ilog.ErrorlnIf(r.Error, q.log)
	return &story, r.Error
}

package repository

import (
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
)

type StoryQuery interface {
	CreateStory(newStory *entities.Story) (*entities.Story, error)
	DeleteStory(storyId string) error
	Exists(storyId string) bool
	GetStoriesInfo(accountId string) (entities.StoryInfoList, error)
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

func (q *storyQuery) GetStoriesInfo(accountId string) (entities.StoryInfoList, error) {
	q.log.Debugf("Fetching all story info for Account{id=%s}", accountId)

	var info entities.StoryInfoList
	r := Database.
		Table("stories").
		Select("stories.id, stories.name, stories.description, stories.status, stories.created_at, "+
			"stories.updated_at, count(tasks.id) task_count").
		Joins("LEFT JOIN users ON users.account_id = ?", accountId).
		Joins("LEFT JOIN tasks ON stories.id = tasks.story_id").
		Where("users.account_id = ?", accountId).
		Group("stories.id").
		Order("stories.created_at DESC").
		Find(&info)

	ilog.ErrorlnIf(r.Error, q.log)
	return info, r.Error
}

func (q *storyQuery) GetStoryById(accountId, storyId string) (*entities.Story, error) {
	q.log.Debugf("Fetching story with Account{id=%s} & Story{id=%s}", accountId, storyId)

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
	q.log.Debugf("Creating Story{name=%s}", newStory.Name)

	result := Database.Create(&newStory)
	ilog.ErrorlnIf(result.Error, q.log)

	return newStory, result.Error
}

func (q *storyQuery) Exists(storyId string) bool {
	q.log.Debugf("Checking if Story{id=%s} exists", storyId)

	var story entities.Story
	r := Database.First(&story, "id = ?", storyId)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.RowsAffected == 1
}

func (q *storyQuery) UpdateStory(story *entities.Story) error {
	q.log.Debugf("Updating Story{id=%s}", story.ID)

	r := Database.Save(&story)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.Error
}

func (q *storyQuery) DeleteStory(storyId string) error {
	q.log.Debugf("Deleting Story{id=%s}", storyId)

	var deletedStory entities.Story
	r := Database.Where("id = ?", storyId).Delete(&deletedStory)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.Error
}

func (q *storyQuery) getStoryByIdOnly(storyId string) (*entities.Story, error) {
	q.log.Debugf("Fetching Story{id=%s}", storyId)

	var story entities.Story
	r := Database.First(&story, "id = ?", storyId)

	ilog.ErrorlnIf(r.Error, q.log)
	return &story, r.Error
}

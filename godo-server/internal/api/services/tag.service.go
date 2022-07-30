package services

import (
	"errors"
	ehand "godo/internal/api/errorhandler"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type TagService interface {
	CreateTag(newTag entities.Tag) (*entities.Tag, error)
	DeleteTag(tagId uint, projectId string) (*entities.Tag, error)
	GetTagById(tagId uint, projectId string) (*entities.Tag, error)
	UpdateTag(newTag entities.Tag) (*entities.Tag, error)
	AddToTask(tagId uint, taskId string) error
	RemoveFromTask(taskId string, tagId uint) error
}

type tagService struct {
	log   ilog.StdLogger
	query repository.TagQuery
}

func NewTagService(query repository.TagQuery, logger ilog.StdLogger) TagService {
	return &tagService{log: logger, query: query}
}

func (t *tagService) GetTagById(tagId uint, projectId string) (*entities.Tag, error) {
	tag, err := t.query.GetTagById(tagId, projectId)
	if err != nil {
		t.log.Debugf("Tag with tagId %d not found", tagId)
		return nil, ehand.ErrorTagNotFound
	}

	return tag, nil
}

func (t *tagService) CreateTag(newTag entities.Tag) (*entities.Tag, error) {
	created, err := t.query.CreateTag(newTag)

	if err != nil {
		t.log.Infof("Could not create Tag: ", err)
		return nil, ehand.ErrorTagNotCreated
	}

	return created, nil
}

func (t *tagService) UpdateTag(newTag entities.Tag) (*entities.Tag, error) {
	updated, err := t.query.UpdateTag(newTag)
	if err != nil {
		return nil, ehand.ErrorTagNotUpdated
	}

	return updated, nil
}

func (t *tagService) DeleteTag(tagId uint, projectId string) (*entities.Tag, error) {
	// Ensure the tag exists
	_, err := t.GetTagById(tagId, projectId)
	if err != nil {
		return nil, err
	}

	deleted, err := t.query.DeleteTag(tagId)
	if err != nil {
		return nil, errors.New("the tag could not be deleted")
	}

	return deleted, nil
}

func (t *tagService) AddToTask(tagId uint, taskId string) error {
	exists := t.query.Exists(tagId)
	if !exists {
		return ehand.ErrorTagNotFound
	}

	err := t.query.AddTagToTask(taskId, tagId)
	if err != nil {
		return err
	}

	return nil
}

func (t *tagService) RemoveFromTask(taskId string, tagId uint) error {
	exists := t.query.Exists(tagId)
	if !exists {
		return ehand.ErrorTagNotFound
	}

	err := t.query.RemoveTagFromTask(taskId, tagId)
	if err != nil {
		return err
	}

	return nil
}

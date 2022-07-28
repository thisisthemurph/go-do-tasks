package services

import (
	"godo/internal/api"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type TagService interface {
	GetTagById(tagId uint, projectId string) (*entities.Tag, error)
	CreateTag(newTag entities.Tag) (*entities.Tag, error)
	UpdateTag(newTag entities.Tag) (*entities.Tag, error)
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
		return nil, api.ErrorTagNotFound
	}

	return tag, nil
}

func (t *tagService) CreateTag(newTag entities.Tag) (*entities.Tag, error) {
	created, err := t.query.CreateTag(newTag)

	if err != nil {
		t.log.Infof("Could not create Tag: ", err)
		return nil, api.ErrorTagNotCreated
	}

	return created, nil
}

func (t *tagService) UpdateTag(newTag entities.Tag) (*entities.Tag, error) {
	updated, err := t.query.UpdateTag(newTag)
	if err != nil {
		return nil, api.ErrorTagNotUpdated
	}

	return updated, nil
}

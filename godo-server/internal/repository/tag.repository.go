package repository

import (
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
)

type TagQuery interface {
	Exists(tagId uint) bool
	GetTagById(tagId uint, projectId string) (*entities.Tag, error)
	CreateTag(newTag entities.Tag) (*entities.Tag, error)
	UpdateTag(newTag entities.Tag) (*entities.Tag, error)
}

type tagQuery struct {
	log ilog.StdLogger
}

func (d *dao) NewTagQuery(logger ilog.StdLogger) TagQuery {
	return &tagQuery{log: logger}
}

func (q *tagQuery) Exists(tagId uint) bool {
	q.log.Infof("Checking if Tag with tagId %d exists", tagId)

	var tag entities.Tag
	r := Database.First(&tag, tagId)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.RowsAffected == 1
}

func (q *tagQuery) GetTagById(tagId uint, projectId string) (*entities.Tag, error) {
	q.log.Infof("Fetching Tag with tagId %d and projectId %s", tagId, projectId)

	var tag entities.Tag
	err := Database.Find(&tag, "id = ? AND project_id = ?", tagId, projectId).Error
	ilog.ErrorlnIf(err, q.log)

	return &tag, err
}

func (q *tagQuery) CreateTag(newTag entities.Tag) (*entities.Tag, error) {
	q.log.Infof("Creating Tag with name %d", newTag.Name)

	r := Database.Create(&newTag)
	ilog.ErrorlnIf(r.Error, q.log)

	return &newTag, r.Error
}

func (q *tagQuery) UpdateTag(newTag entities.Tag) (*entities.Tag, error) {
	q.log.Infof("Updating Tag with tagId %d", newTag.ID)

	err := Database.Save(newTag).Error
	if err != nil {
		q.log.Error("Could not update Tag: ", err)
		return nil, err
	}

	return &newTag, nil
}

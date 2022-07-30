package repository

import (
	"errors"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
)

type TagQuery interface {
	Exists(tagId uint) bool
	GetTagById(tagId uint, projectId string) (*entities.Tag, error)
	CreateTag(newTag entities.Tag) (*entities.Tag, error)
	UpdateTag(newTag entities.Tag) (*entities.Tag, error)
	DeleteTag(tagId uint) (*entities.Tag, error)
	AddTagToTask(taskId string, tagId uint) error
	RemoveTagFromTask(taskId string, tagId uint) error
}

type tagQuery struct {
	log ilog.StdLogger
}

func (d *dao) NewTagQuery(logger ilog.StdLogger) TagQuery {
	return &tagQuery{log: logger}
}

func (q *tagQuery) Exists(tagId uint) bool {
	q.log.Debugf("Checking if Tag with tagId %d exists", tagId)

	var tag entities.Tag
	r := Database.First(&tag, tagId)
	ilog.ErrorlnIf(r.Error, q.log)

	return r.RowsAffected == 1
}

func (q *tagQuery) GetTagById(tagId uint, projectId string) (*entities.Tag, error) {
	q.log.Debugf("Fetching Tag with tagId %d and projectId %s", tagId, projectId)

	var tag entities.Tag
	err := Database.Find(&tag, "id = ? AND project_id = ?", tagId, projectId).Error
	ilog.ErrorlnIf(err, q.log)

	return &tag, err
}

func (q *tagQuery) CreateTag(newTag entities.Tag) (*entities.Tag, error) {
	q.log.Debugf("Creating Tag with name %d", newTag.Name)

	r := Database.Create(&newTag)
	ilog.ErrorlnIf(r.Error, q.log)

	return &newTag, r.Error
}

func (q *tagQuery) UpdateTag(newTag entities.Tag) (*entities.Tag, error) {
	q.log.Debugf("Updating Tag with tagId %d", newTag.ID)

	err := Database.Save(newTag).Error
	if err != nil {
		q.log.Error("Could not update Tag: ", err)
		return nil, err
	}

	return &newTag, nil
}

func (q *tagQuery) DeleteTag(tagId uint) (*entities.Tag, error) {
	q.log.Debugf("Deleting tag with tagId %d", tagId)

	// Delete the related task_tag data
	err := Database.Exec("DELETE FROM task_tags WHERE task_tags.tag_id = ?", tagId).Error
	if err != nil {
		q.log.Error("could not delete task_tags: ", err)
		return nil, err
	}

	// Delete the actual tag
	var deleted entities.Tag
	err = Database.Where("id = ?", tagId).Delete(&deleted).Error
	if err != nil {
		q.log.Error("could not delete Task: ", err)
		return nil, err
	}

	return &deleted, nil
}

func (q *tagQuery) AddTagToTask(taskId string, tagId uint) error {
	q.log.Debugf("Adding tag with tagId %d to task with taskId %s", tagId, taskId)

	r := Database.Exec("INSERT INTO task_tags (task_id, tag_id) VALUES (?, ?)", taskId, tagId)
	if r.Error != nil {
		q.log.Error(r.Error)
		return r.Error
	}

	if r.RowsAffected == 0 {
		q.log.Error("issue relating tag with task - no rows affected")
		return errors.New("issue relating tag with task")
	}

	return nil
}

func (q *tagQuery) RemoveTagFromTask(taskId string, tagId uint) error {
	q.log.Debugf("Removing tag tagId={id=%d} from Task{id=%s}", tagId, taskId)

	err := Database.Exec("DELETE FROM task_tags WHERE task_id = ? AND tag_id = ?", taskId, tagId).Error
	if err != nil {
		return err
	}

	return nil
}

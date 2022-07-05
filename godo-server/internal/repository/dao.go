package repository

import (
	"github.com/jinzhu/gorm"
)

type DAO interface {
	NewStoryQuery() StoryQuery
	NewProjectQuery() ProjectQuery
}

type dao struct{}

var Database *gorm.DB

func NewDAO() DAO {
	Database = GetDatabase()
	return &dao{}
}

func (d *dao) NewStoryQuery() StoryQuery {
	return &storyQuery{}
}

func (d *dao) NewProjectQuery() ProjectQuery {
	return &projectQuery{}
}

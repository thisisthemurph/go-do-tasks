package repository

import (
	"github.com/jinzhu/gorm"
)

type DAO interface {
	NewStoryQuery() StoryQuery
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

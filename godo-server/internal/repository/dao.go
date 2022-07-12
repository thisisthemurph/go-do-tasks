package repository

import (
	"godo/internal/helper/ilog"

	"github.com/jinzhu/gorm"
)

type DAO interface {
	NewStoryQuery(logger ilog.StdLogger) StoryQuery
	NewProjectQuery(logger ilog.StdLogger) ProjectQuery
	NewApiUserQuery(logger ilog.StdLogger) ApiUserQuery
}

type dao struct {
	log ilog.StdLogger
}

var Database *gorm.DB

func NewDAO(logger ilog.StdLogger) DAO {
	Database = GetDatabase(logger)
	return &dao{log: logger}
}

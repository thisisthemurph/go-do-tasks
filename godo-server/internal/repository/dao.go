package repository

import (
	"godo/internal/helper/ilog"

	"github.com/jinzhu/gorm"
)

type DAO interface {
	NewAccountQuery(logger ilog.StdLogger) AccountQuery
	NewApiUserQuery(logger ilog.StdLogger) ApiUserQuery
	NewStoryQuery(logger ilog.StdLogger) StoryQuery
	NewProjectQuery(logger ilog.StdLogger) ProjectQuery
}

type dao struct {
	log ilog.StdLogger
}

var Database *gorm.DB

func NewDAO(logger ilog.StdLogger) DAO {
	Database = GetDatabase(logger)
	return &dao{log: logger}
}

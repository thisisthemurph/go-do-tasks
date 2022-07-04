package entities

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Base struct {
	ID        string
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (base *Base) BeforeCreate(scope *gorm.Scope) error {
 id := uuid.NewV4()
 return scope.SetColumn("ID", id.String())
}

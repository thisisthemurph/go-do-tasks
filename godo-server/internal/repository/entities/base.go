package entities

import (
	"encoding/json"
	"io"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Base struct {
	ID string `json:"id"`
}

type Bases []*Base

type TimestampBase struct {
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.NewV4()
	return scope.SetColumn("ID", id.String())
}

func (b *Bases) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(b)
}

package entities

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"godo/internal/repository/enums"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

type Project struct {
	Base

	Name        string              `json:"name" gorm:"not null" validate:"required,min=1,max=40"`
	Description string              `json:"description" validate:"max=280"`
	Status      enums.ProjectStatus `json:"-" gorm:"type:smallint;default:0;not null" validate:"gte=0"`
	StatusValue string              `json:"status" gorm:"-:all"`
	CreatorId   uint                `json:"-"`
	Creator     User                `json:"creator" gorm:"foreignKey:CreatorId"`
	Stories     []Story             `json:"stories"`

	TimestampBase
}

type ProjectList []*Project
type ProjectKey struct{}

func (p *Project) ToString() string {
	return fmt.Sprintf("Project{Name=%v}", p.Name)
}

// FromJSON TODO: Remove
func (p *Project) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	err := e.Decode(p)
	if err != nil {
		log.Println("Project: there was an error decoding the JSON:", err)
	}

	return err
}

// FromHttpRequest TODO: Remove
func (p Project) FromHttpRequest(r *http.Request) {
	p = r.Context().Value(ProjectKey{}).(Project)
}

// Validate TODO: Remove
func (p *Project) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (p *Project) AfterFind(tx *gorm.DB) {
	p.StatusValue = p.Status.String()
}

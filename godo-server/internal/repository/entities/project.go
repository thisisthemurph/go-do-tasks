package entities

import (
	"encoding/json"
	"fmt"
	"godo/internal/auth"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

type Project struct {
	Base

	Name        string    `json:"name" validate:"required,min=1,max=40"`
	Description string    `json:"description" validate:"max=280"`
	CreatorId   string    `json:"-"`
	Creator     auth.User `gorm:"foreignKey:CreatorId"`
	Stories     []Story
	Tags        []Tag
}

type ProjectList []*Project
type ProjectKey struct{}

func (p *Project) ToString() string {
	return fmt.Sprintf("Project{Name=%v}", p.Name)
}

func (p *Project) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	err := e.Decode(p)
	if err != nil {
		log.Println("Project: there was an error decoding the JSON:", err)
	}

	return err
}

func (p Project) FromHttpRequest(r *http.Request) {
	p = r.Context().Value(ProjectKey{}).(Project)
}

func (p *Project) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

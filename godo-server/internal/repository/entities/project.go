package entities

import (
	"encoding/json"
	"fmt"
	"io"
)

type Project struct {
	Base

	Name        string `json:"name"`
	Description string `json:"description"`

	CreatorId string `json:"-"`
	Creator   Person `gorm:"foreignKey:CreatorId"`
	Stories   []Story
	Tags      []Tag
}

func (p *Project) ToString() string {
	return fmt.Sprintf("Project{Name=%v}", p.Name)
}

func (p *Project) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

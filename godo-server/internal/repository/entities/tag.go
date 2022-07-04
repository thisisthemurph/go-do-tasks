package entities

import (
	"fmt"
)

type Tag struct {
	Base
	ProjectId	string

	Name	string
}

func (t Tag) ToString() string {
	return fmt.Sprintf("Tag{Name=%v}", t.Name)
}
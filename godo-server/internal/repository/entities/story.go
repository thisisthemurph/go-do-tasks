package entities

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Story struct {
	gorm.Model

	// ID		uint
	Name	string
	Tasks 	[]Task
}

func (s Story) ToString() string {
	return fmt.Sprintf("Task{Name=%v}", s.Name)
}
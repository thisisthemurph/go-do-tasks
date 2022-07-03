package entities

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	
	Name        string
	Description string
	Type        int
	Status      int

	StoryID int
}
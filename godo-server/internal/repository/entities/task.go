package entities

type Task struct {
	Base
	StoryId string
	Story   Story

	Name        string
	Description string
	Type        int
	Status      int

	CreatorId string
	Creator   Person `gorm:"foreignKey:CreatorId"`
	Tags      []Tag  `gorm:"many2many:task_tags;"`
}
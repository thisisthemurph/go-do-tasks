package entities

type Task struct {
	Base
	StoryId string
	Story   Story `json:"-"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int
	Status      int

	CreatorId string `json:"-"`
	Creator   Person `gorm:"foreignKey:CreatorId"`
	Tags      []Tag  `gorm:"many2many:task_tags;"`
}

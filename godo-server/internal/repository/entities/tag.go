package entities

type Tag struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	Name      string  `json:"name"`
	ProjectId string  `json:"-"`
	Project   Project `json:"-"`
}

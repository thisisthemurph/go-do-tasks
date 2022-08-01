package entities

//
type Account struct {
	Base

	Name  string `json:"name" validate:"required" gorm:"not null"`
	Email string `json:"email" validate:"required" gorm:"not null"`

	TimestampBase
}

type AccountKey struct{}

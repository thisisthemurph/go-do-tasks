package entities

import "github.com/go-playground/validator"

type Account struct {
	Base

	Name  string `json:"name" validate:"required" gorm:"not null"`
	Email string `json:"email" validate:"required" gorm:"not null"`

	TimestampBase
}

type AccountKey struct{}

func (a *Account) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}

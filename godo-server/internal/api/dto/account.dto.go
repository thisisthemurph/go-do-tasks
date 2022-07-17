package dto

import "github.com/go-playground/validator"

type NewAccountDto struct {
	Name         string `json:"name" validation:"required"`
	UserName     string `json:"user_name" validation:"required"`
	UserUsername string `json:"user_username" validation:"required"`
	UserEmail    string `json:"user_email" validation:"required"`
	Password     string `json:"password" validate:"required"`
}

func (na *NewAccountDto) Validate() error {
	validate := validator.New()
	return validate.Struct(na)
}

package validate

import (
	"github.com/go-playground/validator"
)

func Struct(s any) error {
	v := validator.New()
	return v.Struct(s)
}

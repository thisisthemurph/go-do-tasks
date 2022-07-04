package entities

import (
	"fmt"
)

type Person struct {
	Base
	Name	string `json:"name"`
}

func (p Person) ToString() string {
	return fmt.Sprintf("Person{Name=%v}", p.Name)
}
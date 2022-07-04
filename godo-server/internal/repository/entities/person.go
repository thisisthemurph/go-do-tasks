package entities

import (
	"fmt"
)

type Person struct {
	Base
	Name	string
}

func (p Person) ToString() string {
	return fmt.Sprintf("Person{Name=%v}", p.Name)
}
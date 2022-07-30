package enums

import "fmt"

type ProjectStatus uint8

const (
	Open ProjectStatus = iota
	Closed
)

func (p ProjectStatus) String() string {
	switch p {
	case Open:
	case Closed:
		return "Closed"
	}

	return "Open"
}

func (p ProjectStatus) Print() {
	fmt.Println("ProjectStatus: ", p.String())
}

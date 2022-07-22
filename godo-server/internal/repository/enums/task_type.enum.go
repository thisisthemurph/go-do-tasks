package enums

import "fmt"

type TaskType uint8

const (
	Task TaskType = iota
	Bug
	Test
)

func (t TaskType) String() string {
	switch t {
	case Task:
	case Bug:
		return "Bug"
	case Test:
		return "Test"
	}

	return "Task"
}

func (t TaskType) Print() {
	fmt.Println("ProjectStatus: ", t.String())
}

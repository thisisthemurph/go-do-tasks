package enums

import "fmt"

type StoryStatus uint8

const (
	New StoryStatus = iota
	InProgress
	Complete
)

func (s StoryStatus) String() string {
	switch s {
	case New:
	case InProgress:
		return "In Progress"
	case Complete:
		return "Complete"
	}

	return "New"
}

func (s StoryStatus) Print() {
	fmt.Println("ProjectStatus: ", s.String())
}

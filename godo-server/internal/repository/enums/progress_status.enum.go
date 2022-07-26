package enums

import "fmt"

type ProgressStatus uint8

const (
	New ProgressStatus = iota
	InProgress
	Complete
)

func (s ProgressStatus) String() string {
	switch s {
	case New:
	case InProgress:
		return "In Progress"
	case Complete:
		return "Complete"
	}

	return "New"
}

func (s ProgressStatus) Print() {
	fmt.Println("ProgressStatus: ", s.String())
}

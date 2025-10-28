package vo

import "fmt"

var _ fmt.Stringer = (Status)(0)

type Status int

const (
	_ Status = iota
	Draft
	Serial
	NolongerUpdate
	Completed
)

var statusString = map[Status]string{
	Draft:          `Draft`,
	Serial:         `Serial`,
	NolongerUpdate: `NolongerUpdate`,
	Completed:      `Completed`,
}

func (s Status) String() string {
	str, ok := statusString[s]
	if !ok {
		return `<unknown>`
	}
	return str
}

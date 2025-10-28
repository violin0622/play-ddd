package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*TeacherHired)(nil)

type TeacherHired struct {
	id        ulid.ULID
	teacherID ulid.ULID
	payload   any
	at        time.Time
}

func NewTeacherHired(tid ulid.ULID, payload any) TeacherHired {
	return TeacherHired{
		id:        ulid.Make(),
		at:        time.Now(),
		teacherID: tid,
		payload:   payload,
	}
}

func (t TeacherHired) AggID() ulid.ULID     { return t.teacherID }
func (t TeacherHired) AggKind() string      { return `Teacher` }
func (t TeacherHired) EmittedAt() time.Time { return t.at }
func (t TeacherHired) ID() ulid.ULID        { return t.id }
func (t TeacherHired) Kind() string         { return `TeacherHired` }
func (t TeacherHired) Payload() any         { return t.payload }
func (t TeacherHired) String() string       { return formatEvent(t) }

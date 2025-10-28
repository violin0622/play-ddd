package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*TeacherRetired)(nil)

type TeacherRetired struct {
	at            time.Time
	id, teacherID ulid.ULID
}

func NewTeacherRetired(tid ulid.ULID) TeacherRetired {
	return TeacherRetired{at: time.Now(), teacherID: tid}
}

func (t TeacherRetired) AggID() ulid.ULID     { return t.teacherID }
func (t TeacherRetired) AggKind() string      { return `Teacher` }
func (t TeacherRetired) ID() ulid.ULID        { return t.id }
func (t TeacherRetired) Payload() any         { return nil }
func (t TeacherRetired) EmittedAt() time.Time { return t.at }
func (t TeacherRetired) Kind() string         { return `TeacherRetired` }
func (t TeacherRetired) String() string       { return formatEvent(t) }

package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*TeacherDismissed)(nil)

type TeacherDismissed struct {
	at            time.Time
	id, teacherID ulid.ULID
}

func NewTeacherDismissed(id ulid.ULID) TeacherDismissed {
	return TeacherDismissed{at: time.Now(), teacherID: id}
}

func (t TeacherDismissed) AggID() ulid.ULID     { return t.teacherID }
func (t TeacherDismissed) AggKind() string      { return `Teacher` }
func (t TeacherDismissed) EmittedAt() time.Time { return t.at }
func (t TeacherDismissed) ID() ulid.ULID        { return t.id }
func (t TeacherDismissed) Kind() string         { return `TeacherDismissed` }
func (t TeacherDismissed) Payload() any         { return nil }
func (t TeacherDismissed) String() string       { return formatEvent(t) }

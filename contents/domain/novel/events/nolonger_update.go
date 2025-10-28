package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*NolongerUpdate)(nil)

type NolongerUpdate struct {
	id  ID
	aid ID
	at  time.Time
}

func NewNolongerUpdate(aid ID) NolongerUpdate {
	return NolongerUpdate{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
	}
}

func (t NolongerUpdate) AggID() ID            { return t.aid }
func (t NolongerUpdate) AggKind() string      { return `Novel` }
func (t NolongerUpdate) EmittedAt() time.Time { return t.at }
func (t NolongerUpdate) ID() ID               { return t.id }
func (t NolongerUpdate) Kind() string         { return `NolongerUpdate` }
func (t NolongerUpdate) String() string       { return formatEvent(t) }

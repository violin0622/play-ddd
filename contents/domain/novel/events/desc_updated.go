package events

import (
	"time"

	"play-ddd/contents/domain/novel/vo"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*TagsUpdated)(nil)

type DescUpdated struct {
	id  ID
	aid ID
	at  time.Time

	Desc vo.Description
}

func NewDescUpdated(aid ID, desc vo.Description) DescUpdated {
	return DescUpdated{
		id:   ulid.Make(),
		at:   time.Now(),
		aid:  aid,
		Desc: desc,
	}
}

func (t DescUpdated) AggID() ID            { return t.aid }
func (t DescUpdated) AggKind() string      { return `Novel` }
func (t DescUpdated) EmittedAt() time.Time { return t.at }
func (t DescUpdated) ID() ID               { return t.id }
func (t DescUpdated) Kind() string         { return `DescUpdated` }
func (t DescUpdated) String() string       { return formatEvent(t) }

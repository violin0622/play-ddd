package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*NovelPublished)(nil)

type NovelPublished struct {
	id  ID
	aid ID
	at  time.Time
}

func NewNovelPublished(aid ID) NovelPublished {
	return NovelPublished{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
	}
}

func (t NovelPublished) AggID() ID                { return t.aid }
func (t NovelPublished) AggKind() string          { return `Novel` }
func (t NovelPublished) EmittedAt() time.Time     { return t.at }
func (t NovelPublished) ID() ID                   { return t.id }
func (t NovelPublished) Kind() string             { return `NovelPublished` }
func (t NovelPublished) String() string           { return formatEvent(t) }
func (t NovelPublished) Payload() ([]byte, error) { return emptyPayload() }

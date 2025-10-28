package events

import (
	"time"

	"play-ddd/contents/domain/novel/vo"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*NewChapterUploaded)(nil)

type NewChapterUploaded struct {
	id  ID
	aid ID
	at  time.Time

	Chapter vo.Chapter
}

func NewNewChapterUploaded(aid ID, c vo.Chapter) NewChapterUploaded {
	return NewChapterUploaded{
		id:      ulid.Make(),
		at:      time.Now(),
		aid:     aid,
		Chapter: c,
	}
}

func (t NewChapterUploaded) AggID() ID            { return t.aid }
func (t NewChapterUploaded) AggKind() string      { return `Novel` }
func (t NewChapterUploaded) EmittedAt() time.Time { return t.at }
func (t NewChapterUploaded) ID() ID               { return t.id }
func (t NewChapterUploaded) Kind() string         { return `NewChapterUploaded` }
func (t NewChapterUploaded) String() string       { return formatEvent(t) }

package chapter

import (
	"time"

	"github.com/oklog/ulid/v2"

	"play-ddd/common"
	"play-ddd/contents/domain/novel/vo"
)

var formatEvent = common.FormatEvent[ID, ID]

var _ Event = (*ChapterRevised)(nil)

type ChapterRevised struct {
	id           ID
	aid          ID
	at           time.Time
	seq          int
	title        string
	mainContent  string
	extraContent string
}

func NewChapterRevised(aid ID, c vo.Chapter) ChapterRevised {
	return ChapterRevised{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
	}
}

func (t ChapterRevised) AggID() ID            { return t.aid }
func (t ChapterRevised) AggKind() string      { return `Chapter` }
func (t ChapterRevised) EmittedAt() time.Time { return t.at }
func (t ChapterRevised) ID() ID               { return t.id }
func (t ChapterRevised) Kind() string         { return `ChapterRevised` }
func (t ChapterRevised) String() string       { return formatEvent(t) }

var _ Event = ChapterUploaded{}

type ChapterUploaded struct {
	id           ID
	aid          ID
	at           time.Time
	seq          int
	title        string
	mainContent  string
	extraContent string
	wordCount    int
}

func (t ChapterUploaded) AggID() ID            { return t.aid }
func (t ChapterUploaded) AggKind() string      { return `Chapter` }
func (t ChapterUploaded) EmittedAt() time.Time { return t.at }
func (t ChapterUploaded) ID() ID               { return t.id }
func (t ChapterUploaded) Kind() string         { return `ChapterUploaded` }
func (t ChapterUploaded) String() string       { return formatEvent(t) }

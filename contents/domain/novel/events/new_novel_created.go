package events

import (
	"time"

	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain/novel/vo"
)

var _ Event = (*NovelCreated)(nil)

// NovelCreated is init event of novel agg.
type NovelCreated struct {
	id  ID
	aid ID
	at  time.Time

	AuthorID ID
	Title    vo.Title
	Desc     vo.Description
	Tags     []vo.Tag
	Category vo.Category
	Cover    vo.Cover
	// FirstChapter vo.Chapter
}

func NewNovelCreated(
	aid ID,
	autherID ID,
	title vo.Title,
	cover vo.Cover,
	desc vo.Description,
	category vo.Category,
	tags []vo.Tag,
	// c vo.Chapter,
) NovelCreated {
	return NovelCreated{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,

		AuthorID: autherID,
		Title:    title,
		Desc:     desc,
		Tags:     tags,
		Category: category,
		Cover:    cover,
		// FirstChapter: c,
	}
}

func (t NovelCreated) AggID() ID            { return t.aid }
func (t NovelCreated) AggKind() string      { return `Novel` }
func (t NovelCreated) EmittedAt() time.Time { return t.at }
func (t NovelCreated) ID() ID               { return t.id }
func (t NovelCreated) Kind() string         { return `NovelCreated` }
func (t NovelCreated) String() string       { return formatEvent(t) }

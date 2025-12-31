package events

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain/novel/vo"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
)

// var _ fromPB[novelv1.NovelCreated] = NovelCreated{}

// NovelCreated is init event of novel agg.
type NovelCreated struct {
	id  ID        `json:"-"`
	aid ID        `json:"-"`
	at  time.Time `json:"-"`

	AuthorID ID
	Title    vo.Title
	Desc     vo.Description
	Tags     []vo.Tag
	Category vo.Category
	Cover    vo.Cover
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

func NewNovelCreatedV2(
	aid ID,
	authorID ID,
	title vo.Title,
	cover vo.Cover,
	desc vo.Description,
	category vo.Category,
	tags []vo.Tag,
) NovelCreatedV2 {
	tag := make([]string, len(tags))
	for i := range tags {
		tag[i] = string(tags[i])
	}

	return NovelCreatedV2{
		id:      ulid.Make(),
		aid:     aid,
		at:      time.Now(),
		kind:    "created",
		aggKind: "novel",
		payload: &novelv1.NovelCreated{
			AuthorId:    ulidpb.From(authorID),
			Title:       string(title),
			Description: string(desc),
			Tags:        tag,
			Category:    string(category),
		},
	}
}

type NovelCreatedV2 = event[novelv1.NovelCreated, *novelv1.NovelCreated]

func (t NovelCreated) AggID() ID                { return t.aid }
func (t NovelCreated) AggKind() string          { return `Novel` }
func (t NovelCreated) EmittedAt() time.Time     { return t.at }
func (t NovelCreated) ID() ID                   { return t.id }
func (t NovelCreated) Kind() string             { return `NovelCreated` }
func (t NovelCreated) String() string           { return formatEvent(t) }
func (t NovelCreated) Payload() ([]byte, error) { return json.Marshal(t) }

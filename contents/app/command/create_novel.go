package command

import (
	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain/novel"
)

type CreateNovel struct {
	AuthorID ulid.ULID
	Title    string
	Desc     string
	Category string
	Tags     []string
}

type UploadChapter struct {
	ID        novel.ID
	Title     string
	WordCount int
}

type UpdateTags struct{}

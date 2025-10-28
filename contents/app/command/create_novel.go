package command

import (
	"play-ddd/contents/domain/novel"

	"github.com/oklog/ulid/v2"
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

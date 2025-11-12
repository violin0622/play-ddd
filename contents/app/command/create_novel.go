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
	NovelID                   novel.ID
	Title                     string
	MainContent, ExtraContent string
}

type UpdateNovelInfo struct {
	NovelID novel.ID
	Desc    string
	Tags    []string
}

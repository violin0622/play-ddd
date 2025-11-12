package vo

import (
	"cmp"
	"time"
)

var SentinelChapter = Chapter{}

type Chapter struct {
	// ID         ulid.ULID
	Title      string
	Sequence   int
	WordCount  int
	UploadedAt time.Time
	UpdatedAt  time.Time
}

func ChapterCmp(a, b Chapter) int {
	return cmp.Compare(a.Sequence, b.Sequence)
}

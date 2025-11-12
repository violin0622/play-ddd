package app

import (
	"context"

	"github.com/go-logr/logr"

	"play-ddd/contents/domain/chapter"
	"play-ddd/contents/domain/novel"
)

type CommandHandler struct {
	novelFact novel.Factory
	cf        chapter.Factory
	repo      Repo
	log       logr.Logger
}

func NewCommandHandler(
	repo Repo,
	log logr.Logger,
) CommandHandler {
	return CommandHandler{
		repo: repo,
		log:  log,
	}
}

type QueryHandler struct {
	repo Repo
	log  logr.Logger
}

type Repo interface {
	Chapter() ChapterRepo
	Novel() NovelRepo
	Tx(func(Repo) error) error
}

type ChapterRepo interface {
	chapter.EventRepo
	Get(context.Context, chapter.ID) (chapter.Chapter, error)
	Save(context.Context, chapter.Chapter) error
}

type NovelRepo interface {
	novel.EventRepo
	Get(context.Context, novel.ID) (novel.Novel, error)
	Save(context.Context, novel.Novel) error
}

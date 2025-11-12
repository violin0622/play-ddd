package chapter

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"
)

type Factory struct {
	er  EventRepo
	log logr.Logger
}

func NewFactory(er EventRepo, log logr.Logger) Factory {
	return Factory{er: er, log: log}
}

func (f Factory) WithEventRepo(er EventRepo) Factory {
	return NewFactory(er, f.log)
}

func (cf Factory) UploadChapter(
	ctx context.Context,
	seq int,
	title, mainContent, exContent string) (
	Chapter, error,
) {
	if seq < 1 {
		return Chapter{}, fmt.Errorf(`invalid sequence`)
	}

	c := Chapter{
		id:           ulid.Make(),
		seq:          seq,
		title:        title,
		mainContent:  mainContent,
		extraContent: exContent,
		er:           cf.er,
	}

	err := cf.er.Append(ctx, ChapterUploaded{
		id:           ulid.Make(),
		aid:          c.id,
		at:           Now(),
		seq:          seq,
		title:        c.title,
		mainContent:  c.mainContent,
		extraContent: c.extraContent,
		wordCount:    countWords(c.mainContent),
	})
	if err != nil {
		return Chapter{}, fmt.Errorf(`upload chapter: %w`, err)
	}
	return c, nil
}

package novel

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"
	"github.com/samber/mo"

	ev "play-ddd/contents/domain/novel/events"
	"play-ddd/contents/domain/novel/vo"
)

// If not found, implements should return NotfoundError.
type Query interface {
	ByAuthorAndTitle(context.Context, ID, string) mo.Result[Novel]
}

type Factory struct {
	query Query
	lg    logr.Logger
	er    EventRepo
}

func NewFactory(
	log logr.Logger,
) Factory {
	return Factory{
		lg: log,
	}
}

func (f Factory) WithQuery(q Query) Factory {
	nf := f
	nf.query = q
	return nf
}

func (f Factory) WithEventRepo(er EventRepo) Factory {
	nf := f
	nf.er = er
	return nf
}

// 创建小说分两步：
// 1. 创建小说对象并锁定标题，此时为草案状态；
// 2. 上传第一个章节. 变为连载中；
// 3. 草案状态的小说对外不可见.
// 4. 描述，标签都不是必须的.
func (f Factory) Create(
	ctx context.Context,
	authorID ID,
	title, desc, category string,
	tags []string,
) mo.Result[Novel] {
	var err error
	novel := Novel{
		id:        ulid.Make(),
		authorID:  authorID,
		category:  vo.Category(category),
		desc:      vo.Description(desc),
		s:         vo.Draft,
		er:        f.er,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	if novel.toc, err = vo.NewTOC(); err != nil {
		return mo.Err[Novel](fmt.Errorf(`create new novel: %w`, err))
	}

	if err := f.checkTitleUnique(ctx, authorID, title); err != nil {
		return mo.Err[Novel](fmt.Errorf(`create new novel: %w`, err))
	}

	if novel.title, err = vo.NewTitle(title); err != nil {
		return mo.Err[Novel](fmt.Errorf(`create new novel: %w`, err))
	}

	for i := range tags {
		novel.tags = append(novel.tags, vo.Tag(tags[i]))
	}

	e := ev.NewNovelCreatedV2(
		novel.id,
		authorID,
		novel.title,
		novel.cover,
		novel.desc,
		novel.category,
		novel.tags)
	if err := f.er.Append(ctx, e); err != nil {
		return mo.Err[Novel](fmt.Errorf(`create new novel: %w`, err))
	}

	return mo.Ok(novel)
}

// New returns a novel instance for replay.
func (f Factory) New() Novel {
	toc, _ := vo.NewTOC()
	return Novel{er: f.er, lg: f.lg, toc: toc}
}

// Restore is only intend to be used by infra/repo.
func (f Factory) Restore(
	id, authorID ulid.ULID,
	title, category, desc string,
	tags []string,
	status int,
	cover url.URL,
	chapters []vo.Chapter,
	wordCount int,
	updatedAt, createdAt time.Time) (
	n Novel, err error,
) {
	defer wrapOnError(`resotre novel from repo`, &err)

	n.id = id
	n.authorID = authorID
	n.category = vo.Category(category)
	n.desc = vo.Description(desc)
	n.s = vo.Status(status)
	n.createdAt = createdAt
	n.updatedAt = updatedAt
	n.wordCount = wordCount

	if n.toc, err = vo.NewTOC(chapters...); err != nil {
		return n, err
	}

	if n.title, err = vo.NewTitle(title); err != nil {
		return n, err
	}

	for i := range tags {
		n.tags = append(n.tags, vo.Tag(tags[i]))
	}

	return n, n.checkInvariants()
}

func wrapOnError(msg string, err *error) {
	if err != nil && *err != nil {
		*err = fmt.Errorf(msg+`: %w`, *err)
	}
}

func (f Factory) checkTitleUnique(
	ctx context.Context,
	authorID ID,
	title string,
) error {
	if f.query == nil {
		return fmt.Errorf(`nil query repo`)
	}
	res := f.query.ByAuthorAndTitle(ctx, authorID, title)
	var notfound NotfoundError

	if errors.As(res.Error(), &notfound) {
		return nil
	}

	if res.IsError() {
		return fmt.Errorf(`check title unique: %w`, res.Error())
	}

	return ErrTitleAlreadyExist
}

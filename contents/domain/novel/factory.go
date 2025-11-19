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
type query interface {
	ByAuthorAndTitle(context.Context, ID, string) mo.Result[Novel]
}

type Factory struct {
	query query
	lg    logr.Logger
	er    EventRepo
	repo  Repo
}

func NewFactory(
	er EventRepo,
	log logr.Logger,
	repo Repo,
) Factory {
	return Factory{
		er:   er,
		lg:   log,
		repo: repo,
	}
}

func (f Factory) WithRepo(repo Repo) Factory {
	f = NewFactory(f.er, f.lg, f.repo)
	f.repo = repo
	return f
}

func (f Factory) WithEventRepo(er EventRepo) Factory {
	f = NewFactory(f.er, f.lg, f.repo)
	f.er = er
	return f
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
		id:       ulid.Make(),
		authorID: authorID,
		category: vo.Category(category),
		desc:     vo.Description(desc),
		s:        vo.Draft,
		er:       f.er,
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

	e := ev.NewNovelCreated(
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

func (f Factory) UploadFirstChapter(
	ctx context.Context, id ID, title string, wc int,
) error {
	novel, err := f.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf(`upload first chapter: %w`, err)
	}

	if novel.s != vo.Draft {
		return fmt.Errorf(`novel is not a draft`)
	}

	if err = novel.UploadNewChapter(ctx, title, wc); err != nil {
		return fmt.Errorf(`upload first chapter: %w`, err)
	}

	if err = f.repo.Save(ctx, novel); err != nil {
		return fmt.Errorf(`upload first chapter: %w`, err)
	}

	if err = f.er.Append(ctx, ev.NewNovelPublished(novel.id)); err != nil {
		return fmt.Errorf(`upload first chapter: %w`, err)
	}

	return nil
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

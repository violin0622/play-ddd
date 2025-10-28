package novel

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"play-ddd/common"
	"play-ddd/contents/domain/novel/vo"

	ev "play-ddd/contents/domain/novel/events"
)

var _ Aggregate = (*Novel)(nil)

// Novel is the root of 'novel' aggregate.
// Invariants:
// 1. 必须有作者 (novel)
// 2. 标题在同作者所有作品中必须唯一 (service)
// 3. 至少要有一个章节 (novel)
// 4. 章节编号不可重复 (novel)
// 5. 最多10个标签 (novel)
// 6. desc 不超过500字 (desc)
// 7. 总字数 = 每章字数之和 (novel)
// 8. 更新时间 = 最新章节更新时间 (novel)
// 9. 章节数 = 总章节数 (novel)
//
// 1. 开新书状态初始为 Serial
// 2. 长时间不更新自动变为 NolongerUpdate
// 3. 断更后上传新章节重新变为 Serial
// 4. 由作者修改为 Completed
// 5. Completed 状态不能再上传或撤回章节，但可以修改现有章节。
//
// 1. 只有作者可以更新小说
// 2. 只有管理员可以下架小说
// 3. 下架后不可被编辑
type Novel struct {
	id        ID
	title     vo.Title
	category  vo.Category
	authorID  ID
	desc      vo.Description
	tags      []vo.Tag
	s         vo.Status
	cover     vo.Cover
	toc       vo.TOC
	wordCount int
	updatedAt time.Time
	createdAt time.Time

	er EventRepo
}

func New(
	er EventRepo,
) Novel {
	return Novel{
		er:    er,
		toc:   vo.TOC{Chapters: []vo.Chapter{vo.SentinelChapter}},
		cover: vo.NoCover,
	}
}

func (n Novel) ID() ID                      { return n.id }
func (n Novel) Kind() string                { return `Novel` }
func (n Novel) WordCounts() int             { return n.wordCount }
func (n Novel) Title() vo.Title             { return n.title }
func (n Novel) AuthorID() ID                { return n.authorID }
func (n Novel) Description() vo.Description { return n.desc }
func (n Novel) Tags() []vo.Tag              { return n.tags }
func (n Novel) Category() vo.Category       { return n.category }
func (n Novel) Cover() vo.Cover             { return n.cover }
func (n Novel) Status() vo.Status           { return n.s }
func (n Novel) TOC() vo.TOC                 { return n.toc }
func (n Novel) UpdatedAt() time.Time        { return n.updatedAt }
func (n Novel) CreatedAt() time.Time        { return n.createdAt }

func (n *Novel) ReplayEvents(es ...Event) error {
	if len(es) == 0 {
		return nil
	}

	if _, ok := es[0].(ev.NovelCreated); !ok {
		return common.ErrInitialEvent
	}

	for i := range es {
		if err := n.applyEvent(es[i]); err != nil {
			return fmt.Errorf(`replay events: %w`, err)
		}
	}

	return n.checkInvariants()
}

func (n *Novel) AppendTags(ctx context.Context, tags ...string) error {
	prev := slices.Clone(n.tags)
	for i := range tags {
		n.tags = append(n.tags, vo.Tag(tags[i]))
	}

	return n.finish(ctx, ev.NewTagsUpdated(n.id, prev, n.tags))
}

func (n *Novel) RemoveTags(ctx context.Context, tags ...string) error {
	prev := slices.Clone(n.tags)
	for i := range tags {
		n.tags = slices.DeleteFunc(
			n.tags,
			func(t vo.Tag) bool { return t == vo.Tag(tags[i]) },
		)
	}

	return n.finish(ctx, ev.NewTagsUpdated(n.id, prev, n.tags))
}

func (n *Novel) UpdateDescription(ctx context.Context, d string) error {
	desc, err := vo.NewDescription(d)
	if err != nil {
		return err
	}

	n.desc = desc.MustGet()
	return n.finish(ctx, ev.NewDescUpdated(n.id, n.desc))
}

func (n *Novel) UploadNewChapter(
	ctx context.Context, title string, wc int,
) error {
	now := Now()
	c := vo.Chapter{
		Title:      title,
		Sequence:   len(n.toc.Chapters),
		WordCount:  wc,
		UploadedAt: now,
		UpdatedAt:  now,
	}
	n.toc.Chapters = append(n.toc.Chapters, c)
	n.s = vo.Serial
	n.updatedAt = now

	return n.finish(ctx, ev.NewNewChapterUploaded(n.id, c))
}

func (n *Novel) WithdrawChapter(ctx context.Context) error {
	seq := len(n.toc.Chapters) - 1
	e := ev.NewChapterWithdrawed(n.id, n.toc.Chapters[seq])
	if err := n.imposeChapterWithdrawed(e); err != nil {
		return err
	}

	return n.finish(ctx, e)
}

func (n *Novel) ReviseChapter(
	ctx context.Context, title string, wc, seq int,
) error {
	if seq < 1 || seq > len(n.toc.Chapters) {
		return fmt.Errorf(`invalid chapter sequence`)
	}

	c := n.toc.Chapters[seq]
	c.Title = title
	c.UpdatedAt = Now()
	c.WordCount = wc

	e := ev.NewChapterRevised(n.id, c)
	if err := n.imposeChapterRevised(e); err != nil {
		return err
	}

	return n.finish(ctx, e)
}

func (n *Novel) Complete(ctx context.Context) error {
	n.s = vo.Completed
	return n.finish(ctx, ev.NewCompleted(n.id))
}

func (n *Novel) NolongerUpdate(ctx context.Context) error {
	n.s = vo.NolongerUpdate
	return n.finish(ctx, ev.NewCompleted(n.id))
}

// finish chech invariants and emit events.
func (n Novel) finish(ctx context.Context, e ...Event) error {
	if err := n.checkInvariants(); err != nil {
		return newInvariantsBrokenError(err)
	}

	return n.er.Append(ctx, e...)
}

func (n Novel) checkInvariants() error {
	invarants := []func() error{
		n.autherMustExist,
		n.chapterSequenceOrdered,
		n.noRepeatedOrLostedChapterSequence,
		n.atMost10Tags,
		n.totalWordCountIsSummitOfChapterWordCount,
		n.updatedAtNewstChapterUpload,
	}

	for i := range invarants {
		if err := invarants[i](); err != nil {
			return err
		}
	}

	return nil
}

func (n Novel) autherMustExist() error {
	if n.authorID.IsZero() {
		return errors.New(`authorID must exist`)
	}
	return nil
}

func (n Novel) updatedAtNewstChapterUpload() error {
	last := len(n.toc.Chapters) - 1
	if n.toc.Chapters[last].UploadedAt.Equal(n.updatedAt) {
		return nil
	}
	return errors.New(`novel's update time should be identical to last chapter's update time`)
}

func (n Novel) chapterSequenceOrdered() error {
	if slices.IsSortedFunc(n.toc.Chapters, vo.ChapterCmp) {
		return nil
	}
	return errors.New(`chapter sequence not ordered`)
}

func (n Novel) noRepeatedOrLostedChapterSequence() error {
	for seq := range len(n.toc.Chapters) {
		if n.toc.Chapters[seq].Sequence != seq {
			return errors.New(`chapter sequence missed or repeated`)
		}
	}

	return nil
}

func (n Novel) totalWordCountIsSummitOfChapterWordCount() error {
	sum := 0
	for _, c := range n.toc.Chapters {
		sum += c.WordCount
	}

	if sum != n.wordCount {
		return fmt.Errorf(
			`novel's word count %d should be summit of every chapter's word count %d`,
			n.wordCount, sum)
	}

	return nil
}

func (n Novel) atMost10Tags() error {
	if len(n.tags) > 10 {
		return errors.New(`at most 10 tags allowed`)
	}

	return nil
}

func (n *Novel) applyEvent(e ev.Event) error {
	switch e := e.(type) {
	case ev.NovelCreated:
		return n.imposeNovelCreated(e)
	case ev.NovelPublished:
		return nil
	case ev.TagsUpdated:
		return n.imposeTagsUpdated(e)
	case ev.DescUpdated:
		return n.imposeDescUpdated(e)
	case ev.NewChapterUploaded:
		return n.imposeNewChapterUploaded(e)
	case ev.ChapterRevised:
		return n.imposeChapterRevised(e)
	case ev.ChapterWithdrawed:
		return n.imposeChapterWithdrawed(e)
	case ev.NolongerUpdate:
		return n.imposeNolongerUpdate(e)
	case ev.Completed:
		return n.imposeCompleted(e)
	default:
		return common.ErrUnknownEventKind(e.Kind())
	}
}

func (n *Novel) imposeNovelCreated(e ev.NovelCreated) error {
	if n == nil {
		return errors.New(`nil pointer`)
	}

	n.id = e.AggID()
	n.title = e.Title
	n.authorID = e.AuthorID
	n.tags = e.Tags
	n.cover = e.Cover
	n.s = vo.Draft
	n.category = e.Category
	// n.toc.Chapters = append(n.toc.Chapters, e.FirstChapter)
	// n.wordCount = e.FirstChapter.WordCount
	n.createdAt = e.EmittedAt()
	// n.updatedAt = e.FirstChapter.UploadedAt

	return nil
}

func (n *Novel) imposeTagsUpdated(e ev.TagsUpdated) error {
	n.tags = e.CurrentTags
	return nil
}

func (n *Novel) imposeDescUpdated(e ev.DescUpdated) error {
	n.desc = e.Desc
	return nil
}

func (n *Novel) imposeCompleted(ev.Completed) error {
	n.s = vo.Completed
	return nil
}

func (n *Novel) imposeNolongerUpdate(ev.NolongerUpdate) error {
	n.s = vo.NolongerUpdate
	return nil
}

func (n *Novel) imposeNewChapterUploaded(e ev.NewChapterUploaded) error {
	if n == nil {
		return errors.New(`nil pointer`)
	}

	if n.s == vo.Completed {
		return ErrMutateCompletedNovel
	}

	n.toc.Chapters = append(n.toc.Chapters, e.Chapter)
	n.wordCount += e.Chapter.WordCount
	n.s = vo.Serial
	n.updatedAt = e.Chapter.UploadedAt

	return nil
}

func (n *Novel) imposeChapterRevised(e ev.ChapterRevised) error {
	if n == nil {
		return errors.New(`nil pointer`)
	}

	seq := e.RevisedChapter.Sequence
	if !n.isSequenceValid(seq) {
		return errors.New(`invalid chapter sequence`)
	}

	c := &n.toc.Chapters[seq]
	c.Title = e.RevisedChapter.Title
	c.UpdatedAt = e.RevisedChapter.UpdatedAt
	delta := c.WordCount - e.RevisedChapter.WordCount
	c.WordCount = e.RevisedChapter.WordCount
	n.wordCount -= delta

	return nil
}

func (n Novel) isSequenceValid(seq int) bool {
	return seq < 1 || seq >= len(n.toc.Chapters)
}

func (n *Novel) imposeChapterWithdrawed(ev.ChapterWithdrawed) error {
	if n == nil {
		return errors.New(`nil pointer`)
	}

	if n.s == vo.Completed {
		return ErrMutateCompletedNovel
	}

	seq := len(n.toc.Chapters) - 1
	if seq <= 1 {
		return errors.New(`can not withdraw first chapter`)
	}

	c := n.toc.Chapters[seq]
	n.toc.Chapters = n.toc.Chapters[:seq]
	n.wordCount -= c.WordCount
	n.updatedAt = n.toc.Chapters[seq-1].UploadedAt

	return nil
}

// Now is exported only for test, do not modify elsewhere.
var Now = time.Now

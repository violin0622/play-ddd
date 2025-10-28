package novel_test

import (
	"context"
	"time"

	"play-ddd/common"
	"play-ddd/contents/domain/novel"
	"play-ddd/contents/domain/novel/events"
	"play-ddd/contents/domain/novel/vo"
	"play-ddd/contents/infra/eventstore/fake"

	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ novel.EventRepo = (*nopEventRepo)(nil)

type nopEventRepo struct{}

func (nopEventRepo) Append(context.Context, ...novel.Event) error { return nil }
func (nopEventRepo) Fetch(context.Context, novel.ID) ([]novel.Event, error) {
	return nil, nil
}

var _ = Describe(`NovelReplayEvents`, func() {
	var n novel.Novel
	var err error
	id, autherID := ulid.Make(), ulid.Make()
	// events
	var (
		novelCreated = events.NewNovelCreated(
			id,
			autherID,
			vo.Title(`测试小说`),
			vo.NoCover,
			vo.Description(`初始描述`),
			vo.Category(`玄幻`),
			[]vo.Tag{vo.Tag(`修仙`), vo.Tag(`穿越`)})
		descUpdated = events.NewDescUpdated(id, vo.Description(`更新后的描述`))
		tagsUpdated = events.NewTagsUpdated(id,
			[]vo.Tag{vo.Tag(`修仙`), vo.Tag(`穿越`)},
			[]vo.Tag{vo.Tag(`修仙`), vo.Tag(`穿越`), vo.Tag(`系统`)})

		chapterUploaded = events.NewNewChapterUploaded(id, vo.Chapter{
			Title:      "第二章",
			Sequence:   2,
			WordCount:  1500,
			UploadedAt: time.Now(),
			UpdatedAt:  time.Now(),
		})

		chapterRevised = events.NewChapterRevised(id, vo.Chapter{
			Title:      "第一章（修订版）",
			Sequence:   1,
			WordCount:  1200,
			UploadedAt: time.Now(),
			UpdatedAt:  time.Now(),
		})
		nolongerUpdate = events.NewNolongerUpdate(id)
		completed      = events.NewCompleted(id)
	)
	eventsHist := []events.Event{
		novelCreated,
		descUpdated,
		tagsUpdated,
		chapterUploaded,
		chapterRevised,
		nolongerUpdate,
		completed,
	}

	When(`first event is NovelCreated`, func() {
		It(`should accept`, func() {
			n = novel.New(nopEventRepo{})
			err = n.ReplayEvents(novelCreated)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	When(`first event is not NovelCreated`, func() {
		DescribeTable(`should reject invalid initial events`,
			func(e novel.Event) {
				n = novel.New(nopEventRepo{})
				err = n.ReplayEvents(e)
				Ω(err).Should(MatchError(common.ErrInitialEvent))
			},
			Entry(`DescUpdated`, descUpdated),
			Entry(`TagsUpdated`, tagsUpdated),
			Entry(`NewChapterUploaded`, chapterUploaded),
			Entry(`ChapterRevised`, chapterRevised),
			Entry(`Completed`, completed),
			Entry(`NolongerUpdate`, nolongerUpdate),
		)
	})

	When(`replaying multiple events starting with NovelCreated`, func() {
		n := novel.New(nopEventRepo{})
		err := n.ReplayEvents(eventsHist...)

		It(`accepts`, func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It(`should maintain correct basic properties`, func() {
			Ω(n.ID()).Should(Equal(id))
			Ω(n.AuthorID()).Should(Equal(autherID))
			Ω(n.Title()).Should(Equal(vo.Title(`测试小说`)))
			Ω(n.Description()).Should(Equal(vo.Description(`更新后的描述`)))
			Ω(n.Category()).Should(Equal(vo.Category(`玄幻`)))
			Ω(n.Cover()).Should(Equal(vo.NoCover))
		})

		It(`should maintain correct status and tags`, func() {
			Ω(n.Status()).Should(Equal(vo.Completed))
			expectedTags := []vo.Tag{vo.Tag(`修仙`), vo.Tag(`穿越`), vo.Tag(`系统`)}
			Ω(n.Tags()).Should(Equal(expectedTags))
		})

		It(`should maintain correct chapter information`, func() {
			toc := n.TOC()
			Ω(len(toc.Chapters)).Should(Equal(3)) // 哨兵章节 + 2个实际章节
			Ω(toc.Chapters[1].Title).Should(Equal("第一章（修订版）"))
			Ω(toc.Chapters[1].Sequence).Should(Equal(1))
			Ω(toc.Chapters[1].WordCount).Should(Equal(1200))
			Ω(toc.Chapters[2].Title).Should(Equal("第二章"))
			Ω(toc.Chapters[2].Sequence).Should(Equal(2))
			Ω(toc.Chapters[2].WordCount).Should(Equal(1500))
		})

		It(`should maintain correct word count and timestamps`, func() {
			Ω(n.WordCounts()).Should(Equal(2700)) // 1200 + 1500
			Ω(n.CreatedAt()).ShouldNot(BeZero())
			Ω(n.UpdatedAt()).ShouldNot(BeZero())
		})
	})
})

var _ = Describe(``, func() {
	var _ novel.EventRepo = fake.New[novel.ID, novel.ID]()
})

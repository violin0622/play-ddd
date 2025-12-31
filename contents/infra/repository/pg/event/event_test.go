package event_test

import (
	"testing"

	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	ne "play-ddd/contents/domain/novel/events"
	"play-ddd/contents/domain/novel/vo"
	"play-ddd/contents/infra/repository/pg"
	eventrepo "play-ddd/contents/infra/repository/pg/event"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Event Suite")
}

var db *gorm.DB

var _ = BeforeSuite(func() {
	var err error
	db, err = pg.InitDB(pg.DSN)
	Ω(err).ShouldNot(HaveOccurred())
})

// var _ = Describe(`EventRepo`, func() {
// 	var (
// 		novelID   = ulid.Make()
// 		chapterID = ulid.Make()
// 		desc      = vo.Description(`A New Description.`)
// 		e1        = ne.NewDescUpdated(novelID, desc)
// 		e2        = ne.NewCompleted(novelID)
// 		e3        = chapter.NewChapterRevised(chapterID, vo.Chapter{})
// 	)
//
// 	BeforeEach(func() {
// 		Ω(
// 			db.Unscoped().Where(`1=1`).Delete(&eventrepo.Event{}).Error,
// 		).Should(Succeed())
//
// 		repo := eventrepo.New(db)
// 		Ω(repo.Append(context.TODO(), e1, e2, e3)).Should(Succeed())
// 	})
//
// 	It(`fetch events`, func(ctx SpecContext) {
// 		repo := eventrepo.New(db)
// 		es, err := repo.Fetch(ctx, novelID)
// 		Ω(err).ShouldNot(HaveOccurred())
// 		Ω(es).Should(HaveLen(2))
//
// 		// 验证第一个事件
// 		ev1 := es[0].(ne.DescUpdated)
// 		Ω(ev1.ID()).Should(BeEquivalentTo(e1.ID()))
// 		Ω(ev1.AggID()).Should(BeEquivalentTo(novelID))
// 		Ω(ev1.Kind()).Should(Equal(`DescUpdated`))
// 		Ω(ev1.AggKind()).Should(Equal(`Novel`))
// 		Ω(ev1.Desc).Should(BeEquivalentTo(desc))
// 		Ω(
// 			ev1.EmittedAt(),
// 		).Should(BeTemporally("~", e1.EmittedAt(), time.Second))
//
// 		// 验证第二个事件
// 		ev2 := es[1].(ne.Completed)
// 		Ω(ev2.ID()).Should(BeEquivalentTo(e2.ID()))
// 		Ω(ev2.AggID()).Should(BeEquivalentTo(novelID))
// 		Ω(ev2.Kind()).Should(Equal(`Completed`))
// 		Ω(ev2.AggKind()).Should(Equal(`Novel`))
// 		Ω(
// 			ev2.EmittedAt(),
// 		).Should(BeTemporally("~", e2.EmittedAt(), time.Second))
// 	})
// })

var _ = Describe(`FromDomain`, func() {
	var (
		id, authID   = ulid.Make(), ulid.Make()
		desc         = vo.Description(`A New Description.`)
		title        = vo.Title(`My Title`)
		category     = vo.Category(`Story`)
		tags         = []vo.Tag{`Hot`}
		strTags      = []string{`Hot`}
		novelCreated = ne.NewNovelCreatedV2(
			id,
			authID,
			title,
			vo.Cover{},
			desc,
			category,
			tags,
		)
		event = eventrepo.Event{}
		a, _  = anypb.New(&novelv1.NovelCreated{
			AuthorId:    ulidpb.From(authID),
			Title:       string(title),
			Description: string(desc),
			Tags:        strTags,
			Category:    string(category),
		})
		expectPayload, _ = protojson.Marshal(a)
	)

	It(`can convert from novel.Event2`, func() {
		err := event.FromDomain(novelCreated)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(event.Payload).Should(BeEquivalentTo(expectPayload))
	})
})

package novel

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"gorm.io/gorm"

	"play-ddd/contents/app"
	"play-ddd/contents/domain"
	"play-ddd/contents/domain/novel"
	"play-ddd/contents/domain/novel/vo"
	"play-ddd/contents/infra/repository/pg/datatypes/ts"
	dtulid "play-ddd/contents/infra/repository/pg/datatypes/ulid"
)

var _ app.NovelRepo = novelRepo{}

type ID = dtulid.ULID

type novelRepo struct {
	db   *gorm.DB
	fact novel.Factory
}

// Append implements app.NovelRepo.
func (p novelRepo) Append(
	ctx context.Context,
	es ...domain.Event[novel.ID, novel.ID],
) error {
	return p.db.Save(es).Error
}

// Get implements app.NovelRepo.
func (p novelRepo) Get(
	ctx context.Context,
	id novel.ID,
) (novel novel.Novel, err error) {
	var n Novel
	if err = p.db.Take(&n, dtulid.ULID(id)).Error; err != nil {
		return novel, err
	}

	fmt.Printf("get novel: %+v\n", n)
	fmt.Printf("%+v\n", intoSlice[Chapter, vo.Chapter](n.TOC))
	fmt.Println(n.TOC[0].intoDomain())

	return p.fact.Restore(
		n.ID.Into(),
		n.AuthorID.Into(),
		n.Title,
		n.Category,
		n.Description,
		n.Tags,
		n.Status,
		url.URL{},
		intoSlice[Chapter, vo.Chapter](n.TOC),
		n.WordCount,
		n.UpdatedAt.Into(),
		n.CreatedAt.Into(),
	)
}

// Save implements app.NovelRepo.
func (nr novelRepo) Save(ctx context.Context, n novel.Novel) error {
	novel := fromDomain[novel.Novel, Novel](n)
	return nr.db.WithContext(ctx).Debug().Create(&novel).Error
}

func New(db *gorm.DB) novelRepo { return novelRepo{db: db} }

func Init(db *gorm.DB) error {
	if err := db.AutoMigrate(Novel{}); err != nil {
		return fmt.Errorf(`init repo: automigrate: %w`, err)
	}

	return nil
}

type Novel struct {
	ID          ID `gorm:"primaryKey"`
	CreatedAt   ts.Timestamp
	UpdatedAt   ts.Timestamp
	DeletedAt   ts.Timestamp
	AuthorID    ID
	Title       string
	Category    string
	Description string
	Tags        []string  `gorm:"serializer:json"`
	TOC         []Chapter `gorm:"serializer:json"`
	Status      int
	WordCount   int
}

func (n *Novel) fromDomain(novel novel.Novel) {
	tags := make([]string, len(novel.Tags()))
	for i, t := range novel.Tags() {
		tags[i] = string(t)
	}

	n.ID = dtulid.ULID(novel.ID())
	n.CreatedAt = ts.From(novel.CreatedAt())
	n.UpdatedAt = ts.From(novel.UpdatedAt())
	n.AuthorID = dtulid.ULID(novel.AuthorID())
	n.Title = string(novel.Title())
	n.Description = string(novel.Description())
	n.Category = string(novel.Category())
	n.Status = int(novel.Status())
	n.Tags = tags
	n.TOC = fromSlice[vo.Chapter, Chapter](novel.Chapters())
	n.WordCount = novel.WordCount()
}

type Chapter struct {
	Title      string
	Sequence   int
	WordCount  int
	UploadedAt int64
	UpdatedAt  int64
}

func (c *Chapter) fromDomain(vc vo.Chapter) {
	if c == nil {
		return
	}

	c.Title = vc.Title
	c.Sequence = vc.Sequence
	c.WordCount = vc.WordCount
	c.UpdatedAt = vc.UpdatedAt.UnixMilli()
	c.UploadedAt = vc.UploadedAt.UnixMilli()
}

func (c *Chapter) intoDomain() (vc vo.Chapter) {
	if c == nil {
		return vc
	}

	vc.Title = c.Title
	vc.Sequence = c.Sequence
	vc.WordCount = c.WordCount
	vc.UpdatedAt = time.UnixMilli(c.UpdatedAt)
	vc.UploadedAt = time.UnixMilli(c.UploadedAt)
	return vc
}

func fromDomain[A, M any, P fromPtr[A, M]](a A) (b M) {
	P(&b).fromDomain(a)
	return b
}

// fromPtr 约束类型必须是指向 B 的指针，且该指针类型实现了 fromDomain(A) 方法
type fromPtr[A any, M any] interface {
	fromDomain(A)
	*M
}

func fromSlice[A, M any, Mp fromPtr[A, M]](as []A) (bs []M) {
	if len(as) == 0 {
		return nil
	}
	bs = make([]M, len(as))
	for i := range as {
		Mp(&bs[i]).fromDomain(as[i])
	}
	return bs
}

// intoPtr constraints *M can convert to A
// 'A' stands for Aggregate in domain, and 'M' stands for Model mapping to repo.
type intoPtr[M any, A any] interface {
	intoDomain() A
	*M
}

func intoSlice[M, A any, P intoPtr[M, A]](ms []M) (as []A) {
	if len(ms) == 0 {
		return nil
	}

	as = make([]A, len(ms))
	for i := range as {
		as[i] = P(&ms[i]).intoDomain()
	}
	return as
}

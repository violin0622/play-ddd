package novel

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/samber/mo"
	"gorm.io/gorm"

	"play-ddd/contents/domain/novel"
	"play-ddd/contents/domain/novel/vo"
	"play-ddd/contents/infra/repository/pg/convert"
)

var _ novel.Query = novelRepo{}

// ByAuthorAndTitle implements novel.Query.
func (p novelRepo) ByAuthorAndTitle(
	ctx context.Context, authorID novel.AuthorID, title string,
) mo.Result[novel.Novel] {
	var n Novel
	err := p.db.
		WithContext(ctx).
		Where(Novel{AuthorID: ID(authorID)}).
		Where(nonDeleted).
		Take(&n).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return mo.Err[novel.Novel](novel.NewNotfoundError(err))
	}

	if err != nil {
		return mo.Err[novel.Novel](err)
	}

	dn, err := p.restore(n)
	if err != nil {
		return mo.Err[novel.Novel](err)
	}

	return mo.Ok(dn)
}

// func (p novelRepo) batchRestore(
// 	models ...Novel,
// ) (domains []novel.Novel, err error) {
// 	for _, model := range models {
// 		domain, err := p.restore(model)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		domains = append(domains, domain)
// 	}
//
// 	return domains, nil
// }

func (p novelRepo) restore(model Novel) (domain novel.Novel, err error) {
	chapters, err := convert.SliceIntoDomain[Chapter, vo.Chapter](model.TOC)
	if err != nil {
		return domain, fmt.Errorf(`restore: %w`, err)
	}

	domain, err = p.fact.Restore(
		model.ID.Into(),
		model.AuthorID.Into(),
		model.Title,
		model.Category,
		model.Description,
		model.Tags,
		model.Status,
		url.URL{},
		chapters,
		model.WordCount,
		model.UpdatedTs.Into(),
		model.CreatedTs.Into(),
	)

	return domain, err
}

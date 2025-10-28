package app

import (
	"context"

	"play-ddd/contents/app/command"
	"play-ddd/contents/domain/novel"
)

func (s *service) CreateNovel(
	ctx context.Context,
	cmd command.CreateNovel,
) (novel.ID, error) {
	r := s.novelFact.Create(
		ctx,
		cmd.AuthorID,
		cmd.Title,
		cmd.Desc,
		cmd.Category,
		cmd.Tags)

	if r.IsError() {
		return novel.ZeroID, r.Error()
	}

	return r.MustGet().ID(), nil
}

// func (s *service) UploadFirstChapter(
// 	ctx context.Context,
// ) error {
// 	return s.novelFact.UploadFirstChapter()
// }

func (s *service) UploadChapter(
	ctx context.Context,
	cmd command.UploadChapter,
) error {
	return s.novelRepo.Update(
		ctx,
		cmd.ID,
		func(ctx context.Context, n *novel.Novel) error {
			return n.UploadNewChapter(ctx, cmd.Title, cmd.WordCount)
		})
}

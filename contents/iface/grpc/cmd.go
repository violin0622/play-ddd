package grpc

import (
	"context"

	"play-ddd/contents/app"
	"play-ddd/contents/app/command"
	contv1 "play-ddd/proto/gen/go/contents/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
)

func NewCmdService(
	ah app.CommandHandler,
) cmdServiceServer {
	return cmdServiceServer{
		ah: ah,
	}
}

type cmdServiceServer struct {
	_  contv1.UnimplementedCmdServiceServer
	ah app.CommandHandler
}

func (s *cmdServiceServer) CreateNovel(
	ctx context.Context, req *contv1.CreateNovelRequest,
) (rep *contv1.CreateNovelResponse, err error,
) {
	cmd := command.CreateNovel{
		AuthorID: req.GetAuthorId().AsULID(),
		Title:    req.GetTitle(),
		Desc:     req.GetDesc(),
		Category: req.GetCategory(),
		Tags:     req.GetTags(),
	}

	id, err := s.ah.CreateNovel(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &contv1.CreateNovelResponse{Id: ulidpb.FromULID(id)}, nil
}

func (s *cmdServiceServer) UpdateNovel(
	ctx context.Context, req *contv1.UpdateNovelRequest) (
	rep *contv1.UpdateNovelResponse, err error,
) {
	cmd := command.UpdateNovelInfo{
		Tags: req.GetTags(),
		Desc: req.GetDesc(),
	}

	if err := s.ah.UpdateNovel(ctx, cmd); err != nil {
		return nil, err
	}

	return &contv1.UpdateNovelResponse{}, nil
}

func (s *cmdServiceServer) UploadChapter(
	ctx context.Context, req *contv1.UploadChapterRequest) (
	rep *contv1.UploadChapterResponse, err error,
) {
	id, err := s.ah.UploadChapter(ctx, command.UploadChapter{
		NovelID:      req.GetNovelId().AsULID(),
		Title:        req.GetTitle(),
		MainContent:  req.GetMainContent(),
		ExtraContent: req.GetExtraContent(),
	})
	if err != nil {
		return nil, err
	}

	return &contv1.UploadChapterResponse{Id: ulidpb.FromULID(id)}, nil
}

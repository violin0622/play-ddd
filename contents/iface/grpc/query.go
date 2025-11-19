package grpc

import (
	"play-ddd/contents/app"
	contv1 "play-ddd/proto/gen/go/contents/v1"
)

type Query struct {
	contv1.UnimplementedQueryServiceServer
}

func NewQueryService(
	ah app.QueryHandler,
) Query {
	return Query{}
}

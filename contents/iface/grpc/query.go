package grpc

import contv1 "play-ddd/proto/gen/go/contents/v1"

type Query struct {
	_ contv1.UnimplementedQueryServiceServer
}

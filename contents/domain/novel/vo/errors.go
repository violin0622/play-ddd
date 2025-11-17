package vo

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNoContent = status.New(codes.OutOfRange, `no content`).Err()
var ErrInvalidChapterSequence = status.New(codes.OutOfRange, `chapter sequence out of range`).Err()

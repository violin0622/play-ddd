package nop

import (
	"golang.org/x/time/rate"

	requestlimter "play-ddd/contents/infra/server/requestlimiter"
)

type nop struct{}

func New() nop { return nop{} }

func (nop) Set(r rate.Limit, b int) {}

func (nop) Request(
	requestlimter.RequestInfo,
) (r requestlimter.Result) {
	return r
}

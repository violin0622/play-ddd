package ratelimiter

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	// Allow is alias for AllowN(time.Now(),1)
	Allow() bool

	// Wait is alias for WaitN(ctx, 1)
	Wait(context.Context) error

	AllowN(time.Time, int) bool

	WaitN(context.Context, int) error
}

var NewTokenBucket = rate.NewLimiter

type composite struct {
	limiters []RateLimiter
}

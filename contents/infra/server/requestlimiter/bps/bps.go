package bps

import (
	"time"

	"golang.org/x/time/rate"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"

	"play-ddd/contents/infra/server/requestlimiter"
)

type bps struct {
	rl  rate.Limiter
	api string
}

func (l *bps) Set(r rate.Limit, b int) {
	l.rl.SetBurst(b)
	l.rl.SetLimit(r)
}

func (l *bps) Request(ri requestlimter.RequestInfo) (r requestlimter.Result) {
	if l.rl.AllowN(time.Now(), int(ri.PS.InBytes)) {
		return r
	}

	r.Deny = true
	r.Details = append(r.Details, &epb.QuotaFailure{
		Violations: []*epb.QuotaFailure_Violation{{
			Subject:     `Total BPS`,
			Description: `Server has overall BPS upper limit.`,
			ApiService:  l.api,
			QuotaValue:  int64(l.rl.Limit()),
		}},
	})

	return r
}

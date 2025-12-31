package chapter

import (
	"github.com/oklog/ulid/v2"

	"play-ddd/common"
	"play-ddd/contents/domain"
)

type (
	ID        = ulid.ULID
	Event     = common.Event[ID, ID]
	EventRepo = domain.EventRepo[ID, ID]
	Aggregate = common.Aggregate[ID]
	Repo      = common.AggregateRepo[ID, Chapter]
)

var ZeroID = ulid.Zero

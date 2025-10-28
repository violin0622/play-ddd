package novel

import (
	"play-ddd/contents/domain"

	"github.com/oklog/ulid/v2"
)

type (
	ID        = ulid.ULID
	Event     = domain.Event[ID, ID]
	EventRepo = domain.EventRepo[ID, ID]
	Aggregate = domain.Aggregate[ID]
	Repo      = domain.AggregateRepo[ID, Novel]
)

var ZeroID = ulid.Zero

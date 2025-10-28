package novel

import (
	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain"
)

type (
	ID        = ulid.ULID
	Event     = domain.Event[ID, ID]
	EventRepo = domain.EventRepo[ID, ID]
	Aggregate = domain.Aggregate[ID]
	Repo      = domain.AggregateRepo[ID, Novel]
)

var ZeroID = ulid.Zero

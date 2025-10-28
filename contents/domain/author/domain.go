package author

import (
	"play-ddd/contents/domain"

	"github.com/oklog/ulid/v2"
)

type (
	Event     = domain.Event[ulid.ULID, ulid.ULID]
	EventRepo = domain.EventRepo[ulid.ULID, ulid.ULID]
	Aggregate = domain.Aggregate[ulid.ULID]
)

package author

import (
	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain"
)

type (
	Event     = domain.Event[ulid.ULID, ulid.ULID]
	EventRepo = domain.EventRepo[ulid.ULID, ulid.ULID]
	Aggregate = domain.Aggregate[ulid.ULID]
)

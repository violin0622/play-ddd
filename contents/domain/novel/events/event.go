package events

import (
	"fmt"

	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain"
)

type (
	ID    = ulid.ULID
	Event = domain.Event[ID, ID]
)

func formatEvent(e Event) string {
	return fmt.Sprintf(`%s[%s]: %s[%s] @%s`,
		e.AggKind(), e.AggID(), e.Kind(), e.ID(), e.EmittedAt())
}

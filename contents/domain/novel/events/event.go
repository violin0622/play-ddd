package events

import (
	"fmt"

	"play-ddd/contents/domain"

	"github.com/oklog/ulid/v2"
)

type (
	ID    = ulid.ULID
	Event = domain.Event[ID, ID]
)

func formatEvent(e Event) string {
	return fmt.Sprintf(`%s[%s]: %s[%s] @%s`,
		e.AggKind(), e.AggID(), e.Kind(), e.ID(), e.EmittedAt())
}

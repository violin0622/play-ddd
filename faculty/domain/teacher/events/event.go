package events

import (
	"fmt"

	"github.com/oklog/ulid/v2"

	"play-ddd/faculty/domain"
)

type Event = domain.Event[ulid.ULID, ulid.ULID]

func formatEvent(e Event) string {
	return fmt.Sprintf(`%s[%s]: %s[%s] @%s`,
		e.AggKind(), e.AggID(), e.Kind(), e.ID(), e.EmittedAt())
}
